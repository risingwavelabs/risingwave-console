package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/auth"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

type (
	TradeType   string
	TradeStatus string
	DdWorkEvent string
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrDatabaseNotFound    = errors.New("database not found")
	ErrClusterNotFound     = errors.New("cluster not found")
)

const (
	ExpireDuration             = 2 * time.Minute
	DefaultMaxRetries          = 3
	RefreshTokenExpireDuration = 14 * 24 * time.Hour
)

type ServiceInterface interface {
	// Create a new user and its default organization
	CreateNewUser(ctx context.Context, username, password string) error

	// SignIn authenticates a user and returns credentials
	SignIn(ctx context.Context, params apigen.SignInRequest) (*apigen.Credentials, error)

	// RefreshToken refreshes an authentication token using a refresh token
	RefreshToken(ctx context.Context, userID int32, refreshToken string) (*apigen.Credentials, error)

	// Cluster management
	CreateCluster(ctx context.Context, params apigen.ClusterCreate, orgID int32) (*apigen.Cluster, error)

	// GetCluster gets a cluster by its ID
	GetCluster(ctx context.Context, id int32) (*apigen.Cluster, error)

	// ListClusters lists all clusters in an organization
	ListClusters(ctx context.Context, orgID int32) ([]apigen.Cluster, error)

	// UpdateCluster updates a cluster
	UpdateCluster(ctx context.Context, id int32, params apigen.ClusterCreate) (*apigen.Cluster, error)

	// DeleteCluster deletes a cluster
	DeleteCluster(ctx context.Context, id int32) error

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
	TestDatabaseConnection(ctx context.Context, params apigen.TestConnectionPayload) (*apigen.TestConnectionResult, error)

	// QueryDatabase executes a query on a database
	QueryDatabase(ctx context.Context, id int32, params apigen.QueryRequest, orgID int32) (*apigen.QueryResponse, error)

	// GetDDLProgress gets the progress of DDL operations
	GetDDLProgress(ctx context.Context, id int32, orgID int32) ([]apigen.DDLProgress, error)

	// CancelDDLProgress cancels a DDL operation
	CancelDDLProgress(ctx context.Context, id int32, ddlID string, orgID int32) error
}

type Service struct {
	m    model.ModelInterface
	auth auth.AuthInterface
	sqlm sql.SQLConnectionManegerInterface

	now                 func() time.Time
	generateHashAndSalt func(password string) (string, string, error)
}

func NewService(cfg *config.Config, m model.ModelInterface, auth auth.AuthInterface, sqlm sql.SQLConnectionManegerInterface) ServiceInterface {
	return &Service{
		m:                   m,
		now:                 time.Now,
		generateHashAndSalt: utils.GenerateHashAndSalt,
		auth:                auth,
		sqlm:                sqlm,
	}
}

func (s *Service) SignIn(ctx context.Context, params apigen.SignInRequest) (*apigen.Credentials, error) {
	user, err := s.m.GetUserByName(ctx, params.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.Wrapf(ErrUserNotFound, "user %s not found", params.Name)
		}
		return nil, errors.Wrapf(err, "failed to get user by name")
	}
	input, err := utils.HashPassword(params.Password, user.PasswordSalt)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hash password")
	}
	if input != user.PasswordHash {
		return nil, ErrInvalidPassword
	}
	token, err := s.auth.CreateToken(user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}
	refreshToken, jwtToken, err := s.auth.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}
	err = s.m.UpsertRefreshToken(ctx, querier.UpsertRefreshTokenParams{
		UserID: user.ID,
		Token:  refreshToken,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upsert refresh token")
	}

	return &apigen.Credentials{
		AccessToken:  token,
		RefreshToken: jwtToken,
		TokenType:    apigen.Bearer,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, userID int32, refreshToken string) (*apigen.Credentials, error) {
	originalRefreshToken, err := s.m.GetRefreshToken(ctx, querier.GetRefreshTokenParams{
		UserID: userID,
		Token:  refreshToken,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get refresh token")
	}
	if originalRefreshToken.UpdatedAt.Add(RefreshTokenExpireDuration).Before(s.now()) {
		return nil, ErrRefreshTokenExpired
	}
	user, err := s.m.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get user by id: %d", userID)
	}

	newRefreshToken, jwtToken, err := s.auth.CreateRefreshToken(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate refresh token")
	}
	if err = s.m.UpsertRefreshToken(ctx, querier.UpsertRefreshTokenParams{
		UserID: userID,
		Token:  newRefreshToken,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to upsert refresh token")
	}
	accessToken, err := s.auth.CreateToken(user, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create token")
	}

	return &apigen.Credentials{
		AccessToken:  accessToken,
		RefreshToken: jwtToken,
		TokenType:    apigen.Bearer,
	}, nil
}

func (s *Service) CreateNewUser(ctx context.Context, username, password string) error {
	salt, hash, err := s.generateHashAndSalt(password)
	if err != nil {
		return errors.Wrapf(err, "failed to generate hash and salt")
	}
	return s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
		org, err := txm.CreateOrganization(ctx, fmt.Sprintf("%s's Org", username))
		if err != nil {
			return errors.Wrapf(err, "failed to create organization")
		}
		user, err := txm.CreateUser(ctx, querier.CreateUserParams{
			Name:           username,
			PasswordHash:   hash,
			PasswordSalt:   salt,
			OrganizationID: org.ID,
		})
		if err := txm.CreateOrganizationOwner(ctx, querier.CreateOrganizationOwnerParams{
			UserID:         user.ID,
			OrganizationID: org.ID,
		}); err != nil {
			return errors.Wrapf(err, "failed to create organization owner")
		}
		if err != nil {
			return errors.Wrapf(err, "failed to create user")
		}
		return nil
	})
}

