package meta

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v68/github"
	"github.com/risingwavelabs/wavekit/internal/config"
)

const risectlFileName = "risectl"

// RisectlManager is a manager for risectl.
type RisectlManager struct {
	risectlDir string
	noInternet bool

	cache []string
	mu    sync.RWMutex
	exp   time.Time
}

func NewRisectlManager(cfg *config.Config) (RisectlManagerInterface, error) {
	risectlDir := cfg.RisectlDir
	if risectlDir == "" {
		risectlDir = filepath.Join(os.Getenv("HOME"), ".risectl")
		log.Default().Printf("Using default risectl dir: %s", risectlDir)
	}

	_, err := os.Stat(risectlDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(risectlDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create risectl dir %s: %w", risectlDir, err)
			}
		} else {
			return nil, fmt.Errorf("failed to stat risectl dir %s: %w", risectlDir, err)
		}
	}

	return &RisectlManager{
		risectlDir: risectlDir,
		noInternet: cfg.NoInternet,
	}, nil
}

func (m *RisectlManager) NewConn(ctx context.Context, version string, host string, port int32) (RisectlConn, error) {
	if version == "" {
		return nil, fmt.Errorf("version is required")
	}

	path := filepath.Join(m.risectlDir, version, risectlFileName)

	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Default().Printf("Downloading risectl binary for version %s", version)
			if err := m.downloadRisectl(ctx, version); err != nil {
				return nil, fmt.Errorf("failed to download risectl binary %s: %w", path, err)
			}
		} else {
			return nil, fmt.Errorf("failed to stat risectl binary %s: %w", path, err)
		}
	}

	return &RisectlConnection{
		version:     version,
		risectlPath: path,
		endpoint:    fmt.Sprintf("http://%s:%d", host, port),
	}, nil
}

func (m *RisectlManager) ListVersions(ctx context.Context) ([]string, error) {
	if m.noInternet {
		return m.listVersionsNoInternet()
	}

	return m.listVersions(ctx)
}

func (m *RisectlManager) listVersionsNoInternet() ([]string, error) {
	dirs, err := os.ReadDir(m.risectlDir)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if dir.IsDir() {
			versions = append(versions, dir.Name())
		}
	}

	return versions, nil
}

func (m *RisectlManager) listVersions(ctx context.Context) ([]string, error) {
	if time.Now().Before(m.exp) {
		m.mu.RLock()
		defer m.mu.RUnlock()
		return m.cache, nil
	}
	client := github.NewClient(nil)

	releases, _, err := client.Repositories.ListReleases(ctx, "risingwavelabs", "risingwave", nil)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(releases))
	for _, release := range releases {
		if release.TagName != nil {
			versions = append(versions, *release.TagName)
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache = versions
	m.exp = time.Now().Add(5 * time.Minute)

	return versions, nil
}

func (m *RisectlManager) downloadRisectl(ctx context.Context, version string) error {
	client := github.NewClient(nil)

	// Get the release by tag
	release, _, err := client.Repositories.GetReleaseByTag(ctx, "risingwavelabs", "risingwave", version)
	if err != nil {
		return fmt.Errorf("failed to get release %s: %w", version, err)
	}

	// Find the risectl asset for x86_64
	var assetID int64
	var assetName string
	for _, asset := range release.Assets {
		if asset.GetName() != "" && strings.HasPrefix(asset.GetName(), "risectl") && strings.Contains(asset.GetName(), "x86_64") {
			assetID = asset.GetID()
			assetName = asset.GetName()
			break
		}
	}
	if assetID == 0 {
		return fmt.Errorf("no risectl x86_64 asset found for version %s", version)
	}

	// Create version directory if it doesn't exist
	versionDir := filepath.Join(m.risectlDir, version)
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", versionDir, err)
	}

	// Download the asset
	rc, _, err := client.Repositories.DownloadReleaseAsset(ctx, "risingwavelabs", "risingwave", assetID, http.DefaultClient)
	if err != nil {
		return fmt.Errorf("failed to download asset %s: %w", assetName, err)
	}
	defer rc.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(rc)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzipReader)

	// Find and extract the risectl binary
	outPath := filepath.Join(versionDir, risectlFileName)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return fmt.Errorf("risectl binary not found in archive")
		}
		if err != nil {
			return fmt.Errorf("failed to read tar archive: %w", err)
		}

		// Look for the risectl binary in the archive
		if strings.HasSuffix(header.Name, "risectl") {
			// Create the output file
			out, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", outPath, err)
			}
			defer out.Close()

			// Copy the content
			if _, err := io.Copy(out, tarReader); err != nil {
				return fmt.Errorf("failed to save asset to %s: %w", outPath, err)
			}
			return nil
		}
	}
}
