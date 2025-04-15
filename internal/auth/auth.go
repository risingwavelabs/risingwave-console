package auth

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/macaroons"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
)

const (
	ContextKeyUserID = iota
	ContextKeyOrgID
	ContextKeyMacaroon
)

const (
	TimeoutRefreshToken = time.Hour * 24 * 30
	TimeoutAccessToken  = time.Minute * 30
)

var ErrUserIdentityNotExist = errors.New("user identity not exists")

type User struct {
	ID             int32
	OrganizationID int32
	AccessRules    map[string]struct{}
}

type AuthInterface interface {
	Authfunc(c *fiber.Ctx, rules ...string) error

	// CreateToken creates a new JWT token for the given user with specified access rules
	CreateToken(ctx context.Context, user *querier.User, rules []string) (int64, string, error)

	// CreateRefreshToken returns a refresh token
	CreateRefreshToken(ctx context.Context, accessKeyID int64, userID int32) (string, error)

	// ParseRefreshToken parses the given refresh token and returns the user ID
	ParseRefreshToken(ctx context.Context, refreshToken string) (int32, error)
}

type Auth struct {
	macaroons macaroons.MacaroonManagerInterface
}

// Ensure AuthService implements AuthServiceInterface
var _ AuthInterface = (*Auth)(nil)

func NewAuth(macaroons macaroons.MacaroonManagerInterface) (AuthInterface, error) {
	return &Auth{
		macaroons: macaroons,
	}, nil
}

func (a *Auth) Authfunc(c *fiber.Ctx, rules ...string) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("missing authorization header")
	}

	// Remove "Bearer " prefix if present
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	token, err := a.macaroons.Parse(c.Context(), tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	for _, caveat := range token.Caveats() {
		if err := caveat.Validate(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
	}

	return nil
}

func (a *Auth) CreateToken(ctx context.Context, user *querier.User, rules []string) (int64, string, error) {
	token, err := a.macaroons.CreateToken(ctx, []macaroons.Caveat{
		NewUserContextCaveat(user.ID, user.OrganizationID),
	}, TimeoutAccessToken)
	if err != nil {
		return 0, "", errors.Wrap(err, "failed to create macaroon token")
	}
	return token.KeyID(), token.StringToken(), nil
}

func (a *Auth) CreateRefreshToken(ctx context.Context, accessKeyID int64, userID int32) (string, error) {
	token, err := a.macaroons.CreateToken(ctx, []macaroons.Caveat{
		NewRefreshOnlyCaveat(userID, accessKeyID),
	}, TimeoutRefreshToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to create macaroon token")
	}
	return token.StringToken(), nil
}

func (a *Auth) ParseRefreshToken(ctx context.Context, refreshToken string) (int32, error) {
	token, err := a.macaroons.Parse(ctx, refreshToken)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse macaroon token")
	}

	for _, caveat := range token.Caveats() {
		if caveat.Type() == CaveatRefreshOnly {
			roc, ok := caveat.(*RefreshOnlyCaveat)
			if !ok {
				return 0, errors.Errorf("caveat is not a RefreshOnlyCaveat even though it has type %s", CaveatRefreshOnly)
			}
			return roc.UserID, nil
		}
	}

	return 0, errors.New("no userID found in refresh token")
}

func GetUserID(c *fiber.Ctx) (int32, error) {
	userID, ok := c.Locals(ContextKeyUserID).(int32)
	if !ok {
		return 0, ErrUserIdentityNotExist
	}
	return userID, nil
}

func GetOrgID(c *fiber.Ctx) (int32, error) {
	orgID, ok := c.Locals(ContextKeyOrgID).(int32)
	if !ok {
		return 0, ErrUserIdentityNotExist
	}
	return orgID, nil
}
