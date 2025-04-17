package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	mock_http "github.com/risingwavelabs/wavekit/internal/conn/http/mock"
	mock_meta "github.com/risingwavelabs/wavekit/internal/conn/meta/mock"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
	"github.com/risingwavelabs/wavekit/internal/task"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExecuteAutoBackup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		orgID                = int32(201)
		clusterVersion       = "v2.2.1"
		clusterID            = int32(101)
		clusterHost          = "localhost"
		clusterPort          = int32(9000)
		snapshotID           = int64(1)
		retentionDurationRaw = "3d"
		retentionDuration, _ = utils.ParseDuration(retentionDurationRaw)
		currTime             = time.Now()
	)

	model := model.NewMockModelInterface(ctrl)
	risectlm := mock_meta.NewMockRisectlManagerInterface(ctrl)
	risectlcm := mock_meta.NewMockRisectlConn(ctrl)
	taskstore := task.NewMockTaskStoreInterface(ctrl)

	model.EXPECT().GetClusterByID(gomock.Any(), clusterID).Return(&querier.Cluster{
		ID:             clusterID,
		Version:        clusterVersion,
		Host:           clusterHost,
		MetaPort:       clusterPort,
		OrganizationID: orgID,
	}, nil)

	risectlm.EXPECT().NewConn(gomock.Any(), clusterVersion, clusterHost, clusterPort).Return(risectlcm, nil)
	risectlcm.EXPECT().MetaBackup(gomock.Any()).Return(snapshotID, nil)

	c := &modelctx.ModelCtx{
		ModelInterface: model,
		Context:        context.Background(),
	}

	handler := &TaskHandler{
		risectlm:  risectlm,
		taskstore: taskstore,
		now:       func() time.Time { return currTime },
	}

	model.EXPECT().CreateSnapshot(gomock.Any(), querier.CreateSnapshotParams{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
	}).Return(nil)

	taskstore.EXPECT().PushTask(
		gomock.Any(),
		apigen.TaskSpec{
			Type: apigen.DeleteSnapshot,
			DeleteSnapshot: &apigen.TaskSpecDeleteSnapshot{
				ClusterID:  clusterID,
				SnapshotID: snapshotID,
			},
		},
		task.MatchTaskOpt(task.OrgID(orgID)),
		task.MatchTaskOpt(task.AlwaysRetryOnFailure("10m")),
		task.MatchTaskOpt(task.StartedAt(currTime.Add(retentionDuration))),
	).Return(int32(1), nil)

	err := handler.ExecuteAutoBackup(c, apigen.TaskSpecAutoBackup{
		ClusterID:         clusterID,
		RetentionDuration: retentionDurationRaw,
	})
	require.NoError(t, err)
}

func TestExecuteAutoDiagnostic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		orgID                = int32(201)
		clusterID            = int32(101)
		clusterHost          = "localhost"
		clusterHttpPort      = int32(9000)
		diagnose             = "diagnose"
		currTime             = time.Now()
		diagnosticID         = int32(301)
		retentionDurationRaw = "3d"
		retentionDuration, _ = utils.ParseDuration(retentionDurationRaw)
	)

	metahttp := mock_http.NewMockMetaHttpManagerInterface(ctrl)
	model := model.NewMockModelInterface(ctrl)
	taskstore := task.NewMockTaskStoreInterface(ctrl)

	model.EXPECT().GetClusterByID(gomock.Any(), clusterID).Return(&querier.Cluster{
		ID:             clusterID,
		Host:           clusterHost,
		HttpPort:       clusterHttpPort,
		OrganizationID: orgID,
	}, nil)

	metahttp.
		EXPECT().
		GetDiagnose(gomock.Any(), fmt.Sprintf("http://%s:%d", clusterHost, clusterHttpPort)).
		Return(diagnose, nil)

	model.EXPECT().CreateClusterDiagnostic(gomock.Any(), querier.CreateClusterDiagnosticParams{
		ClusterID: clusterID,
		Content:   diagnose,
	}).Return(&querier.ClusterDiagnostic{
		ID: diagnosticID,
	}, nil)

	taskstore.EXPECT().PushTask(
		gomock.Any(),
		apigen.TaskSpec{
			Type: apigen.DeleteClusterDiagnostic,
			DeleteClusterDiagnostic: &apigen.TaskSpecDeleteClusterDiagnostic{
				ClusterID:    clusterID,
				DiagnosticID: diagnosticID,
			},
		},
		task.MatchTaskOpt(task.OrgID(orgID)),
		task.MatchTaskOpt(task.AlwaysRetryOnFailure("10m")),
		task.MatchTaskOpt(task.StartedAt(currTime.Add(retentionDuration))),
	).Return(int32(1), nil)

	c := &modelctx.ModelCtx{
		ModelInterface: model,
		Context:        context.Background(),
	}

	handler := &TaskHandler{
		taskstore: taskstore,
		metahttp:  metahttp,
		now:       func() time.Time { return currTime },
	}

	err := handler.ExecuteAutoDiagnostic(c, apigen.TaskSpecAutoDiagnostic{
		ClusterID:         clusterID,
		RetentionDuration: retentionDurationRaw,
	})
	require.NoError(t, err)
}

func TestExecuteDeleteClusterDiagnostic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		diagID = int32(301)
	)

	model := model.NewMockModelInterface(ctrl)

	model.EXPECT().DeleteClusterDiagnostic(gomock.Any(), diagID).Return(nil)

	c := &modelctx.ModelCtx{
		ModelInterface: model,
		Context:        context.Background(),
	}

	handler := &TaskHandler{}

	err := handler.ExecuteDeleteClusterDiagnostic(c, apigen.TaskSpecDeleteClusterDiagnostic{
		DiagnosticID: diagID,
	})
	require.NoError(t, err)
}

func TestExecuteDeleteSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		clusterID      = int32(101)
		snapshotID     = int64(1)
		clusterHost    = "localhost"
		clusterPort    = int32(9000)
		clusterVersion = "v2.2.1"
	)

	model := model.NewMockModelInterface(ctrl)
	risectlm := mock_meta.NewMockRisectlManagerInterface(ctrl)
	risectlcm := mock_meta.NewMockRisectlConn(ctrl)

	model.EXPECT().GetClusterByID(gomock.Any(), clusterID).Return(&querier.Cluster{
		ID:       clusterID,
		Host:     clusterHost,
		MetaPort: clusterPort,
		Version:  clusterVersion,
	}, nil)

	risectlm.EXPECT().NewConn(gomock.Any(), clusterVersion, clusterHost, clusterPort).Return(risectlcm, nil)
	risectlcm.EXPECT().DeleteSnapshot(gomock.Any(), snapshotID).Return(nil)

	model.EXPECT().DeleteSnapshot(gomock.Any(), querier.DeleteSnapshotParams{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
	}).Return(nil)

	c := &modelctx.ModelCtx{
		ModelInterface: model,
		Context:        context.Background(),
	}

	handler := &TaskHandler{
		risectlm: risectlm,
	}

	err := handler.ExecuteDeleteSnapshot(c, apigen.TaskSpecDeleteSnapshot{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
	})
	require.NoError(t, err)
}
