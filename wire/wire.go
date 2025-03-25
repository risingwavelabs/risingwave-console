//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/risingwavelabs/wavekit/internal/apps/server"
	"github.com/risingwavelabs/wavekit/internal/auth"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/conn/metricsstore"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/controller"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/service"
)

func InitializeServer() (*server.Server, error) {
	wire.Build(
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
	)
	return nil, nil
}
