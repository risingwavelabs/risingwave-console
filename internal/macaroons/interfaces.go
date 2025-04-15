package macaroons

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type KeyStore interface {
	// Create creates a new key and returns the keyID.
	Create(ctx context.Context, key []byte, ttl time.Duration) (int64, error)

	// Get returns the key for the given keyID. returns ErrKeyNotFound if the key is not found.
	Get(ctx context.Context, keyID int64) ([]byte, error)

	// Delete deletes the key for the given keyID. returns ErrKeyNotFound if the key is not found.
	Delete(ctx context.Context, keyID int64) error
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
	CreateToken(ctx context.Context, caveats []Caveat, ttl time.Duration) (*Macaroon, error)

	Parse(ctx context.Context, token string) (*Macaroon, error)

	InvalidateToken(ctx context.Context, token string) error
}
