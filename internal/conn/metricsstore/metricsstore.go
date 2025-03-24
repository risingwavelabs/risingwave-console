package metricsstore

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	prom_model "github.com/prometheus/common/model"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model"
)

var ErrNoMetricsStoreFound = errors.New("no metrics store found")

// PrometheusConn is an API wrapper for prometheus,
// the connection is only established when the query is made.
type MetricsConn interface {
	GetMaterializedViewThroughput(ctx context.Context) (prom_model.Matrix, error)
}

type MetricsManager struct {
	model         model.ModelInterface
	defaultLabels map[string]string
}

func NewMetricsManager(m model.ModelInterface, cfg *config.Config) (*MetricsManager, error) {
	defaultLabels := make(map[string]string)

	for _, label := range strings.Split(cfg.Prometheus.DefaultLabels, ",") {
		parts := strings.Split(label, "=")
		if len(parts) == 2 {
			defaultLabels[parts[0]] = parts[1]
		}
	}

	return &MetricsManager{
		model:         m,
		defaultLabels: defaultLabels,
	}, nil
}

func (m *MetricsManager) GetMetricsConn(ctx context.Context, clusterID int32) (MetricsConn, error) {
	metricsStore, err := m.model.GetMetricsStore(ctx, clusterID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}
	if metricsStore.Spec.Prometheus != nil {
		return NewPrometheusConn(metricsStore.Spec.Prometheus.Endpoint, m.defaultLabels)
	}

	return nil, ErrNoMetricsStoreFound
}
