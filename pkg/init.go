package pkg

import (
	"context"
	"time"

	"github.com/cloudcarver/anchor/pkg/app"
	"github.com/cloudcarver/anchor/pkg/taskcore/worker"
	"github.com/gofiber/fiber/v2"
	"github.com/risingwavelabs/wavekit/pkg/config"
	"github.com/risingwavelabs/wavekit/pkg/service"
	"github.com/risingwavelabs/wavekit/pkg/zgen/apigen"
)

type App struct {
	anchorApp *app.Application
}

func (a *App) Start() error {
	return a.anchorApp.Start()
}

func NewApp(anchorApp *app.Application, cfg *config.Config, plugin *Plugin, initService *service.InitService) (*App, error) {
	plugin.Plug(anchorApp)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := initService.Init(ctx, cfg, anchorApp); err != nil {
		return nil, err
	}

	return &App{
		anchorApp: anchorApp,
	}, nil
}

type Plugin struct {
	serverInterface apigen.ServerInterface
	validator       apigen.Validator
	taskHandler     worker.TaskHandler
}

func NewPlugin(serverInterface apigen.ServerInterface, validator apigen.Validator, taskHandler worker.TaskHandler) *Plugin {
	return &Plugin{
		serverInterface: serverInterface,
		validator:       validator,
		taskHandler:     taskHandler,
	}
}

func (p *Plugin) Plug(anchorApp *app.Application) {
	p.PlugToFiberApp(anchorApp.GetServer().GetApp())
	p.PlugToWorker(anchorApp.GetWorker())
}

func (p *Plugin) PlugToFiberApp(fiberApp *fiber.App) {
	apigen.RegisterHandlersWithOptions(fiberApp, apigen.NewXMiddleware(p.serverInterface, p.validator), apigen.FiberServerOptions{
		BaseURL:     "/api/v1",
		Middlewares: []apigen.MiddlewareFunc{},
	})
}

func (p *Plugin) PlugToWorker(worker *worker.Worker) {
	worker.RegisterTaskHandler(p.taskHandler)
}
