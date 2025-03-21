package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang/snappy"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/prompb"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"

	_ "github.com/marcboeker/go-duckdb/v2"
)

func main() {
	app := &cli.App{
		Name:  "promduck",
		Usage: "Serve the file exported by promdump as a Prometheus remote_read endpoint",
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "Serve a Prometheus dump file as a remote_read endpoint",
				Action: runServe,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "rebuild",
						Aliases: []string{"r"},
						Usage:   "Rebuild the store, otherwise it will use the existing store",
					},
					&cli.StringFlag{
						Name:     "filename",
						Aliases:  []string{"f"},
						Usage:    "Input filename (compressed Prometheus dump)",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "out",
						Aliases: []string{"o"},
						Usage:   "Output directory",
						Value:   "",
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to serve on",
						Value:   39090,
					},
					&cli.StringFlag{
						Name:    "host",
						Aliases: []string{"H"},
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

// runServe implements the 'serve' command to serve a prometheus dump file
func runServe(c *cli.Context) error {
	inputFilename := c.String("filename")
	if inputFilename == "" {
		return fmt.Errorf("filename is required")
	}

	outDir := c.String("out")
	if outDir == "" {
		outDir = fmt.Sprintf("%s/.promdump", os.Getenv("HOME"))
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create output directory")
	}

	port := c.Int("port")
	if port == 0 {
		port = 9090 // Default port
	}

	address := c.String("address")

	// Set up the HTTP server
	addr := fmt.Sprintf("%s:%d", address, port)

	db, err := sql.Open("duckdb", "")
	if err != nil {
		return errors.Wrap(err, "failed to open duckdb")
	}
	defer db.Close()

	store := &Store{
		db:             db,
		inputFilename:  inputFilename,
		outputFilename: filepath.Join(outDir, "metrics.parquet"),
	}

	if err := store.init(context.Background(), c.Bool("rebuild")); err != nil {
		return errors.Wrap(err, "failed to init store")
	}

	// Handler for remote_read API
	http.HandleFunc("/api/v1/read", func(w http.ResponseWriter, r *http.Request) {
		compressed, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reqProtobuf, err := snappy.Decode(nil, compressed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// resolve json
		var req prompb.ReadRequest
		if err := proto.Unmarshal(reqProtobuf, protoadapt.MessageV2Of(&req)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := store.remoteRead(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resProtobuf, err := proto.Marshal(protoadapt.MessageV2Of(data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/x-protobuf")
		w.Header().Set("Content-Encoding", "snappy")
		if _, err := w.Write(snappy.Encode(nil, resProtobuf)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Basic health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the HTTP server
	fmt.Printf("Serving Prometheus data from file %s on %s\n", inputFilename, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return errors.Wrap(err, "server error")
	}

	return nil
}

type Store struct {
	db             *sql.DB
	inputFilename  string
	outputFilename string
}

func (s *Store) init(ctx context.Context, rebuild bool) error {
	fmt.Println("Initializing store...", s.outputFilename)
	if _, err := os.Stat(s.outputFilename); !rebuild && err == nil {
		fmt.Println("Store already exists, skipping. Add -r to rebuild the store if needed")
		return nil
	}

	if _, err := s.db.ExecContext(ctx, fmt.Sprintf(
		`COPY (
			SELECT 
				metric, 
				(value[1]::DOUBLE * 1000)::BIGINT AS timestamp, 
				value[2]::VARCHAR AS val 
			FROM read_json_auto('%s'), LATERAL UNNEST(values) AS t(value)
		) TO '%s' (FORMAT 'parquet');`,
		s.inputFilename,
		s.outputFilename,
	)); err != nil {
		return fmt.Errorf("failed to convert %s to parquet file %s: %v", s.inputFilename, s.outputFilename, err)
	}
	return nil
}

func (s *Store) remoteRead(ctx context.Context, querys prompb.ReadRequest) (*prompb.ReadResponse, error) {
	//query
	var res = prompb.ReadResponse{}
	for _, query := range querys.Queries {
		startTime := query.StartTimestampMs
		endTime := query.EndTimestampMs
		matchers := query.Matchers

		stmt := fmt.Sprintf(
			`
			SELECT 
				metric, 
				JSON_GROUP_ARRAY(json_data) AS time_series
			FROM (
				SELECT 
					metric, 
					JSON('{ "timestamp": ' || timestamp || ', "value": ' || val || ' }') AS json_data
				FROM read_parquet('%s') WHERE %d <= timestamp AND timestamp <= %d`,
			s.outputFilename,
			startTime,
			endTime,
		)

		for _, matcher := range matchers {
			switch matcher.Type {
			case prompb.LabelMatcher_EQ:
				stmt += fmt.Sprintf(" AND metric['%s'] = '%s'", matcher.Name, matcher.Value)
			case prompb.LabelMatcher_NEQ:
				stmt += fmt.Sprintf(" AND metric['%s'] != '%s'", matcher.Name, matcher.Value)
			case prompb.LabelMatcher_RE:
				if matcher.Value == "" || matcher.Value == "()" || matcher.Value == ".*" {
					continue
				}
				stmt += fmt.Sprintf(" AND regexp_matches(metric['%s'], '%s')", matcher.Name, matcher.Value)
			case prompb.LabelMatcher_NRE:
				stmt += fmt.Sprintf(" AND NOT regexp_matches(metric['%s'], '%s')", matcher.Name, matcher.Value)
			}
		}

		stmt += "\nORDER BY metric, timestamp) sorted_data\nGROUP BY metric"

		fmt.Println(stmt)

		rows, err := s.db.QueryContext(ctx, stmt)
		if err != nil {
			return nil, fmt.Errorf("failed to query: %v", err)
		}
		defer rows.Close()

		timeSeries := []*prompb.TimeSeries{}

		for rows.Next() {
			var (
				metric     map[any]any
				rawSamples []any
			)

			if err := rows.Scan(&metric, &rawSamples); err != nil {
				return nil, fmt.Errorf("failed to scan row: %v", err)
			}

			labels := []prompb.Label{}
			for k, v := range metric {
				labels = append(labels, prompb.Label{
					Name:  k.(string),
					Value: v.(string),
				})
			}

			samples := []prompb.Sample{}
			for _, sample := range rawSamples {
				sampleMap := sample.(map[string]any)
				timestamp := int64(sampleMap["timestamp"].(float64))
				value := sampleMap["value"].(string)
				floatValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse value: %v", err)
				}

				samples = append(samples, prompb.Sample{
					Timestamp: timestamp,
					Value:     floatValue,
				})
			}

			timeSeries = append(timeSeries, &prompb.TimeSeries{
				Labels:  labels,
				Samples: samples,
			})
		}

		res.Results = append(res.Results, &prompb.QueryResult{
			Timeseries: timeSeries,
		})
	}

	return &res, nil
}
