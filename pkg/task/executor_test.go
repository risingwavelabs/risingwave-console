package task

import (
	"context"
	"fmt"
	"testing"
	"time"

	mock_http "github.com/risingwavelabs/risingwave-console/pkg/conn/http/mock"
	mock_meta "github.com/risingwavelabs/risingwave-console/pkg/conn/meta/mock"
	"github.com/risingwavelabs/risingwave-console/pkg/utils"
	"github.com/risingwavelabs/risingwave-console/pkg/zcore/model"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/querier"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/taskgen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/cloudcarver/anchor/pkg/taskcore"
	anchor_apigen "github.com/cloudcarver/anchor/pkg/zgen/apigen"
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

	model.EXPECT().GetClusterByID(gomock.Any(), clusterID).Return(&querier.Cluster{
		ID:       clusterID,
		Version:  clusterVersion,
		Host:     clusterHost,
		MetaPort: clusterPort,
		OrgID:    orgID,
	}, nil)

	risectlm.EXPECT().NewConn(gomock.Any(), clusterVersion, clusterHost, clusterPort).Return(risectlcm, nil)
	risectlcm.EXPECT().MetaBackup(gomock.Any()).Return(snapshotID, nil)

	taskRunner := taskgen.NewMockTaskRunner(ctrl)

	executor := &TaskExecutor{
		risectlm:   risectlm,
		now:        func() time.Time { return currTime },
		taskRunner: taskRunner,
		model:      model,
	}

	model.EXPECT().CreateClusterSnapshot(gomock.Any(), querier.CreateClusterSnapshotParams{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
		Name:       fmt.Sprintf("auto-backup-%s", currTime.Format("2006-01-02-15-04-05")),
	}).Return(nil)

	taskRunner.EXPECT().RunDeleteSnapshot(
		gomock.Any(),
		&taskgen.DeleteSnapshotParameters{
			ClusterID:  clusterID,
			SnapshotID: snapshotID,
		},
		taskcore.Eq(func(task *anchor_apigen.Task) error {
			task.StartedAt = utils.Ptr(currTime.Add(retentionDuration))
			return nil
		}),
	).Return(int32(1), nil)

	err := executor.ExecuteAutoBackup(context.Background(), &taskgen.AutoBackupParameters{
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
	taskRunner := taskgen.NewMockTaskRunner(ctrl)

	model.EXPECT().GetClusterByID(gomock.Any(), clusterID).Return(&querier.Cluster{
		ID:       clusterID,
		Host:     clusterHost,
		HttpPort: clusterHttpPort,
		OrgID:    orgID,
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

	taskRunner.EXPECT().RunDeleteClusterDiagnostic(
		gomock.Any(),
		&taskgen.DeleteClusterDiagnosticParameters{
			ClusterID:    clusterID,
			DiagnosticID: diagnosticID,
		},
		taskcore.Eq(func(task *anchor_apigen.Task) error {
			task.StartedAt = utils.Ptr(currTime.Add(retentionDuration))
			return nil
		}),
	).Return(int32(1), nil)

	executor := &TaskExecutor{
		taskRunner: taskRunner,
		metahttp:   metahttp,
		now:        func() time.Time { return currTime },
		model:      model,
	}

	err := executor.ExecuteAutoDiagnostic(context.Background(), &taskgen.AutoDiagnosticParameters{
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

	executor := &TaskExecutor{
		model: model,
	}

	err := executor.ExecuteDeleteClusterDiagnostic(context.Background(), &taskgen.DeleteClusterDiagnosticParameters{
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

	model.EXPECT().DeleteClusterSnapshot(gomock.Any(), querier.DeleteClusterSnapshotParams{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
	}).Return(nil)

	executor := &TaskExecutor{
		risectlm: risectlm,
		model:    model,
	}

	err := executor.ExecuteDeleteSnapshot(context.Background(), &taskgen.DeleteSnapshotParameters{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
	})
	require.NoError(t, err)
}
