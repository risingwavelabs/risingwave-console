package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
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
	token, err := s.auth.CreateToken(user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}
	refreshToken, jwtToken, err := s.auth.CreateRefreshToken(user.ID)
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
		RefreshToken: jwtToken,
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

	newRefreshToken, jwtToken, err := s.auth.CreateRefreshToken(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}
	if err = s.m.UpsertRefreshToken(ctx, querier.UpsertRefreshTokenParams{
		UserID: userID,
		Token:  newRefreshToken,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to upsert refresh token")
	}
	accessToken, err := s.auth.CreateToken(user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}

	return &apigen.Credentials{
		AccessToken:  accessToken,
		RefreshToken: jwtToken,
		TokenType:    apigen.Bearer,
	}, nil
}

func (s *Service) CreateNewUser(ctx context.Context, username, password string) (int32, error) {
	salt, hash, err := s.generateHashAndSalt(password)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to generate hash and salt")
	}
	var orgID int32
	if err := s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
		org, err := txm.CreateOrganization(ctx, fmt.Sprintf("%s's Org", username))
		if err != nil {
			return errors.Wrapf(err, "failed to create organization")
		}
		user, err := txm.CreateUser(ctx, querier.CreateUserParams{
			Name:           username,
			PasswordHash:   hash,
			PasswordSalt:   salt,
			OrganizationID: org.ID,
		})
		if err := txm.CreateOrganizationOwner(ctx, querier.CreateOrganizationOwnerParams{
			UserID:         user.ID,
			OrganizationID: org.ID,
		}); err != nil {
			return errors.Wrapf(err, "failed to create organization owner")
		}
		if err != nil {
			return errors.Wrapf(err, "failed to create user")
		}
		orgID = org.ID
		return nil
	}); err != nil {
		return 0, err
	}
	return orgID, nil
}
