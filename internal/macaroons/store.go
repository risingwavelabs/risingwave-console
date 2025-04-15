package macaroons

import (
	"context"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/task"
	"github.com/risingwavelabs/wavekit/internal/utils"

	"github.com/risingwavelabs/wavekit/internal/apigen"
)

type Store struct {
	model     model.ModelInterface
	taskStore task.TaskStoreInterface
}

func NewStore(model model.ModelInterface, taskStore task.TaskStoreInterface) KeyStore {
	return &Store{
		model:     model,
		taskStore: taskStore,
	}
}

func (s *Store) Create(ctx context.Context, key []byte, ttl time.Duration) (int64, error) {
	var ret int64
	if err := s.model.RunTransaction(ctx, func(txm model.ModelInterface) error {
		keyID, err := txm.CreateOpaqueKey(ctx, key)
		if err != nil {
			return errors.Wrap(err, "failed to create key")
		}
		ret = keyID
		if ttl > 0 {
			if _, err := txm.CreateTask(ctx, querier.CreateTaskParams{
				Attributes: apigen.TaskAttributes{
					RetryPolicy: &apigen.TaskRetryPolicy{
						Interval:             "30m",
						AlwaysRetryOnFailure: utils.Ptr(true),
					},
				},
				Spec: apigen.TaskSpec{
					Type: apigen.DeleteOpaqueKey,
					DeleteOpaqueKey: &apigen.TaskSpecDeleteOpaqueKey{
						KeyID: keyID,
					},
				},
			}); err != nil {
				return errors.Wrap(err, "failed to create task")
			}
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return ret, nil
}

func (s *Store) Get(ctx context.Context, keyID int64) ([]byte, error) {
	key, err := s.model.GetOpaqueKey(ctx, keyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrKeyNotFound
		}
		return nil, errors.Wrap(err, "failed to get key")
	}

	return key, nil
}

func (s *Store) Delete(ctx context.Context, keyID int64) error {
	err := s.model.DeleteOpaqueKey(ctx, keyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrKeyNotFound
		}
		return errors.Wrap(err, "failed to delete key")
	}
	return nil
}
