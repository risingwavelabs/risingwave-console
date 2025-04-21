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
	"github.com/risingwavelabs/wavekit/internal/globalctx"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"github.com/risingwavelabs/wavekit/internal/service"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("server")

type Server struct {
	app             *fiber.App
	host            string
	port            int
	auth            auth.AuthInterface
	globalCtx       *globalctx.GlobalContext
	serverInterface apigen.ServerInterface
	validator       apigen.Validator
}

func NewServer(
	cfg *config.Config,
	globalCtx *globalctx.GlobalContext,
	auth auth.AuthInterface,
	initSvc *service.InitService,
	serverInterface apigen.ServerInterface,
	validator apigen.Validator,
) (*Server, error) {
	// create fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
		BodyLimit:    50 * 1024 * 1024, // 50MB
	})

	var port = 8020
	if cfg.Port != 0 {
		port = cfg.Port
	} else {
		log.Infof("Using default port: %d", port)
	}

	var host = "localhost"
	if cfg.Host != "" {
		host = cfg.Host
	} else {
		log.Infof("Using default host: %s", host)
	}

	s := &Server{
		app:             app,
		host:            host,
		port:            port,
		auth:            auth,
		serverInterface: serverInterface,
		globalCtx:       globalCtx,
		validator:       validator,
	}

	s.registerMiddleware()

	apigen.RegisterHandlersWithOptions(s.app, apigen.NewXMiddleware(s.serverInterface, s.validator), apigen.FiberServerOptions{
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

	s.app.Get("/config.js", func(c *fiber.Ctx) error {
		endpoint := fmt.Sprintf("http://%s:%d/api/v1", s.host, s.port)
		c.Set("Content-Type", "application/javascript")
		return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("window.APP_ENDPOINT = '%s';", endpoint))
	})

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
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("token", fmt.Sprintf("%v", c.Get("Authorization"))),
			zap.String("request-id", c.Locals(requestid.ConfigDefault.ContextKey).(string)),
			zap.Float32("latency-ms", float32(end.Sub(start).Milliseconds())),
			zap.String("body", utils.TruncateString(string(c.Response().Body()), 512)),
			zap.Error(err),
		)
		return err
	})

}

func (s *Server) Listen() error {
	// Create a channel to receive shutdown signal
	shutdownChan := make(chan error)

	// Start the server in a goroutine
	go func() {
		if err := s.app.Listen(fmt.Sprintf(":%d", s.port)); err != nil {
			shutdownChan <- err
		}
	}()

	// Wait for either context cancellation or server error
	select {
	case err := <-shutdownChan:
		return err
	case <-s.globalCtx.Context().Done():
		log.Info("shutting down server due to context cancellation")
		return s.app.Shutdown()
	}
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) GetApp() *fiber.App {
	return s.app
}
