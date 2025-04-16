package task

import (
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
)

type TaskStoreInterface interface {
	PushTask(c *modelctx.ModelCtx, spec apigen.TaskSpec, opts ...TaskOpt) (int32, error)

	CreateCronJob(c *modelctx.ModelCtx, timeoutDuration *string, orgID *int32, cronExpression string, specType apigen.TaskSpec) (int32, error)

	UpdateCronJob(c *modelctx.ModelCtx, taskID int32, timeoutDuration *string, orgID *int32, cronExpression string, specType apigen.TaskSpec) error

	PauseCronJob(c *modelctx.ModelCtx, taskID int32) error

	ResumeCronJob(c *modelctx.ModelCtx, taskID int32) error
}
