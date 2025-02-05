package meta

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewRisectlManager(t *testing.T) {
	cfg := &config.Config{
		RisectlDir: "",
	}

	manager, err := NewRisectlManager(cfg)
	require.NoError(t, err)

	m, ok := manager.(*RisectlManager)
	require.True(t, ok)

	require.Equal(t, filepath.Join(os.Getenv("HOME"), ".risectl"), m.risectlDir)
}
