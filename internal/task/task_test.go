package task

import (
	"context"
	"testing"
	"time"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	mock_meta "github.com/risingwavelabs/wavekit/internal/conn/meta/mock"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	mock_task "github.com/risingwavelabs/wavekit/internal/task/mock"
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
	taskstore := mock_task.NewMockTaskStoreInterface(ctrl)

	model.EXPECT().GetClusterByID(gomock.Any(), clusterID).Return(&querier.Cluster{
		ID:             clusterID,
		Version:        clusterVersion,
		Host:           clusterHost,
		MetaPort:       clusterPort,
		OrganizationID: orgID,
	}, nil)

	risectlm.EXPECT().NewConn(gomock.Any(), clusterVersion, clusterHost, clusterPort).Return(risectlcm, nil)
	risectlcm.EXPECT().MetaBackup(gomock.Any()).Return(snapshotID, nil)

	executor := &TaskExecutor{
		model:     model,
		risectlm:  risectlm,
		taskstore: taskstore,
		now:       func() time.Time { return currTime },
	}

	model.EXPECT().CreateSnapshot(gomock.Any(), querier.CreateSnapshotParams{
		ClusterID:  clusterID,
		SnapshotID: snapshotID,
	}).Return(nil)

	taskstore.EXPECT().CreateTask(gomock.Any(), &orgID, apigen.TaskSpec{
		Type: apigen.DeleteSnapshot,
		DeleteSnapshot: &apigen.TaskSpecDeleteSnapshot{
			ClusterID:  clusterID,
			SnapshotID: snapshotID,
		},
	}, utils.Ptr(currTime.Add(retentionDuration))).Return(int32(1), nil)

	err := executor.ExecuteAutoBackup(context.Background(), apigen.TaskSpecAutoBackup{
		ClusterID:         clusterID,
		RetentionDuration: retentionDurationRaw,
	})
	require.NoError(t, err)
}
