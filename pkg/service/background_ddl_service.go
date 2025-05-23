package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/pkg/zgen/apigen"
)

const getDDLProgressSQL = `SELECT * FROM rw_ddl_progress ORDER BY initialized_at DESC`

func (s *Service) GetDDLProgress(ctx context.Context, id int32, orgID int32) ([]apigen.DDLProgress, error) {
	conn, err := s.sqlm.GetConn(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get database connection")
	}

	result, err := conn.Query(ctx, getDDLProgressSQL, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get DDL progress")
	}

	progress := []apigen.DDLProgress{}
	for _, row := range result.Rows {
		item := apigen.DDLProgress{
			ID:        row["ddl_id"].(int64),
			Statement: row["ddl_statement"].(string),
			Progress:  row["progress"].(string),
		}
		if t, ok := row["initialized_at"].(time.Time); ok {
			item.InitializedAt = &t
		}
		progress = append(progress, item)
	}
	return progress, nil
}

func (s *Service) CancelDDLProgress(ctx context.Context, id int32, ddlID int64, orgID int32) error {
	conn, err := s.sqlm.GetConn(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to get database connection")
	}

	_, err = conn.Query(ctx, fmt.Sprintf("CANCEL JOB %d", ddlID), false)
	if err != nil {
		return errors.Wrapf(err, "failed to cancel DDL progress")
	}

	return nil
}
