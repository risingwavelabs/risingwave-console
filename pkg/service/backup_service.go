package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudcarver/anchor/pkg/taskcore"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/risingwave-console/pkg/zcore/model"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/apigen"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/querier"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/taskgen"
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
		ID:    id,
		OrgID: orgID,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to get cluster")
	}

	orgSettings, err := s.m.GetOrgSettings(ctx, orgID)
	if err != nil {
		return errors.Wrapf(err, "failed to get organization")
	}
	cronExpression := fmt.Sprintf("CRON_TZ=%s %s", orgSettings.Timezone, params.CronExpression)

	c, err := s.m.GetAutoBackupConfig(ctx, cluster.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No existing auto backup config, create a new one and a new cron job
			if err := s.m.RunTransactionWithTx(ctx, func(tx pgx.Tx, txm model.ModelInterface) error {
				taskID, err := s.taskRunner.RunAutoBackupWithTx(ctx, tx, &taskgen.AutoBackupParameters{
					ClusterID:         cluster.ID,
					RetentionDuration: params.RetentionDuration,
				}, taskcore.WithCronjob(cronExpression))
				if err != nil {
					return errors.Wrapf(err, "failed to create cron job")
				}

				if err := txm.CreateAutoBackupConfig(ctx, querier.CreateAutoBackupConfigParams{
					ClusterID: cluster.ID,
					TaskID:    taskID,
					Enabled:   true,
				}); err != nil {
					return errors.Wrapf(err, "failed to create auto backup config")
				}
				return nil
			}); err != nil {
				return errors.Wrapf(err, "failed to create new cluster auto backup config")
			}
			return nil
		}
		return errors.Wrapf(err, "failed to get auto backup config")
	}

	if err := s.m.RunTransactionWithTx(ctx, func(tx pgx.Tx, txm model.ModelInterface) error {
		txTaskstore := s.taskstore.WithTx(tx)

		if !params.Enabled {
			if err := txTaskstore.PauseCronJob(ctx, c.TaskID); err != nil {
				return errors.Wrapf(err, "failed to pause cron job")
			}
		} else {
			if err := txTaskstore.ResumeCronJob(ctx, c.TaskID); err != nil {
				return errors.Wrapf(err, "failed to resume cron job")
			}
		}

		if err := txm.UpdateAutoBackupConfig(ctx, querier.UpdateAutoBackupConfigParams{
			ClusterID: cluster.ID,
			Enabled:   params.Enabled,
		}); err != nil {
			return errors.Wrapf(err, "failed to update auto backup config")
		}

		taskParams := taskgen.AutoBackupParameters{
			ClusterID:         cluster.ID,
			RetentionDuration: params.RetentionDuration,
		}

		spec, err := taskParams.Marshal()
		if err != nil {
			return errors.Wrapf(err, "failed to marshal task parameters")
		}

		if err := txTaskstore.UpdateCronJob(ctx, c.TaskID, cronExpression, spec); err != nil {
			return errors.Wrapf(err, "failed to update cron job")
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to update cluster auto backup config")
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
	task, err := s.anchorSvc.GetTaskByID(ctx, c.TaskID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get task")
	}

	var params taskgen.AutoBackupParameters
	if err := params.Parse(task.Spec.Payload); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal task spec")
	}

	return &apigen.AutoBackupConfig{
		Enabled:           c.Enabled,
		CronExpression:    task.Attributes.Cronjob.CronExpression,
		RetentionDuration: params.RetentionDuration,
	}, nil
}
