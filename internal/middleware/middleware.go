package middleware

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/risingwavelabs/wavekit/internal/service"
)

const (
	UserContextKey      = "user"
	AuthorizationHeader = "Authorization"
)

var (
	ErrUserIdentityNotExist = errors.New("user identity not exists")
)

type Middleware struct {
	authService service.AuthServiceInterface
}

func NewMiddleware(authService service.AuthServiceInterface) *Middleware {
	return &Middleware{
		authService: authService,
	}
}

func (m *Middleware) Auth(c *fiber.Ctx) error {
	authHeader := c.Get(AuthorizationHeader)
	if authHeader == "" {
		return c.Status(401).SendString("missing authorization header")
	}

	// Remove "Bearer " prefix if present
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	user, err := m.authService.ValidateToken(tokenString)
	if err != nil {
		return c.Status(401).SendString(err.Error())
	}

	c.Locals(UserContextKey, user)
	return c.Next()
}

func (m *Middleware) CheckRules(c *fiber.Ctx, rules []string, rejectRules []string) error {
	user, err := GetUser(c)
	if err != nil {
		return c.Status(403).SendString(err.Error())
	}
	for _, rule := range rules {
		if _, ok := user.AccessRules[rule]; !ok {
			return c.Status(403).SendString(fmt.Sprintf("没有权限，需要访问规则%s", rule))
		}
	}
	for _, rule := range rejectRules {
		if _, ok := user.AccessRules[rule]; ok {
			return c.Status(403).SendString(fmt.Sprintf("没有权限，因为带有访问规则%s", rule))
		}
	}
	return nil
}

func GetUser(c *fiber.Ctx) (*service.User, error) {
	user, ok := c.Locals(UserContextKey).(*service.User)
	if !ok {
		return nil, ErrUserIdentityNotExist
	}
	return user, nil
}
