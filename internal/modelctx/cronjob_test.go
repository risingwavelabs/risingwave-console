package modelctx

import (
	"context"
	"testing"
	"time"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/stretchr/testify/assert"
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
		taskID           = int32(101)
		timeoutDuration  = "5m"
	)

	mockModel := model.NewMockModelInterface(ctrl)
	mockModel.EXPECT().CreateTask(ctx, querier.CreateTaskParams{
		Attributes: apigen.TaskAttributes{
			OrgID: &orgID,
			Cronjob: &apigen.TaskCronjob{
				CronExpression: cronExpression,
			},
			Timeout: utils.Ptr(timeoutDuration),
		},
		Spec:      taskSpec,
		StartedAt: &expectedNextTime,
		Status:    string(apigen.Pending),
	}).Return(&querier.Task{
		ID: taskID,
	}, nil)

	mc := &ModelContext{
		model: mockModel,
		now: func() time.Time {
			return currentTime
		},
	}
	id, err := mc.CreateCronJob(ctx, utils.Ptr(timeoutDuration), &orgID, cronExpression, taskSpec)
	require.NoError(t, err)
	assert.Equal(t, taskID, id)
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
		timeoutDuration  = "5m"
	)

	mockModel := model.NewMockModelInterface(ctrl)
	mockModel.EXPECT().UpdateTask(ctx, querier.UpdateTaskParams{
		ID: taskID,
		Attributes: apigen.TaskAttributes{
			OrgID: &orgID,
			Cronjob: &apigen.TaskCronjob{
				CronExpression: cronExpression,
			},
			Timeout: utils.Ptr(timeoutDuration),
		},
		Spec:      taskSpec,
		StartedAt: &expectedNextTime,
	})

	mc := &ModelContext{
		model: mockModel,
		now: func() time.Time {
			return currentTime
		},
	}
	err := mc.UpdateCronJob(ctx, taskID, utils.Ptr(timeoutDuration), &orgID, cronExpression, taskSpec)
	require.NoError(t, err)
}

func TestPauseCronJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	var (
		taskID = int32(1)
	)

	mockModel := model.NewMockModelInterface(ctrl)
	mockModel.EXPECT().UpdateTaskStatus(ctx, querier.UpdateTaskStatusParams{
		ID:     taskID,
		Status: string(apigen.Paused),
	}).Return(nil)

	mc := &ModelContext{
		model: mockModel,
	}
	err := mc.PauseCronJob(ctx, taskID)
	require.NoError(t, err)
}

func TestResumeCronJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	var (
		taskID = int32(1)
	)

	mockModel := model.NewMockModelInterface(ctrl)
	mockModel.EXPECT().UpdateTaskStatus(ctx, querier.UpdateTaskStatusParams{
		ID:     taskID,
		Status: string(apigen.Pending),
	}).Return(nil)

	mc := &ModelContext{
		model: mockModel,
	}
	err := mc.ResumeCronJob(ctx, taskID)
	require.NoError(t, err)
}
