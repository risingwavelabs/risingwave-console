package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

type (
	TradeType   string
	TradeStatus string
	DdWorkEvent string
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

const (
	ExpireDuration             = 2 * time.Minute
	DefaultMaxRetries          = 3
	RefreshTokenExpireDuration = 14 * 24 * time.Hour
)

type ServiceInterface interface {
	// SignIn authenticates a user and returns credentials
	SignIn(ctx context.Context, params apigen.SignInRequest) (*apigen.Credentials, error)

	// RefreshToken refreshes an authentication token using a refresh token
	RefreshToken(ctx context.Context, userID int32, refreshToken string) (*apigen.Credentials, error)
}

type Service struct {
	m           model.ModelInterface
	authService AuthServiceInterface

	now                 func() time.Time
	generateHashAndSalt func(password string) (string, string, error)
}

func NewService(cfg *config.Config, m model.ModelInterface, authService AuthServiceInterface) ServiceInterface {
	return &Service{
		m:                   m,
		now:                 time.Now,
		generateHashAndSalt: utils.GenerateHashAndSalt,
		authService:         authService,
	}
}

func (s *Service) SignIn(ctx context.Context, params apigen.SignInRequest) (*apigen.Credentials, error) {
	user, err := s.m.GetUserByName(ctx, params.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrapf(ErrUserNotFound, "user %s not found", params.Name)
		}
		return nil, errors.Wrapf(err, "failed to get user by name")
	}
	input, err := utils.HashPassword(params.Password, user.PasswordSalt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hash password")
	}
	if input != user.PasswordHash {
		return nil, ErrInvalidPassword
	}
	token, err := s.authService.CreateToken(user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}
	refreshToken, err := s.authService.GenerateRefreshToken()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}
	err = s.m.UpsertRefreshToken(ctx, querier.UpsertRefreshTokenParams{
		UserID: user.ID,
		Token:  refreshToken,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert refresh token")
	}

	return &apigen.Credentials{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    apigen.Bearer,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, userID int32, refreshToken string) (*apigen.Credentials, error) {
	originalRefreshToken, err := s.m.GetRefreshToken(ctx, querier.GetRefreshTokenParams{
		UserID: userID,
		Token:  refreshToken,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get refresh token")
	}
	if originalRefreshToken.UpdatedAt.Add(RefreshTokenExpireDuration).Before(s.now()) {
		return nil, ErrRefreshTokenExpired
	}
	user, err := s.m.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user by id: %d", userID)
	}

	newRefreshToken, err := s.authService.GenerateRefreshToken()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}
	if err = s.m.UpsertRefreshToken(ctx, querier.UpsertRefreshTokenParams{
		UserID: userID,
		Token:  newRefreshToken,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to upsert refresh token")
	}
	accessToken, err := s.authService.CreateToken(user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}

	return &apigen.Credentials{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    apigen.Bearer,
	}, nil
}
