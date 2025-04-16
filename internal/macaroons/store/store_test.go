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
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/modelctx"
	"github.com/risingwavelabs/wavekit/internal/task"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	model := model.NewExtendedMockModelInterface(ctrl)
	taskStore := task.NewMockTaskStoreInterface(ctrl)

	var (
		ctx      = context.Background()
		ttl      = 1 * time.Hour
		key      = []byte("test")
		userID   = int32(201)
		currTime = time.Now()
		keyID    = int64(101)
		taskID   = int32(201)
	)

	c := &modelctx.ModelCtx{
		ModelInterface: model,
		Context:        ctx,
	}

	model.EXPECT().CreateOpaqueKey(gomock.Any(), gomock.Any()).Return(keyID, nil)
	taskStore.EXPECT().PushTask(
		c,
		apigen.TaskSpec{
			Type: apigen.DeleteOpaqueKey,
			DeleteOpaqueKey: &apigen.TaskSpecDeleteOpaqueKey{
				KeyID: keyID,
			},
		},
		task.MatchTaskOpt(task.StartedAt(currTime.Add(ttl))),
		task.MatchTaskOpt(task.AlwaysRetryOnFailure("30m")),
	).Return(taskID, nil)

	store := &Store{
		model:     model,
		taskStore: taskStore,
		now:       func() time.Time { return currTime },
	}

	ret, err := store.Create(ctx, userID, key, ttl)
	require.NoError(t, err)
	require.Equal(t, keyID, ret)
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		ctx   = context.Background()
		keyID = int64(101)
	)

	var testCases = []struct {
		name string
		err  error
	}{
		{
			name: "success",
			err:  nil,
		},
		{
			name: "no row",
			err:  pgx.ErrNoRows,
		},
		{
			name: "error",
			err:  errors.New("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := model.NewExtendedMockModelInterface(ctrl)

			store := &Store{
				model: model,
			}

			if tc.err == nil {
				model.EXPECT().DeleteOpaqueKey(gomock.Any(), keyID).Return(nil)
			} else {
				model.EXPECT().DeleteOpaqueKey(gomock.Any(), keyID).Return(tc.err)
			}

			err := store.Delete(ctx, keyID)
			if tc.err == nil {
				require.NoError(t, err)
			} else if tc.err == pgx.ErrNoRows {
				require.ErrorIs(t, err, ErrKeyNotFound)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var testCases = []struct {
		name string
		err  error
	}{
		{
			name: "success",
			err:  nil,
		},
		{
			name: "no row",
			err:  pgx.ErrNoRows,
		},
		{
			name: "error",
			err:  errors.New("error"),
		},
	}

	var (
		ctx   = context.Background()
		keyID = int64(101)
		key   = []byte("test")
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := model.NewExtendedMockModelInterface(ctrl)

			store := &Store{
				model: model,
			}

			if tc.err == nil {
				model.EXPECT().GetOpaqueKey(gomock.Any(), keyID).Return(key, nil)
			} else {
				model.EXPECT().GetOpaqueKey(gomock.Any(), keyID).Return(nil, tc.err)
			}

			key, err := store.Get(ctx, keyID)
			if tc.err == nil {
				require.NoError(t, err)
				require.Equal(t, key, key)
			} else if tc.err == pgx.ErrNoRows {
				require.ErrorIs(t, err, ErrKeyNotFound)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestDeleteUserKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var testCases = []struct {
		name string
		err  error
	}{
		{
			name: "success",
			err:  nil,
		},
		{
			name: "no row",
			err:  pgx.ErrNoRows,
		},
		{
			name: "error",
			err:  errors.New("error"),
		},
	}

	var (
		ctx    = context.Background()
		userID = int32(201)
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := model.NewExtendedMockModelInterface(ctrl)

			store := &Store{
				model: model,
			}

			if tc.err == nil {
				model.EXPECT().DeleteOpaqueKeys(gomock.Any(), userID).Return(nil)
			} else {
				model.EXPECT().DeleteOpaqueKeys(gomock.Any(), userID).Return(tc.err)
			}

			err := store.DeleteUserKeys(ctx, userID)
			if tc.err == nil {
				require.NoError(t, err)
			} else if tc.err == pgx.ErrNoRows {
				require.ErrorIs(t, err, ErrKeyNotFound)
			} else {
				require.Error(t, err)
			}
		})
	}
}
