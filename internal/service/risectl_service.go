package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
)

func (s *Service) getRisectlConn(ctx context.Context, id int32) (meta.RisectlConn, error) {
	cluster, err := s.m.GetClusterByID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	return s.risectlm.NewConn(ctx, cluster.Version, cluster.Host, cluster.MetaPort)
}

func (s *Service) RunRisectlCommand(ctx context.Context, id int32, params apigen.RisectlCommand, orgID int32) (*apigen.RisectlCommandResult, error) {
	conn, err := s.getRisectlConn(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get risectl connection")
	}

	stdout, stderr, exitCode, err := conn.Run(ctx, params.Args...)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return &apigen.RisectlCommandResult{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: int32(exitCode),
		Err:      errMsg,
	}, nil
}
