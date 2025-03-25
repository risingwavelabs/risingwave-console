package metricsstore

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prom_model "github.com/prometheus/common/model"
	"go.uber.org/zap"

	"github.com/risingwavelabs/wavekit/internal/logger"
)

var log = logger.NewLogAgent("metricsstore")

var ErrPrometheusEndpointNotFound = errors.New("prometheus endpoint not found")

type PrometheusConn struct {
	v1api           v1.API
	defaultSelector string
}

func NewPrometheusConn(endpoint string, defaultSelector string) (*PrometheusConn, error) {
	client, err := api.NewClient(api.Config{
		Address: endpoint,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create prometheus client")
	}

	return &PrometheusConn{
		v1api:           v1.NewAPI(client),
		defaultSelector: defaultSelector,
	}, nil
}

func (c *PrometheusConn) GetMaterializedViewThroughput(ctx context.Context) (prom_model.Matrix, error) {
	rate := "1m"
	query := fmt.Sprintf(`sum(rate(%s[%s])) by (table_id) * on(table_id) group_left(table_name) group(%s) by (table_id, table_name)`,
		metricWithLabelSelector("stream_mview_input_row_count", c.defaultSelector),
		rate,
		metricWithLabelSelector("table_info", c.defaultSelector))
	result, warnings, err := c.v1api.QueryRange(ctx, query, v1.Range{
		Start: time.Now().Add(-1 * time.Minute),
		End:   time.Now(),
		Step:  time.Second * 5,
	})

	log.Info("prometheus query", zap.String("query", query), zap.String("defaultSelector", c.defaultSelector))

	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		return nil, errors.New(strings.Join(warnings, "\n"))
	}

	// Check if the result is a matrix type
	if result.Type() == prom_model.ValMatrix {
		return result.(prom_model.Matrix), nil
	}

	return nil, errors.New("result is not a matrix")
}

func metricWithLabelSelector(metric string, selector string) string {
	if selector == "" {
		return metric
	}
	return fmt.Sprintf("%s{%s}", metric, selector)
}
