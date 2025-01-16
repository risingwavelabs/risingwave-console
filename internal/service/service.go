package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/auth"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/conn/risectl"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"golang.org/x/mod/semver"
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
	ListClusters(ctx context.Context, orgID int32) ([]apigen.Cluster, error)

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

	ListClusterVersions(ctx context.Context) ([]string, error)

	// ClusterSnapshot management
	CreateClusterSnapshot(ctx context.Context, id int32, name string, orgID int32) (*apigen.Snapshot, error)

	// ListClusterSnapshots lists all snapshots of a cluster
	ListClusterSnapshots(ctx context.Context, id int32, orgID int32) ([]apigen.Snapshot, error)

	DeleteClusterSnapshot(ctx context.Context, id int32, snapshotID int64, orgID int32) error

	TestClusterConnection(ctx context.Context, params apigen.TestClusterConnectionPayload, orgID int32) (*apigen.TestClusterConnectionResult, error)
}

type Service struct {
	m        model.ModelInterface
	auth     auth.AuthInterface
	sqlm     sql.SQLConnectionManegerInterface
	risectlm risectl.RisectlManagerInterface

	now                 func() time.Time
	generateHashAndSalt func(password string) (string, string, error)
}

func NewService(cfg *config.Config, m model.ModelInterface, auth auth.AuthInterface, sqlm sql.SQLConnectionManegerInterface, risectlm risectl.RisectlManagerInterface) ServiceInterface {
	return &Service{
		m:                   m,
		now:                 time.Now,
		generateHashAndSalt: utils.GenerateHashAndSalt,
		auth:                auth,
		sqlm:                sqlm,
		risectlm:            risectlm,
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

func (s *Service) CreateNewUser(ctx context.Context, username, password string) (int32, error) {
	salt, hash, err := s.generateHashAndSalt(password)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to generate hash and salt")
	}
	var orgID int32
	if err := s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
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
		orgID = org.ID
		return nil
	}); err != nil {
		return 0, err
	}
	return orgID, nil
}

func (s *Service) CreateCluster(ctx context.Context, params apigen.ClusterCreate, orgID int32) (*apigen.Cluster, error) {
	cluster, err := s.m.CreateCluster(ctx, querier.CreateClusterParams{
		OrganizationID: orgID,
		Name:           params.Name,
		Host:           params.Host,
		SqlPort:        params.SqlPort,
		MetaPort:       params.MetaPort,
		HttpPort:       params.HttpPort,
		Version:        params.Version,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster")
	}

	return &apigen.Cluster{
		ID:             cluster.ID,
		OrganizationID: cluster.OrganizationID,
		Name:           cluster.Name,
		Host:           cluster.Host,
		Version:        cluster.Version,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		HttpPort:       cluster.HttpPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) GetCluster(ctx context.Context, id int32, orgID int32) (*apigen.Cluster, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
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
		Version:        cluster.Version,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		HttpPort:       cluster.HttpPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) ListClusters(ctx context.Context, orgID int32) ([]apigen.Cluster, error) {
	clusters, err := s.m.ListOrgClusters(ctx, orgID)
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
			Version:        cluster.Version,
			SqlPort:        cluster.SqlPort,
			MetaPort:       cluster.MetaPort,
			HttpPort:       cluster.HttpPort,
			CreatedAt:      cluster.CreatedAt,
			UpdatedAt:      cluster.UpdatedAt,
		}
	}
	return result, nil
}

func (s *Service) UpdateCluster(ctx context.Context, id int32, params apigen.ClusterCreate, orgID int32) (*apigen.Cluster, error) {
	cluster, err := s.m.UpdateOrgCluster(ctx, querier.UpdateOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
		Name:           params.Name,
		Host:           params.Host,
		Version:        params.Version,
		SqlPort:        int32(params.SqlPort),
		MetaPort:       int32(params.MetaPort),
		HttpPort:       int32(params.HttpPort),
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
		Version:        cluster.Version,
		SqlPort:        cluster.SqlPort,
		MetaPort:       cluster.MetaPort,
		HttpPort:       cluster.HttpPort,
		CreatedAt:      cluster.CreatedAt,
		UpdatedAt:      cluster.UpdatedAt,
	}, nil
}

