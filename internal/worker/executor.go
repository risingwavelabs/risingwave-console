package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/robfig/cron/v3"
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

func (e *Executor) ExecuteAutoBackup(ctx context.Context, spec apigen.TaskSpecAutoBackup) error {
	cluster, err := e.model.GetClusterByID(ctx, spec.ClusterID)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster")
	}
	tz, err := e.model.GetTimeZone(ctx, cluster.OrganizationID)
	if err != nil {
		return errors.Wrap(err, "failed to get timezone")
	}
	config, err := e.model.GetAutoBackupConfig(ctx, cluster.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get auto backup config")
	}

	// schedule the next task
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cron, err := parser.Parse(fmt.Sprintf("TZ=%s %s", tz, config.CronExpression))
	if err != nil {
		return errors.Wrapf(err, "failed to parse cron expression: %s", config.CronExpression)
	}
	next := cron.Next(time.Now())
	if _, err := e.model.CreateTask(ctx, querier.CreateTaskParams{
		Spec: apigen.TaskSpec{
			Type:       apigen.AutoBackup,
			AutoBackup: &spec,
		},
		Status:    string(apigen.Pending),
		StartedAt: &next,
	}); err != nil {
		return errors.Wrap(err, "failed to create task")
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

	// record the snapshot ID
	if err := e.model.CreateSnapshot(ctx, querier.CreateSnapshotParams{
		ClusterID:  cluster.ID,
		SnapshotID: snapshotID,
	}); err != nil {
		return errors.Wrap(err, "failed to create snapshot")
	}

	// clean up old snapshots
	snapshots, err := e.model.ListSnapshots(ctx, cluster.ID)
	if err != nil {
		return errors.Wrap(err, "failed to list snapshots")
	}
	for _, snapshot := range snapshots[config.KeepLast:] {
		log.Info("deleting snapshot", zap.Int64("snapshot_id", snapshot.SnapshotID), zap.Time("created_at", snapshot.CreatedAt))
		if err := conn.DeleteSnapshot(ctx, snapshot.SnapshotID); err != nil {
			return errors.Wrap(err, "failed to delete snapshot")
		}
		if err := e.model.DeleteSnapshot(ctx, querier.DeleteSnapshotParams{
			ClusterID:  cluster.ID,
			SnapshotID: snapshot.SnapshotID,
		}); err != nil {
			return errors.Wrap(err, "failed to delete snapshot")
		}
	}
	return nil
}

func (e *Executor) ExecuteAutoDiagnostic(ctx context.Context, spec apigen.TaskSpecAutoDiagnostic) error {
	return nil
}
