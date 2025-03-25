package service

import (
	"context"

	"github.com/risingwavelabs/wavekit/internal/apigen"
)

func (s *Service) CreateMetricsStore(ctx context.Context, req apigen.MetricsStoreCreate) (*apigen.MetricsStore, error) {
	return nil, nil
}

func (s *Service) DeleteMetricsStore(ctx context.Context, id int32, force bool) error {
	return nil
}

func (s *Service) GetMetricsStore(ctx context.Context, id int32) (*apigen.MetricsStore, error) {
	return nil, nil
}

func (s *Service) ListMetricsStores(ctx context.Context) ([]*apigen.MetricsStore, error) {
	return nil, nil
}

func (s *Service) UpdateMetricsStore(ctx context.Context, id int32, req apigen.MetricsStoreCreate) (*apigen.MetricsStore, error) {
	return nil, nil
}
