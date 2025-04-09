package task

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
	taskstore := NewMockTaskStoreInterface(ctrl)

	model.EXPECT().GetClusterByID(gomock.Any(), clusterID).Return(&querier.Cluster{
		ID:             clusterID,
		Version:        clusterVersion,
		Host:           clusterHost,
		MetaPort:       clusterPort,
		OrganizationID: orgID,
	}, nil)

	risectlm.EXPECT().NewConn(gomock.Any(), clusterVersion, clusterHost, clusterPort).Return(risectlcm, nil)
	risectlcm.EXPECT().MetaBackup(gomock.Any()).Return(snapshotID, nil)

	executor := &TaskHandler{
		model:     model,
		risectlm:  risectlm,
		taskstore: taskstore,
		now:       func() time.Time { return currTime },
	}

	model.EXPECT().CreateSnapshot(gomock.Any(), querier.CreateSnapshotParams{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
	}).Return(nil)

	taskstore.EXPECT().CreateTask(gomock.Any(), CreateTaskParams{
		OrgID: &orgID,
		Spec: apigen.TaskSpec{
			Type: apigen.DeleteSnapshot,
			DeleteSnapshot: &apigen.TaskSpecDeleteSnapshot{
				ClusterID:  clusterID,
				SnapshotID: snapshotID,
			},
		},
		StartedAt:            utils.Ptr(currTime.Add(retentionDuration)),
		AlwaysRetryOnFailure: true,
		RetryInterval:        "10m",
	}).Return(int32(1), nil)

	err := executor.ExecuteAutoBackup(context.Background(), apigen.TaskSpecAutoBackup{
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
	taskstore := NewMockTaskStoreInterface(ctrl)

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

	taskstore.EXPECT().CreateTask(gomock.Any(), CreateTaskParams{
		OrgID: &orgID,
		Spec: apigen.TaskSpec{
			Type: apigen.DeleteClusterDiagnostic,
			DeleteClusterDiagnostic: &apigen.TaskSpecDeleteClusterDiagnostic{
				ClusterID:    clusterID,
				DiagnosticID: diagnosticID,
			},
		},
		StartedAt:            utils.Ptr(currTime.Add(retentionDuration)),
		AlwaysRetryOnFailure: true,
		RetryInterval:        "10m",
	}).Return(int32(1), nil)

	executor := &TaskHandler{
		model:     model,
		taskstore: taskstore,
		now:       func() time.Time { return currTime },
		metahttp:  metahttp,
	}

	err := executor.ExecuteAutoDiagnostic(context.Background(), apigen.TaskSpecAutoDiagnostic{
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

	executor := &TaskHandler{
		model: model,
	}

	err := executor.ExecuteDeleteClusterDiagnostic(context.Background(), apigen.TaskSpecDeleteClusterDiagnostic{
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

	executor := &TaskHandler{
		model:    model,
		risectlm: risectlm,
	}

	err := executor.ExecuteDeleteSnapshot(context.Background(), apigen.TaskSpecDeleteSnapshot{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
	})
	require.NoError(t, err)
}
