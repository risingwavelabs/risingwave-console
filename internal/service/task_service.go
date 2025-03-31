package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/robfig/cron/v3"
)

func (s *Service) CreateCronJob(ctx context.Context, orgID *int32, cronExpression string, specType apigen.TaskSpec) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cron, err := parser.Parse(cronExpression)
	if err != nil {
		return errors.Wrapf(err, "failed to parse cron expression")
	}
	nextTime := cron.Next(s.now())

	taskAttributes := apigen.TaskAttributes{
		OrgID: orgID,
		Cronjob: &apigen.TaskCronjob{
			CronExpression: cronExpression,
		},
	}

	_, err = s.m.CreateTask(ctx, querier.CreateTaskParams{
		Attributes: taskAttributes,
		Spec:       specType,
		StartedAt:  &nextTime,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create task")
	}
	return nil
}

func (s *Service) UpdateCronJob(ctx context.Context, taskID int32, orgID *int32, cronExpression string, specType apigen.TaskSpec) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cron, err := parser.Parse(cronExpression)
	if err != nil {
		return errors.Wrapf(err, "failed to parse cron expression")
	}
	nextTime := cron.Next(s.now())

	taskAttributes := apigen.TaskAttributes{
		OrgID: orgID,
		Cronjob: &apigen.TaskCronjob{
			CronExpression: cronExpression,
		},
	}

	if err := s.m.UpdateTask(ctx, querier.UpdateTaskParams{
		ID:         taskID,
		Attributes: taskAttributes,
		StartedAt:  &nextTime,
		Spec:       specType,
	}); err != nil {
		return errors.Wrapf(err, "failed to update task")
	}
	return nil
}
