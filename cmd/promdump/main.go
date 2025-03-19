package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/prompb"
	"github.com/urfave/cli/v2"

	"database/sql"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/risingwavelabs/wavekit/pkg/promdump"
)

const localStorePath = "metrics.parquet"

// runDump implements the 'dump' command to dump Prometheus data to a file
func runDump(c *cli.Context) error {
	filename := c.String("filename")
	if filename == "" {
		return fmt.Errorf("filename is required")
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

	fmt.Printf("Dumping Prometheus data from %s to file %s\n", endpoint, filename)
	fmt.Printf("Time range: %s to %s with step %s\n", start.Format(time.RFC3339), end.Format(time.RFC3339), step)

	// Create dump options
	opt := &promdump.DumpOpt{
		Endpoint:      endpoint,
		Start:         start,
		End:           end,
		Step:          step,
		QueryInterval: queryInterval,
	}

	// Execute the dump
	err = promdump.DumpPromToFile(context.Background(), opt, filename)
	if err != nil {
		return errors.Wrap(err, "failed to dump prometheus data")
	}

	fmt.Printf("Successfully dumped prometheus data to %s\n", filename)
	return nil
}

// runServe implements the 'serve' command to serve a prometheus dump file
func runServe(c *cli.Context) error {
	filename := c.String("filename")
	if filename == "" {
		return fmt.Errorf("filename is required")
	}

	port := c.Int("port")
	if port == 0 {
		port = 9090 // Default port
	}

	address := c.String("address")

	// Set up the HTTP server
	addr := fmt.Sprintf("%s:%d", address, port)
	fmt.Printf("Serving Prometheus data from file %s on %s\n", filename, addr)

	db, err := sql.Open("duckdb", "promfile.db")
	if err != nil {
		return errors.Wrap(err, "failed to open duckdb")
	}
	defer db.Close()

	// Handler for remote_read API
	http.HandleFunc("/api/v1/read", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement actual handling of remote_read requests using promdump.DecompressPromdumpFile
		fmt.Println("Received remote_read request")

		// This is just a placeholder that will need to be implemented
		// You would need to:
		// 1. Decompress the file using promdump.DecompressPromdumpFile
		// 2. Parse the requested time range from the request
		// 3. Find matching series in the decompressed data
		// 4. Format and return the response

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success","data":[]}`))
	})

	// Basic health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the HTTP server
	if err := http.ListenAndServe(addr, nil); err != nil {
		return errors.Wrap(err, "server error")
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:  "promfile",
		Usage: "Dump Prometheus data to a file, and serve the file as a Prometheus remote_read endpoint",
		Commands: []*cli.Command{
			{
				Name:   "dump",
				Usage:  "Dump Prometheus data to a file",
				Action: runDump,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "filename",
						Aliases:  []string{"f"},
						Usage:    "Output filename",
						Required: true,
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
						Value: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
					},
					&cli.StringFlag{
						Name:  "end",
						Usage: "End time (RFC3339 format)",
						Value: time.Now().Format(time.RFC3339),
					},
					&cli.DurationFlag{
						Name:  "step",
						Usage: "Query step interval",
						Value: time.Minute,
					},
					&cli.DurationFlag{
						Name:  "query-interval",
						Usage: "Interval between queries to avoid overloading Prometheus",
						Value: 0,
					},
				},
			},
			{
				Name:   "serve",
				Usage:  "Serve a Prometheus dump file as a remote_read endpoint",
				Action: runServe,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "filename",
						Aliases:  []string{"f"},
						Usage:    "Input filename (compressed Prometheus dump)",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to serve on",
						Value:   39090,
					},
					&cli.StringFlag{
						Name:    "host",
						Aliases: []string{"h"},
						Usage:   "Host to bind to",
						Value:   "0.0.0.0",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type Store struct {
	db             *sql.DB
	inputFilename  string
	outputFilename string
}

func (s *Store) init(ctx context.Context) error {
	absPath, err := filepath.Abs(s.inputFilename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	if _, err = s.db.ExecContext(ctx, fmt.Sprintf(
		`COPY (
			SELECT 
				metric, 
				value[1]::DOUBLE AS timestamp, 
				value[2]::VARCHAR AS val 
			FROM read_json_auto('%s'), LATERAL UNNEST(values) AS t(value)
		) TO '%s' (FORMAT 'parquet');`,
		absPath,
		s.outputFilename,
	)); err != nil {
		return fmt.Errorf("failed to convert %s to parquet file %s: %v", s.inputFilename, s.outputFilename, err)
	}
	return nil
}

func (s *Store) remoteRead(ctx context.Context, querys prompb.ReadRequest) (*prompb.ReadResponse, error) {
	//query
	query := querys.Queries[0]
	startTime := query.StartTimestampMs
	endTime := query.EndTimestampMs
	matchers := query.Matchers

	stmt := fmt.Sprintf(
		`SELECT * FROM '%s' WHERE %d <= timestamp AND timestamp <= %d`,
		localStorePath,
		startTime,
		endTime,
	)

	for _, matcher := range matchers {
		stmt += fmt.Sprintf(" AND metric['%s'] = '%s'", matcher.Name, matcher.Value)
	}

	rows, err := s.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	defer rows.Close()

	return nil, nil
}
