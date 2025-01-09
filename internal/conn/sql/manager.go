package sql

import (
	"context"
	"fmt"
	"sync"

	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

type SQLConnectionManegerInterface interface {
	GetConn(ctx context.Context, databaseID int32) (SQLConnectionInterface, error)
}

type SQLConnectionManager struct {
	mu    sync.RWMutex
	m     model.ModelInterface
	conns map[int32]SQLConnectionInterface
}

func NewSQLConnectionManager(m model.ModelInterface) SQLConnectionManegerInterface {
	return &SQLConnectionManager{
		m:     m,
		conns: make(map[int32]SQLConnectionInterface),
	}
}

func (s *SQLConnectionManager) NewConn(ctx context.Context, databaseID int32) (SQLConnectionInterface, error) {
	databaseInfo, err := s.m.GetDatabaseConnection(ctx, databaseID)
	if err != nil {
		return nil, err
	}

	clusterInfo, err := s.m.GetCluster(ctx, databaseInfo.ClusterID)
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", databaseInfo.Username, utils.UnwrapOrDefault(databaseInfo.Password, ""), clusterInfo.Host, clusterInfo.SqlPort, databaseInfo.Database)

	s.mu.Lock()
	s.conns[databaseID] = &SimpleSQLConnection{
		connStr: connStr,
	}
	s.mu.Unlock()
	return &SimpleSQLConnection{
		connStr: connStr,
	}, nil
}

func (s *SQLConnectionManager) GetConn(ctx context.Context, databaseID int32) (SQLConnectionInterface, error) {
	s.mu.RLock()
	conn, ok := s.conns[databaseID]
	if ok {
		return conn, nil
	}
	s.mu.RUnlock()

	conn, err := s.NewConn(ctx, databaseID)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
