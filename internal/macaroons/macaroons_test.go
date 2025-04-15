package macaroons

import (
	"context"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

type TestCaveat struct {
	Data string
}

func (c *TestCaveat) Encode() (string, error) {
	return c.Data, nil
}

func (c *TestCaveat) Decode(s string) error {
	c.Data = s
	return nil
}

func (c *TestCaveat) Type() string {
	return "test"
}

func (c *TestCaveat) Validate(*fiber.Ctx) error {
	return nil
}

func TestMacaroonManager_CreateMacaroon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	keyStore := NewMockKeyStore(ctrl)
	caveatParser := NewMockCaveatParser(ctrl)

	var (
		keyID   = int64(9527)
		caveats = []Caveat{
			&TestCaveat{Data: "caveat1"},
			&TestCaveat{Data: "caveat2"},
		}
		ttl    = time.Second * 10
		userID = int32(1)
	)

	keyStore.EXPECT().Create(gomock.Any(), userID, []byte("key"), ttl).Return(keyID, nil)
	keyStore.EXPECT().Get(gomock.Any(), keyID).Return([]byte("key"), nil)

	caveatParser.EXPECT().Parse("caveat1").Return(caveats[0], nil)
	caveatParser.EXPECT().Parse("caveat2").Return(caveats[1], nil)

	manager := &MacaroonManager{
		keyStore:     keyStore,
		caveatParser: caveatParser,
		randomKey:    func() ([]byte, error) { return []byte("key"), nil },
	}

	macaroon, err := manager.CreateToken(context.Background(), userID, caveats, ttl)
	require.NoError(t, err)

	parsed, err := manager.Parse(context.Background(), macaroon.StringToken())
	require.NoError(t, err)
	require.Equal(t, keyID, parsed.keyID)
	require.Equal(t, caveats, parsed.caveats)
}
