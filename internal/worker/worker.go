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
	model model.ModelInterface

	getExecutor executorGetter
	getHandler  LifeCycleHandlerGetter

	risectlm *meta.RisectlManager
}

func NewWorker(globalCtx context.Context, model model.ModelInterface, risectlm *meta.RisectlManager) (*Worker, error) {
	w := &Worker{
		model:       model,
		getExecutor: newExecutor,
		getHandler:  newTaskLifeCycleHandler,
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

func taskToAPI(task *querier.Task) apigen.Task {
	return apigen.Task{
		ID:        task.ID,
		CreatedAt: task.CreatedAt,
		Spec:      task.Spec,
		StartedAt: task.StartedAt,
		Status:    apigen.TaskStatus(task.Status),
		UpdatedAt: task.UpdatedAt,
	}
}

func (w *Worker) executeTask(ctx context.Context, model model.ModelInterface, task apigen.Task) error {
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

		// life cycle handler
		lh, err := w.getHandler(txm)
		if err != nil {
			return errors.Wrap(err, "failed to create attribute handler")
		}

		// handle attributes
		if err := lh.HandleAttributes(ctx, task); err != nil {
			return errors.Wrap(err, "failed to handle attributes")
		}

		// run task
		err = w.executeTask(ctx, txm, task)
		if err != nil { // handle failed
			log.Error("error executing task", zap.Int32("task_id", task.ID), zap.Error(err))
			if err := lh.HandleFailed(ctx, task, err); err != nil {
				return errors.Wrap(err, "failed to handle failed task")
			}
		} else { // handle completed
			if err := lh.HandleCompleted(ctx, task); err != nil {
				log.Error("error handling completed task", zap.Int32("task_id", task.ID), zap.Error(err))
				return errors.Wrap(err, "failed to handle completed task")
			}
			log.Info("task completed", zap.Int32("task_id", task.ID))
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
