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
)

type KeyStore interface {
	// Create creates a new key and returns the keyID.
	Create(ctx context.Context, userID int32, key []byte, ttl time.Duration) (int64, error)

	// Get returns the key for the given keyID. returns ErrKeyNotFound if the key is not found.
	Get(ctx context.Context, keyID int64) ([]byte, error)

	// Delete deletes the key for the given keyID. returns ErrKeyNotFound if the key is not found.
	Delete(ctx context.Context, keyID int64) error

	// DeleteUserKeys deletes all keys for the given userID.
	DeleteUserKeys(ctx context.Context, userID int32) error
}
