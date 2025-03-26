package worker

import (
	"context"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
)

type Worker struct {
	// Unique name for the worker
	Name string

	model model.ModelInterface
}

func NewWorker(name string, model model.ModelInterface) *Worker {
	return &Worker{
		Name:  name,
		model: model,
	}
}

func (w *Worker) PullTask(ctx context.Context) (*apigen.Task, error) {
	return nil, nil
}
