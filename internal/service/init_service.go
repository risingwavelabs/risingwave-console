package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model"
)

type InitService struct {
	m model.ModelInterface
	s ServiceInterface
}

func NewInitService(m model.ModelInterface, s ServiceInterface) *InitService {
	return &InitService{
		m: m,
		s: s,
	}
}

func (s *InitService) Init(ctx context.Context, cfg *config.Config) error {
	// remove the root user if it is not set in the config
	if cfg.Root == nil {
		if err := s.m.DeleteUserByName(ctx, "root"); err != nil {
			return errors.Wrapf(err, "failed to delete root user")
		}
		return nil
	}

	// create the root user
	if err := s.s.CreateNewUser(ctx, "root", cfg.Root.Password); err != nil {
		return errors.Wrapf(err, "failed to create root user")
	}
	return nil
}
