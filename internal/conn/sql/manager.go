package sql

import (
	"context"
	"fmt"

	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/zcore/model"
)

type SQLConnectionManegerInterface interface {
	GetConn(ctx context.Context, databaseID int32) (SQLConnectionInterface, error)
}

type SQLConnectionManager struct {
	m model.ModelInterface
}

func NewSQLConnectionManager(m model.ModelInterface) SQLConnectionManegerInterface {
	return &SQLConnectionManager{
		m: m,
	}
}

func (s *SQLConnectionManager) GetConn(ctx context.Context, databaseID int32) (SQLConnectionInterface, error) {
	databaseInfo, err := s.m.GetDatabaseConnectionByID(ctx, databaseID)
	if err != nil {
		return nil, err
	}

	clusterInfo, err := s.m.GetClusterByID(ctx, databaseInfo.ClusterID)
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", databaseInfo.Username, utils.UnwrapOrDefault(databaseInfo.Password, ""), clusterInfo.Host, clusterInfo.SqlPort, databaseInfo.Database)

	return &SimpleSQLConnection{
		connStr: connStr,
	}, nil
}
