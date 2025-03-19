package service

import (
	"context"

	prom_model "github.com/prometheus/common/model"
)

func (s *Service) GetMaterializedViewThroughput(ctx context.Context, clusterID int32) (prom_model.Matrix, error) {
	conn, err := s.promm.GetPrometheusConn(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	return conn.GetMaterializedViewThroughput(ctx)
}
