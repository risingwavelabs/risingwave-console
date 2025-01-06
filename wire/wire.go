//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/risingwavelabs/wavekit/internal/apps/server"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/controller"
	"github.com/risingwavelabs/wavekit/internal/middleware"
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
		middleware.NewMiddleware,
		service.NewInitService,
		service.NewAuthService,
	)
	return nil, nil
}
