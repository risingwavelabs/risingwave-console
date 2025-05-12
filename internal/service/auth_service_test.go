package service

import (
	"context"
	"testing"

	"github.com/cloudcarver/anchor/pkg/auth"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/zcore/model"
	"github.com/risingwavelabs/wavekit/internal/zgen/apigen"
	"github.com/risingwavelabs/wavekit/internal/zgen/querier"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModel := model.NewMockModelInterface(ctrl)
	mockAuth := auth.NewMockAuthInterface(ctrl)

	service := &Service{
		m:    mockModel,
		auth: mockAuth,
	}

	ctx := context.Background()
	userID := int32(1)
	username := "testuser"
	password := "testpass"
	salt := "salt"

	hash, err := utils.HashPassword(password, salt)
	require.NoError(t, err)
	keyID := int64(123)
	accessToken := "access_token"
	refreshToken := "refresh_token"

	user := &querier.User{
		ID:           userID,
		Name:         username,
		PasswordHash: hash,
		PasswordSalt: salt,
	}

	testCases := []struct {
		name          string
		params        apigen.SignInRequest
		setupMock     func()
		expectedError error
	}{
		{
			name: "successful sign in",
			params: apigen.SignInRequest{
				Name:     username,
				Password: password,
			},
			setupMock: func() {
				mockModel.EXPECT().GetUserByName(ctx, username).Return(user, nil)
				mockAuth.EXPECT().InvalidateUserTokens(ctx, userID).Return(nil)
				mockAuth.EXPECT().CreateToken(ctx, user.ID, nil).Return(keyID, accessToken, nil)
				mockAuth.EXPECT().CreateRefreshToken(ctx, keyID, userID).Return(refreshToken, nil)
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			params: apigen.SignInRequest{
				Name:     username,
				Password: password,
			},
			setupMock: func() {
				mockModel.EXPECT().GetUserByName(ctx, username).Return(nil, pgx.ErrNoRows)
			},
			expectedError: ErrUserNotFound,
		},
		{
			name: "invalid password",
			params: apigen.SignInRequest{
				Name:     username,
				Password: "wrongpass",
			},
			setupMock: func() {
				mockModel.EXPECT().GetUserByName(ctx, username).Return(user, nil)
			},
			expectedError: ErrInvalidPassword,
		},
		{
			name: "failed to invalidate tokens",
			params: apigen.SignInRequest{
				Name:     username,
				Password: password,
			},
			setupMock: func() {
				mockModel.EXPECT().GetUserByName(ctx, username).Return(user, nil)
				mockAuth.EXPECT().InvalidateUserTokens(ctx, userID).Return(errors.New("invalidation failed"))
			},
			expectedError: errors.New("failed to invalidate user tokens"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			creds, err := service.SignIn(ctx, tc.params)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
				require.Nil(t, creds)
			} else {
				require.NoError(t, err)
				require.NotNil(t, creds)
				require.Equal(t, accessToken, creds.AccessToken)
				require.Equal(t, refreshToken, creds.RefreshToken)
				require.Equal(t, apigen.Bearer, creds.TokenType)
			}
		})
	}
}

func TestService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModel := model.NewMockModelInterface(ctrl)
	mockAuth := auth.NewMockAuthInterface(ctrl)

	service := &Service{
		m:    mockModel,
		auth: mockAuth,
	}

	ctx := context.Background()
	userID := int32(1)
	keyID := int64(123)
	accessToken := "new_access_token"
	refreshToken := "new_refresh_token"

	user := &querier.User{
		ID:   userID,
		Name: "testuser",
	}

	testCases := []struct {
		name          string
		userID        int32
		refreshToken  string
		setupMock     func()
		expectedError error
	}{
		{
			name:         "successful token refresh",
			userID:       userID,
			refreshToken: "old_refresh_token",
			setupMock: func() {
				mockModel.EXPECT().GetUser(ctx, userID).Return(user, nil)
				mockAuth.EXPECT().InvalidateUserTokens(ctx, userID).Return(nil)
				mockAuth.EXPECT().CreateToken(ctx, user.ID, nil).Return(keyID, accessToken, nil)
				mockAuth.EXPECT().CreateRefreshToken(ctx, keyID, userID).Return(refreshToken, nil)
			},
			expectedError: nil,
		},
		{
			name:         "user not found",
			userID:       userID,
			refreshToken: "old_refresh_token",
			setupMock: func() {
				mockModel.EXPECT().GetUser(ctx, userID).Return(nil, pgx.ErrNoRows)
			},
			expectedError: errors.New("failed to get user by id"),
		},
		{
			name:         "failed to invalidate tokens",
			userID:       userID,
			refreshToken: "old_refresh_token",
			setupMock: func() {
				mockModel.EXPECT().GetUser(ctx, userID).Return(user, nil)
				mockAuth.EXPECT().InvalidateUserTokens(ctx, userID).Return(errors.New("invalidation failed"))
			},
			expectedError: errors.New("failed to invalidate user tokens"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			creds, err := service.RefreshToken(ctx, tc.userID, tc.refreshToken)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
				require.Nil(t, creds)
			} else {
				require.NoError(t, err)
				require.NotNil(t, creds)
				require.Equal(t, accessToken, creds.AccessToken)
				require.Equal(t, refreshToken, creds.RefreshToken)
				require.Equal(t, apigen.Bearer, creds.TokenType)
			}
		})
	}
}

func TestService_CreateNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockModel := model.NewMockModelInterfaceWithTransaction(ctrl)

	service := &Service{
		m:                   mockModel,
		generateHashAndSalt: utils.GenerateHashAndSalt,
	}

	ctx := context.Background()
	username := "newuser"
	password := "newpass"
	orgID := int32(1)
	userID := int32(2)

	testCases := []struct {
		name          string
		username      string
		password      string
		setupMock     func()
		expectedError error
	}{
		{
			name:     "successful user creation",
			username: username,
			password: password,
			setupMock: func() {
				org := &querier.Organization{ID: orgID}
				user := &querier.User{ID: userID}

				mockModel.EXPECT().CreateOrganization(gomock.Any(), gomock.Any()).Return(org, nil)
				mockModel.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(user, nil)
				mockModel.EXPECT().CreateOrganizationOwner(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "failed to create organization",
			username: username,
			password: password,
			setupMock: func() {
				mockModel.EXPECT().CreateOrganization(gomock.Any(), gomock.Any()).Return(nil, errors.New("org creation failed"))
			},
			expectedError: errors.New("failed to create organization"),
		},
		{
			name:     "failed to create user",
			username: username,
			password: password,
			setupMock: func() {
				org := &querier.Organization{ID: orgID}
				mockModel.EXPECT().CreateOrganization(gomock.Any(), gomock.Any()).Return(org, nil)
				mockModel.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("user creation failed"))
			},
			expectedError: errors.New("failed to create user"),
		},
		{
			name:     "failed to create organization owner",
			username: username,
			password: password,
			setupMock: func() {
				org := &querier.Organization{ID: orgID}
				user := &querier.User{ID: userID}
				mockModel.EXPECT().CreateOrganization(gomock.Any(), gomock.Any()).Return(org, nil)
				mockModel.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(user, nil)
				mockModel.EXPECT().CreateOrganizationOwner(gomock.Any(), gomock.Any()).Return(errors.New("owner creation failed"))
			},
			expectedError: errors.New("failed to create organization owner"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			createdOrgID, err := service.CreateNewUser(ctx, tc.username, tc.password)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
				require.Equal(t, int32(0), createdOrgID)
			} else {
				require.NoError(t, err)
				require.Equal(t, orgID, createdOrgID)
			}
		})
	}
}
