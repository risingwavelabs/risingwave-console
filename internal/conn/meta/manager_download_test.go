//go:build !ut
// +build !ut

package meta

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadRisectl(t *testing.T) {
	// Create a temporary subdirectory instead of using the root temp dir
	tempDir, err := os.MkdirTemp("", "risectl-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	risectlManager := &RisectlManager{
		risectlDir: tempDir,
	}

	err = risectlManager.downloadRisectl(context.Background(), "v2.2.1")
	require.NoError(t, err)

	risectlPath := filepath.Join(tempDir, "v2.2.1", "risectl")
	require.FileExists(t, risectlPath)

	out, err := exec.CommandContext(context.Background(), risectlPath, "help").Output()
	require.NoError(t, err)
	require.Contains(t, string(out), "Usage:")
}
