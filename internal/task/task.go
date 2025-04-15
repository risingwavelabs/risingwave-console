package task

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/conn/http"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/worker"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("worker")

type TaskHandler struct {
	model     model.ModelInterface
	risectlm  meta.RisectlManagerInterface
	taskstore TaskStoreInterface
	metahttp  http.MetaHttpManagerInterface
	now       func() time.Time
}

type TaskStoreInterface interface {
	CreateTask(ctx context.Context, params CreateTaskParams) (int32, error)
}

type TaskStore struct {
	model model.ModelInterface
}

func NewTaskStore(model model.ModelInterface) TaskStoreInterface {
	return &TaskStore{model: model}
}

func NewTaskHandler(model model.ModelInterface, risectlm meta.RisectlManagerInterface, taskstore TaskStoreInterface, metahttp http.MetaHttpManagerInterface) worker.TaskHandler {
	return &TaskHandler{
		model:     model,
		risectlm:  risectlm,
		taskstore: taskstore,
		now:       time.Now,
		metahttp:  metahttp,
	}
}

type CreateTaskParams struct {
	OrgID                *int32
	Spec                 apigen.TaskSpec
	StartedAt            *time.Time
	AlwaysRetryOnFailure bool
	RetryInterval        string
}

func (s *TaskStore) CreateTask(ctx context.Context, params CreateTaskParams) (int32, error) {
	attributes := apigen.TaskAttributes{
		OrgID: params.OrgID,
	}
	if params.AlwaysRetryOnFailure {
		_, err := time.ParseDuration(params.RetryInterval)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse retry interval")
		}
		attributes.RetryPolicy = &apigen.TaskRetryPolicy{
			AlwaysRetryOnFailure: utils.Ptr(true),
			Interval:             params.RetryInterval,
		}
	}
	task, err := s.model.CreateTask(ctx, querier.CreateTaskParams{
		Attributes: attributes,
		Spec:       params.Spec,
		StartedAt:  params.StartedAt,
		Status:     string(apigen.Pending),
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to create task")
	}
	return task.ID, nil
}

func (e *TaskHandler) HandleTask(ctx context.Context, task apigen.Task) error {
	switch task.Spec.Type {
	case apigen.AutoBackup:
		if task.Spec.AutoBackup == nil {
			return fmt.Errorf("auto backup spec is nil")
		}
		return e.ExecuteAutoBackup(ctx, *task.Spec.AutoBackup)
	case apigen.AutoDiagnostic:
		if task.Spec.AutoDiagnostic == nil {
			return fmt.Errorf("auto diagnostic spec is nil")
		}
		return e.ExecuteAutoDiagnostic(ctx, *task.Spec.AutoDiagnostic)
	case apigen.DeleteClusterDiagnostic:
		if task.Spec.DeleteClusterDiagnostic == nil {
			return fmt.Errorf("delete cluster diagnostic spec is nil")
		}
		return e.ExecuteDeleteClusterDiagnostic(ctx, *task.Spec.DeleteClusterDiagnostic)
	case apigen.DeleteSnapshot:
		if task.Spec.DeleteSnapshot == nil {
			return fmt.Errorf("delete snapshot spec is nil")
		}
		return e.ExecuteDeleteSnapshot(ctx, *task.Spec.DeleteSnapshot)
	case apigen.DeleteOpaqueKey:
		if task.Spec.DeleteOpaqueKey == nil {
			return fmt.Errorf("delete opaque key spec is nil")
		}
		return e.ExecuteDeleteOpaqueKey(ctx, *task.Spec.DeleteOpaqueKey)
	default:
		return fmt.Errorf("unknown task type: %s", task.Spec.Type)
	}
}

func (e *TaskHandler) ExecuteAutoBackup(ctx context.Context, spec apigen.TaskSpecAutoBackup) error {
	cluster, err := e.model.GetClusterByID(ctx, spec.ClusterID)
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
	if err := e.model.CreateSnapshot(ctx, querier.CreateSnapshotParams{
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
	taskID, err := e.taskstore.CreateTask(ctx, CreateTaskParams{
		OrgID: &cluster.OrganizationID,
		Spec: apigen.TaskSpec{
			Type: apigen.DeleteSnapshot,
			DeleteSnapshot: &apigen.TaskSpecDeleteSnapshot{
				ClusterID:  cluster.ID,
				SnapshotID: snapshotID,
			},
		},
		StartedAt:            utils.Ptr(e.now().Add(retentionDuration)),
		AlwaysRetryOnFailure: true,
		RetryInterval:        "10m",
	})
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

func (e *TaskHandler) ExecuteAutoDiagnostic(ctx context.Context, spec apigen.TaskSpecAutoDiagnostic) error {
	cluster, err := e.model.GetClusterByID(ctx, spec.ClusterID)
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
	retentionDuration, err := utils.ParseDuration(spec.RetentionDuration)
	if err != nil {
		return errors.Wrap(err, "failed to parse retention duration")
	}
	taskID, err := e.taskstore.CreateTask(ctx, CreateTaskParams{
		OrgID: &cluster.OrganizationID,
		Spec: apigen.TaskSpec{
			Type: apigen.DeleteClusterDiagnostic,
			DeleteClusterDiagnostic: &apigen.TaskSpecDeleteClusterDiagnostic{
				ClusterID:    cluster.ID,
				DiagnosticID: diag.ID,
			},
		},
		StartedAt:            utils.Ptr(e.now().Add(retentionDuration)),
		AlwaysRetryOnFailure: true,
		RetryInterval:        "10m",
	})
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

func (e *TaskHandler) ExecuteDeleteClusterDiagnostic(ctx context.Context, spec apigen.TaskSpecDeleteClusterDiagnostic) error {
	if err := e.model.DeleteClusterDiagnostic(ctx, spec.DiagnosticID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("cluster diagnostic not found, skipping delete", zap.Int32("diagnostic_id", spec.DiagnosticID))
			return nil
		}
		return errors.Wrap(err, "failed to delete cluster diagnostic")
	}
	return nil
}

func (e *TaskHandler) ExecuteDeleteSnapshot(ctx context.Context, spec apigen.TaskSpecDeleteSnapshot) error {
	cluster, err := e.model.GetClusterByID(ctx, spec.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}

	conn, err := e.risectlm.NewConn(ctx, cluster.Version, cluster.Host, cluster.MetaPort)
	if err != nil {
		return errors.Wrap(err, "failed to get risectl connection")
	}

	if err := conn.DeleteSnapshot(ctx, spec.SnapshotID); err != nil {
		return errors.Wrapf(err, "failed to delete snapshot in risingwave, snapshot_id: %d", spec.SnapshotID)
	}

	if err := e.model.DeleteSnapshot(ctx, querier.DeleteSnapshotParams{
		ClusterID:  cluster.ID,
		SnapshotID: spec.SnapshotID,
	}); err != nil {
		return errors.Wrapf(err, "failed to delete snapshot in database, cluster_name: %s, cluster_id: %d, snapshot_id: %d", cluster.Name, cluster.ID, spec.SnapshotID)
	}

	return nil
}

func (e *TaskHandler) ExecuteDeleteOpaqueKey(ctx context.Context, spec apigen.TaskSpecDeleteOpaqueKey) error {
	if err := e.model.DeleteOpaqueKey(ctx, spec.KeyID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("opaque key not found, skipping delete", zap.Int64("key_id", spec.KeyID))
			return nil
		}
		return errors.Wrapf(err, "failed to delete opaque key in database, key_id: %d", spec.KeyID)
	}
	return nil
}
