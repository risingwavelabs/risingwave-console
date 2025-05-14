package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/zgen/apigen"
	"github.com/risingwavelabs/wavekit/internal/zgen/querier"
)

func (s *Service) TestDatabaseConnection(ctx context.Context, params apigen.TestDatabaseConnectionPayload, orgID int32) (*apigen.TestDatabaseConnectionResult, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:    params.ClusterID,
		OrgID: orgID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", params.Username, utils.UnwrapOrDefault(params.Password, ""), cluster.Host, cluster.SqlPort, params.Database)

	_, err = sql.Query(ctx, connStr, "SELECT 1", false)
	if err != nil {
		return &apigen.TestDatabaseConnectionResult{
			Success: false,
			Result:  err.Error(),
		}, nil
	}

	return &apigen.TestDatabaseConnectionResult{
		Success: true,
		Result:  "Connection successful",
	}, nil
}

func (s *Service) QueryDatabase(ctx context.Context, id int32, params apigen.QueryRequest, orgID int32, backgroundDDL bool) (*apigen.QueryResponse, error) {
	db, err := s.m.GetOrgDatabaseByID(ctx, querier.GetOrgDatabaseByIDParams{
		ID:    id,
		OrgID: orgID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrDatabaseNotFound
		}
		if errors.Is(err, sql.ErrQueryFailed) {
			return &apigen.QueryResponse{
				Error: utils.Ptr(err.Error()),
			}, nil
		}
		return nil, errors.Wrapf(err, "failed to get database connection")
	}

	conn, err := s.sqlm.GetConn(ctx, db.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get database connection")
	}

	result, err := conn.Query(ctx, params.Query, backgroundDDL)
	if err != nil {
		if errors.Is(err, sql.ErrQueryFailed) {
			return &apigen.QueryResponse{
				Error: utils.Ptr(err.Error()),
			}, nil
		}
		return nil, errors.Wrapf(err, "failed to query database")
	}

	columns := make([]apigen.Column, len(result.Columns))
	for i, column := range result.Columns {
		columns[i] = apigen.Column{
			Name: column.Name,
			Type: column.Type,
		}
	}

	return &apigen.QueryResponse{
		Columns: columns,
		Rows:    result.Rows,
	}, nil
}
