package worker

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/worker/mock"
	"go.uber.org/mock/gomock"
)

type RunTaskTest struct {
	errExecute error
}

func TestRunTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		taskID         int32 = 1
		clusterID      int32 = 1
		autoBackupSpec       = apigen.TaskSpecAutoBackup{
			ClusterID: clusterID,
		}
		taskSpec = apigen.TaskSpec{
			Type:       apigen.AutoBackup,
			AutoBackup: &autoBackupSpec,
		}
		taskStatus = apigen.Pending
		task       = apigen.Task{
			ID:     taskID,
			Spec:   taskSpec,
			Status: taskStatus,
		}
	)

	testCases := []RunTaskTest{
		{
			errExecute: nil,
		},
		{
			errExecute: errors.New("buckethead"),
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("errExecute: %v", testCase.errExecute), func(t *testing.T) {
			mockModel := model.NewExtendedMockModelInterface(ctrl)
			mockExecutor := mock.NewMockExecutorInterface(ctrl)
			mockLifeCycleHandler := mock.NewMockTaskLifeCycleHandlerInterface(ctrl)

			worker := &Worker{
				model: mockModel,
				getExecutor: func(m model.ModelInterface, risectlm *meta.RisectlManager) ExecutorInterface {
					return mockExecutor
				},
				getHandler: func(txm model.ModelInterface) (TaskLifeCycleHandlerInterface, error) {
					return mockLifeCycleHandler, nil
				},
			}
			// pull task
			mockModel.EXPECT().PullTask(gomock.Any()).Return(&querier.Task{
				ID:     taskID,
				Spec:   taskSpec,
				Status: string(taskStatus),
			}, nil)

			// called by handle attributes
			mockLifeCycleHandler.EXPECT().HandleAttributes(gomock.Any(), task).Return(nil)

			// executor run business logic
			mockExecutor.EXPECT().ExecuteAutoBackup(gomock.Any(), autoBackupSpec).Return(testCase.errExecute)

			if testCase.errExecute != nil {
				mockLifeCycleHandler.EXPECT().HandleFailed(gomock.Any(), task, testCase.errExecute).Return(nil)
			} else {
				mockLifeCycleHandler.EXPECT().HandleCompleted(gomock.Any(), task).Return(nil)
			}

			worker.runTask(context.Background())
		})
	}
}
