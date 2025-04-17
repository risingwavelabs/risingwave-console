package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/macaroons"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuth_Authfunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacaroons := macaroons.NewMockMacaroonManagerInterface(ctrl)
	auth, err := NewAuth(mockMacaroons)
	require.NoError(t, err)

	// Test token
	testToken := "test_token"
	testBearerToken := "Bearer " + testToken

	testCases := []struct {
		name           string
		authHeader     string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func() {},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			authHeader: testToken,
			setupMock: func() {
				mockMacaroons.EXPECT().Parse(gomock.Any(), testToken).Return(nil, macaroons.ErrMalformedToken)
			},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "bearer token prefix",
			authHeader: testBearerToken,
			setupMock: func() {
				mockMacaroons.EXPECT().Parse(gomock.Any(), testToken).Return(nil, macaroons.ErrMalformedToken)
			},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "caveat validation error",
			authHeader: testToken,
			setupMock: func() {
				mockCaveat := macaroons.NewMockCaveat(ctrl)
				mockCaveat.EXPECT().Encode().Return("caveat", nil)

				macaroon, err := macaroons.CreateMacaroon(123, []byte("key"), []macaroons.Caveat{mockCaveat})
				require.NoError(t, err)

				mockMacaroons.EXPECT().Parse(gomock.Any(), testToken).Return(macaroon, nil)
				mockCaveat.EXPECT().Validate(gomock.Any()).Return(errors.New("caveat validation error"))

			},
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:       "successful authorization",
			authHeader: testToken,
			setupMock: func() {
				mockCaveat := macaroons.NewMockCaveat(ctrl)
				mockCaveat.EXPECT().Encode().Return("caveat", nil)
				macaroon, err := macaroons.CreateMacaroon(123, []byte("key"), []macaroons.Caveat{mockCaveat})
				require.NoError(t, err)

				mockMacaroons.EXPECT().Parse(gomock.Any(), testToken).Return(macaroon, nil)
				mockCaveat.EXPECT().Validate(gomock.Any()).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test Fiber app and set up error handling
			app := fiber.New(fiber.Config{
				ErrorHandler: utils.ErrorHandler,
			})

			// Add a test route
			app.Use(func(c *fiber.Ctx) error {
				// Call auth.Authfunc
				err := auth.Authfunc(c)
				if err != nil {
					return err
				}
				return nil
			})

			// Set up mock expectations
			tc.setupMock()

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			// Execute request
			resp, err := app.Test(req)
			require.NoError(t, err)

			// Verify status code
			require.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

func TestAuth_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacaroons := macaroons.NewMockMacaroonManagerInterface(ctrl)
	auth, err := NewAuth(mockMacaroons)
	require.NoError(t, err)

	ctx := context.Background()
	userID := int32(1)
	orgID := int32(10)
	keyID := int64(123)

	user := &querier.User{
		ID:             userID,
		OrganizationID: orgID,
	}

	macaroon, err := macaroons.CreateMacaroon(123, []byte("key"), nil)
	require.NoError(t, err)
	testCases := []struct {
		name          string
		user          *querier.User
		rules         []string
		setupMock     func()
		expectedKeyID int64
		expectedToken string
		expectedError error
	}{
		{
			name:  "successful token creation",
			user:  user,
			rules: []string{"read", "write"},
			setupMock: func() {
				mockMacaroons.EXPECT().CreateToken(
					gomock.Any(),
					userID,
					gomock.Any(), // Here we expect a UserContextCaveat but it's difficult to match in tests
					TimeoutAccessToken,
				).Return(macaroon, nil)
			},
			expectedKeyID: keyID,
			expectedToken: macaroon.StringToken(),
			expectedError: nil,
		},
		{
			name:  "token creation failure",
			user:  user,
			rules: []string{"read"},
			setupMock: func() {
				mockMacaroons.EXPECT().CreateToken(
					gomock.Any(),
					userID,
					gomock.Any(),
					TimeoutAccessToken,
				).Return(nil, errors.New("token creation failed"))
			},
			expectedKeyID: 0,
			expectedToken: "",
			expectedError: errors.New("failed to create macaroon token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			gotKeyID, gotToken, err := auth.CreateToken(ctx, tc.user, tc.rules)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
				require.Equal(t, tc.expectedKeyID, gotKeyID)
				require.Equal(t, tc.expectedToken, gotToken)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedKeyID, gotKeyID)
				require.Equal(t, tc.expectedToken, gotToken)
			}
		})
	}
}

