package app

import (
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/metrics"
	"github.com/risingwavelabs/wavekit/internal/server"
	"github.com/risingwavelabs/wavekit/internal/worker"
)

type Application struct {
	server        *server.Server
	prometheus    *metrics.MetricsServer
	worker        *worker.Worker
	disableWorker bool
	debugServer   *DebugServer
}

func NewApplication(cfg *config.Config, server *server.Server, prometheus *metrics.MetricsServer, worker *worker.Worker, debugServer *DebugServer) *Application {
	return &Application{
		server:        server,
		prometheus:    prometheus,
		worker:        worker,
		disableWorker: cfg.Worker.Disable,
		debugServer:   debugServer,
	}
}

func (a *Application) Start() error {
	go a.debugServer.Start()
	go a.prometheus.Start()
	if !a.disableWorker {
		go a.worker.Start()
	}
	return a.server.Listen()
}

func (a *Application) GetServer() *server.Server {
	return a.server
}
