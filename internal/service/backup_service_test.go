package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cloudcarver/anchor/pkg/taskcore"
	"github.com/jackc/pgx/v5"
	"github.com/risingwavelabs/wavekit/internal/zcore/model"
	"github.com/risingwavelabs/wavekit/internal/zgen/apigen"
	"github.com/risingwavelabs/wavekit/internal/zgen/querier"
	"github.com/risingwavelabs/wavekit/internal/zgen/taskgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		tz                = "America/New_York"
		cronExpression    = "0 0 0 * * *"
		retentionDuration = "10m"
		taskParams        = taskgen.AutoBackupParameters{
			ClusterID:         clusterID,
			RetentionDuration: retentionDuration,
		}
	)

	spec, err := json.Marshal(taskParams)
	require.NoError(t, err)

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
			taskstore := taskcore.NewMockTaskStoreInterfaceWithTx(ctrl)
			taskRunner := taskgen.NewMockTaskRunner(ctrl)

			service := &Service{
				m:          mockModel,
				taskstore:  taskstore,
				taskRunner: taskRunner,
			}

			mockModel.EXPECT().GetOrgCluster(gomock.Any(), querier.GetOrgClusterParams{
				ID:    clusterID,
				OrgID: orgID,
			}).Return(&querier.Cluster{
				ID: clusterID,
			}, nil)

			mockModel.EXPECT().GetOrgSettings(gomock.Any(), orgID).Return(&querier.OrgSetting{
				OrgID:    orgID,
				Timezone: tz,
			}, nil)

			mockModel.EXPECT().GetAutoBackupConfig(gomock.Any(), clusterID).Return(tc.cfg, tc.err)

			if tc.err == nil { // UPDATE
				if tc.enabled {
					taskstore.EXPECT().ResumeCronJob(gomock.Any(), taskID).Return(nil)
				} else {
					taskstore.EXPECT().PauseCronJob(gomock.Any(), taskID).Return(nil)
				}
				taskstore.EXPECT().UpdateCronJob(gomock.Any(), taskID, fmt.Sprintf("CRON_TZ=%s %s", tz, cronExpression), spec).Return(nil)
				mockModel.EXPECT().UpdateAutoBackupConfig(gomock.Any(), querier.UpdateAutoBackupConfigParams{
					ClusterID: clusterID,
					Enabled:   tc.enabled,
				}).Return(nil)
			} else { // CREATEs
				taskRunner.EXPECT().RunAutoBackupWithTx(gomock.Any(), gomock.Any(), &taskgen.AutoBackupParameters{
					ClusterID:         clusterID,
					RetentionDuration: retentionDuration,
				}, taskcore.Eq(taskcore.WithCronjob(fmt.Sprintf("CRON_TZ=%s %s", tz, cronExpression)))).Return(taskID, nil)
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
