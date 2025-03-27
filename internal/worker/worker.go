package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("worker")

type Worker struct {
	model       model.ModelInterface
	getExecutor executorGetter

	risectlm *meta.RisectlManager
}

func NewWorker(globalCtx context.Context, model model.ModelInterface, risectlm *meta.RisectlManager) (*Worker, error) {
	w := &Worker{
		model:       model,
		getExecutor: newExecutor,
		risectlm:    risectlm,
	}

	go func() {
		for {
			select {
			case <-globalCtx.Done():
				return
			case <-time.Tick(1 * time.Second):
				if err := w.runTask(globalCtx); err != nil {
					log.Error("error running task", zap.Error(err))
				}
			}
		}
	}()

	return w, nil
}

func taskToAPI(task *querier.Task) *apigen.Task {
	return &apigen.Task{
		ID:        task.ID,
		CreatedAt: task.CreatedAt,
		Spec:      task.Spec,
		StartedAt: task.StartedAt,
		Status:    apigen.TaskStatus(task.Status),
		UpdatedAt: task.UpdatedAt,
	}
}

func (w *Worker) executeTask(ctx context.Context, model model.ModelInterface, task *apigen.Task) error {
	executor := w.getExecutor(model, w.risectlm)

	switch task.Spec.Type {
	case apigen.AutoBackup:
		if task.Spec.AutoBackup == nil {
			return fmt.Errorf("auto backup spec is nil")
		}
		return executor.ExecuteAutoBackup(ctx, *task.Spec.AutoBackup)
	case apigen.AutoDiagnostic:
		if task.Spec.AutoDiagnostic == nil {
			return fmt.Errorf("auto diagnostic spec is nil")
		}
		return executor.ExecuteAutoDiagnostic(ctx, *task.Spec.AutoDiagnostic)
	default:
		return fmt.Errorf("unknown task type: %s", task.Spec.Type)
	}
}

func (w *Worker) runTask(ctx context.Context) error {
	if err := w.model.RunTransaction(ctx, func(txm model.ModelInterface) error {
		qtask, err := txm.PullTask(ctx)
		if err != nil {
			return err
		}
		task := taskToAPI(qtask)

		log.Info("executing task", zap.Int32("task_id", task.ID), zap.Any("task", task))

		if err := w.executeTask(ctx, txm, task); err != nil {
			log.Error("error executing task", zap.Int32("task_id", task.ID), zap.Error(err))

			if err := txm.UpdateTaskStatus(ctx, querier.UpdateTaskStatusParams{
				ID:     task.ID,
				Status: string(apigen.Failed),
			}); err != nil {
				return errors.Wrap(err, "update task status")
			}

			if _, err := txm.InsertEvent(ctx, apigen.EventSpec{
				Type: apigen.TaskError,
				TaskError: &apigen.EventTaskError{
					TaskID: task.ID,
					Error:  err.Error(),
				},
			}); err != nil {
				return errors.Wrap(err, "insert task error event")
			}

			return nil
		}

		if err := txm.UpdateTaskStatus(ctx, querier.UpdateTaskStatusParams{
			ID:     task.ID,
			Status: string(apigen.Completed),
		}); err != nil {
			return errors.Wrap(err, "update task status")
		}

		log.Info("task completed", zap.Int32("task_id", task.ID))

		return nil
	}); err != nil {
		return err
	}

	return nil
}