func TestAuth_CreateRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacaroons := macaroons.NewMockMacaroonManagerInterface(ctrl)
	auth, err := NewAuth(mockMacaroons)
	require.NoError(t, err)

	ctx := context.Background()
	userID := int32(1)
	accessKeyID := int64(123)

	macaroon, err := macaroons.CreateMacaroon(0, []byte("key"), nil)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		userID        int32
		accessKeyID   int64
		setupMock     func()
		expectedToken string
		expectedError error
	}{
		{
			name:        "successful refresh token creation",
			userID:      userID,
			accessKeyID: accessKeyID,
			setupMock: func() {

				mockMacaroons.EXPECT().CreateToken(
					gomock.Any(),
					userID,
					gomock.Any(), // Expect RefreshOnlyCaveat but difficult to match in tests
					TimeoutRefreshToken,
				).Return(macaroon, nil)
			},
			expectedToken: macaroon.StringToken(),
			expectedError: nil,
		},
		{
			name:        "refresh token creation failure",
			userID:      userID,
			accessKeyID: accessKeyID,
			setupMock: func() {
				mockMacaroons.EXPECT().CreateToken(
					gomock.Any(),
					userID,
					gomock.Any(),
					TimeoutRefreshToken,
				).Return(nil, errors.New("token creation failed"))
			},
			expectedToken: "",
			expectedError: errors.New("failed to create macaroon token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			gotToken, err := auth.CreateRefreshToken(ctx, tc.accessKeyID, tc.userID)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
				require.Equal(t, tc.expectedToken, gotToken)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToken, gotToken)
			}
		})
	}
}

func TestAuth_ParseRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacaroons := macaroons.NewMockMacaroonManagerInterface(ctrl)
	auth, err := NewAuth(mockMacaroons)
	require.NoError(t, err)

	ctx := context.Background()
	userID := int32(1)

	refreshCaveat := NewRefreshOnlyCaveat(userID, 456)
	macaroon, err := macaroons.CreateMacaroon(0, []byte("key"), []macaroons.Caveat{refreshCaveat})
	require.NoError(t, err)

	noRefreshCaveat := macaroons.NewMockCaveat(ctrl)
	noRefreshCaveat.EXPECT().Encode().Return("caveat", nil)
	noRefreshCaveat.EXPECT().Type().Return("not_refresh")
	noRefreshMacaroon, err := macaroons.CreateMacaroon(0, []byte("key"), []macaroons.Caveat{noRefreshCaveat})
	require.NoError(t, err)

	testCases := []struct {
		name           string
		refreshToken   string
		setupMock      func()
		expectedUserID int32
		expectedError  error
	}{
		{
			name:         "successful refresh token parsing",
			refreshToken: macaroon.StringToken(),
			setupMock: func() {
				mockMacaroons.EXPECT().Parse(gomock.Any(), macaroon.StringToken()).Return(macaroon, nil)
			},
			expectedUserID: userID,
			expectedError:  nil,
		},
		{
			name:         "parse failure",
			refreshToken: macaroon.StringToken(),
			setupMock: func() {
				mockMacaroons.EXPECT().Parse(gomock.Any(), macaroon.StringToken()).Return(nil, errors.New("parse failed"))
			},
			expectedUserID: 0,
			expectedError:  errors.New("failed to parse macaroon token"),
		},
		{
			name:         "no refresh caveat",
			refreshToken: noRefreshMacaroon.StringToken(),
			setupMock: func() {
				mockMacaroons.EXPECT().Parse(gomock.Any(), noRefreshMacaroon.StringToken()).Return(noRefreshMacaroon, nil)
			},
			expectedUserID: 0,
			expectedError:  errors.New("no userID found in refresh token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			gotUserID, err := auth.ParseRefreshToken(ctx, tc.refreshToken)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
				require.Equal(t, tc.expectedUserID, gotUserID)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUserID, gotUserID)
			}
		})
	}
}

