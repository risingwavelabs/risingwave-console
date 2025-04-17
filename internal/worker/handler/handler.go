package handler

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/conn/http"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
	"github.com/risingwavelabs/wavekit/internal/task"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/worker"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("task handler")

type TaskHandler struct {
	risectlm meta.RisectlManagerInterface

	taskstore task.TaskStoreInterface

	metahttp http.MetaHttpManagerInterface

	now func() time.Time
}

func NewTaskHandler(risectlm meta.RisectlManagerInterface, taskstore task.TaskStoreInterface, metahttp http.MetaHttpManagerInterface) worker.TaskHandler {
	return &TaskHandler{
		risectlm:  risectlm,
		taskstore: taskstore,
		now:       time.Now,
		metahttp:  metahttp,
	}
}

func (e *TaskHandler) HandleTask(c *modelctx.ModelCtx, task apigen.Task) error {
	switch task.Spec.Type {
	case apigen.AutoBackup:
		if task.Spec.AutoBackup == nil {
			return fmt.Errorf("auto backup spec is nil")
		}
		return e.ExecuteAutoBackup(c, *task.Spec.AutoBackup)
	case apigen.AutoDiagnostic:
		if task.Spec.AutoDiagnostic == nil {
			return fmt.Errorf("auto diagnostic spec is nil")
		}
		return e.ExecuteAutoDiagnostic(c, *task.Spec.AutoDiagnostic)
	case apigen.DeleteClusterDiagnostic:
		if task.Spec.DeleteClusterDiagnostic == nil {
			return fmt.Errorf("delete cluster diagnostic spec is nil")
		}
		return e.ExecuteDeleteClusterDiagnostic(c, *task.Spec.DeleteClusterDiagnostic)
	case apigen.DeleteSnapshot:
		if task.Spec.DeleteSnapshot == nil {
			return fmt.Errorf("delete snapshot spec is nil")
		}
		return e.ExecuteDeleteSnapshot(c, *task.Spec.DeleteSnapshot)
	case apigen.DeleteOpaqueKey:
		if task.Spec.DeleteOpaqueKey == nil {
			return fmt.Errorf("delete opaque key spec is nil")
		}
		return e.ExecuteDeleteOpaqueKey(c, *task.Spec.DeleteOpaqueKey)
	default:
		return fmt.Errorf("unknown task type: %s", task.Spec.Type)
	}
}

