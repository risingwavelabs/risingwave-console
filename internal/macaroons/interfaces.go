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

package macaroons

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CaveatParser interface {
	Parse(string) (Caveat, error)
}

type Caveat interface {
	Encode() (string, error)

	Decode(string) error

	Type() string

	Validate(*fiber.Ctx) error
}

type MacaroonManagerInterface interface {
	CreateToken(ctx context.Context, userID int32, caveats []Caveat, ttl time.Duration) (*Macaroon, error)

	Parse(ctx context.Context, token string) (*Macaroon, error)

	InvalidateUserTokens(ctx context.Context, userID int32) error
}
