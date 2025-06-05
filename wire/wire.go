//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/risingwavelabs/risingwave-console/pkg"
	"github.com/risingwavelabs/risingwave-console/pkg/config"
	"github.com/risingwavelabs/risingwave-console/pkg/conn/http"
	"github.com/risingwavelabs/risingwave-console/pkg/conn/meta"
	"github.com/risingwavelabs/risingwave-console/pkg/conn/metricsstore"
	"github.com/risingwavelabs/risingwave-console/pkg/conn/sql"
	"github.com/risingwavelabs/risingwave-console/pkg/controller"
	"github.com/risingwavelabs/risingwave-console/pkg/service"
	"github.com/risingwavelabs/risingwave-console/pkg/task"
	"github.com/risingwavelabs/risingwave-console/pkg/zcore/injection"
	"github.com/risingwavelabs/risingwave-console/pkg/zcore/model"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/taskgen"

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
