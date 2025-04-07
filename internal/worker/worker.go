package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/globalctx"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"github.com/risingwavelabs/wavekit/internal/metrics"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("worker")

const maxTaskTimeout = 1 * time.Hour

type ExecutorInterface interface {
	ExecuteAutoBackup(ctx context.Context, spec apigen.TaskSpecAutoBackup) error
	ExecuteAutoDiagnostic(ctx context.Context, spec apigen.TaskSpecAutoDiagnostic) error
}

type Worker struct {
	model model.ModelInterface

	getHandler LifeCycleHandlerGetter

	executor ExecutorInterface

	globalCtx *globalctx.GlobalContext
}

func NewWorker(globalCtx *globalctx.GlobalContext, model model.ModelInterface, executor ExecutorInterface) (*Worker, error) {
	w := &Worker{
		model:      model,
		getHandler: newTaskLifeCycleHandler,
		globalCtx:  globalCtx,
		executor:   executor,
	}

	return w, nil
}

func taskToAPI(task *querier.Task) apigen.Task {
	return apigen.Task{
		ID:         task.ID,
		CreatedAt:  task.CreatedAt,
		Spec:       task.Spec,
		StartedAt:  task.StartedAt,
		Status:     apigen.TaskStatus(task.Status),
		UpdatedAt:  task.UpdatedAt,
		Attributes: task.Attributes,
	}
}

func (w *Worker) Start() {
	for {
		select {
		case <-w.globalCtx.Context().Done():
			return
		case <-time.Tick(1 * time.Second):
			go func() {
				metrics.WorkerGoroutines.Inc()
				defer metrics.WorkerGoroutines.Dec()
				if err := w.runTask(w.globalCtx.Context()); err != nil {
					log.Error("error running task", zap.Error(err))
				}
			}()
		}
	}
}

func (w *Worker) executeTask(ctx context.Context, task apigen.Task) error {
	switch task.Spec.Type {
	case apigen.AutoBackup:
		if task.Spec.AutoBackup == nil {
			return fmt.Errorf("auto backup spec is nil")
		}
		return w.executor.ExecuteAutoBackup(ctx, *task.Spec.AutoBackup)
	case apigen.AutoDiagnostic:
		if task.Spec.AutoDiagnostic == nil {
			return fmt.Errorf("auto diagnostic spec is nil")
		}
		return w.executor.ExecuteAutoDiagnostic(ctx, *task.Spec.AutoDiagnostic)
	default:
		return fmt.Errorf("unknown task type: %s", task.Spec.Type)
	}
}

func (w *Worker) runTask(parentCtx context.Context) error {
	if err := w.model.RunTransaction(parentCtx, func(txm model.ModelInterface) error {
		qtask, err := txm.PullTask(parentCtx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}
		task := taskToAPI(qtask)

		timeout := maxTaskTimeout
		if task.Attributes.Timeout != nil {
			timeout, err = time.ParseDuration(*task.Attributes.Timeout)
			if err != nil {
				return errors.Wrap(err, "failed to parse timeout")
			}
		}
		if timeout > maxTaskTimeout {
			timeout = maxTaskTimeout
		}
		ctx, cancel := context.WithTimeout(parentCtx, timeout)
		defer cancel()

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
		err = w.executeTask(ctx, task)
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
