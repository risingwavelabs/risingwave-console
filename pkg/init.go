package pkg

import (
	"context"
	"time"

	"github.com/cloudcarver/anchor/pkg/app"
	"github.com/cloudcarver/anchor/pkg/taskcore/worker"
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

func NewApp(anchorApp *app.Application, cfg *config.Config, serverInterface apigen.ServerInterface, validator apigen.Validator, taskHandler worker.TaskHandler, initService *service.InitService) (*App, error) {
	anchorApp.GetWorker().RegisterTaskHandler(taskHandler)

	apigen.RegisterHandlersWithOptions(anchorApp.GetServer().GetApp(), apigen.NewXMiddleware(serverInterface, validator), apigen.FiberServerOptions{
		BaseURL:     "/api/v1",
		Middlewares: []apigen.MiddlewareFunc{},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := initService.Init(ctx, cfg, anchorApp); err != nil {
		return nil, err
	}

	return &App{
		anchorApp: anchorApp,
	}, nil
}
