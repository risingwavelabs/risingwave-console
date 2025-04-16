package worker

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/globalctx"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"github.com/risingwavelabs/wavekit/internal/metrics"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("worker")

const maxTaskTimeout = 1 * time.Hour

type TaskHandler interface {
	HandleTask(c *modelctx.ModelCtx, task apigen.Task) error
}

type Worker struct {
	model model.ModelInterface

	getHandler LifeCycleHandlerGetter

	globalCtx *globalctx.GlobalContext

	taskHandler TaskHandler
}

func NewWorker(globalCtx *globalctx.GlobalContext, model model.ModelInterface, taskHandler TaskHandler) (*Worker, error) {
	w := &Worker{
		model:       model,
		getHandler:  newTaskLifeCycleHandler,
		globalCtx:   globalCtx,
		taskHandler: taskHandler,
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
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-w.globalCtx.Context().Done():
			return
		case <-ticker.C:
			go func() {
				metrics.WorkerGoroutines.Inc()
				defer metrics.WorkerGoroutines.Dec()
				if err := w.runTask(w.globalCtx.Context()); err != nil {
					metrics.RunTaskErrors.Inc()
					log.Error("error running task", zap.Error(err))
				}
			}()
		}
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

		metrics.PulledTasks.Inc()

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
		err = w.taskHandler.HandleTask(modelctx.NewModelctx(ctx, txm), task)
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