func (e *TaskHandler) ExecuteAutoBackup(c *modelctx.ModelCtx, spec apigen.TaskSpecAutoBackup) error {
	cluster, err := c.GetClusterByID(c, spec.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}

	// run meta backup
	conn, err := e.risectlm.NewConn(c, cluster.Version, cluster.Host, cluster.MetaPort)
	if err != nil {
		return errors.Wrap(err, "failed to get risectl connection")
	}
	snapshotID, err := conn.MetaBackup(c)
	if err != nil {
		return errors.Wrap(err, "failed to get meta backup")
	}
	log.Info(
		"auto backup task created",
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("snapshot_id", fmt.Sprintf("%d", snapshotID)),
	)

	// record the snapshot ID
	if err := c.CreateSnapshot(c, querier.CreateSnapshotParams{
		ClusterID:  cluster.ID,
		SnapshotID: snapshotID,
	}); err != nil {
		return errors.Wrap(err, "failed to create snapshot")
	}

	// create a task to delete the snapshot after the retention duration
	retentionDuration, err := utils.ParseDuration(spec.RetentionDuration)
	if err != nil {
		return errors.Wrap(err, "failed to parse retention duration")
	}
	taskID, err := e.taskstore.PushTask(
		c,
		apigen.TaskSpec{
			Type: apigen.DeleteSnapshot,
			DeleteSnapshot: &apigen.TaskSpecDeleteSnapshot{
				ClusterID:  cluster.ID,
				SnapshotID: snapshotID,
			},
		},
		task.OrgID(cluster.OrganizationID),
		task.AlwaysRetryOnFailure("10m"),
		task.StartedAt(e.now().Add(retentionDuration)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	log.Info(
		"auto delete snapshot task created",
		zap.Int32("task_id", taskID),
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("snapshot_id", fmt.Sprintf("%d", snapshotID)),
		zap.String("retention_duration", spec.RetentionDuration),
	)

	return nil
}

func (e *TaskHandler) ExecuteAutoDiagnostic(c *modelctx.ModelCtx, spec apigen.TaskSpecAutoDiagnostic) error {
	cluster, err := c.GetClusterByID(c, spec.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}

	// run diagnostics
	content, err := e.metahttp.GetDiagnose(c, fmt.Sprintf("http://%s:%d", cluster.Host, cluster.HttpPort))
	if err != nil {
		return errors.Wrap(err, "failed to get diagnose")
	}
	diag, err := c.CreateClusterDiagnostic(c, querier.CreateClusterDiagnosticParams{
		ClusterID: cluster.ID,
		Content:   content,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create cluster diagnostic")
	}
	log.Info(
		"cluster diagnostic created",
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("diagnostic_id", fmt.Sprintf("%d", diag.ID)),
	)

	// create a task to delete the cluster diagnostic after the retention duration
	retentionDuration, err := utils.ParseDuration(spec.RetentionDuration)
	if err != nil {
		return errors.Wrap(err, "failed to parse retention duration")
	}
	taskID, err := e.taskstore.PushTask(
		c,
		apigen.TaskSpec{
			Type: apigen.DeleteClusterDiagnostic,
			DeleteClusterDiagnostic: &apigen.TaskSpecDeleteClusterDiagnostic{
				ClusterID:    cluster.ID,
				DiagnosticID: diag.ID,
			},
		},
		task.OrgID(cluster.OrganizationID),
		task.AlwaysRetryOnFailure("10m"),
		task.StartedAt(e.now().Add(retentionDuration)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}
	log.Info(
		"auto delete cluster diagnostic task created",
		zap.Int32("task_id", taskID),
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("diagnostic_id", fmt.Sprintf("%d", diag.ID)),
		zap.String("retention_duration", spec.RetentionDuration),
	)

	return nil
}

func (e *TaskHandler) ExecuteDeleteClusterDiagnostic(c *modelctx.ModelCtx, spec apigen.TaskSpecDeleteClusterDiagnostic) error {
	if err := c.DeleteClusterDiagnostic(c, spec.DiagnosticID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("cluster diagnostic not found, skipping delete", zap.Int32("diagnostic_id", spec.DiagnosticID))
			return nil
		}
		return errors.Wrap(err, "failed to delete cluster diagnostic")
	}
	return nil
}

func (e *TaskHandler) ExecuteDeleteSnapshot(c *modelctx.ModelCtx, spec apigen.TaskSpecDeleteSnapshot) error {
	cluster, err := c.GetClusterByID(c, spec.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}

	conn, err := e.risectlm.NewConn(c, cluster.Version, cluster.Host, cluster.MetaPort)
	if err != nil {
		return errors.Wrap(err, "failed to get risectl connection")
	}

	if err := conn.DeleteSnapshot(c, spec.SnapshotID); err != nil {
		return errors.Wrapf(err, "failed to delete snapshot in risingwave, snapshot_id: %d", spec.SnapshotID)
	}

	if err := c.DeleteSnapshot(c, querier.DeleteSnapshotParams{
		ClusterID:  cluster.ID,
		SnapshotID: spec.SnapshotID,
	}); err != nil {
		return errors.Wrapf(err, "failed to delete snapshot in database, cluster_name: %s, cluster_id: %d, snapshot_id: %d", cluster.Name, cluster.ID, spec.SnapshotID)
	}

	return nil
}

func (e *TaskHandler) ExecuteDeleteOpaqueKey(c *modelctx.ModelCtx, spec apigen.TaskSpecDeleteOpaqueKey) error {
	if err := c.DeleteOpaqueKey(c, spec.KeyID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("opaque key not found, skipping delete", zap.Int64("key_id", spec.KeyID))
			return nil
		}
		return errors.Wrapf(err, "failed to delete opaque key in database, key_id: %d", spec.KeyID)
	}
	return nil
}
