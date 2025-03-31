package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
)

func (s *Service) CreateClusterDiagnostic(ctx context.Context, id int32, orgID int32) (*apigen.DiagnosticData, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	conn, err := s.getMetaHttpConn(ctx, cluster.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get meta http connection")
	}
	content, err := conn.GetDiagnose(ctx)
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
			if err := s.Self().CreateCronJob(ctx, &orgID, params.CronExpression, taskSpec); err != nil {
				return errors.Wrapf(err, "failed to create cron job")
			}
			return nil
		}
		return errors.Wrapf(err, "failed to get auto diagnostics config")
	}

	if err := s.Self().UpdateCronJob(ctx, c.TaskID, &orgID, params.CronExpression, taskSpec); err != nil {
		return errors.Wrapf(err, "failed to update cron job")
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
