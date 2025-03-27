package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/robfig/cron/v3"
)

type executorGetter = func(model model.ModelInterface) ExecutorInterface

type ExecutorInterface interface {
	ExecuteAutoBackup(ctx context.Context, spec apigen.TaskSpecAutoBackup) error
	ExecuteAutoDiagnostic(ctx context.Context, spec apigen.TaskSpecAutoDiagnostic) error
}

type Executor struct {
	model model.ModelInterface
}

func newExecutor(model model.ModelInterface) ExecutorInterface {
	return &Executor{
		model: model,
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

	// TODO: meta backup and delete older snapshots

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
	return nil
}

func (e *Executor) ExecuteAutoDiagnostic(ctx context.Context, spec apigen.TaskSpecAutoDiagnostic) error {
	return nil
}
