package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/conn/http"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"go.uber.org/zap"
)

type executorGetter = func(model model.ModelInterface, risectlm *meta.RisectlManager) ExecutorInterface

type ExecutorInterface interface {
	ExecuteAutoBackup(ctx context.Context, spec apigen.TaskSpecAutoBackup) error
	ExecuteAutoDiagnostic(ctx context.Context, spec apigen.TaskSpecAutoDiagnostic) error
}

type Executor struct {
	model    model.ModelInterface
	risectlm *meta.RisectlManager
}

func newExecutor(model model.ModelInterface, risectlm *meta.RisectlManager) ExecutorInterface {
	return &Executor{
		model:    model,
		risectlm: risectlm,
	}
}

func (e *Executor) createTask(ctx context.Context, orgID *int32, spec apigen.TaskSpec, startedAt *time.Time) (*querier.Task, error) {
	return e.model.CreateTask(ctx, querier.CreateTaskParams{
		Attributes: apigen.TaskAttributes{
			OrgID: orgID,
		},
		Spec:      spec,
		StartedAt: startedAt,
		Status:    string(apigen.Pending),
	})
}

func (e *Executor) ExecuteAutoBackup(ctx context.Context, spec apigen.TaskSpecAutoBackup) error {
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
	retentionDuration, err := time.ParseDuration(spec.RetentionDuration)
	if err != nil {
		return errors.Wrap(err, "failed to parse retention duration")
	}
	task, err := e.createTask(ctx, &cluster.OrganizationID, apigen.TaskSpec{
		Type: apigen.DeleteSnapshot,
		DeleteSnapshot: &apigen.TaskSpecDeleteSnapshot{
			ClusterID:  cluster.ID,
			SnapshotID: snapshotID,
		},
	}, utils.Ptr(time.Now().Add(retentionDuration)))
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	log.Info(
		"auto delete snapshot task created",
		zap.Int32("task_id", task.ID),
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("snapshot_id", fmt.Sprintf("%d", snapshotID)),
		zap.String("retention_duration", spec.RetentionDuration),
	)

	return nil
}

func (e *Executor) ExecuteAutoDiagnostic(ctx context.Context, spec apigen.TaskSpecAutoDiagnostic) error {
	cluster, err := e.model.GetClusterByID(ctx, spec.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}

	// run diagnostics
	conn := http.NewMetaHttpConnection(fmt.Sprintf("http://%s:%d", cluster.Host, cluster.HttpPort))
	content, err := conn.GetDiagnose(ctx)
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
	retentionDuration, err := time.ParseDuration(spec.RetentionDuration)
	if err != nil {
		return errors.Wrap(err, "failed to parse retention duration")
	}
	task, err := e.createTask(ctx, &cluster.OrganizationID, apigen.TaskSpec{
		Type: apigen.DeleteClusterDiagnostic,
		DeleteClusterDiagnostic: &apigen.TaskDeleteClusterDiagnostic{
			ClusterID:    cluster.ID,
			DiagnosticID: diag.ID,
		},
	}, utils.Ptr(time.Now().Add(retentionDuration)))
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}
	log.Info(
		"auto delete cluster diagnostic task created",
		zap.Int32("task_id", task.ID),
		zap.String("cluster_id", fmt.Sprintf("%d", cluster.ID)),
		zap.String("diagnostic_id", fmt.Sprintf("%d", diag.ID)),
		zap.String("retention_duration", spec.RetentionDuration),
	)

	return nil
}
