package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/auth"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/conn/meta"
	"github.com/risingwavelabs/wavekit/internal/conn/prom"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/utils"

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
	// Create a new user and its default organization
	CreateNewUser(ctx context.Context, username, password string) (int32, error)

	// SignIn authenticates a user and returns credentials
	SignIn(ctx context.Context, params apigen.SignInRequest) (*apigen.Credentials, error)

	// RefreshToken refreshes an authentication token using a refresh token
	RefreshToken(ctx context.Context, userID int32, refreshToken string) (*apigen.Credentials, error)

	// Cluster management
	CreateCluster(ctx context.Context, params apigen.ClusterCreate, orgID int32) (*apigen.Cluster, error)

	// GetCluster gets a cluster by its ID
	GetCluster(ctx context.Context, id int32, orgID int32) (*apigen.Cluster, error)

	// ListClusters lists all clusters in an organization
	ListClusters(ctx context.Context, orgID int32) ([]*apigen.Cluster, error)

	// UpdateCluster updates a cluster
	UpdateCluster(ctx context.Context, id int32, params apigen.ClusterCreate, orgID int32) (*apigen.Cluster, error)

	// DeleteCluster deletes a cluster
	DeleteCluster(ctx context.Context, id int32, cascade bool, orgID int32) error

	// Database management
	CreateDatabase(ctx context.Context, params apigen.DatabaseConnectInfo, orgID int32) (*apigen.Database, error)

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
}

type Service struct {
	m        model.ModelInterface
	auth     auth.AuthInterface
	sqlm     sql.SQLConnectionManegerInterface
	risectlm meta.RisectlManagerInterface
	promm    prom.PromManagerInterface

	now                 func() time.Time
	generateHashAndSalt func(password string) (string, string, error)
}

func NewService(
	cfg *config.Config,
	m model.ModelInterface,
	auth auth.AuthInterface,
	sqlm sql.SQLConnectionManegerInterface,
	risectlm meta.RisectlManagerInterface,
	promm prom.PromManagerInterface,
) ServiceInterface {
	return &Service{
		m:                   m,
		now:                 time.Now,
		generateHashAndSalt: utils.GenerateHashAndSalt,
		auth:                auth,
		sqlm:                sqlm,
		risectlm:            risectlm,
		promm:               promm,
	}
}
