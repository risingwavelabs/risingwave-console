package middleware

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("middleware")

func (m *Middleware) ErrorHandler(c *fiber.Ctx, err error) error {
	// default 500
	var code = fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	rid := c.Locals(requestid.ConfigDefault.ContextKey)

	if code == fiber.StatusInternalServerError {
		log.Error(fmt.Sprintf("unexpected error, request-id: %v", rid), zap.Error(err))
		return c.Status(code).SendString(fmt.Sprintf("unexpected error, request-id: %v", rid))
	}

	// Return status code with error message
	return c.Status(code).SendString(err.Error())
}
