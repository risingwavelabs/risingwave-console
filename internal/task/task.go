package task

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
)

type TaskStore struct {
	now func() time.Time
}

func NewTaskStore() TaskStoreInterface {
	return &TaskStore{
		now: time.Now,
	}
}

type CreateTaskParams struct {
	OrgID                *int32
	Spec                 apigen.TaskSpec
	StartedAt            *time.Time
	AlwaysRetryOnFailure bool
	RetryInterval        string
}

type CreateScheduledTaskParams struct {
	OrgID                *int32
	Spec                 apigen.TaskSpec
	StartedAt            *time.Time
	AlwaysRetryOnFailure bool
}

func (s *TaskStore) PushTask(c *modelctx.ModelCtx, spec apigen.TaskSpec, opts ...TaskOpt) (int32, error) {
	task := &apigen.Task{
		Spec: spec,
	}
	for _, opt := range opts {
		if err := opt.fn(task); err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("failed to apply task opt %s", opt.name))
		}
	}

	createdTask, err := c.CreateTask(c, querier.CreateTaskParams{
		Attributes: task.Attributes,
		Spec:       task.Spec,
		StartedAt:  task.StartedAt,
		Status:     string(apigen.Pending),
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to create task")
	}
	return createdTask.ID, nil
}
