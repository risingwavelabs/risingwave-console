//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/risingwavelabs/wavekit/internal/app"
	"github.com/risingwavelabs/wavekit/internal/auth"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/conn/http"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/conn/metricsstore"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/controller"
	"github.com/risingwavelabs/wavekit/internal/globalctx"
	"github.com/risingwavelabs/wavekit/internal/macaroons"
	"github.com/risingwavelabs/wavekit/internal/metrics"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/server"
	"github.com/risingwavelabs/wavekit/internal/service"
	"github.com/risingwavelabs/wavekit/internal/task"
	"github.com/risingwavelabs/wavekit/internal/worker"
	"github.com/risingwavelabs/wavekit/internal/worker/handler"
)

func InitializeApplication() (*app.Application, error) {
	wire.Build(
		app.NewApplication,
		config.NewConfig,
		service.NewService,
		controller.NewController,
		model.NewModel,
		server.NewServer,
		service.NewInitService,
		auth.NewAuth,
		sql.NewSQLConnectionManager,
		meta.NewRisectlManager,
		metricsstore.NewMetricsManager,
		metrics.NewMetricsServer,
		app.NewDebugServer,
		globalctx.New,
		worker.NewWorker,
		task.NewTaskStore,
		http.NewMetaHttpManager,
		macaroons.NewMacaroonManager,
		macaroons.NewStore,
		auth.NewCaveatParser,
		handler.NewTaskHandler,
	)
	return nil, nil
}
