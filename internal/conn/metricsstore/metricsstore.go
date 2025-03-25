package metricsstore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	prom_model "github.com/prometheus/common/model"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model"
)

var ErrMetricsStoreNotSupported = errors.New("Metrics store not supported")

// PrometheusConn is an API wrapper for prometheus,
// the connection is only established when the query is made.
type MetricsConn interface {
	GetMaterializedViewThroughput(ctx context.Context) (prom_model.Matrix, error)
}

type MetricsManager struct {
	model model.ModelInterface
}

func NewMetricsManager(m model.ModelInterface, cfg *config.Config) (*MetricsManager, error) {
	return &MetricsManager{
		model: m,
	}, nil
}

func (m *MetricsManager) GetMetricsConn(ctx context.Context, clusterID int32) (MetricsConn, error) {
	metricsStore, err := m.model.GetMetricsStore(ctx, clusterID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}
	selectors := ""
	if metricsStore.DefaultLabels != nil {
		for i, label := range metricsStore.DefaultLabels {
			if i > 0 {
				selectors += ","
			}
			var op string
			switch label.Op {
			case apigen.EQ:
				op = "="
			case apigen.NEQ:
				op = "!="
			case apigen.RE:
				op = "=~"
			case apigen.NRE:
				op = "!~"
			}
			selectors += fmt.Sprintf("%s%s\"%s\"", label.Key, op, label.Value)
		}
	}
	if metricsStore.Spec.Prometheus != nil {
		return NewPrometheusConn(metricsStore.Spec.Prometheus.Endpoint, selectors)
	}

	return nil, ErrMetricsStoreNotSupported
}
