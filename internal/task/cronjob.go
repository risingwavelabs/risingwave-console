package task

import (
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
	"github.com/robfig/cron/v3"
)

func (s *TaskStore) CreateCronJob(c *modelctx.ModelCtx, timeoutDuration *string, orgID *int32, cronExpression string, specType apigen.TaskSpec) (int32, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cron, err := parser.Parse(cronExpression)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse cron expression")
	}
	nextTime := cron.Next(s.now())

	if timeoutDuration != nil { // validate if it is set
		_, err = time.ParseDuration(*timeoutDuration)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to parse timeout duration")
		}
	}

	taskAttributes := apigen.TaskAttributes{
		OrgID: orgID,
		Cronjob: &apigen.TaskCronjob{
			CronExpression: cronExpression,
		},
		Timeout: timeoutDuration,
	}

	task, err := c.CreateTask(c, querier.CreateTaskParams{
		Attributes: taskAttributes,
		Spec:       specType,
		StartedAt:  &nextTime,
		Status:     string(apigen.Pending),
	})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create task")
	}
	return task.ID, nil
}

func (s *TaskStore) UpdateCronJob(c *modelctx.ModelCtx, taskID int32, timeoutDuration *string, orgID *int32, cronExpression string, specType apigen.TaskSpec) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cron, err := parser.Parse(cronExpression)
	if err != nil {
		return errors.Wrapf(err, "failed to parse cron expression")
	}
	nextTime := cron.Next(s.now())

	if timeoutDuration != nil { // validate if it is set
		_, err = time.ParseDuration(*timeoutDuration)
		if err != nil {
			return errors.Wrapf(err, "failed to parse timeout duration")
		}
	}

	taskAttributes := apigen.TaskAttributes{
		OrgID: orgID,
		Cronjob: &apigen.TaskCronjob{
			CronExpression: cronExpression,
		},
		Timeout: timeoutDuration,
	}

	if err := c.UpdateTask(c, querier.UpdateTaskParams{
		ID:         taskID,
		Attributes: taskAttributes,
		StartedAt:  &nextTime,
		Spec:       specType,
	}); err != nil {
		return errors.Wrapf(err, "failed to update task")
	}
	return nil
}

func (s *TaskStore) PauseCronJob(c *modelctx.ModelCtx, taskID int32) error {
	if err := c.UpdateTaskStatus(c, querier.UpdateTaskStatusParams{
		ID:     taskID,
		Status: string(apigen.Paused),
	}); err != nil {
		return errors.Wrapf(err, "failed to pause task")
	}
	return nil
}

func (s *TaskStore) ResumeCronJob(c *modelctx.ModelCtx, taskID int32) error {
	if err := c.UpdateTaskStatus(c, querier.UpdateTaskStatusParams{
		ID:     taskID,
		Status: string(apigen.Pending),
	}); err != nil {
		return errors.Wrapf(err, "failed to resume task")
	}
	return nil
}
