package internal

import (
	"context"
	"time"

	anchor_app "github.com/cloudcarver/anchor/pkg/app"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/service"
)

func Init(cfg *config.Config, initService *service.InitService) func(anchorApp *anchor_app.Application) error {
	return func(anchorApp *anchor_app.Application) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := initService.Init(ctx, cfg, anchorApp); err != nil {
			return err
		}

		return nil
	}
}
