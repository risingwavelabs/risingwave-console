package prom

import (
	"fmt"
	"strings"
)

type Operation struct {
	Type OperationType
	Args string
}

type OperationType string

const (
	OpRate  OperationType = "rate"
	OpRange OperationType = "range"
)

type QueryBuilder struct {
	defaultLabels map[string]string
}

func NewQueryBuilder(defaultLabels map[string]string) *QueryBuilder {
	labels := make(map[string]string)
	if defaultLabels != nil {
		for k, v := range defaultLabels {
			labels[k] = v
		}
	}
	return &QueryBuilder{
		defaultLabels: labels,
	}
}

// NewQuery creates a new QueryBuilder instance. Note that
// the operation is not thread safe.
func (q *QueryBuilder) NewQuery(metric string) *Query {
	return &Query{
		metric: metric,
		labels: make(map[string]string),
	}
}

// QueryBuilder helps construct PromQL queries
type Query struct {
	metric     string
	labels     map[string]string
	operations []Operation
}

// WithLabel adds a label filter
func (q *Query) WithLabel(key, value string) *Query {
	q.labels[key] = value
	return q
}

// Range adds a time range
func (q *Query) Range(rng string) *Query {
	q.operations = append(q.operations, Operation{
		Type: OpRange,
		Args: rng,
	})
	return q
}

// Rate wraps the query in a rate function
func (q *Query) Rate() *Query {
	q.operations = append(q.operations, Operation{
		Type: OpRate,
		Args: "",
	})
	return q
}

// Build constructs the final PromQL query string
func (q *Query) Build() string {
	labelStr := ""
	if len(q.labels) > 0 {
		var parts []string
		for k, v := range q.labels {
			parts = append(parts, fmt.Sprintf(`%s="%s"`, k, v))
		}
		labelStr = "{" + strings.Join(parts, ", ") + "}"
	}

	query := fmt.Sprintf("%s%s", q.metric, labelStr)

	// Apply operations in order
	for _, op := range q.operations {
		switch op.Type {
		case OpRange:
			query = fmt.Sprintf("%s[%s]", query, op.Args)
		case OpRate:
			query = fmt.Sprintf("rate(%s)", query)
		}
	}

	return query
}