func TestAuth_InvalidateUserTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacaroons := macaroons.NewMockMacaroonManagerInterface(ctrl)
	auth, err := NewAuth(mockMacaroons)
	require.NoError(t, err)

	ctx := context.Background()
	userID := int32(1)

	testCases := []struct {
		name          string
		userID        int32
		setupMock     func()
		expectedError error
	}{
		{
			name:   "successful invalidation",
			userID: userID,
			setupMock: func() {
				mockMacaroons.EXPECT().InvalidateUserTokens(gomock.Any(), userID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "invalidation failure",
			userID: userID,
			setupMock: func() {
				mockMacaroons.EXPECT().InvalidateUserTokens(gomock.Any(), userID).Return(errors.New("invalidation failed"))
			},
			expectedError: errors.New("invalidation failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			err := auth.InvalidateUserTokens(ctx, tc.userID)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	userID := int32(1)

	testCases := []struct {
		name           string
		setupContext   func(*fiber.Ctx)
		expectedUserID int32
		expectedError  error
	}{
		{
			name: "successful user ID retrieval",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(ContextKeyUserID, userID)
			},
			expectedUserID: userID,
			expectedError:  nil,
		},
		{
			name: "user ID not in context",
			setupContext: func(c *fiber.Ctx) {
				// Do not set user ID
			},
			expectedUserID: 0,
			expectedError:  ErrUserIdentityNotExist,
		},
		{
			name: "user ID wrong type",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(ContextKeyUserID, "not an int32")
			},
			expectedUserID: 0,
			expectedError:  ErrUserIdentityNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test Fiber app
			app := fiber.New()

			// Add a test route
			app.Get("/test", func(c *fiber.Ctx) error {
				// Set up context
				tc.setupContext(c)

				// Call the test function
				gotUserID, err := GetUserID(c)

				// Verify results
				if tc.expectedError != nil {
					require.Error(t, err)
					require.Equal(t, tc.expectedError, err)
					require.Equal(t, tc.expectedUserID, gotUserID)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedUserID, gotUserID)
				}

				return c.SendStatus(fiber.StatusOK)
			})

			// Send request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusOK, resp.StatusCode)
		})
	}
}

func TestGetOrgID(t *testing.T) {
	orgID := int32(10)

	testCases := []struct {
		name          string
		setupContext  func(*fiber.Ctx)
		expectedOrgID int32
		expectedError error
	}{
		{
			name: "successful org ID retrieval",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(ContextKeyOrgID, orgID)
			},
			expectedOrgID: orgID,
			expectedError: nil,
		},
		{
			name: "org ID not in context",
			setupContext: func(c *fiber.Ctx) {
				// Do not set organization ID
			},
			expectedOrgID: 0,
			expectedError: ErrUserIdentityNotExist,
		},
		{
			name: "org ID wrong type",
			setupContext: func(c *fiber.Ctx) {
				c.Locals(ContextKeyOrgID, "not an int32")
			},
			expectedOrgID: 0,
			expectedError: ErrUserIdentityNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test Fiber app
			app := fiber.New()

			// Add a test route
			app.Get("/test", func(c *fiber.Ctx) error {
				// Set up context
				tc.setupContext(c)

				// Call the test function
				gotOrgID, err := GetOrgID(c)

				// Verify results
				if tc.expectedError != nil {
					require.Error(t, err)
					require.Equal(t, tc.expectedError, err)
					require.Equal(t, tc.expectedOrgID, gotOrgID)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedOrgID, gotOrgID)
				}

				return c.SendStatus(fiber.StatusOK)
			})

			// Send request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, fiber.StatusOK, resp.StatusCode)
		})
	}
}

func TestNewAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacaroons := macaroons.NewMockMacaroonManagerInterface(ctrl)

	auth, err := NewAuth(mockMacaroons)
	require.NoError(t, err)
	require.NotNil(t, auth)
}
