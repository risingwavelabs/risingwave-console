package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
)

func (s *Service) CreateClusterSnapshot(ctx context.Context, id int32, name string, orgID int32) (*apigen.Snapshot, error) {
	conn, err := s.getRisectlConn(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get risectl connection")
	}

	snapshotID, err := conn.MetaBackup(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create snapshot")
	}

	if err := s.m.CreateClusterSnapshot(ctx, querier.CreateClusterSnapshotParams{
		ClusterID:  id,
		SnapshotID: snapshotID,
		Name:       name,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to create snapshot")
	}

	return &apigen.Snapshot{
		ID:        snapshotID,
		Name:      name,
		ClusterID: id,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Service) ListClusterSnapshots(ctx context.Context, id int32, orgID int32) ([]apigen.Snapshot, error) {
	snapshots, err := s.m.ListClusterSnapshots(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list cluster snapshots")
	}

	result := make([]apigen.Snapshot, len(snapshots))
	for i, snapshot := range snapshots {
		result[i] = apigen.Snapshot{
			ID:        snapshot.SnapshotID,
			Name:      snapshot.Name,
			ClusterID: snapshot.ClusterID,
			CreatedAt: snapshot.CreatedAt,
		}
	}

	return result, nil
}

func (s *Service) DeleteClusterSnapshot(ctx context.Context, id int32, snapshotID int64, orgID int32) error {
	conn, err := s.getRisectlConn(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to get risectl connection")
	}

	if err := conn.DeleteSnapshot(ctx, snapshotID); err != nil {
		return errors.Wrapf(err, "failed to delete snapshot")
	}

	return nil
}

func (s *Service) UpdateClusterAutoBackupConfig(ctx context.Context, id int32, params apigen.AutoBackupConfig, orgID int32) error {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to get cluster")
	}

	if err := s.m.UpsertAutoBackupConfig(ctx, querier.UpsertAutoBackupConfigParams{
		ClusterID:      cluster.ID,
		Enabled:        params.Enabled,
		CronExpression: params.CronExpression,
		KeepLast:       params.KeepLast,
	}); err != nil {
		return errors.Wrapf(err, "failed to update auto backup config")
	}

	return nil
}

func (s *Service) GetClusterAutoBackupConfig(ctx context.Context, id int32, orgID int32) (*apigen.AutoBackupConfig, error) {
	c, err := s.m.GetAutoBackupConfig(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &apigen.AutoBackupConfig{
				Enabled: false,
			}, nil
		}
		return nil, errors.Wrapf(err, "failed to get auto backup config")
	}

	return &apigen.AutoBackupConfig{
		Enabled:        c.Enabled,
		CronExpression: c.CronExpression,
		KeepLast:       c.KeepLast,
	}, nil
}
