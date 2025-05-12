package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/conn/sql"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/zgen/apigen"
	"github.com/risingwavelabs/wavekit/internal/zgen/querier"
)

func (s *Service) ImportDatabase(ctx context.Context, params apigen.DatabaseConnectInfo, orgID int32) (*apigen.Database, error) {
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