func (s *Service) DeleteCluster(ctx context.Context, id int32, cascade bool, orgID int32) error {
	if cascade {
		return s.deleteClusterCacasde(ctx, id, orgID)
	}
	return s.deleteClusterNonCacasde(ctx, id, orgID)
}

func (s *Service) deleteClusterCacasde(ctx context.Context, id int32, orgID int32) error {
	return s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
		if err := txm.DeleteAllOrgDatabaseConnectionsByClusterID(ctx, querier.DeleteAllOrgDatabaseConnectionsByClusterIDParams{
			ClusterID:      id,
			OrganizationID: orgID,
		}); err != nil {
			return errors.Wrapf(err, "failed to delete associated database connections")
		}

		if err := txm.DeleteOrgCluster(ctx, querier.DeleteOrgClusterParams{
			ID:             id,
			OrganizationID: orgID,
		}); err != nil {
			return errors.Wrapf(err, "failed to delete cluster")
		}
		return nil
	})
}

func (s *Service) deleteClusterNonCacasde(ctx context.Context, id int32, orgID int32) error {
	dbConnections, err := s.m.GetAllOrgDatabseConnectionsByClusterID(ctx, querier.GetAllOrgDatabseConnectionsByClusterIDParams{
		ClusterID:      id,
		OrganizationID: orgID,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to get database connections")
	}
	names := make([]string, len(dbConnections))
	for i, db := range dbConnections {
		names[i] = db.Name
	}

	if len(dbConnections) > 0 {
		return errors.Wrapf(ErrClusterHasDatabaseConnections, "cluster has %d database connections: %s", len(dbConnections), strings.Join(names, ", "))
	}

	err = s.m.DeleteOrgCluster(ctx, querier.DeleteOrgClusterParams{
		ID:             id,
		OrganizationID: orgID,
	})
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

func (s *Service) getConnStr(ctx context.Context, db *querier.DatabaseConnection) (string, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             db.ClusterID,
		OrganizationID: db.OrganizationID,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to get cluster")
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", db.Username, utils.UnwrapOrDefault(db.Password, ""), cluster.Host, cluster.SqlPort, db.Database), nil
}

const getRelationsSQL = `SELECT 
    rw_relations.id            AS relation_id,
    rw_schemas.name            AS schema, 
    rw_relations.name          AS relation_name, 
    rw_relations.relation_type AS relation_type, 
    rw_columns.name            AS column_name,
    rw_columns.data_type       AS column_type,
    rw_columns.is_primary_key  AS is_primary_key,
	rw_columns.is_hidden       AS is_hidden
FROM rw_columns
JOIN rw_relations ON rw_relations.id = rw_columns.relation_id
JOIN rw_schemas   ON rw_schemas.id = rw_relations.schema_id
`

const getRwDependSQL = `SELECT * FROM rw_depend`

func (s *Service) getDb(ctx context.Context, id int32, orgID int32) (*querier.DatabaseConnection, error) {
	db, err := s.m.GetOrgDatabaseByID(ctx, querier.GetOrgDatabaseByIDParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrDatabaseNotFound
		}
		return nil, errors.Wrapf(err, "failed to get database")
	}
	return db, nil
}

func (s *Service) GetDatabase(ctx context.Context, id int32, orgID int32) (*apigen.Database, error) {
	db, err := s.getDb(ctx, id, orgID)
	if err != nil {
		return nil, err
	}

	connStr, err := s.getConnStr(ctx, db)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get connection string")
	}

	result, err := sql.Query(ctx, connStr, getRelationsSQL, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query database")
	}

	data := make(map[string]map[string]apigen.Relation)

	idToDepends := make(map[int32][]int32)
	depend, err := sql.Query(ctx, connStr, getRwDependSQL, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query database")
	}
	for _, row := range depend.Rows {
		objid := row["objid"].(int32)
		refobjid := row["refobjid"].(int32)
		if _, ok := idToDepends[objid]; !ok {
			idToDepends[objid] = []int32{}
		}
		idToDepends[objid] = append(idToDepends[objid], refobjid)
	}

	for _, row := range result.Rows {
		schemaName := row["schema"].(string)
		if _, ok := data[schemaName]; !ok {
			data[schemaName] = make(map[string]apigen.Relation)
		}
		schema := data[schemaName]

		relationName := row["relation_name"].(string)
		if _, ok := schema[relationName]; !ok {
			schema[relationName] = apigen.Relation{
				ID:           row["relation_id"].(int32),
				Name:         row["relation_name"].(string),
				Type:         apigen.RelationType(row["relation_type"].(string)),
				Columns:      []apigen.Column{},
				Dependencies: idToDepends[row["relation_id"].(int32)],
			}
		}
		relation := schema[relationName]
		relation.Columns = append(relation.Columns, apigen.Column{
			Name:         row["column_name"].(string),
			Type:         row["column_type"].(string),
			IsPrimaryKey: row["is_primary_key"].(bool),
			IsHidden:     row["is_hidden"].(bool),
		})
		schema[relationName] = relation
		data[schemaName] = schema
	}

	schemas := []apigen.Schema{}
	for schemaName, schema := range data {
		s := apigen.Schema{
			Name:      schemaName,
			Relations: []apigen.Relation{},
		}
		for _, relation := range schema {
			s.Relations = append(s.Relations, relation)
		}
		schemas = append(schemas, s)
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
		Schemas:        &schemas,
	}, nil
}

