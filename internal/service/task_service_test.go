package service

import (
	"context"
	"testing"
	"time"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateCronJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	var (
		orgID            = int32(1)
		cronExpression   = "0 0 * * *"
		taskSpec         = apigen.TaskSpec{}
		currentTime      = time.Date(2025, 3, 31, 12, 0, 0, 0, time.UTC)
		expectedNextTime = time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	)

	mockModel := model.NewMockModelInterface(ctrl)
	mockModel.EXPECT().CreateTask(ctx, querier.CreateTaskParams{
		Attributes: apigen.TaskAttributes{
			OrgID: &orgID,
			Cronjob: &apigen.TaskCronjob{
				CronExpression: cronExpression,
			},
		},
		Spec:      taskSpec,
		StartedAt: &expectedNextTime,
	})

	taskService := &Service{
		m: mockModel,
		now: func() time.Time {
			return currentTime
		},
	}
	err := taskService.CreateCronJob(ctx, &orgID, cronExpression, taskSpec)
	require.NoError(t, err)
}

func TestUpdateCronJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	var (
		taskID           = int32(1)
		orgID            = int32(1)
		cronExpression   = "0 0 * * *"
		taskSpec         = apigen.TaskSpec{}
		currentTime      = time.Date(2025, 3, 31, 12, 0, 0, 0, time.UTC)
		expectedNextTime = time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	)

	mockModel := model.NewMockModelInterface(ctrl)
	mockModel.EXPECT().UpdateTask(ctx, querier.UpdateTaskParams{
		ID: taskID,
		Attributes: apigen.TaskAttributes{
			OrgID: &orgID,
			Cronjob: &apigen.TaskCronjob{
				CronExpression: cronExpression,
			},
		},
		Spec:      taskSpec,
		StartedAt: &expectedNextTime,
	})

	taskService := &Service{
		m: mockModel,
		now: func() time.Time {
			return currentTime
		},
	}
	err := taskService.UpdateCronJob(ctx, taskID, &orgID, cronExpression, taskSpec)
	require.NoError(t, err)
}
