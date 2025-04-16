package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

const defaultDiagnosticTaskTimeout = "30m"

func (s *Service) CreateClusterDiagnostic(ctx context.Context, id int32, orgID int32) (*apigen.DiagnosticData, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	content, err := s.metahttp.GetDiagnose(ctx, fmt.Sprintf("http://%s:%d", cluster.Host, cluster.HttpPort))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get diagnose")
	}
	diag, err := s.m.CreateClusterDiagnostic(ctx, querier.CreateClusterDiagnosticParams{
		ClusterID: cluster.ID,
		Content:   content,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster diagnostic")
	}
	return &apigen.DiagnosticData{
		ID:        diag.ID,
		CreatedAt: diag.CreatedAt,
		Content:   diag.Content,
	}, nil
}

func (s *Service) ListClusterDiagnostics(ctx context.Context, id int32, orgID int32) ([]apigen.DiagnosticData, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	diagnostics, err := s.m.ListClusterDiagnostics(ctx, cluster.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list cluster diagnostics")
	}

	result := make([]apigen.DiagnosticData, len(diagnostics))
	for i, diagnostic := range diagnostics {
		result[i] = apigen.DiagnosticData{
			ID:        diagnostic.ID,
			CreatedAt: diagnostic.CreatedAt,
		}
	}
	return result, nil
}

func (s *Service) GetClusterDiagnostic(ctx context.Context, id int32, diagnosticID int32, orgID int32) (*apigen.DiagnosticData, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	diagnostic, err := s.m.GetClusterDiagnostic(ctx, diagnosticID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster diagnostic")
	}

	if diagnostic.ClusterID != cluster.ID {
		return nil, ErrDiagnosticNotFound
	}

	return &apigen.DiagnosticData{
		ID:        diagnostic.ID,
		CreatedAt: diagnostic.CreatedAt,
		Content:   diagnostic.Content,
	}, nil
}

func (s *Service) UpdateClusterAutoDiagnosticConfig(ctx context.Context, id int32, params apigen.AutoDiagnosticConfig, orgID int32) error {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to get cluster")
	}

	taskSpec := apigen.TaskSpec{
		Type: apigen.AutoDiagnostic,
		AutoDiagnostic: &apigen.TaskSpecAutoDiagnostic{
			ClusterID:         cluster.ID,
			RetentionDuration: params.RetentionDuration,
		},
	}

	c, err := s.m.GetAutoDiagnosticsConfig(ctx, cluster.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No existing auto backup config, create a new one and a new cron job
			if err := s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
				c := modelctx.NewModelctx(ctx, txm)
				taskID, err := s.taskstore.CreateCronJob(c, utils.Ptr(defaultDiagnosticTaskTimeout), &orgID, params.CronExpression, taskSpec)
				if err != nil {
					return errors.Wrapf(err, "failed to create cron job")
				}
				if err := txm.CreateAutoDiagnosticsConfig(ctx, querier.CreateAutoDiagnosticsConfigParams{
					ClusterID: cluster.ID,
					TaskID:    taskID,
					Enabled:   true,
				}); err != nil {
					return errors.Wrapf(err, "failed to create auto diagnostics config")
				}
				return nil
			}); err != nil {
				return errors.Wrapf(err, "failed to create new cluster auto diagnostics config")
			}
			return nil
		}
		return errors.Wrapf(err, "failed to get auto diagnostics config")
	}

	if err := s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
		mc := modelctx.NewModelctx(ctx, txm)
		if !params.Enabled {
			if err := s.taskstore.PauseCronJob(mc, c.TaskID); err != nil {
				return errors.Wrapf(err, "failed to pause cron job")
			}
		} else {
			if err := s.taskstore.ResumeCronJob(mc, c.TaskID); err != nil {
				return errors.Wrapf(err, "failed to resume cron job")
			}
		}
		if err := txm.UpdateAutoDiagnosticsConfig(ctx, querier.UpdateAutoDiagnosticsConfigParams{
			ClusterID: cluster.ID,
			Enabled:   params.Enabled,
		}); err != nil {
			return errors.Wrapf(err, "failed to update auto diagnostics config")
		}
		if err := s.taskstore.UpdateCronJob(mc, c.TaskID, utils.Ptr(defaultDiagnosticTaskTimeout), &orgID, params.CronExpression, taskSpec); err != nil {
			return errors.Wrapf(err, "failed to update cron job")
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to update cluster auto diagnostics config")
	}
	return nil
}

func (s *Service) GetClusterAutoDiagnosticConfig(ctx context.Context, id int32, orgID int32) (*apigen.AutoDiagnosticConfig, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}
	c, err := s.m.GetAutoDiagnosticsConfig(ctx, cluster.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &apigen.AutoDiagnosticConfig{
				Enabled: false,
			}, nil
		}
	}
	task, err := s.m.GetTaskByID(ctx, c.TaskID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get task")
	}

	return &apigen.AutoDiagnosticConfig{
		Enabled:           c.Enabled,
		CronExpression:    task.Attributes.Cronjob.CronExpression,
		RetentionDuration: task.Spec.AutoDiagnostic.RetentionDuration,
	}, nil
}
