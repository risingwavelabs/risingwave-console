package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudcarver/anchor/pkg/taskcore"
	"github.com/jackc/pgx/v5"
	mock_http "github.com/risingwavelabs/wavekit/pkg/conn/http/mock"
	"github.com/risingwavelabs/wavekit/pkg/zcore/model"
	"github.com/risingwavelabs/wavekit/pkg/zgen/apigen"
	"github.com/risingwavelabs/wavekit/pkg/zgen/querier"
	"github.com/risingwavelabs/wavekit/pkg/zgen/taskgen"
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
		tz                = "America/New_York"
	)

	type testCase struct {
		name    string
		err     error
		cfg     *querier.AutoDiagnosticsConfig
		enabled bool
	}

	testCases := []testCase{
		{
			name: "existing auto diagnostics config",
			err:  nil,
			cfg: &querier.AutoDiagnosticsConfig{
				TaskID: taskID,
			},
		},
		{
			name: "no existing auto diagnostics config",
			err:  pgx.ErrNoRows,
			cfg:  nil,
		},
		{
			name: "resume auto diagnostics config",
			err:  nil,
			cfg: &querier.AutoDiagnosticsConfig{
				TaskID: taskID,
			},
			enabled: true,
		},
		{
			name: "pause auto diagnostics config",
			err:  nil,
			cfg: &querier.AutoDiagnosticsConfig{
				TaskID: taskID,
			},
			enabled: false,
		},
	}

	for _, tc := range testCases {
		mockModel := model.NewMockModelInterfaceWithTransaction(ctrl)
		taskRunner := taskgen.NewMockTaskRunner(ctrl)
		taskstore := taskcore.NewMockTaskStoreInterfaceWithTx(ctrl)

		service := &Service{
			m:          mockModel,
			taskRunner: taskRunner,
			taskstore:  taskstore,
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

		mockModel.EXPECT().GetAutoDiagnosticsConfig(gomock.Any(), clusterID).Return(tc.cfg, tc.err)

		if tc.err == nil { // UPDATE
			if tc.enabled {
				taskstore.EXPECT().ResumeCronJob(gomock.Any(), taskID).Return(nil)
			} else {
				taskstore.EXPECT().PauseCronJob(gomock.Any(), taskID).Return(nil)
			}
			taskParams := taskgen.AutoDiagnosticParameters{
				ClusterID:         clusterID,
				RetentionDuration: retentionDuration,
			}
			spec, err := taskParams.Marshal()
			if err != nil {
				t.Fatalf("failed to marshal task parameters: %v", err)
			}

			taskstore.EXPECT().UpdateCronJob(gomock.Any(), taskID, fmt.Sprintf("CRON_TZ=%s %s", tz, cronExpression), spec).Return(nil)
			mockModel.EXPECT().UpdateAutoDiagnosticsConfig(gomock.Any(), querier.UpdateAutoDiagnosticsConfigParams{
				ClusterID: clusterID,
				Enabled:   tc.enabled,
			}).Return(nil)
		} else { // CREATE
			taskRunner.EXPECT().RunAutoDiagnosticWithTx(gomock.Any(), gomock.Any(), &taskgen.AutoDiagnosticParameters{
				ClusterID:         clusterID,
				RetentionDuration: retentionDuration,
			}, taskcore.Eq(taskcore.WithCronjob(fmt.Sprintf("CRON_TZ=%s %s", tz, cronExpression)))).Return(taskID, nil)

			mockModel.EXPECT().CreateAutoDiagnosticsConfig(gomock.Any(), querier.CreateAutoDiagnosticsConfigParams{
				ClusterID: clusterID,
				TaskID:    taskID,
				Enabled:   true,
			}).Return(nil)
		}
		err := service.UpdateClusterAutoDiagnosticConfig(ctx, clusterID, apigen.AutoDiagnosticConfig{
			CronExpression:    cronExpression,
			RetentionDuration: retentionDuration,
			Enabled:           tc.enabled,
		}, orgID)
		assert.NoError(t, err)
	}
}

func TestCreateClusterDiagnostic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx               = context.Background()
		orgID             = int32(201)
		clusterID         = int32(101)
		diagnosticID      = int32(301)
		clusterHost       = "localhost"
		clusterHttpPort   = int32(9000)
		diagnosticContent = "diagnostic content"
	)

	metahttp := mock_http.NewMockMetaHttpManagerInterface(ctrl)
	model := model.NewMockModelInterfaceWithTransaction(ctrl)
	service := &Service{
		m:        model,
		metahttp: metahttp,
	}

	model.EXPECT().GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:    clusterID,
		OrgID: orgID,
	}).Return(&querier.Cluster{
		ID:       clusterID,
		Host:     clusterHost,
		HttpPort: clusterHttpPort,
		OrgID:    orgID,
	}, nil)

	metahttp.
		EXPECT().
		GetDiagnose(ctx, fmt.Sprintf("http://%s:%d", clusterHost, clusterHttpPort)).
		Return(diagnosticContent, nil)

	model.EXPECT().CreateClusterDiagnostic(ctx, querier.CreateClusterDiagnosticParams{
		ClusterID: clusterID,
		Content:   diagnosticContent,
	}).Return(&querier.ClusterDiagnostic{
		ID:      diagnosticID,
		Content: diagnosticContent,
	}, nil)

	diagnostic, err := service.CreateClusterDiagnostic(ctx, clusterID, orgID)
	assert.NoError(t, err)
	assert.Equal(t, diagnosticID, diagnostic.ID)
	assert.Equal(t, diagnosticContent, diagnostic.Content)
}
