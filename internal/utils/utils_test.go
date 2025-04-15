package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieveFromJSON(t *testing.T) {
	s := `{"type": "test", "value": "test"}`

	ret, err := RetrieveFromJSON[string](s, "type")

	assert.NoError(t, err)
	assert.Equal(t, "test", *ret)
}
