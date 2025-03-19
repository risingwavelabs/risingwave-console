package promdump

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDump(t *testing.T) {
	err := DumpPromToFileWithCallback(context.Background(), &DumpOpt{
		Endpoint: "http://localhost:9500",
		Start:    time.Now().Add(-12 * time.Hour),
		End:      time.Now(),
		Step:     5 * time.Second,
	}, "test/test.ndjson.gz", nil)
	require.NoError(t, err)
}

func TestDecompressPromdumpFile(t *testing.T) {
	raw, err := DecompressPromdumpFile(context.Background(), "test.json.gz", "test", "", 64*1024*1024)
	require.NoError(t, err)
	fmt.Println(raw)
}