func (s *Service) CreateCluster(ctx context.Context, params apigen.ClusterCreate, orgID int32) (*apigen.Cluster, error) {
	cluster, err := s.m.CreateCluster(ctx, querier.CreateClusterParams{
		OrganizationID: orgID,
		Name:           params.Name,
		Host:           params.Host,
		SqlPort:        int32(params.SqlPort),
		MetaPort:       int32(params.MetaPort),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster")
	}

	return &apigen.Cluster{
		ID:             cluster.ID,
		OrganizationID: cluster.OrganizationID,
		Name:           cluster.Name,
		Host:           cluster.Host,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) GetCluster(ctx context.Context, id int32) (*apigen.Cluster, error) {
	cluster, err := s.m.GetCluster(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrClusterNotFound
		}
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	return &apigen.Cluster{
		ID:             cluster.ID,
		OrganizationID: cluster.OrganizationID,
		Name:           cluster.Name,
		Host:           cluster.Host,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) ListClusters(ctx context.Context, orgID int32) ([]apigen.Cluster, error) {
	clusters, err := s.m.ListClusters(ctx, orgID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to list clusters")
	}

	result := make([]apigen.Cluster, len(clusters))
	for i, cluster := range clusters {
		result[i] = apigen.Cluster{
			ID:             cluster.ID,
			OrganizationID: cluster.OrganizationID,
			Name:           cluster.Name,
			Host:           cluster.Host,
			SqlPort:        cluster.SqlPort,
			MetaPort:       cluster.MetaPort,
			CreatedAt:      cluster.CreatedAt,
			UpdatedAt:      cluster.UpdatedAt,
		}
	}
	return result, nil
}

func (s *Service) UpdateCluster(ctx context.Context, id int32, params apigen.ClusterCreate) (*apigen.Cluster, error) {
	cluster, err := s.m.UpdateCluster(ctx, querier.UpdateClusterParams{
		ID:       id,
		Name:     params.Name,
		Host:     params.Host,
		SqlPort:  int32(params.SqlPort),
		MetaPort: int32(params.MetaPort),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrClusterNotFound
		}
		return nil, errors.Wrapf(err, "failed to update cluster")
	}

	return &apigen.Cluster{
		ID:             cluster.ID,
		OrganizationID: cluster.OrganizationID,
		Name:           cluster.Name,
		Host:           cluster.Host,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) DeleteCluster(ctx context.Context, id int32) error {
	err := s.m.DeleteCluster(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete cluster")
	}
	return nil
}

func (s *Service) CreateDatabase(ctx context.Context, params apigen.DatabaseConnectInfo, orgID int32) (*apigen.Database, error) {
	cluster, err := s.m.CreateDatabaseConnection(ctx, querier.CreateDatabaseConnectionParams{
		ClusterID:      params.ClusterID,
		Name:           params.Name,
		Username:       params.Username,
		Password:       params.Password,
		Database:       params.Database,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create database")
	}

	return &apigen.Database{
		ID:             cluster.ID,
		Name:           cluster.Name,
		ClusterID:      cluster.ClusterID,
		OrganizationID: cluster.OrganizationID,
		Username:       cluster.Username,
		Password:       cluster.Password,
		Database:       cluster.Database,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) GetDatabase(ctx context.Context, id int32, orgID int32) (*apigen.Database, error) {
	db, err := s.m.GetDatabaseConnection(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrDatabaseNotFound
		}
		return nil, errors.Wrapf(err, "failed to get database")
	}

	return &apigen.Database{
		ID:             db.ID,
		Name:           db.Name,
		ClusterID:      db.ClusterID,
		OrganizationID: db.OrganizationID,
		Username:       db.Username,
		Password:       db.Password,
		Database:       db.Database,
		CreatedAt:      db.CreatedAt,
		UpdatedAt:      db.UpdatedAt,
	}, nil
}

func (s *Service) ListDatabases(ctx context.Context, orgID int32) ([]apigen.Database, error) {
	dbs, err := s.m.ListDatabaseConnections(ctx, orgID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to list databases")
	}

	result := make([]apigen.Database, len(dbs))
	for i, db := range dbs {
		result[i] = apigen.Database{
			ID:             db.ID,
			Name:           db.Name,
			ClusterID:      db.ClusterID,
			OrganizationID: db.OrganizationID,
			Username:       db.Username,
			Password:       db.Password,
			Database:       db.Database,
			CreatedAt:      db.CreatedAt,
			UpdatedAt:      db.UpdatedAt,
		}
	}
	return result, nil
}

func (s *Service) UpdateDatabase(ctx context.Context, id int32, params apigen.DatabaseConnectInfo, orgID int32) (*apigen.Database, error) {
	db, err := s.m.UpdateDatabaseConnection(ctx, querier.UpdateDatabaseConnectionParams{
		ID:             id,
		ClusterID:      params.ClusterID,
		Name:           params.Name,
		Username:       params.Username,
		Password:       params.Password,
		Database:       params.Database,
		OrganizationID: orgID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrDatabaseNotFound
		}
		return nil, errors.Wrapf(err, "failed to update database")
	}

	return &apigen.Database{
		ID:             db.ID,
		Name:           db.Name,
		ClusterID:      db.ClusterID,
		Database:       db.Database,
		OrganizationID: db.OrganizationID,
		Username:       db.Username,
		Password:       db.Password,
		CreatedAt:      db.CreatedAt,
		UpdatedAt:      db.UpdatedAt,
	}, nil
}

func (s *Service) DeleteDatabase(ctx context.Context, id int32, orgID int32) error {
	err := s.m.DeleteDatabaseConnection(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete database")
	}
	return nil
}

func (s *Service) TestDatabaseConnection(ctx context.Context, params apigen.TestConnectionPayload) (*apigen.TestConnectionResult, error) {
	cluster, err := s.m.GetCluster(ctx, params.ClusterID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", params.Username, utils.UnwrapOrDefault(params.Password, ""), cluster.Host, cluster.SqlPort, params.Database)

	_, err = sql.Query(ctx, connStr, "SELECT 1")
	if err != nil {
		return &apigen.TestConnectionResult{
			Success: false,
			Result:  err.Error(),
		}, nil
	}

	return &apigen.TestConnectionResult{
		Success: true,
		Result:  "Connection successful",
	}, nil
}

func (s *Service) QueryDatabase(ctx context.Context, id int32, params apigen.QueryRequest, orgID int32) (*apigen.QueryResponse, error) {
	db, err := s.m.GetUserDatabaseByID(ctx, querier.GetUserDatabaseByIDParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrDatabaseNotFound
		}
		if errors.Is(err, sql.ErrQueryFailed) {
			return &apigen.QueryResponse{
				Error: utils.Ptr(err.Error()),
			}, nil
		}
		return nil, errors.Wrapf(err, "failed to get database connection")
	}

	conn, err := s.sqlm.GetConn(ctx, db.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get database connection")
	}

	result, err := conn.Query(ctx, params.Query)
	if err != nil {
		if errors.Is(err, sql.ErrQueryFailed) {
			return &apigen.QueryResponse{
				Error: utils.Ptr(err.Error()),
			}, nil
		}
		return nil, errors.Wrapf(err, "failed to query database")
	}

	columns := make([]apigen.Column, len(result.Columns))
	for i, column := range result.Columns {
		columns[i] = apigen.Column{
			Name: column.Name,
			Type: column.Type,
		}
	}

	return &apigen.QueryResponse{

		Columns: columns,
		Rows:    result.Rows,
	}, nil
}

func (s *Service) GetDDLProgress(ctx context.Context, id int32, orgID int32) ([]apigen.DDLProgress, error) {
	// TODO: Implement actual DDL progress tracking
	return []apigen.DDLProgress{
		{
			ID:            1,
			Statement:     "CREATE MATERIALIZED VIEW ...",
			Progress:      "50%",
			InitializedAt: s.now().Round(time.Hour),
		},
		{
			ID:            1,
			Statement:     "CREATE MATERIALIZED VIEW ...",
			Progress:      "100%",
			InitializedAt: s.now().Round(time.Hour),
		},
	}, nil
}

func (s *Service) CancelDDLProgress(ctx context.Context, id int32, ddlID string, orgID int32) error {
	// TODO: Implement actual DDL cancellation
	return nil
}
