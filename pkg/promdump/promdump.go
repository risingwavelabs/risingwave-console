package promdump

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prom_model "github.com/prometheus/common/model"
)

type DumpOpt struct {
	Endpoint      string
	Start         time.Time
	End           time.Time
	Step          time.Duration
	QueryInterval time.Duration
}

func dump(ctx context.Context, opt *DumpOpt, cb func(name string, value prom_model.Matrix) error) error {
	client, err := api.NewClient(api.Config{
		Address: opt.Endpoint,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create prometheus client")
	}
	v1api := v1.NewAPI(client)

	// get all metric names
	labelValues, warnings, err := v1api.LabelValues(ctx, "__name__", []string{}, opt.Start, opt.End)
	if err != nil {
		return errors.Wrapf(err, "failed to get label values")
	}

	if len(warnings) > 0 {
		return errors.Errorf("warnings: %v", warnings)
	}
	for _, labelValue := range labelValues {
		v, warnings, err := v1api.QueryRange(ctx, string(labelValue), v1.Range{
			Start: opt.Start,
			End:   opt.End,
			Step:  opt.Step,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to query range")
		}
		if len(warnings) > 0 {
			return errors.Errorf("warnings: %v", warnings)
		}
		matrix, ok := v.(prom_model.Matrix)
		if !ok {
			return errors.New("value is not a matrix")
		}
		if cb != nil {
			if err := cb(string(labelValue), matrix); err != nil {
				return errors.Wrapf(err, "failed to run callback")
			}
		}
		if opt.QueryInterval > 0 {
			time.Sleep(opt.QueryInterval)
		}
	}
	return nil
}

func DumpPromToFileWithCallback(ctx context.Context, opt *DumpOpt, filename string, cb func(name string, value prom_model.Matrix) error) error {
	buf := bytes.NewBuffer(nil)
	w := gzip.NewWriter(buf)

	isFirstItem := true
	if err := dump(ctx, opt, func(name string, value prom_model.Matrix) error {
		if !isFirstItem {
			if _, err := w.Write([]byte("\n")); err != nil {
				return errors.Wrapf(err, "failed to write comma")
			}
		} else {
			isFirstItem = false
		}

		for _, series := range value {
			raw, err := json.Marshal(series)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal value")
			}
			if _, err := w.Write(raw); err != nil {
				return errors.Wrapf(err, "failed to write value")
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				return errors.Wrapf(err, "failed to write newline")
			}
		}

		if cb != nil {
			if err := cb(name, value); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to dump")
	}
	w.Close()
	return os.WriteFile(filename, buf.Bytes(), 0644)
}

func DumpPromToFile(ctx context.Context, opt *DumpOpt, filename string) error {
	return DumpPromToFileWithCallback(ctx, opt, filename, nil)
}

// DecompressPromdumpFile decompresses a prometheus dump file and saves the uncompressed
// data to a temporary file. It returns the path to the temporary file. If outDir is not
// provided, it will be saved to the system's default temporary directory.
func DecompressPromdumpFile(ctx context.Context, filename string, outDir string, outFile string, bufSize int) (string, error) {
	compressedFile, err := os.Open(filename)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open compressed file")
	}
	defer compressedFile.Close()

	if outFile == "" {
		outFile = "promdump.ndjson"
	}

	tempFile, err := os.CreateTemp(outDir, outFile)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create temporary file")
	}
	defer tempFile.Close()

	zlibReader, err := gzip.NewReader(compressedFile)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", errors.Wrapf(err, "failed to create zlib reader")
	}
	defer zlibReader.Close()

	buffer := make([]byte, bufSize)
	for {
		n, err := zlibReader.Read(buffer)
		if n > 0 {
			if _, err := tempFile.Write(buffer[:n]); err != nil {
				os.Remove(tempFile.Name())
				return "", errors.Wrapf(err, "failed to write to temporary file")
			}
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			os.Remove(tempFile.Name())
			return "", errors.Wrapf(err, "failed to read from compressed file")
		}
	}

	return tempFile.Name(), nil
}
