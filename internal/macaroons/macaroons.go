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
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/macaroons/store"
)

var (
	ErrMalformedToken   = errors.New("malformed token")
	ErrInvalidSignature = errors.New("invalid signature")
)

type Macaroon struct {
	keyID   int64
	caveats []Caveat

	signature []byte

	encodedToken *string
}

func (m *Macaroon) StringToken() string {
	return *m.encodedToken
}

func (m *Macaroon) Caveats() []Caveat {
	return m.caveats
}

func (m *Macaroon) KeyID() int64 {
	return m.keyID
}

type MacaroonManager struct {
	keyStore     store.KeyStore
	caveatParser CaveatParser

	randomKey func() ([]byte, error)
}

func NewMacaroonManager(keyStore store.KeyStore, caveatParser CaveatParser) MacaroonManagerInterface {
	return &MacaroonManager{
		keyStore:     keyStore,
		caveatParser: caveatParser,
		randomKey:    randomKey,
	}
}

func (m *MacaroonManager) CreateToken(ctx context.Context, userID int32, caveats []Caveat, ttl time.Duration) (*Macaroon, error) {
	key, err := m.randomKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random key")
	}
	keyID, err := m.keyStore.Create(ctx, userID, key, ttl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key")
	}

	return CreateMacaroon(keyID, key, caveats)
}

func CreateMacaroon(keyID int64, key []byte, caveats []Caveat) (*Macaroon, error) {
	encodedKeyID := base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(keyID, 10)))
	token := encodedKeyID

	encodedCaveats := make([]string, len(caveats))
	for i, caveat := range caveats {
		encodedCaveat, err := caveat.Encode()
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode caveat")
		}
		encodedCaveats[i] = encodedCaveat
		token += "." + encodedCaveat
	}

	signature := chainedHmac(key, encodedKeyID, encodedCaveats)

	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	token += "." + encodedSignature

	return &Macaroon{
		keyID:        keyID,
		caveats:      caveats,
		signature:    signature,
		encodedToken: &token,
	}, nil
}

func (m *MacaroonManager) Parse(ctx context.Context, token string) (*Macaroon, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, errors.Wrap(ErrMalformedToken, "token must contain at least 2 parts")
	}
	// decode nounce and keyID
	header, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.Wrap(ErrMalformedToken, "failed to decode header")
	}
	keyID, err := strconv.ParseInt(string(header), 10, 64)
	if err != nil {
		return nil, errors.Wrap(ErrMalformedToken, "failed to convert keyID to int")
	}
	// decode signature
	signature, err := base64.StdEncoding.DecodeString(parts[len(parts)-1])
	if err != nil {
		return nil, errors.Wrap(ErrMalformedToken, "failed to decode signature")
	}

	// verify signature
	key, err := m.keyStore.Get(ctx, keyID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key")
	}

	calculatedSignature := chainedHmac(key, parts[0], parts[1:len(parts)-1])

	if !hmac.Equal(signature, calculatedSignature) {
		return nil, ErrInvalidSignature
	}

	// decode caveats
	caveats := make([]Caveat, len(parts)-2)
	for i, part := range parts[1 : len(parts)-1] {
		caveat, err := m.caveatParser.Parse(part)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse caveat")
		}
		caveats[i] = caveat
	}

	return &Macaroon{
		keyID:     keyID,
		caveats:   caveats,
		signature: signature,
	}, nil
}

func (m *MacaroonManager) InvalidateUserTokens(ctx context.Context, userID int32) error {
	if err := m.keyStore.DeleteUserKeys(ctx, userID); err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return nil
		}
		return errors.Wrap(err, "failed to delete user keys")
	}
	return nil
}

func chainedHmac(key []byte, encodedKeyID string, encodedCaveats []string) []byte {
	parts := []string{encodedKeyID}
	parts = append(parts, encodedCaveats...)

	hmac := hmac.New(sha256.New, key)
	for _, part := range parts {
		hmac.Write([]byte(part))
	}
	return hmac.Sum(nil)
}

func randomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random key")
	}
	return key, nil
}
