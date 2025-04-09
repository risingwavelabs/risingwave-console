package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleCronjob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txm := model.NewExtendedMockModelInterface(ctrl)

	tz := "Asia/Shanghai"
	location, err := time.LoadLocation(tz)
	require.NoError(t, err)

	var (
		orgID = int32(1)

		currTime = time.Date(2025, 3, 27, 0, 0, 1, 0, location)
		cronExpr = "0 0 * * *"
		nextTime = time.Date(2025, 3, 28, 0, 0, 0, 0, location)
		taskID   = int32(1)
	)

	handler := &TaskLifeCycleHandler{
		txm: txm,
		now: func() time.Time {
			return currTime
		},
	}

	txm.EXPECT().GetOrganization(context.Background(), orgID).Return(&querier.Organization{
		Timezone: tz,
	}, nil)

	txm.EXPECT().UpdateTaskStartedAt(context.Background(), querier.UpdateTaskStartedAtParams{
		ID:        taskID,
		StartedAt: utils.Ptr(nextTime),
	}).Return(nil)

	task := apigen.Task{
		ID: taskID,
		Attributes: apigen.TaskAttributes{
			OrgID: &orgID,
			Cronjob: &apigen.TaskCronjob{
				CronExpression: cronExpr,
			},
		},
	}

	err = handler.handleCronjob(context.Background(), task)
	require.NoError(t, err)
}

func TestHandleCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		taskID = int32(1)
	)

	txm := model.NewExtendedMockModelInterface(ctrl)

	handler := &TaskLifeCycleHandler{
		txm: txm,
	}

	task := apigen.Task{
		ID: taskID,
	}

	txm.EXPECT().InsertEvent(context.Background(), apigen.EventSpec{
		Type: apigen.TaskCompleted,
		TaskCompleted: &apigen.EventTaskCompleted{
			TaskID: taskID,
		},
	}).Return(&querier.Event{}, nil)

	txm.EXPECT().UpdateTaskStatus(context.Background(), querier.UpdateTaskStatusParams{
		ID:     taskID,
		Status: string(apigen.Completed),
	}).Return(nil)

	err := handler.HandleCompleted(context.Background(), task)
	require.NoError(t, err)
}

func TestHandleFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		taskID = int32(1)
		err    = errors.New("test error")
	)

	txm := model.NewExtendedMockModelInterface(ctrl)

	handler := &TaskLifeCycleHandler{
		txm: txm,
	}

	task := apigen.Task{
		ID: taskID,
	}

	txm.EXPECT().InsertEvent(context.Background(), apigen.EventSpec{
		Type: apigen.TaskError,
		TaskError: &apigen.EventTaskError{
			TaskID: taskID,
			Error:  err.Error(),
		},
	}).Return(&querier.Event{}, nil)

	txm.EXPECT().UpdateTaskStatus(context.Background(), querier.UpdateTaskStatusParams{
		ID:     taskID,
		Status: string(apigen.Failed),
	}).Return(nil)

	err = handler.HandleFailed(context.Background(), task, err)
	require.NoError(t, err)
}

func TestHandleFailedWithRetryPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		taskID      = int32(1)
		err         = errors.New("test error")
		intervalRaw = "1h"
		interval, _ = time.ParseDuration(intervalRaw)
		currTime    = time.Now()
	)

	txm := model.NewExtendedMockModelInterface(ctrl)

	handler := &TaskLifeCycleHandler{
		txm: txm,
		now: func() time.Time {
			return currTime
		},
	}

	task := apigen.Task{
		ID: taskID,
		Attributes: apigen.TaskAttributes{
			RetryPolicy: &apigen.TaskRetryPolicy{
				AlwaysRetryOnFailure: utils.Ptr(true),
				Interval:             intervalRaw,
			},
		},
	}

	txm.EXPECT().InsertEvent(context.Background(), apigen.EventSpec{
		Type: apigen.TaskError,
		TaskError: &apigen.EventTaskError{
			TaskID: taskID,
			Error:  err.Error(),
		},
	}).Return(&querier.Event{}, nil)

	txm.EXPECT().UpdateTaskStartedAt(context.Background(), querier.UpdateTaskStartedAtParams{
		ID:        taskID,
		StartedAt: utils.Ptr(currTime.Add(interval)),
	}).Return(nil)

	err = handler.HandleFailed(context.Background(), task, err)
	require.NoError(t, err)

}
