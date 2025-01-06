package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

type InitService struct {
	m model.ModelInterface
}

func NewInitService(m model.ModelInterface) *InitService {
	return &InitService{
		m: m,
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
	hash, salt, err := utils.GenerateHashAndSalt(cfg.Root.Password)
	if err != nil {
		return errors.Wrapf(err, "failed to generate hash and salt")
	}
	_, err = s.m.CreateUser(ctx, querier.CreateUserParams{
		Name:         "root",
		PasswordHash: hash,
		PasswordSalt: salt,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create root user")
	}
	return nil
}
