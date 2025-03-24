package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/pkg/promdump"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "promdump",
		Usage:  "Dump Prometheus data to a *.ndjson.gz file, and serve the file as a Prometheus remote_read endpoint",
		Action: runDump,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "Output filename",
				Value:   "out.ndjson.gz",
			},
			&cli.StringFlag{
				Name:     "endpoint",
				Aliases:  []string{"e"},
				Usage:    "Prometheus endpoint URL",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "start",
				Usage: "Start time (RFC3339 format)",
				Value: time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339),
			},
			&cli.StringFlag{
				Name:  "end",
				Usage: "End time (RFC3339 format)",
				Value: time.Now().Format(time.RFC3339),
			},
			&cli.DurationFlag{
				Name:  "step",
				Usage: "Format: 1s, 1m, 1h, 1d, default is 1s.",
				Value: time.Second,
			},
			&cli.StringFlag{
				Name: "query",
				Usage: "PromQL query to filter time series, e.g. use {risingwave_cluster=\"default\"} " +
					"to dump all time series with the label risingwave_cluster=default. " +
					"If not provided, all time series will be dumped.",
				Value: "",
			},
			&cli.BoolFlag{
				Name:  "plain",
				Usage: "Output in uncompressed NDJSON format",
				Value: false,
			},
			&cli.Float64Flag{
				Name:  "query-ratio",
				Usage: "(0, 1], if OOM, reduce the memory usage in Prometheus instance by this ratio",
				Value: 1,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// runDump implements the 'dump' command to dump Prometheus data to a file
func runDump(c *cli.Context) error {
	outFilename := c.String("out")
	if outFilename == "" {
		return fmt.Errorf("out is required")
	}

	endpoint := c.String("endpoint")
	if endpoint == "" {
		return fmt.Errorf("prometheus endpoint is required")
	}

	startStr := c.String("start")
	endStr := c.String("end")
	step := c.Duration("step")
	queryInterval := c.Duration("query-interval")

	// Parse time strings
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse start time")
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse end time")
	}

	fmt.Printf("Dumping Prometheus data from %s to file %s\n", endpoint, outFilename)
	fmt.Printf("Time range: %s to %s with step %s\n", start.Format(time.RFC3339), end.Format(time.RFC3339), step)

	// Create dump options
	opt := &promdump.DumpOpt{
		Endpoint:      endpoint,
		Start:         start,
		End:           end,
		Step:          step,
		QueryInterval: queryInterval,
		Query:         c.String("query"),
		Plain:         c.Bool("plain"),
		QueryRatio:    c.Float64("query-ratio"),
	}

	// mkdir -p
	outDir := filepath.Dir(outFilename)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	// Execute the dump
	err = promdump.DumpPromToFile(context.Background(), opt, outFilename, true)
	if err != nil {
		return errors.Wrap(err, "failed to dump prometheus data")
	}

	fmt.Printf("Successfully dumped prometheus data to %s\n", outFilename)
	return nil
}
