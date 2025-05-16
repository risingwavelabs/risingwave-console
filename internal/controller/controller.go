package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloudcarver/anchor/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/risingwavelabs/wavekit/internal/caveats"
	"github.com/risingwavelabs/wavekit/internal/conn/metricsstore"
	"github.com/risingwavelabs/wavekit/internal/service"
	"github.com/risingwavelabs/wavekit/internal/utils"
	"github.com/risingwavelabs/wavekit/internal/zgen/apigen"
)

type Controller struct {
	svc  service.ServiceInterface
	auth auth.AuthInterface
}

func NewSeverInterface(s service.ServiceInterface, auth auth.AuthInterface) apigen.ServerInterface {
	return &Controller{
		svc:  s,
		auth: auth,
	}
}

func getOrgID(c *fiber.Ctx) (int32, error) {
	orgID, ok := c.Locals(caveats.ContextKeyOrgID).(int32)
	if !ok {
		return 0, errors.New("missing orgID in request context")
	}
	return orgID, nil
}

func (controller *Controller) ImportCluster(c *fiber.Ctx) error {
	var params apigen.ClusterImport
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	cluster, err := controller.svc.ImportCluster(c.Context(), params, orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(cluster)
}

func (controller *Controller) DeleteCluster(c *fiber.Ctx, id int32, params apigen.DeleteClusterParams) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	err = controller.svc.DeleteCluster(c.Context(), id, utils.UnwrapOrDefault(params.Cascade, false), orgID)
	if err != nil {
		if errors.Is(err, service.ErrClusterHasDatabaseConnections) {
			return c.Status(fiber.StatusConflict).SendString(err.Error())
		}
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (controller *Controller) GetCluster(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	cluster, err := controller.svc.GetCluster(c.Context(), id, orgID)
	if err != nil {
		if errors.Is(err, service.ErrClusterNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(cluster)
}

func (controller *Controller) UpdateCluster(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	var params apigen.ClusterImport
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	cluster, err := controller.svc.UpdateCluster(c.Context(), id, params, orgID)
	if err != nil {
		if errors.Is(err, service.ErrClusterNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(cluster)
}

func (controller *Controller) ListClusters(c *fiber.Ctx) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	clusters, err := controller.svc.ListClusters(c.Context(), orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(clusters)
}

func (controller *Controller) ImportDatabase(c *fiber.Ctx) error {
	var params apigen.DatabaseConnectInfo
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	database, err := controller.svc.ImportDatabase(c.Context(), params, orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(database)
}

func (controller *Controller) DeleteDatabase(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	err = controller.svc.DeleteDatabase(c.Context(), id, orgID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (controller *Controller) GetDatabase(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	database, err := controller.svc.GetDatabase(c.Context(), id, orgID)
	if err != nil {
		if errors.Is(err, service.ErrDatabaseNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(database)
}

func (controller *Controller) UpdateDatabase(c *fiber.Ctx, id int32) error {
	var params apigen.DatabaseConnectInfo
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	database, err := controller.svc.UpdateDatabase(c.Context(), id, params, orgID)
	if err != nil {
		if errors.Is(err, service.ErrDatabaseNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(database)
}

func (controller *Controller) ListDatabases(c *fiber.Ctx) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	databases, err := controller.svc.ListDatabases(c.Context(), orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(databases)
}

func (controller *Controller) GetDDLProgress(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	progress, err := controller.svc.GetDDLProgress(c.Context(), id, orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(progress)
}

func (controller *Controller) CancelDDLProgress(c *fiber.Ctx, id int32, ddlID int64) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	err = controller.svc.CancelDDLProgress(c.Context(), id, ddlID, orgID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (controller *Controller) TestDatabaseConnection(c *fiber.Ctx) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	var params apigen.TestDatabaseConnectionPayload
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result, err := controller.svc.TestDatabaseConnection(c.Context(), params, orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (controller *Controller) QueryDatabase(c *fiber.Ctx, id int32) error {
	var params apigen.QueryRequest
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	result, err := controller.svc.QueryDatabase(c.Context(), id, params, orgID, utils.UnwrapOrDefault(params.BackgroundDDL, false))
	if err != nil {
		if errors.Is(err, service.ErrDatabaseNotFound) {
			return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("database %d not found", id))
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (controller *Controller) CreateClusterSnapshot(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	var params apigen.CreateClusterSnapshotJSONRequestBody
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	snapshot, err := controller.svc.CreateClusterSnapshot(c.Context(), id, params.Name, orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(snapshot)
}

func (controller *Controller) DeleteClusterSnapshot(c *fiber.Ctx, id int32, snapshotId int64) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	err = controller.svc.DeleteClusterSnapshot(c.Context(), id, snapshotId, orgID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (controller *Controller) RestoreClusterSnapshot(c *fiber.Ctx, id int32, snapshotId int64) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) ListClusterSnapshots(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	snapshots, err := controller.svc.ListClusterSnapshots(c.Context(), id, orgID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(snapshots)
}

func (controller *Controller) ListClusterDiagnostics(c *fiber.Ctx, id int32, params apigen.ListClusterDiagnosticsParams) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	diagnostics, err := controller.svc.ListClusterDiagnostics(c.Context(), id, orgID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(diagnostics)
}

func (controller *Controller) GetClusterAutoBackupConfig(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	config, err := controller.svc.GetClusterAutoBackupConfig(c.Context(), id, orgID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(config)
}

func (controller *Controller) UpdateClusterAutoBackupConfig(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	var params apigen.AutoBackupConfig
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = controller.svc.UpdateClusterAutoBackupConfig(c.Context(), id, params, orgID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (controller *Controller) GetClusterAutoDiagnosticConfig(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	config, err := controller.svc.GetClusterAutoDiagnosticConfig(c.Context(), id, orgID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(config)
}

func (controller *Controller) UpdateClusterAutoDiagnosticConfig(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	var params apigen.AutoDiagnosticConfig
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = controller.svc.UpdateClusterAutoDiagnosticConfig(c.Context(), id, params, orgID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (controller *Controller) ListClusterVersions(c *fiber.Ctx) error {
	versions, err := controller.svc.ListClusterVersions(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(versions)
}

func (controller *Controller) TestClusterConnection(c *fiber.Ctx) error {
	var params apigen.TestClusterConnectionPayload
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	conn, err := controller.svc.TestClusterConnection(c.Context(), params, orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(conn)
}

func (controller *Controller) RunRisectlCommand(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	var params apigen.RisectlCommand
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result, err := controller.svc.RunRisectlCommand(c.Context(), id, params, orgID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (controller *Controller) CreateClusterDiagnostic(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	diagnostic, err := controller.svc.CreateClusterDiagnostic(c.Context(), id, orgID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(diagnostic)
}

func (controller *Controller) GetClusterDiagnostic(c *fiber.Ctx, id int32, diagnosticId int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	diagnostic, err := controller.svc.GetClusterDiagnostic(c.Context(), id, diagnosticId, orgID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(diagnostic)
}

func (controller *Controller) GetMaterializedViewThroughput(c *fiber.Ctx, clusterID int32) error {
	throughput, err := controller.svc.GetMaterializedViewThroughput(c.Context(), clusterID)
	if err != nil {
		if errors.Is(err, metricsstore.ErrMetricsStoreNotSupported) {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(throughput)
}

func (controller *Controller) ImportMetricsStore(c *fiber.Ctx) error {
	var req apigen.MetricsStoreImport
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	ms, err := controller.svc.ImportMetricsStore(c.Context(), req, orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(ms)
}

func (controller *Controller) DeleteMetricsStore(c *fiber.Ctx, id int32, params apigen.DeleteMetricsStoreParams) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	if !params.Force {
		clusters, err := controller.svc.ListClustersByMetricsStoreID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		var names []string
		for _, cluster := range clusters {
			names = append(names, cluster.Name)
		}
		if len(names) > 0 {
			return c.Status(fiber.StatusConflict).SendString(fmt.Sprintf("Metrics store is in use by clusters: %s", strings.Join(names, ", ")))
		}
	}

	if err := controller.svc.DeleteMetricsStore(c.Context(), id, orgID, params.Force); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (controller *Controller) GetMetricsStore(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	ms, err := controller.svc.GetMetricsStore(c.Context(), id, orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(ms)
}

func (controller *Controller) ListMetricsStores(c *fiber.Ctx) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	msList, err := controller.svc.ListMetricsStores(c.Context(), orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(msList)
}

func (controller *Controller) UpdateMetricsStore(c *fiber.Ctx, id int32) error {
	orgID, err := getOrgID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("missing orgID in request context")
	}

	var req apigen.MetricsStoreImport
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	ms, err := controller.svc.UpdateMetricsStore(c.Context(), id, req, orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(ms)
}

func (controller *Controller) ListTasks(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (controller *Controller) ListEvents(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (controller *Controller) CreateCluster(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
