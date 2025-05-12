// Code generate by anchor. DO NOT EDIT.
package initapp

import (
	"github.com/cloudcarver/anchor/pkg/app"
	"github.com/cloudcarver/anchor/pkg/taskcore/worker"
	"github.com/risingwavelabs/wavekit/internal/zgen/apigen"
)

type App struct {
	anchorApp *app.Application
}

func NewApp(anchorApp *app.Application, serverInterface apigen.ServerInterface, validator apigen.Validator, taskHandler worker.TaskHandler) *App {
	app := &App{
		anchorApp: anchorApp,
	}

	app.RegisterServerInterface(serverInterface, validator)
	app.RegisterTaskHandler(taskHandler)

	return app
}

func (a *App) Start() error {
	return a.anchorApp.Start()
}

func (a *App) RegisterTaskHandler(taskHandler worker.TaskHandler) {
	a.anchorApp.GetWorker().RegisterTaskHandler(taskHandler)
}

func (a *App) RegisterServerInterface(serverInterface apigen.ServerInterface, validator apigen.Validator) {
	apigen.RegisterHandlersWithOptions(a.anchorApp.GetServer().GetApp(), apigen.NewXMiddleware(serverInterface, validator), apigen.FiberServerOptions{
		BaseURL:     "/api/v1",
		Middlewares: []apigen.MiddlewareFunc{},
	})
}
