package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
	mock_modelctx "github.com/risingwavelabs/wavekit/internal/modelctx/mock"
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
			mockModel := model.NewExtendedMockModelInterface(ctrl)
			mockModelctx := mock_modelctx.NewMockModelContextInterface(ctrl)
			service := &Service{
				m: mockModel,
				modelctx: func(model model.ModelInterface) modelctx.ModelContextInterface {
					return mockModelctx
				},
			}

			mockModel.EXPECT().GetOrgCluster(ctx, querier.GetOrgClusterParams{
				ID:             clusterID,
				OrganizationID: orgID,
			}).Return(&querier.Cluster{
				ID: clusterID,
			}, nil)

			mockModel.EXPECT().GetAutoBackupConfig(ctx, clusterID).Return(tc.cfg, tc.err)

			if tc.err == nil { // UPDATE
				if tc.enabled {
					mockModelctx.EXPECT().ResumeCronJob(ctx, taskID).Return(nil)
				} else {
					mockModelctx.EXPECT().PauseCronJob(ctx, taskID).Return(nil)
				}
				mockModelctx.EXPECT().UpdateCronJob(ctx, taskID, utils.Ptr(defaultBackupTaskTimeout), &orgID, cronExpression, taskSpec).Return(nil)
				mockModel.EXPECT().UpdateAutoBackupConfig(ctx, querier.UpdateAutoBackupConfigParams{
					ClusterID: clusterID,
					Enabled:   tc.enabled,
				}).Return(nil)
			} else { // CREATEs
				mockModelctx.EXPECT().CreateCronJob(ctx, utils.Ptr(defaultBackupTaskTimeout), &orgID, cronExpression, taskSpec).Return(taskID, nil)
				mockModel.EXPECT().CreateAutoBackupConfig(ctx, querier.CreateAutoBackupConfigParams{
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