func (s *Service) ListDatabases(ctx context.Context, orgID int32) ([]apigen.Database, error) {
	dbs, err := s.m.ListOrgDatabaseConnections(ctx, orgID)
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
	db, err := s.m.UpdateOrgDatabaseConnection(ctx, querier.UpdateOrgDatabaseConnectionParams{
		ID:               id,
		ClusterID:        params.ClusterID,
		Name:             params.Name,
		Username:         params.Username,
		Password:         params.Password,
		Database:         params.Database,
		OrganizationID:   orgID,
		OrganizationID_2: orgID,
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
	err := s.m.DeleteOrgDatabaseConnection(ctx, querier.DeleteOrgDatabaseConnectionParams{
		ID:             id,
		OrganizationID: orgID,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to delete database")
	}
	return nil
}

func (s *Service) TestDatabaseConnection(ctx context.Context, params apigen.TestDatabaseConnectionPayload, orgID int32) (*apigen.TestDatabaseConnectionResult, error) {
	cluster, err := s.m.GetOrgCluster(ctx, querier.GetOrgClusterParams{
		ID:             params.ClusterID,
		OrganizationID: orgID,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", params.Username, utils.UnwrapOrDefault(params.Password, ""), cluster.Host, cluster.SqlPort, params.Database)

	_, err = sql.Query(ctx, connStr, "SELECT 1", false)
	if err != nil {
		return &apigen.TestDatabaseConnectionResult{
			Success: false,
			Result:  err.Error(),
		}, nil
	}

	return &apigen.TestDatabaseConnectionResult{
		Success: true,
		Result:  "Connection successful",
	}, nil
}

func (s *Service) QueryDatabase(ctx context.Context, id int32, params apigen.QueryRequest, orgID int32, backgroundDDL bool) (*apigen.QueryResponse, error) {
	db, err := s.m.GetOrgDatabaseByID(ctx, querier.GetOrgDatabaseByIDParams{
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

	result, err := conn.Query(ctx, params.Query, backgroundDDL)
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

const getDDLProgressSQL = `SELECT * FROM rw_ddl_progress ORDER BY initialized_at DESC`

func (s *Service) GetDDLProgress(ctx context.Context, id int32, orgID int32) ([]apigen.DDLProgress, error) {
	conn, err := s.sqlm.GetConn(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get database connection")
	}

	result, err := conn.Query(ctx, getDDLProgressSQL, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get DDL progress")
	}

	progress := []apigen.DDLProgress{}
	for _, row := range result.Rows {
		item := apigen.DDLProgress{
			ID:        row["ddl_id"].(int64),
			Statement: row["ddl_statement"].(string),
			Progress:  row["progress"].(string),
		}
		if t, ok := row["initialized_at"].(time.Time); ok {
			item.InitializedAt = &t
		}
		progress = append(progress, item)
	}
	return progress, nil
}

func (s *Service) CancelDDLProgress(ctx context.Context, id int32, ddlID int64, orgID int32) error {
	conn, err := s.sqlm.GetConn(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to get database connection")
	}

	_, err = conn.Query(ctx, fmt.Sprintf("CANCEL JOB %d", ddlID), false)
	if err != nil {
		return errors.Wrapf(err, "failed to cancel DDL progress")
	}

	return nil
}

func (s *Service) ListClusterVersions(ctx context.Context) ([]string, error) {
	versions, err := s.risectlm.ListVersions(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list cluster versions")
	}
	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare(versions[i], versions[j]) > 0
	})
	return versions, nil
}

func (s *Service) getRisectlConn(ctx context.Context, id int32) (risectl.RisectlConn, error) {
	cluster, err := s.m.GetClusterByID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cluster")
	}

	return s.risectlm.NewConn(ctx, cluster.Version, fmt.Sprintf("http://%s:%d", cluster.Host, cluster.MetaPort))
}

func (s *Service) CreateClusterSnapshot(ctx context.Context, id int32, name string, orgID int32) (*apigen.Snapshot, error) {
	conn, err := s.getRisectlConn(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get risectl connection")
	}

	snapshotID, err := conn.MetaBackup(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create snapshot")
	}

	if err := s.m.CreateClusterSnapshot(ctx, querier.CreateClusterSnapshotParams{
		ClusterID:  id,
		SnapshotID: snapshotID,
		Name:       name,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to create snapshot")
	}

	return &apigen.Snapshot{
		ID:        snapshotID,
		Name:      name,
		ClusterID: id,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Service) ListClusterSnapshots(ctx context.Context, id int32, orgID int32) ([]apigen.Snapshot, error) {
	snapshots, err := s.m.ListClusterSnapshots(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list cluster snapshots")
	}

	result := make([]apigen.Snapshot, len(snapshots))
	for i, snapshot := range snapshots {
		result[i] = apigen.Snapshot{
			ID:        snapshot.SnapshotID,
			Name:      snapshot.Name,
			ClusterID: snapshot.ClusterID,
			CreatedAt: snapshot.CreatedAt,
		}
	}
	return result, nil
}

func (s *Service) DeleteClusterSnapshot(ctx context.Context, id int32, snapshotID int64, orgID int32) error {
	conn, err := s.getRisectlConn(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to get risectl connection")
	}

	if err := conn.DeleteSnapshot(ctx, snapshotID); err != nil {
		return errors.Wrapf(err, "failed to delete snapshot")
	}

	return nil
}

func (s *Service) TestClusterConnection(ctx context.Context, params apigen.TestClusterConnectionPayload, orgID int32) (*apigen.TestClusterConnectionResult, error) {
	errMsg := ""
	if err := utils.TestTCPConnection(ctx, params.Host, params.MetaPort, 5*time.Second); err != nil {
		errMsg += fmt.Sprintf("Failed to connect to meta port: %s\n", err.Error())
	}
	if err := utils.TestTCPConnection(ctx, params.Host, params.SqlPort, 5*time.Second); err != nil {
		errMsg += fmt.Sprintf("Failed to connect to sql port: %s\n", err.Error())
	}
	if err := utils.TestTCPConnection(ctx, params.Host, params.HttpPort, 5*time.Second); err != nil {
		errMsg += fmt.Sprintf("Failed to connect to http port: %s\n", err.Error())
	}

	return &apigen.TestClusterConnectionResult{
		Success: errMsg == "",
		Result:  utils.IfElse(errMsg == "", "Connection successful", errMsg),
	}, nil
}
