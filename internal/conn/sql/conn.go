package sql

import (
	"context"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
)

var ErrQueryFailed = errors.New("query failed")

type Column struct {
	Name string
	Type string
}

type Result struct {
	RowsAffected int64
	Columns      []Column
	Rows         []map[string]any
}

type SQLConnectionInterface interface {
	Query(context.Context, string) (*Result, error)
}

type SimpleSQLConnection struct {
	connStr string
}

func (s *SimpleSQLConnection) Query(ctx context.Context, query string) (*Result, error) {
	return Query(ctx, s.connStr, query)
}

func Query(ctx context.Context, connStr string, query string) (*Result, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(ErrQueryFailed, err.Error())
	}
	defer rows.Close()

	fieldDescs := rows.FieldDescriptions()
	columns := make([]Column, len(fieldDescs))
	for i, d := range fieldDescs {
		columns[i] = Column{
			Name: string(d.Name),
			Type: getDataTypeName(d.DataTypeOID),
		}
	}

	var result []map[string]any
	for rows.Next() {
		values := make([]any, len(fieldDescs))
		scanArgs := make([]any, len(fieldDescs))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(fieldDescs))
		for i, col := range fieldDescs {
			row[string(col.Name)] = values[i]
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(ErrQueryFailed, err.Error())
	}

	return &Result{
		RowsAffected: rows.CommandTag().RowsAffected(),
		Columns:      columns,
		Rows:         result,
	}, nil
}
