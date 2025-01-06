package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

type User struct {
	ID          int32
	AccessRules map[string]struct{}
}

type AuthServiceInterface interface {
	// CreateToken creates a new JWT token for the given user with specified access rules
	CreateToken(user *querier.User, rules []string) (string, error)

	// ValidateToken validates the given token string and returns the user if valid
	ValidateToken(tokenString string) (*User, error)

	// GetJWTSecret returns the JWT secret used for signing tokens
	GetJWTSecret() []byte

	// GenerateRefreshToken generates a new refresh token for the given user
	GenerateRefreshToken() (string, error)
}

type AuthService struct {
	jwtSecret []byte
}

// Ensure AuthService implements AuthServiceInterface
var _ AuthServiceInterface = (*AuthService)(nil)

func NewAuthService(cfg *config.Config) (AuthServiceInterface, error) {
	if len(cfg.Jwt.Secret) == 0 {
		return nil, errors.New("jwt secret is empty")
	}

	return &AuthService{
		jwtSecret: []byte(cfg.Jwt.Secret),
	}, nil
}

func (s *AuthService) CreateToken(user *querier.User, rules []string) (string, error) {
	claims := s.createClaims(user, rules)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
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

func (s *AuthService) createClaims(user *querier.User, accessRules []string) jwt.MapClaims {
	ruleMap := make(map[string]struct{})
	for _, rule := range accessRules {
		ruleMap[rule] = struct{}{}
	}

	return jwt.MapClaims{
		"user": &User{
			ID:          user.ID,
			AccessRules: ruleMap,
		},
		"exp": time.Now().Add(12 * time.Hour).Unix(),
	}
}

func (s *AuthService) GetJWTSecret() []byte {
	return s.jwtSecret
}

func (s *AuthService) GenerateRefreshToken() (string, error) {
	const length = 32
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b), nil
}
