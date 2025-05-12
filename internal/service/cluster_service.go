package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/zcore/model"
	"github.com/risingwavelabs/wavekit/internal/zgen/apigen"
	"github.com/risingwavelabs/wavekit/internal/zgen/querier"
	"golang.org/x/mod/semver"
)

func (s *Service) TestClusterConnection(ctx context.Context, params apigen.TestClusterConnectionPayload, orgID int32) (*apigen.TestClusterConnectionResult, error) {
	errMsg := ""
	if err := utils.TestTCPConnection(ctx, params.Host, params.MetaPort, 5*time.Second); err != nil {
		errMsg += fmt.Sprintf("Failed to connect to meta port: %s\n", err.Error())
	}
	if err := utils.TestTCPConnection(ctx, params.Host, params.SqlPort, 5*time.Second); err != nil {
		errMsg += fmt.Sprintf("Failed to connect to sql port: %s\n", err.Error())
	}
	if err := utils.TestTCPConnection(ctx, params.Host, params.HttpPort, 5*time.Second); err != nil {
		errMsg += fmt.Sprintf("Failed to connect to http port: %s\n", err.Error())
	}

	return &apigen.TestClusterConnectionResult{
		Success: errMsg == "",
		Result:  utils.IfElse(errMsg == "", "Connection successful", errMsg),
	}, nil
}

func (s *Service) ListClusterVersions(ctx context.Context) ([]string, error) {
	versions, err := s.risectlm.ListVersions(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list cluster versions")
	}
	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare(versions[i], versions[j]) > 0
	})
	return versions, nil
}

func (s *Service) ImportCluster(ctx context.Context, params apigen.ClusterImport, orgID int32) (*apigen.Cluster, error) {
	cluster, err := s.m.CreateCluster(ctx, querier.CreateClusterParams{
		OrganizationID: orgID,
		Name:           params.Name,
		Host:           params.Host,
		SqlPort:        params.SqlPort,
		MetaPort:       params.MetaPort,
		HttpPort:       params.HttpPort,
		Version:        params.Version,
		MetricsStoreID: params.MetricsStoreID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster")
	}

	return clusterToApi(cluster), nil
}

func (s *Service) GetCluster(ctx context.Context, id int32, orgID int32) (*apigen.Cluster, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrClusterNotFound
		}
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	return clusterToApi(cluster), nil
}

func (s *Service) ListClusters(ctx context.Context, orgID int32) ([]*apigen.Cluster, error) {
	clusters, err := s.m.ListOrgClusters(ctx, orgID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to list clusters")
	}

	result := make([]*apigen.Cluster, len(clusters))
	for i, cluster := range clusters {
		result[i] = clusterToApi(cluster)
	}
	return result, nil
}

func (s *Service) UpdateCluster(ctx context.Context, id int32, params apigen.ClusterImport, orgID int32) (*apigen.Cluster, error) {
	cluster, err := s.m.UpdateOrgCluster(ctx, querier.UpdateOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
		Name:           params.Name,
		Host:           params.Host,
		Version:        params.Version,
		SqlPort:        int32(params.SqlPort),
		MetaPort:       int32(params.MetaPort),
		HttpPort:       int32(params.HttpPort),
		MetricsStoreID: params.MetricsStoreID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrClusterNotFound
		}
		return nil, errors.Wrapf(err, "failed to update cluster")
	}

	return clusterToApi(cluster), nil
}

func (s *Service) DeleteCluster(ctx context.Context, id int32, cascade bool, orgID int32) error {
	if cascade {
		return s.deleteClusterCacasde(ctx, id, orgID)
	}
	return s.deleteClusterNonCacasde(ctx, id, orgID)
}

func (s *Service) deleteClusterCacasde(ctx context.Context, id int32, orgID int32) error {
	return s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
		if err := txm.DeleteAllOrgDatabaseConnectionsByClusterID(ctx, querier.DeleteAllOrgDatabaseConnectionsByClusterIDParams{
			ClusterID:      id,
			OrganizationID: orgID,
		}); err != nil {
			return errors.Wrapf(err, "failed to delete associated database connections")
		}

		if err := txm.DeleteOrgCluster(ctx, querier.DeleteOrgClusterParams{
			ID:             id,
			OrganizationID: orgID,
		}); err != nil {
			return errors.Wrapf(err, "failed to delete cluster")
		}
		return nil
	})
}

func (s *Service) deleteClusterNonCacasde(ctx context.Context, id int32, orgID int32) error {
	dbConnections, err := s.m.GetAllOrgDatabseConnectionsByClusterID(ctx, querier.GetAllOrgDatabseConnectionsByClusterIDParams{
		ClusterID:      id,
		OrganizationID: orgID,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to get database connections")
	}
	names := make([]string, len(dbConnections))
	for i, db := range dbConnections {
		names[i] = db.Name
	}

	if len(dbConnections) > 0 {
		return errors.Wrapf(ErrClusterHasDatabaseConnections, "cluster has %d database connections: %s", len(dbConnections), strings.Join(names, ", "))
	}

	err = s.m.DeleteOrgCluster(ctx, querier.DeleteOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to delete cluster")
	}
	return nil
}
