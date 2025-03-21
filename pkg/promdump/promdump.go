package promdump

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prom_model "github.com/prometheus/common/model"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

const PrometheusDefaultMaxResolution = 11_000

type QueryCallback func(query string, value prom_model.Matrix, progress float64) error

type DumpOpt struct {
	Endpoint      string
	Start         time.Time
	End           time.Time
	Step          time.Duration
	QueryInterval time.Duration
	Query         string
	Plain         bool
	QueryRatio    float64
}

type QueryRangeChunkParams struct {
	Start time.Time
	End   time.Time
}

func queryChunks(start time.Time, end time.Time, step time.Duration, queryRatio float64) []QueryRangeChunkParams {
	maxDuration := time.Duration(float64(PrometheusDefaultMaxResolution)*queryRatio) * step
	chunks := []QueryRangeChunkParams{}
	for {
		d := end.Sub(start)
		if d < maxDuration {
			chunks = append(chunks, QueryRangeChunkParams{
				Start: start,
				End:   end,
			})
			break
		}
		chunks = append(chunks, QueryRangeChunkParams{
			Start: start.Add(step), // last end is inclusive
			End:   start.Add(maxDuration),
		})
		start = start.Add(maxDuration)
	}
	return chunks
}

func dump(ctx context.Context, opt *DumpOpt, cb QueryCallback) error {
	client, err := api.NewClient(api.Config{
		Address: opt.Endpoint,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create prometheus client")
	}

	v1api := v1.NewAPI(client)

	var queries []string
	if len(opt.Query) > 0 {
		queries = []string{opt.Query}
	} else {
		// get all metric names
		labelValues, warnings, err := v1api.LabelValues(ctx, "__name__", []string{}, opt.Start, opt.End)
		if err != nil {
			return errors.Wrapf(err, "failed to get label values")
		}
		if len(warnings) > 0 {
			return errors.Errorf("warnings: %v", warnings)
		}
		for _, labelValue := range labelValues {
			queries = append(queries, string(labelValue))
		}
	}

	chunks := queryChunks(opt.Start, opt.End, opt.Step, opt.QueryRatio)
	for pi, query := range queries {
		vs, warnings, err := queryFullRange(ctx, v1api, string(query), opt.Step, chunks)
		if err != nil {
			return errors.Wrapf(err, "failed to query range")
		}
		if len(warnings) > 0 {
			return errors.Errorf("warnings: %v", warnings)
		}
		for pj, v := range vs {
			matrix, ok := v.(prom_model.Matrix)
			if !ok {
				return errors.New("value is not a matrix")
			}
			progress := float64(pi+1)/float64(len(queries)) +
				float64(pj+1)/float64(len(vs))*(1/float64(len(queries)))
			if cb != nil {
				if err := cb(string(query), matrix, progress); err != nil {
					return errors.Wrapf(err, "failed to run callback")
				}
			}
		}
	}
	return nil
}

func DumpPromToFileWithCallback(ctx context.Context, opt *DumpOpt, filename string, cb QueryCallback) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to open file")
	}
	defer f.Close()

	var w io.Writer
	if !opt.Plain {
		gw := gzip.NewWriter(f)
		defer gw.Close()
		w = gw
	} else {
		w = f
	}

	isFirstItem := true
	if err := dump(ctx, opt, func(query string, value prom_model.Matrix, progress float64) error {
		if !isFirstItem {
			if _, err := w.Write([]byte("\n")); err != nil {
				return errors.Wrapf(err, "failed to write comma")
			}
		} else {
			isFirstItem = false
		}

		for i, series := range value {
			raw, err := json.Marshal(series)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal value")
			}
			if len(raw) == 0 {
				continue
			}
			if _, err := w.Write(raw); err != nil {
				return errors.Wrapf(err, "failed to write value")
			}
			if i < len(value)-1 {
				if _, err := w.Write([]byte("\n")); err != nil {
					return errors.Wrapf(err, "failed to write newline")
				}
			}
		}

		if cb != nil {
			if err := cb(query, value, progress); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to dump")
	}
	return nil
}

func DumpPromToFile(ctx context.Context, opt *DumpOpt, filename string, showProgress bool) error {
	var lastProgress float64
	if err := DumpPromToFileWithCallback(ctx, opt, filename, func(query string, value prom_model.Matrix, progress float64) error {
		if showProgress {
			if progress-lastProgress > 0.01 {
				// Clear the line and print the progress bar with percentage
				fmt.Printf("\033[2K\rprogress: %s", utils.RenderProgressBar(progress))
				lastProgress = progress
			}
		}
		return nil
	}); err != nil {
		fmt.Println()
		return errors.Wrapf(err, "failed to dump")
	}

	// Clear the line and print final progress
	fmt.Printf("\033[2K\rprogress: %s\n", utils.RenderProgressBar(1.0))
	return nil
}

func queryFullRange(ctx context.Context, v1api v1.API, query string, step time.Duration, chunks []QueryRangeChunkParams, opts ...v1.Option) ([]prom_model.Value, v1.Warnings, error) {
	var vs []prom_model.Value
	var retWarnings v1.Warnings
	for _, chunk := range chunks {
		v, warnings, err := v1api.QueryRange(ctx, query, v1.Range{
			Start: chunk.Start,
			End:   chunk.End,
			Step:  step,
		}, opts...)
		if err != nil {
			return nil, warnings, errors.Wrapf(err, "failed to query range")
		}
		vs = append(vs, v)
		retWarnings = append(retWarnings, warnings...)
	}
	return vs, retWarnings, nil
}
