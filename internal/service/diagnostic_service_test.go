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

func TestUpdateClusterAutoDiagnosticConfig(t *testing.T) {
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
			Type: apigen.AutoDiagnostic,
			AutoDiagnostic: &apigen.TaskSpecAutoDiagnostic{
				ClusterID:         clusterID,
				RetentionDuration: retentionDuration,
			},
		}
	)

	type getAutoBackupConfigRtn struct {
		err error
		cfg *querier.AutoDiagnosticsConfig
	}

	testCases := []getAutoBackupConfigRtn{
		{
			err: nil,
			cfg: &querier.AutoDiagnosticsConfig{
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

		mockModel.EXPECT().GetAutoDiagnosticsConfig(ctx, clusterID).Return(tc.cfg, tc.err)

		if tc.err == nil {
			self.EXPECT().UpdateCronJob(ctx, taskID, &orgID, cronExpression, taskSpec).Return(nil)
		} else {
			self.EXPECT().CreateCronJob(ctx, &orgID, cronExpression, taskSpec).Return(nil)
		}

		err := service.UpdateClusterAutoDiagnosticConfig(ctx, clusterID, apigen.AutoDiagnosticConfig{
			CronExpression:    cronExpression,
			RetentionDuration: retentionDuration,
		}, orgID)
		assert.NoError(t, err)
	}
}
