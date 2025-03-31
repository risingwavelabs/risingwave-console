package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
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

	type getAutoBackupConfigRtn struct {
		err error
		cfg *querier.AutoBackupConfig
	}

	testCases := []getAutoBackupConfigRtn{
		{
			err: nil,
			cfg: &querier.AutoBackupConfig{
				TaskID: taskID,
			},
		},
		{
			err: pgx.ErrNoRows,
			cfg: nil,
		},
	}

	for _, tc := range testCases {
		self := NewMockServiceInterface(ctrl)
		mockModel := model.NewExtendedMockModelInterface(ctrl)
		service := &Service{
			self: self,
			m:    mockModel,
		}

		mockModel.EXPECT().GetOrgCluster(ctx, querier.GetOrgClusterParams{
			ID:             clusterID,
			OrganizationID: orgID,
		}).Return(&querier.Cluster{
			ID: clusterID,
		}, nil)

		mockModel.EXPECT().GetAutoBackupConfig(ctx, clusterID).Return(tc.cfg, tc.err)

		if tc.err == nil {
			self.EXPECT().UpdateCronJob(ctx, taskID, &orgID, cronExpression, taskSpec).Return(nil)
		} else {
			self.EXPECT().CreateCronJob(ctx, &orgID, cronExpression, taskSpec).Return(nil)
		}

		err := service.UpdateClusterAutoBackupConfig(ctx, clusterID, apigen.AutoBackupConfig{
			CronExpression:    cronExpression,
			RetentionDuration: retentionDuration,
		}, orgID)
		assert.NoError(t, err)
	}
}
