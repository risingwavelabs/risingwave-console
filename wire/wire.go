//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/risingwavelabs/wavekit/internal"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/conn/http"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/conn/metricsstore"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/controller"
	"github.com/risingwavelabs/wavekit/internal/service"
	"github.com/risingwavelabs/wavekit/internal/task"
	"github.com/risingwavelabs/wavekit/internal/zcore/initapp"
	"github.com/risingwavelabs/wavekit/internal/zcore/injection"
	"github.com/risingwavelabs/wavekit/internal/zcore/model"
	"github.com/risingwavelabs/wavekit/internal/zgen/taskgen"

	anchor_wire "github.com/cloudcarver/anchor/wire"
)

func InitializeApplication() (*initapp.App, error) {
	wire.Build(
		anchor_wire.InitializeApplication,
		initapp.NewApp,
		injection.InjectAuth,
		injection.InjectTaskStore,
		injection.InjectAnchorSvc,
		config.NewConfig,
		service.NewService,
		service.NewInitService,
		model.NewModel,
		sql.NewSQLConnectionManager,
		meta.NewRisectlManager,
		metricsstore.NewMetricsManager,
		http.NewMetaHttpManager,
		controller.NewValidator,
		controller.NewSeverInterface,
		taskgen.NewTaskHandler,
		taskgen.NewTaskRunner,
		task.NewTaskExecutor,
		internal.Init,
	)
	return nil, nil
}
