package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/auth"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/controller"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"github.com/risingwavelabs/wavekit/internal/service"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("server")

type Server struct {
	app        *fiber.App
	port       int
	auth       auth.AuthInterface
	controller *controller.Controller
}

func NewServer(cfg *config.Config, c *controller.Controller, auth auth.AuthInterface, initSvc *service.InitService) (*Server, error) {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
		BodyLimit:    50 * 1024 * 1024, // 50MB
	})

	var port = 8020
	if cfg.Port != 0 {
		port = cfg.Port
	} else {
		log.Infof("Using default port: %d", port)
	}

	s := &Server{
		app:        app,
		port:       port,
		auth:       auth,
		controller: c,
	}

	s.registerMiddleware()

	apigen.RegisterHandlersWithOptions(s.app, s.controller, apigen.FiberServerOptions{
		BaseURL:     "/api/v1",
		Middlewares: []apigen.MiddlewareFunc{},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := initSvc.Init(ctx, cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to initialize server")
	}

	return s, nil
}

func (s *Server) GetController() *controller.Controller {
	return s.controller
}

func (s *Server) registerMiddleware() {
	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	s.app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(wavekit.StaticFiles),
		PathPrefix:   "web/out",
		NotFoundFile: "404.html",
		Index:        "index.html",
	}))

	s.app.Use(cors.New(cors.Config{}))
	s.app.Use(requestid.New())
	s.app.Use(func(c *fiber.Ctx) error {
		// log
		log.Info(
			"request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("request-id", c.Locals(requestid.ConfigDefault.ContextKey).(string)),
		)
		start := time.Now()
		err := c.Next()
		end := time.Now()
		log.Info(
			"response",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("request-id", c.Locals(requestid.ConfigDefault.ContextKey).(string)),
			zap.Float32("latency-ms", float32(end.Sub(start).Milliseconds())),
			zap.String("body", utils.TruncateString(string(c.Response().Body()), 512)),
			zap.Error(err),
		)
		return err
	})

	apigen.RegisterAuthFunc(s.app, s.auth.Authfunc)
}

func (s *Server) Listen() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.port))
}

func (s *Server) GetApp() *fiber.App {
	return s.app
}

func ErrorHandler(c *fiber.Ctx, err error) error {
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
		zap.L().Info(fmt.Sprintf("unexpected error, request-id: %v", rid), zap.Error(err))
		return c.Status(code).SendString(fmt.Sprintf("unexpected error, request-id: %v", rid))
	}

	// Return status code with error message
	return c.Status(code).SendString(err.Error())
}
