package service

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	root "github.com/risingwavelabs/risingwave-console"
	"github.com/risingwavelabs/risingwave-console/pkg/config"
	"github.com/risingwavelabs/risingwave-console/pkg/utils"
	"github.com/risingwavelabs/risingwave-console/pkg/zcore/model"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/apigen"
	"github.com/risingwavelabs/risingwave-console/pkg/zgen/querier"
	"gopkg.in/yaml.v3"

	anchor_app "github.com/cloudcarver/anchor/pkg/app"
	anchor_svc "github.com/cloudcarver/anchor/pkg/service"
)

type InitService struct {
	m         model.ModelInterface
	anchorSvc anchor_svc.ServiceInterface
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

func NewInitService(m model.ModelInterface, anchor_svc anchor_svc.ServiceInterface) *InitService {
	return &InitService{
		m:         m,
		anchorSvc: anchor_svc,
	}
}

func (s *InitService) Init(ctx context.Context, cfg *config.Config, anchorApp *anchor_app.Application) error {
	// init the static web pages
	anchorApp.GetServer().GetApp().Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(root.StaticFiles),
		PathPrefix:   "web/out",
		NotFoundFile: "404.html",
		Index:        "index.html",
	}))

	anchorApp.GetServer().GetApp().Get("/config.js", func(c *fiber.Ctx) error {
		endpoint := fmt.Sprintf("http://%s:%d/api/v1", anchorApp.GetServer().GetHost(), anchorApp.GetServer().GetPort())
		c.Set("Content-Type", "application/javascript")
		return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("window.APP_ENDPOINT = '%s';", endpoint))
	})

	// init hooks
	anchorApp.GetHooks().RegisterOnOrgCreatedWithTx(func(ctx context.Context, tx pgx.Tx, orgID int32) error {
		txm := s.m.SpawnWithTx(tx)
		if err := txm.CreateOrgSettings(ctx, querier.CreateOrgSettingsParams{
			OrgID:    orgID,
			Timezone: "UTC",
		}); err != nil {
			return errors.Wrapf(err, "failed to create org settings")
		}
		return nil
	})

	// remove the root user if it is not set in the config
	if cfg.Root == nil {
		if err := s.anchorSvc.DeleteUserByName(ctx, "root"); err != nil {
			return errors.Wrapf(err, "failed to delete root user")
		}
		return nil
	}

	// create the root user
	rootPwd := cfg.Root.Password
	if rootPwd == "" {
		rootPwd = "123456"
	}
	orgID, err := s.anchorSvc.CreateNewUser(ctx, "root", rootPwd)
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

		ms, err := s.m.ListMetricsStoresByOrgID(ctx, orgID)
		if err != nil {
			return errors.Wrapf(err, "failed to list metrics stores")
		}
		for _, m := range ms {
			metricsStoreNameToID[m.Name] = m.ID
		}

		for _, metricsStore := range cfg.MetricsStores {
			if _, ok := metricsStoreNameToID[metricsStore.Name]; ok {
				_, err := s.m.UpdateMetricsStore(ctx, querier.UpdateMetricsStoreParams{
					ID:            metricsStoreNameToID[metricsStore.Name],
					OrgID:         orgID,
					Name:          metricsStore.Name,
					Spec:          metricsStore.Spec,
					DefaultLabels: metricsStore.DefaultLabels,
				})
				if err != nil {
					return errors.Wrapf(err, "failed to update metrics store: %s", metricsStore.Name)
				}
				continue
			}
			ms, err := s.m.CreateMetricsStore(ctx, querier.CreateMetricsStoreParams{
				OrgID:         orgID,
				Name:          metricsStore.Name,
				Spec:          metricsStore.Spec,
				DefaultLabels: metricsStore.DefaultLabels,
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
				OrgID:          orgID,
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
				Name:      database.Name,
				OrgID:     orgID,
				ClusterID: clusterID,
				Username:  database.Username,
				Password:  database.Password,
				Database:  database.Database,
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
