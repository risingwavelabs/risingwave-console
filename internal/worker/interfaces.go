package worker

import (
	"context"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
)

type TaskLifeCycleHandlerInterface interface {
	HandleAttributes(ctx context.Context, task apigen.Task) error
	HandleFailed(ctx context.Context, task apigen.Task, err error) error
	HandleCompleted(ctx context.Context, task apigen.Task) error
}

var ErrUnknownTaskType = errors.New("unknown task type")

type TaskHandler interface {
	HandleTask(c *modelctx.ModelCtx, task apigen.Task) error

	RegisterTaskHandler(taskHandler TaskHandler)
}
