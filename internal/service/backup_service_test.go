package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/task"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateClusterAutoBackupConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	var (
		orgID             = int32(201)
		clusterID         = int32(101)
		taskID            = int32(301)
		cronExpression    = "0 0 * * *"
		retentionDuration = "10m"
		taskSpec          = apigen.TaskSpec{
			Type: apigen.AutoBackup,
			AutoBackup: &apigen.TaskSpecAutoBackup{
				ClusterID:         clusterID,
				RetentionDuration: retentionDuration,
			},
		}
	)

	type testCase struct {
		name    string
		err     error
		cfg     *querier.AutoBackupConfig
		enabled bool
	}

	testCases := []testCase{
		{
			name: "existing auto backup config",
			err:  nil,
			cfg: &querier.AutoBackupConfig{
				TaskID: taskID,
			},
		},
		{
			name: "no existing auto backup config",
			err:  pgx.ErrNoRows,
			cfg:  nil,
		},
		{
			name: "resume auto backup config",
			err:  nil,
			cfg: &querier.AutoBackupConfig{
				TaskID: taskID,
			},
			enabled: true,
		},
		{
			name: "pause auto backup config",
			err:  nil,
			cfg: &querier.AutoBackupConfig{
				TaskID: taskID,
			},
			enabled: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockModel := model.NewMockModelInterfaceWithTransaction(ctrl)
			taskstore := task.NewMockTaskStoreInterface(ctrl)
			service := &Service{
				m:         mockModel,
				taskstore: taskstore,
			}

			mockModel.EXPECT().GetOrgCluster(gomock.Any(), querier.GetOrgClusterParams{
				ID:             clusterID,
				OrganizationID: orgID,
			}).Return(&querier.Cluster{
				ID: clusterID,
			}, nil)

			mockModel.EXPECT().GetAutoBackupConfig(gomock.Any(), clusterID).Return(tc.cfg, tc.err)

			if tc.err == nil { // UPDATE
				if tc.enabled {
					taskstore.EXPECT().ResumeCronJob(gomock.Any(), taskID).Return(nil)
				} else {
					taskstore.EXPECT().PauseCronJob(gomock.Any(), taskID).Return(nil)
				}
				taskstore.EXPECT().UpdateCronJob(gomock.Any(), taskID, utils.Ptr(defaultBackupTaskTimeout), &orgID, cronExpression, taskSpec).Return(nil)
				mockModel.EXPECT().UpdateAutoBackupConfig(gomock.Any(), querier.UpdateAutoBackupConfigParams{
					ClusterID: clusterID,
					Enabled:   tc.enabled,
				}).Return(nil)
			} else { // CREATEs
				taskstore.EXPECT().CreateCronJob(gomock.Any(), utils.Ptr(defaultBackupTaskTimeout), &orgID, cronExpression, taskSpec).Return(taskID, nil)
				mockModel.EXPECT().CreateAutoBackupConfig(gomock.Any(), querier.CreateAutoBackupConfigParams{
					ClusterID: clusterID,
					TaskID:    taskID,
					Enabled:   true,
				}).Return(nil)
			}

			err := service.UpdateClusterAutoBackupConfig(ctx, clusterID, apigen.AutoBackupConfig{
				CronExpression:    cronExpression,
				RetentionDuration: retentionDuration,
				Enabled:           tc.enabled,
			}, orgID)
			assert.NoError(t, err)
		})
	}
}
