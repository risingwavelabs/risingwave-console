package service

import (
	"context"
	"time"

	"github.com/cloudcarver/anchor/pkg/auth"
	anchor_svc "github.com/cloudcarver/anchor/pkg/service"
	"github.com/cloudcarver/anchor/pkg/taskcore"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/conn/http"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/conn/metricsstore"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/zcore/model"
	"github.com/risingwavelabs/wavekit/internal/zgen/apigen"
	"github.com/risingwavelabs/wavekit/internal/zgen/taskgen"

	prom_model "github.com/prometheus/common/model"
)

type (
	TradeType   string
	TradeStatus string
	DdWorkEvent string
)

var (
	ErrUserNotFound                  = errors.New("user not found")
	ErrInvalidPassword               = errors.New("invalid password")
	ErrRefreshTokenExpired           = errors.New("refresh token expired")
	ErrDatabaseNotFound              = errors.New("database not found")
	ErrClusterNotFound               = errors.New("cluster not found")
	ErrClusterHasDatabaseConnections = errors.New("cluster has database connections")
	ErrDiagnosticNotFound            = errors.New("diagnostic not found")
)

const (
	ExpireDuration             = 2 * time.Minute
	DefaultMaxRetries          = 3
	RefreshTokenExpireDuration = 14 * 24 * time.Hour
)

type ServiceInterface interface {
	// Cluster management
	ImportCluster(ctx context.Context, params apigen.ClusterImport, orgID int32) (*apigen.Cluster, error)

	// GetCluster gets a cluster by its ID
	GetCluster(ctx context.Context, id int32, orgID int32) (*apigen.Cluster, error)

	// ListClusters lists all clusters in an organization
	ListClusters(ctx context.Context, orgID int32) ([]*apigen.Cluster, error)

	// UpdateCluster updates a cluster
	UpdateCluster(ctx context.Context, id int32, params apigen.ClusterImport, orgID int32) (*apigen.Cluster, error)

	// DeleteCluster deletes a cluster
	DeleteCluster(ctx context.Context, id int32, cascade bool, orgID int32) error

	// Database management
	ImportDatabase(ctx context.Context, params apigen.DatabaseConnectInfo, orgID int32) (*apigen.Database, error)

	// GetDatabase gets a database by its ID and organization ID
	GetDatabase(ctx context.Context, id int32, orgID int32) (*apigen.Database, error)

	// ListDatabases lists all databases in an organization
	ListDatabases(ctx context.Context, orgID int32) ([]apigen.Database, error)

	// UpdateDatabase updates a database
	UpdateDatabase(ctx context.Context, id int32, params apigen.DatabaseConnectInfo, orgID int32) (*apigen.Database, error)

	// DeleteDatabase deletes a database
	DeleteDatabase(ctx context.Context, id int32, orgID int32) error

	// TestDatabaseConnection tests a database connection
	TestDatabaseConnection(ctx context.Context, params apigen.TestDatabaseConnectionPayload, orgID int32) (*apigen.TestDatabaseConnectionResult, error)

	// QueryDatabase executes a query on a database
	QueryDatabase(ctx context.Context, id int32, params apigen.QueryRequest, orgID int32, backgroundDDL bool) (*apigen.QueryResponse, error)

	// GetDDLProgress gets the progress of DDL operations
	GetDDLProgress(ctx context.Context, id int32, orgID int32) ([]apigen.DDLProgress, error)

	// CancelDDLProgress cancels a DDL operation
	CancelDDLProgress(ctx context.Context, id int32, ddlID int64, orgID int32) error

	// ListClusterVersions lists all supported cluster versions via github releases API, if
	// internet is disabled, fallback to list local files.
	ListClusterVersions(ctx context.Context) ([]string, error)

	// ClusterSnapshot management
	CreateClusterSnapshot(ctx context.Context, id int32, name string, orgID int32) (*apigen.Snapshot, error)

	// ListClusterSnapshots lists all snapshots of a cluster
	ListClusterSnapshots(ctx context.Context, id int32, orgID int32) ([]apigen.Snapshot, error)

	// DeleteClusterSnapshot deletes a snapshot from a cluster
	DeleteClusterSnapshot(ctx context.Context, id int32, snapshotID int64, orgID int32) error

	// TestClusterConnection tests the connection to a cluster
	TestClusterConnection(ctx context.Context, params apigen.TestClusterConnectionPayload, orgID int32) (*apigen.TestClusterConnectionResult, error)

	// RunRisectlCommand executes a risectl command on a cluster
	RunRisectlCommand(ctx context.Context, id int32, params apigen.RisectlCommand, orgID int32) (*apigen.RisectlCommandResult, error)

	// GetClusterDiagnostic gets diagnostic information dump for a cluster by ID
	GetClusterDiagnostic(ctx context.Context, id int32, diagnosticID int32, orgID int32) (*apigen.DiagnosticData, error)

	// ListClusterDiagnostics lists all diagnostic information dumps for a cluster
	ListClusterDiagnostics(ctx context.Context, id int32, orgID int32) ([]apigen.DiagnosticData, error)

	// CreateClusterDiagnostic creates a new diagnostic information dump for a cluster
	CreateClusterDiagnostic(ctx context.Context, id int32, orgID int32) (*apigen.DiagnosticData, error)

	// UpdateClusterAutoBackupConfig updates the auto-backup configuration for a cluster
	UpdateClusterAutoBackupConfig(ctx context.Context, id int32, params apigen.AutoBackupConfig, orgID int32) error

	// UpdateClusterAutoDiagnosticConfig updates the auto-diagnostic configuration for a cluster
	UpdateClusterAutoDiagnosticConfig(ctx context.Context, id int32, params apigen.AutoDiagnosticConfig, orgID int32) error

	// GetClusterAutoBackupConfig gets the auto-backup configuration for a cluster
	GetClusterAutoBackupConfig(ctx context.Context, id int32, orgID int32) (*apigen.AutoBackupConfig, error)

	// GetClusterAutoDiagnosticConfig gets the auto-diagnostic configuration for a cluster
	GetClusterAutoDiagnosticConfig(ctx context.Context, id int32, orgID int32) (*apigen.AutoDiagnosticConfig, error)

	// GetMaterializedViewThroughput gets the throughput of materialized views
	GetMaterializedViewThroughput(ctx context.Context, clusterID int32) (prom_model.Matrix, error)

	// ImportMetricsStore creates a new metrics store
	ImportMetricsStore(context.Context, apigen.MetricsStoreImport, int32) (*apigen.MetricsStore, error)

	// DeleteMetricsStore deletes a metrics store
	DeleteMetricsStore(ctx context.Context, id int32, OrgID int32, force bool) error

	// GetMetricsStore gets a metrics store by ID
	GetMetricsStore(ctx context.Context, id int32, OrgID int32) (*apigen.MetricsStore, error)

	// UpdateMetricsStore updates a metrics store
	UpdateMetricsStore(ctx context.Context, id int32, req apigen.MetricsStoreImport, OrgID int32) (*apigen.MetricsStore, error)

	// ListClustersByMetricsStoreID lists all clusters by metrics store ID
	ListMetricsStores(ctx context.Context, OrgID int32) ([]*apigen.MetricsStore, error)

	// ListClustersByMetricsStoreID lists all clusters by metrics store ID
	ListClustersByMetricsStoreID(ctx context.Context, id int32) ([]*apigen.Cluster, error)
}

