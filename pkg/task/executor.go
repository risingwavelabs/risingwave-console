package task

import (
	"context"
	"fmt"
	"time"

	anchor_apigen "github.com/cloudcarver/anchor/pkg/zgen/apigen"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/risingwave-console/pkg/conn/http"
	"github.com/risingwavelabs/risingwave-console/pkg/conn/meta"
	"github.com/risingwavelabs/risingwave-console/pkg/logger"
	"github.com/risingwavelabs/risingwave-console/pkg/utils"
	"github.com/risingwavelabs/risingwave-console/pkg/zcore/model"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/querier"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/taskgen"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("task runner")

type TaskExecutor struct {
	risectlm meta.RisectlManagerInterface

	metahttp http.MetaHttpManagerInterface

	model model.ModelInterface

	taskRunner taskgen.TaskRunner

	now func() time.Time
}

func NewTaskExecutor(taskRunner taskgen.TaskRunner, model model.ModelInterface, risectlm meta.RisectlManagerInterface, metahttp http.MetaHttpManagerInterface) taskgen.ExecutorInterface {
	return &TaskExecutor{
		taskRunner: taskRunner,
		model:      model,
		risectlm:   risectlm,
		now:        time.Now,
		metahttp:   metahttp,
	}
}

func (e *TaskExecutor) ExecuteAutoBackup(ctx context.Context, params *taskgen.AutoBackupParameters) error {
	cluster, err := e.model.GetClusterByID(ctx, params.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}

	// run meta backup
	conn, err := e.risectlm.NewConn(ctx, cluster.Version, cluster.Host, cluster.MetaPort)
	if err != nil {
		return errors.Wrap(err, "failed to get risectl connection")
	}
	snapshotID, err := conn.MetaBackup(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get meta backup")
	}
	log.Info(
		"auto backup task created",
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("snapshot_id", fmt.Sprintf("%d", snapshotID)),
	)

	// record the snapshot ID
	if err := e.model.CreateClusterSnapshot(ctx, querier.CreateClusterSnapshotParams{
		ClusterID:  cluster.ID,
		SnapshotID: snapshotID,
		Name:       fmt.Sprintf("auto-backup-%s", e.now().Format("2006-01-02-15-04-05")),
	}); err != nil {
		return errors.Wrap(err, "failed to create snapshot")
	}

	// create a task to delete the snapshot after the retention duration
	retentionDuration, err := utils.ParseDuration(params.RetentionDuration)
	if err != nil {
		return errors.Wrap(err, "failed to parse retention duration")
	}
	taskID, err := e.taskRunner.RunDeleteSnapshot(ctx, &taskgen.DeleteSnapshotParameters{
		ClusterID:  cluster.ID,
		SnapshotID: snapshotID,
	}, func(task *anchor_apigen.Task) error {
		task.StartedAt = utils.Ptr(e.now().Add(retentionDuration))
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	log.Info(
		"auto delete snapshot task created",
		zap.Int32("task_id", taskID),
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("snapshot_id", fmt.Sprintf("%d", snapshotID)),
		zap.String("retention_duration", params.RetentionDuration),
	)

	return nil
}

func (e *TaskExecutor) ExecuteAutoDiagnostic(ctx context.Context, params *taskgen.AutoDiagnosticParameters) error {
	cluster, err := e.model.GetClusterByID(ctx, params.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}

	// run diagnostics
	content, err := e.metahttp.GetDiagnose(ctx, fmt.Sprintf("http://%s:%d", cluster.Host, cluster.HttpPort))
	if err != nil {
		return errors.Wrap(err, "failed to get diagnose")
	}
	diag, err := e.model.CreateClusterDiagnostic(ctx, querier.CreateClusterDiagnosticParams{
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
	retentionDuration, err := utils.ParseDuration(params.RetentionDuration)
	if err != nil {
		return errors.Wrap(err, "failed to parse retention duration")
	}
	taskID, err := e.taskRunner.RunDeleteClusterDiagnostic(ctx, &taskgen.DeleteClusterDiagnosticParameters{
		ClusterID:    cluster.ID,
		DiagnosticID: diag.ID,
	}, func(task *anchor_apigen.Task) error {
		task.StartedAt = utils.Ptr(e.now().Add(retentionDuration))
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	log.Info(
		"auto delete cluster diagnostic task created",
		zap.Int32("task_id", taskID),
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("diagnostic_id", fmt.Sprintf("%d", diag.ID)),
		zap.String("retention_duration", params.RetentionDuration),
	)

	return nil
}

func (e *TaskExecutor) ExecuteDeleteClusterDiagnostic(ctx context.Context, params *taskgen.DeleteClusterDiagnosticParameters) error {
	if err := e.model.DeleteClusterDiagnostic(ctx, params.DiagnosticID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("cluster diagnostic not found, skipping delete", zap.Int32("diagnostic_id", params.DiagnosticID))
			return nil
		}
		return errors.Wrap(err, "failed to delete cluster diagnostic")
	}
	return nil
}

func (e *TaskExecutor) ExecuteDeleteSnapshot(ctx context.Context, params *taskgen.DeleteSnapshotParameters) error {
	cluster, err := e.model.GetClusterByID(ctx, params.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}

	conn, err := e.risectlm.NewConn(ctx, cluster.Version, cluster.Host, cluster.MetaPort)
	if err != nil {
		return errors.Wrap(err, "failed to get risectl connection")
	}

	if err := conn.DeleteSnapshot(ctx, params.SnapshotID); err != nil {
		return errors.Wrapf(err, "failed to delete snapshot in risingwave, snapshot_id: %d", params.SnapshotID)
	}

	if err := e.model.DeleteClusterSnapshot(ctx, querier.DeleteClusterSnapshotParams{
		ClusterID:  cluster.ID,
		SnapshotID: params.SnapshotID,
	}); err != nil {
		return errors.Wrapf(err, "failed to delete snapshot in database, cluster_name: %s, cluster_id: %d, snapshot_id: %d", cluster.Name, cluster.ID, params.SnapshotID)
	}

	return nil
}
