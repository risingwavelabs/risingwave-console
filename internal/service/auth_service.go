package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

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

	if err := s.auth.InvalidateUserTokens(ctx, user.ID); err != nil {
		return nil, errors.Wrapf(err, "failed to invalidate user tokens")
	}

	keyID, token, err := s.auth.CreateToken(ctx, user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}

	refreshToken, err := s.auth.CreateRefreshToken(ctx, keyID, user.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}

	return &apigen.Credentials{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    apigen.Bearer,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, userID int32, refreshToken string) (*apigen.Credentials, error) {
	user, err := s.m.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user by id: %d", userID)
	}
	if err := s.auth.InvalidateUserTokens(ctx, userID); err != nil {
		return nil, errors.Wrapf(err, "failed to invalidate user tokens")
	}

	keyID, accessToken, err := s.auth.CreateToken(ctx, user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}

	newRefreshToken, err := s.auth.CreateRefreshToken(ctx, keyID, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}

	return &apigen.Credentials{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    apigen.Bearer,
	}, nil
}

func (s *Service) CreateNewUser(ctx context.Context, username, password string) (int32, error) {
	salt, hash, err := s.generateHashAndSalt(password)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to generate hash and salt")
	}

	org, err := s.m.CreateOrganization(ctx, fmt.Sprintf("%s's Org", username))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create organization")
	}

	user, err := s.m.CreateUser(ctx, querier.CreateUserParams{
		Name:           username,
		PasswordHash:   hash,
		PasswordSalt:   salt,
		OrganizationID: org.ID,
	})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create user")
	}

	if err := s.m.CreateOrganizationOwner(ctx, querier.CreateOrganizationOwnerParams{
		UserID:         user.ID,
		OrganizationID: org.ID,
	}); err != nil {
		return 0, errors.Wrapf(err, "failed to create organization owner")
	}

	return org.ID, nil
}
