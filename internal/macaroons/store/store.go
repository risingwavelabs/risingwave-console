// Copyright 2025 RisingWave Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
	"github.com/risingwavelabs/wavekit/internal/task"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Store struct {
	model     model.ModelInterface
	taskStore task.TaskStoreInterface
	now       func() time.Time
}

func NewStore(model model.ModelInterface, taskStore task.TaskStoreInterface) KeyStore {
	return &Store{
		model:     model,
		taskStore: taskStore,
		now:       time.Now,
	}
}

func (s *Store) Create(ctx context.Context, userID int32, key []byte, ttl time.Duration) (int64, error) {
	var ret int64
	if err := s.model.RunTransaction(ctx, func(txm model.ModelInterface) error {
		keyID, err := txm.CreateOpaqueKey(ctx, querier.CreateOpaqueKeyParams{
			UserID: userID,
			Key:    key,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create key")
		}

		ret = keyID

		if ttl > 0 {
			c := modelctx.NewModelctx(ctx, txm)
			if _, err := s.taskStore.PushTask(
				c,
				apigen.TaskSpec{
					Type: apigen.DeleteOpaqueKey,
					DeleteOpaqueKey: &apigen.TaskSpecDeleteOpaqueKey{
						KeyID: keyID,
					},
				},
				task.StartedAt(s.now().Add(ttl)),
				task.AlwaysRetryOnFailure("30m"),
			); err != nil {
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

func (s *Store) DeleteUserKeys(ctx context.Context, userID int32) error {
	err := s.model.DeleteOpaqueKeys(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrKeyNotFound
		}
		return errors.Wrap(err, "failed to delete user keys")
	}
	return nil
}
