//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/risingwavelabs/wavekit/pkg"
	"github.com/risingwavelabs/wavekit/pkg/config"
	"github.com/risingwavelabs/wavekit/pkg/conn/http"
	"github.com/risingwavelabs/wavekit/pkg/conn/meta"
	"github.com/risingwavelabs/wavekit/pkg/conn/metricsstore"
	"github.com/risingwavelabs/wavekit/pkg/conn/sql"
	"github.com/risingwavelabs/wavekit/pkg/controller"
	"github.com/risingwavelabs/wavekit/pkg/service"
	"github.com/risingwavelabs/wavekit/pkg/task"
	"github.com/risingwavelabs/wavekit/pkg/zcore/injection"
	"github.com/risingwavelabs/wavekit/pkg/zcore/model"
	"github.com/risingwavelabs/wavekit/pkg/zgen/taskgen"

	anchor_wire "github.com/cloudcarver/anchor/wire"
)

func InitializeApplication() (*pkg.App, error) {
	wire.Build(
		anchor_wire.InitializeApplication,
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
		pkg.NewApp,
		pkg.NewPlugin,
	)
	return nil, nil
}
