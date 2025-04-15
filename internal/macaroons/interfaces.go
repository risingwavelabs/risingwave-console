package macaroons

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
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
