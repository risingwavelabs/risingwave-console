package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "prompush",
		Usage:  "Push Prometheus data to a remote endpoint",
		Action: runPush,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "filename",
				Aliases:  []string{"f"},
				Usage:    "Input filename",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "vm-endpoint",
				Aliases:  []string{"e"},
				Usage:    "VictoriaMetrics endpoint URL",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "use-legacy-format",
				Usage:    "Use legacy format",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "batch-size",
				Aliases:  []string{"b"},
				Usage:    "Batch size",
				Required: false,
				Value:    1000,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type LegacyFormat struct {
	Metric map[string]string `json:"metric"`
	Values [][]any           `json:"values"`
}

type Item struct {
	Metric     map[string]string `json:"metric"`
	Values     []float64         `json:"values"`
	Timestamps []int64           `json:"timestamps"`
}

func runPush(c *cli.Context) error {
	vmEndpoint := c.String("vm-endpoint")
	useLegacyFormat := c.Bool("use-legacy-format")
	filename := c.String("filename")
	batchSize := c.Int("batch-size")

	if batchSize <= 0 {
		return fmt.Errorf("batch-size must be greater than 0")
	}
	if len(vmEndpoint) == 0 {
		return fmt.Errorf("vm-endpoint is required")
	}
	if len(filename) == 0 {
		return fmt.Errorf("filename is required")
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	pw := NewPushWorker(c.Context, vmEndpoint, batchSize)
	defer pw.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	scanner := bufio.NewScanner(gzReader)
	scanner.Buffer(make([]byte, 1024*1024), 100*1024*1024) // 100MB max line size
	for scanner.Scan() {
		line := scanner.Bytes()

		// Get current position in the compressed file
		currentPos, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("failed to get current file position: %w", err)
		}

		// Display progress based on compressed file position
		fmt.Printf("\033[2K\rprogress: %s", utils.RenderProgressBar(float64(currentPos)/float64(fileSize)))

		if len(line) == 0 {
			continue
		}
		if !useLegacyFormat {
			pw.Push(line)
			continue
		}

		var legacy LegacyFormat
		if err := json.Unmarshal(line, &legacy); err != nil {
			return fmt.Errorf("failed to unmarshal line: %w", err)
		}
		item := Item{
			Metric: legacy.Metric,
		}

		for _, v := range legacy.Values {
			val, err := strconv.ParseFloat(v[1].(string), 64)
			if err != nil {
				return fmt.Errorf("failed to parse timestamp: %w", err)
			}
			item.Timestamps = append(item.Timestamps, int64(1000*v[0].(float64)))
			item.Values = append(item.Values, val)
		}
		l, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal item: %w", err)
		}
		l = append(l, '\n')
		pw.Push(l)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return pw.Flush(c.Context)
}

type PushWorker struct {
	vmEndpoint string
	c          chan []byte
	buf        bytes.Buffer
	cnt        int
	mu         sync.Mutex
}

func NewPushWorker(ctx context.Context, vmEndpoint string, batchSize int) *PushWorker {
	w := &PushWorker{
		vmEndpoint: vmEndpoint,
		c:          make(chan []byte, batchSize),
		cnt:        0,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case line := <-w.c:
				if w.cnt == batchSize {
					if err := w.Flush(ctx); err != nil {
						log.Printf("failed to flush: %s", err)
					}
				}
				w.Append(line)
			}
		}
	}()

	return w
}

func (w *PushWorker) Append(line []byte) {
	w.buf.Write(line)
	w.cnt++
}

func (w *PushWorker) Push(line []byte) {
	w.c <- line
}

func (w *PushWorker) Flush(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buf.Len() == 0 {
		return nil
	}

	// Create a copy of the buffer data to avoid modification during HTTP transfer
	data := make([]byte, w.buf.Len())
	copy(data, w.buf.Bytes())

	// Use http.Client directly instead of http.Post to keep mutex locked during transmission
	req, err := http.NewRequestWithContext(ctx, "POST", w.vmEndpoint+"/api/v1/import", bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/jsonl")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to push metrics")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to push metrics: status=%d body=%s", resp.StatusCode, string(body))
	}

	// reset the buffer
	w.buf.Reset()
	w.cnt = 0
	return nil
}

func (w *PushWorker) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = w.Flush(ctx)
	close(w.c)
}
