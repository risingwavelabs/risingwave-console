package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

const UserContextKey = "user"

var ErrUserIdentityNotExist = errors.New("user identity not exists")

type User struct {
	ID             int32
	OrganizationID int32
	AccessRules    map[string]struct{}
}

type AuthInterface interface {
	Authfunc(c *fiber.Ctx, rules ...string) error

	// CreateToken creates a new JWT token for the given user with specified access rules
	CreateToken(user *querier.User, rules []string) (string, error)

	// ValidateToken validates the given token string and returns the user if valid
	ValidateToken(tokenString string) (*User, error)

	// GetJWTSecret returns the JWT secret used for signing tokens
	GetJWTSecret() []byte

	// CreateRefreshToken returns a refresh token and its JWT token
	CreateRefreshToken(userID int32) (string, string, error)

	// ParseJWTRefreshToken parses the given JWT refresh token and returns the user ID and the refresh token
	ParseJWTRefreshToken(jwtToken string) (int32, string, error)
}

type Auth struct {
	jwtSecret []byte
}

// Ensure AuthService implements AuthServiceInterface
var _ AuthInterface = (*Auth)(nil)

func NewAuth(cfg *config.Config) (AuthInterface, error) {
	if len(cfg.Jwt.Secret) == 0 && !cfg.Jwt.RandomSecret {
		return nil, errors.New("jwt secret is empty")
	}

	return &Auth{
		jwtSecret: []byte(utils.IfElse(cfg.Jwt.RandomSecret, randomString(32), cfg.Jwt.Secret)),
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

	user, err := a.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	c.Locals(UserContextKey, user)

	for _, rule := range rules {
		if _, ok := user.AccessRules[rule]; !ok {
			return c.Status(fiber.StatusForbidden).SendString(fmt.Sprintf("Permission denied, need rule %s", rule))
		}
	}
	return nil
}

func (a *Auth) CreateToken(user *querier.User, rules []string) (string, error) {
	claims := a.createClaims(user, rules)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

func (a *Auth) ValidateToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unexpected error when parsing claims: claims is not jwt.MapClaims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("failed to parse exp")
	}
	if time.Since(time.Unix(int64(exp), 0)) > 0 {
		return nil, errors.New("token is expired")
	}

	var user User
	if err := utils.JSONConvert(claims["user"], &user); err != nil {
		return nil, errors.Wrapf(err, "failed to parse user from claims: %s", utils.TryMarshal(claims["user"]))
	}
	return &user, nil
}

func (a *Auth) createClaims(user *querier.User, accessRules []string) jwt.MapClaims {
	ruleMap := make(map[string]struct{})
	for _, rule := range accessRules {
		ruleMap[rule] = struct{}{}
	}

	return jwt.MapClaims{
		"user": &User{
			ID:             user.ID,
			OrganizationID: user.OrganizationID,
			AccessRules:    ruleMap,
		},
		"exp": time.Now().Add(12 * time.Hour).Unix(),
	}
}

func (a *Auth) GetJWTSecret() []byte {
	return a.jwtSecret
}

func (a *Auth) CreateRefreshToken(userID int32) (string, string, error) {
	refreshToken := randomString(32)
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"token":  refreshToken,
	})
	jwtToken, err := jwt.SignedString(a.jwtSecret)
	if err != nil {
		return "", "", err
	}
	return refreshToken, jwtToken, nil
}

func (a *Auth) ParseJWTRefreshToken(jwtToken string) (int32, string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("unexpected error when parsing claims: claims is not jwt.MapClaims")
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, "", errors.New("unexpected error when parsing userID: userID is not float64")
	}

	refreshToken, ok := claims["token"].(string)
	if !ok {
		return 0, "", errors.New("unexpected error when parsing token: token is not string")
	}
	return int32(userID), refreshToken, nil
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GetUser(c *fiber.Ctx) (*User, error) {
	user, ok := c.Locals(UserContextKey).(*User)
	if !ok {
		return nil, ErrUserIdentityNotExist
	}
	return user, nil
}
