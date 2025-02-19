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
	risectlDir := os.TempDir()
	defer os.RemoveAll(risectlDir)

	risectlManager := &RisectlManager{
		risectlDir: risectlDir,
	}

	err := risectlManager.downloadRisectl(context.Background(), "v2.2.1")
	require.NoError(t, err)

	risectlPath := filepath.Join(risectlDir, "v2.2.1", "risectl")
	require.FileExists(t, risectlPath)

	out, err := exec.CommandContext(context.Background(), risectlPath, "help").Output()
	require.NoError(t, err)
	require.Contains(t, string(out), "Usage:")
}
