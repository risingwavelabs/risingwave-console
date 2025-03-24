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
)

var ErrPrometheusEndpointNotFound = errors.New("prometheus endpoint not found")

type PrometheusConn struct {
	v1api         v1.API
	defaultLabels map[string]string
}

func NewPrometheusConn(endpoint string, defaultLabels map[string]string) (*PrometheusConn, error) {
	client, err := api.NewClient(api.Config{
		Address: endpoint,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create prometheus client")
	}

	return &PrometheusConn{
		v1api:         v1.NewAPI(client),
		defaultLabels: defaultLabels,
	}, nil
}

func (c *PrometheusConn) GetMaterializedViewThroughput(ctx context.Context) (prom_model.Matrix, error) {
	rate := "1m"
	query := fmt.Sprintf(`sum(rate(%s[%s])) by (table_id) * on(table_id) group_left(table_name) group(%s) by (table_id, table_name)`,
		metricWithLabels("stream_mview_input_row_count", c.defaultLabels),
		rate,
		metricWithLabels("table_info", c.defaultLabels))
	result, warnings, err := c.v1api.QueryRange(ctx, query, v1.Range{
		Start: time.Now().Add(-1 * time.Minute),
		End:   time.Now(),
		Step:  time.Second * 5,
	})
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

func metricWithLabels(metric string, labels map[string]string) string {
	labelStr := ""
	for k, v := range labels {
		labelStr += fmt.Sprintf("%s=%s,", k, v)
	}
	if labelStr == "" {
		return metric
	}
	return fmt.Sprintf("%s{%s}", metric, labelStr)
}
