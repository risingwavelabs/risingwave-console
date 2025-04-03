package modelctx

import (
	"context"
	"time"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
)

type ModelContextInterface interface {
	CreateCronJob(ctx context.Context, timeoutDuration *string, orgID *int32, cronExpression string, specType apigen.TaskSpec) (int32, error)
	UpdateCronJob(ctx context.Context, taskID int32, timeoutDuration *string, orgID *int32, cronExpression string, specType apigen.TaskSpec) error
	PauseCronJob(ctx context.Context, taskID int32) error
	ResumeCronJob(ctx context.Context, taskID int32) error
}

type ModelContext struct {
	model model.ModelInterface
	now   func() time.Time
}

func NewModelctx(model model.ModelInterface, now func() time.Time) ModelContextInterface {
	return &ModelContext{
		model: model,
		now:   now,
	}
}
