package prom

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type PromConn interface {
	GetMaterializedViewThroughput(ctx context.Context) (float64, error)
}

type PrometheusConn struct {
	v1api        v1.API
	queryBuilder *QueryBuilder
}

func NewPrometheusConn(endpoint string, defaultLabels map[string]string) (PromConn, error) {
	client, err := api.NewClient(api.Config{
		Address: endpoint,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create prometheus client")
	}
	v1api := v1.NewAPI(client)
	return &PrometheusConn{v1api: v1api, queryBuilder: NewQueryBuilder(defaultLabels)}, nil
}

func (c *PrometheusConn) GetMaterializedViewThroughput(ctx context.Context) (float64, error) {
	query := c.queryBuilder.
		NewQuery("risingwave_materialized_view_throughput").
		WithLabel("view_name", "mv_1").
		Range("1m").
		Rate()
	c.v1api.QueryRange(ctx, query, v1.Range{
		Start: time.Now().Add(-1 * time.Minute),
		End:   time.Now(),
		Step:  time.Second * 10,
	})
}
