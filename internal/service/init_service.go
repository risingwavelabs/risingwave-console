package service

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"gopkg.in/yaml.v3"
)

type InitService struct {
	m model.ModelInterface
	s ServiceInterface
}

type ClusterConnections struct {
	Host     string `yaml:"host" validate:"required,hostname_rfc1123"`
	SqlPort  int32  `yaml:"sqlPort" validate:"required,min=1,max=65535"`
	MetaPort int32  `yaml:"metaPort" validate:"required,min=1,max=65535"`
	HttpPort int32  `yaml:"httpPort" validate:"required,min=1,max=65535"`
}

type Cluster struct {
	Name         string              `yaml:"name" validate:"required"`
	Version      string              `yaml:"version" validate:"required"`
	Connections  *ClusterConnections `yaml:"connections" validate:"required"`
	MetricsStore string              `yaml:"metricsStore" validate:"required"`
}

type Database struct {
	Name     string  `yaml:"name" validate:"required"`
	Cluster  string  `yaml:"cluster" validate:"required"`
	Username string  `yaml:"username" validate:"required"`
	Password *string `yaml:"password"`
	Database string  `yaml:"database" validate:"required"`
}

type Query struct {
	Name      *string `yaml:"name"`
	Statement string  `yaml:"statement" validate:"required"`
}

type InitConfig struct {
	Clusters      []Cluster             `yaml:"clusters"`
	Databases     []Database            `yaml:"databases"`
	Queries       []Query               `yaml:"queries"`
	MetricsStores []apigen.MetricsStore `yaml:"metricsStores"`
}

func NewInitService(m model.ModelInterface, s ServiceInterface) *InitService {
	return &InitService{
		m: m,
		s: s,
	}
}

func (s *InitService) Init(ctx context.Context, cfg *config.Config) error {
	// remove the root user if it is not set in the config
	if cfg.Root == nil {
		if err := s.m.DeleteUserByName(ctx, "root"); err != nil {
			return errors.Wrapf(err, "failed to delete root user")
		}
		return nil
	}

	// create the root user
	rootPwd := cfg.Root.Password
	if rootPwd == "" {
		rootPwd = "123456"
	}
	orgID, err := s.s.CreateNewUser(ctx, "root", rootPwd)
	if err != nil {
		return errors.Wrapf(err, "failed to create root user")
	}

	// init the database
	if cfg.Init != "" {
		raw, err := os.ReadFile(cfg.Init)
		if err != nil {
			return errors.Wrapf(err, "failed to read init file: %s", cfg.Init)
		}
		var initCfg InitConfig
		if err := yaml.Unmarshal(raw, &initCfg); err != nil {
			return errors.Wrapf(err, "failed to unmarshal init file: %s", cfg.Init)
		}
		if err := validator.New().Struct(initCfg); err != nil {
			return errors.Wrapf(err, "failed to validate init file: %s", cfg.Init)
		}
		if err := s.initDatabase(ctx, &initCfg, orgID); err != nil {
			return errors.Wrapf(err, "failed to init database")
		}
	}

	return nil
}

func (s *InitService) initDatabase(ctx context.Context, cfg *InitConfig, orgID int32) error {
	if err := s.m.RunTransaction(ctx, func(txm model.ModelInterface) error {
		clusterNameToID := make(map[string]int32)
		metricsStoreNameToID := make(map[string]int32)

		for _, metricsStore := range cfg.MetricsStores {
			ms, err := s.m.CreateMetricsStore(ctx, querier.CreateMetricsStoreParams{
				OrganizationID: orgID,
				Name:           metricsStore.Name,
				Spec:           metricsStore.Spec,
				DefaultLabels:  metricsStore.DefaultLabels,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to create metrics store: %s", metricsStore.Name)
			}
			metricsStoreNameToID[metricsStore.Name] = ms.ID
		}

		for _, cluster := range cfg.Clusters {
			if cluster.Connections == nil {
				return errors.New("cluster connections is required")
			}
			msid, ok := metricsStoreNameToID[cluster.MetricsStore]
			cluster, err := s.m.InitCluster(ctx, querier.InitClusterParams{
				OrganizationID: orgID,
				Name:           cluster.Name,
				Host:           cluster.Connections.Host,
				SqlPort:        cluster.Connections.SqlPort,
				MetaPort:       cluster.Connections.MetaPort,
				HttpPort:       cluster.Connections.HttpPort,
				Version:        cluster.Version,
				MetricsStoreID: utils.IfElse(ok, &msid, nil),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to create cluster: %s", cluster.Name)
			}
			clusterNameToID[cluster.Name] = cluster.ID
		}

		for _, database := range cfg.Databases {
			if _, ok := clusterNameToID[database.Cluster]; !ok {
				return errors.Errorf("cluster %s not found", database.Cluster)
			}
			clusterID := clusterNameToID[database.Cluster]
			if _, err := s.m.InitDatabaseConnection(ctx, querier.InitDatabaseConnectionParams{
				Name:           database.Name,
				OrganizationID: orgID,
				ClusterID:      clusterID,
				Username:       database.Username,
				Password:       database.Password,
				Database:       database.Database,
			}); err != nil {
				return errors.Wrapf(err, "failed to init cluster: %s", database.Cluster)
			}
		}
		return nil
	}); err != nil {
		return errors.Wrapf(err, "failed to run transaction")
	}
	return nil
}
