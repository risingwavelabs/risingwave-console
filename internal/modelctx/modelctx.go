package modelctx

import (
	"context"

	"github.com/risingwavelabs/wavekit/internal/model"
)

type ModelCtx struct {
	model.ModelInterface
	context.Context
}

func NewModelctx(ctx context.Context, model model.ModelInterface) *ModelCtx {
	return &ModelCtx{
		ModelInterface: model,
		Context:        ctx,
	}
}
