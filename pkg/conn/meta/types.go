package meta

import "context"

type RisectlConn interface {
	RunCombined(ctx context.Context, args ...string) (string, int, error)
	Run(ctx context.Context, args ...string) (string, string, int, error)
	MetaBackup(ctx context.Context) (int64, error)
	DeleteSnapshot(ctx context.Context, snapshotID int64) error
}

type RisectlManagerInterface interface {
	ListVersions(ctx context.Context) ([]string, error)
	NewConn(ctx context.Context, version string, host string, metaPort int32) (RisectlConn, error)
}