type Service struct {
	m                  model.ModelInterface
	auth               auth.AuthInterface
	sqlm               sql.SQLConnectionManegerInterface
	risectlm           meta.RisectlManagerInterface
	metahttp           http.MetaHttpManagerInterface
	metricsConnManager *metricsstore.MetricsManager
	taskRunner         taskgen.TaskRunner
	taskstore          taskcore.TaskStoreInterface
	anchorSvc          anchor_svc.ServiceInterface

	now                 func() time.Time
	generateHashAndSalt func(password string) (string, string, error)
}

func NewService(
	cfg *config.Config,
	m model.ModelInterface,
	auth auth.AuthInterface,
	sqlm sql.SQLConnectionManegerInterface,
	risectlm meta.RisectlManagerInterface,
	metricsConnManager *metricsstore.MetricsManager,
	metahttp http.MetaHttpManagerInterface,
	taskRunner taskgen.TaskRunner,
	taskstore taskcore.TaskStoreInterface,
	anchorSvc anchor_svc.ServiceInterface,
) (ServiceInterface, error) {
	s := &Service{
		m:                   m,
		now:                 time.Now,
		generateHashAndSalt: utils.GenerateHashAndSalt,
		auth:                auth,
		sqlm:                sqlm,
		risectlm:            risectlm,
		metahttp:            metahttp,
		metricsConnManager:  metricsConnManager,
		taskRunner:          taskRunner,
		taskstore:           taskstore,
		anchorSvc:           anchorSvc,
	}
	return s, nil
}
