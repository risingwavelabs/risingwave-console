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
	"github.com/risingwavelabs/wavekit/internal/middleware"
	"github.com/risingwavelabs/wavekit/internal/service"
)

type Server struct {
	app        *fiber.App
	port       int
	auth       auth.AuthInterface
	controller *controller.Controller
}

func NewServer(cfg *config.Config, c *controller.Controller, auth auth.AuthInterface, initSvc *service.InitService) (*Server, error) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		BodyLimit:    50 * 1024 * 1024, // 50MB
	})

	s := &Server{
		app:        app,
		port:       cfg.Port,
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
	s.app.Use(middleware.NewLogger())

	apigen.RegisterAuthFunc(s.app, s.auth.Authfunc)
}

func (s *Server) Listen() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.port))
}

func (s *Server) GetApp() *fiber.App {
	return s.app
}
