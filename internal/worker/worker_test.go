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
			ctx := context.Background()
			mockModel := model.NewExtendedMockModelInterface(ctrl)
			mockExecutor := NewMockExecutorInterface(ctrl)

			worker := &Worker{
				model: mockModel,
				getExecutor: func(m model.ModelInterface, risectlm *meta.RisectlManager) ExecutorInterface {
					return mockExecutor
				},
			}

			mockModel.EXPECT().PullTask(ctx).Return(&querier.Task{
				ID:     taskID,
				Spec:   taskSpec,
				Status: string(taskStatus),
			}, nil)

			mockExecutor.EXPECT().ExecuteAutoBackup(ctx, autoBackupSpec).Return(testCase.errExecute)

			if testCase.errExecute != nil {
				mockModel.EXPECT().UpdateTaskStatus(ctx, querier.UpdateTaskStatusParams{
					ID:     taskID,
					Status: string(apigen.Failed),
				}).Return(nil)

				mockModel.EXPECT().InsertEvent(ctx, apigen.EventSpec{
					Type: apigen.TaskError,
					TaskError: &apigen.EventTaskError{
						TaskID: taskID,
						Error:  testCase.errExecute.Error(),
					},
				}).Return(&querier.Event{}, nil)
			} else {
				mockModel.EXPECT().UpdateTaskStatus(ctx, querier.UpdateTaskStatusParams{
					ID:     taskID,
					Status: string(apigen.Completed),
				}).Return(nil)
			}

			worker.runTask(ctx)
		})
	}
}
