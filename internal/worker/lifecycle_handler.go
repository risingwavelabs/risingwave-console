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

type LifeCycleHandlerGetter = func(txm model.ModelInterface) (TaskLifeCycleHandlerInterface, error)

type TaskLifeCycleHandlerInterface interface {
	HandleAttributes(ctx context.Context, task apigen.Task) error
	HandleFailed(ctx context.Context, task apigen.Task, err error) error
	HandleCompleted(ctx context.Context, task apigen.Task) error
}

type TaskLifeCycleHandler struct {
	txm model.ModelInterface
	now func() time.Time
}

func newTaskLifeCycleHandler(txm model.ModelInterface) (TaskLifeCycleHandlerInterface, error) {
	if !txm.InTransaction() {
		return nil, errors.Errorf("task life cycle handler must be used within a transaction")
	}
	return &TaskLifeCycleHandler{
		txm: txm,
		now: time.Now,
	}, nil
}

func (a *TaskLifeCycleHandler) HandleAttributes(ctx context.Context, task apigen.Task) error {
	if task.Attributes.Cronjob != nil {
		return a.handleCronjob(ctx, task)
	}
	return nil
}

func (a *TaskLifeCycleHandler) HandleFailed(ctx context.Context, task apigen.Task, err error) error {
	// the event must be reported
	if _, err := a.txm.InsertEvent(ctx, apigen.EventSpec{
		Type: apigen.TaskError,
		TaskError: &apigen.EventTaskError{
			TaskID: task.ID,
			Error:  err.Error(),
		},
	}); err != nil {
		return errors.Wrap(err, "insert task error event")
	}

	// cronjob should be run again anyway, no need to update status
	if task.Attributes.Cronjob != nil {
		return nil
	}

	// update task status
	if err := a.txm.UpdateTaskStatus(ctx, querier.UpdateTaskStatusParams{
		ID:     task.ID,
		Status: string(apigen.Failed),
	}); err != nil {
		return errors.Wrap(err, "update task status")
	}
	return nil
}

func (a *TaskLifeCycleHandler) HandleCompleted(ctx context.Context, task apigen.Task) error {
	if err := a.txm.UpdateTaskStatus(ctx, querier.UpdateTaskStatusParams{
		ID:     task.ID,
		Status: string(apigen.Completed),
	}); err != nil {
		return errors.Wrap(err, "update task status")
	}
	return nil
}

func (a *TaskLifeCycleHandler) handleCronjob(ctx context.Context, task apigen.Task) error {
	if task.Attributes.Cronjob == nil {
		return errors.Errorf("")
	}
	cronjob := task.Attributes.Cronjob

	tz := "UTC"
	if task.Attributes.OrgID != nil {
		org, err := a.txm.GetOrganization(ctx, *task.Attributes.OrgID)
		if err != nil {
			return errors.Wrapf(err, "get organization by id: %d", *task.Attributes.OrgID)
		}
		tz = org.Timezone
	}

	// schedule next task
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cron, err := parser.Parse(fmt.Sprintf("TZ=%s %s", tz, cronjob.CronExpression))
	if err != nil {
		return errors.Wrapf(err, "failed to parse cron expression: %s", cronjob.CronExpression)
	}
	nextTime := cron.Next(a.now())
	if err := a.txm.UpdateTaskStartedAt(ctx, querier.UpdateTaskStartedAtParams{
		ID:        task.ID,
		StartedAt: &nextTime,
	}); err != nil {
		return errors.Wrap(err, "failed to create task")
	}
	return nil
}
