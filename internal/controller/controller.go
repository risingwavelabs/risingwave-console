package controller

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/auth"
	"github.com/risingwavelabs/wavekit/internal/conn/prom"
	"github.com/risingwavelabs/wavekit/internal/service"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

type Controller struct {
	svc  service.ServiceInterface
	auth auth.AuthInterface
}

var _ apigen.ServerInterface = &Controller{}

func NewController(
	s service.ServiceInterface,
	auth auth.AuthInterface,
) *Controller {
	return &Controller{
		svc:  s,
		auth: auth,
	}
}

func (controller *Controller) SignIn(c *fiber.Ctx) error {
	var params apigen.SignInRequest
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	credentials, err := controller.svc.SignIn(c.Context(), params)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPassword) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(credentials)
}

func (controller *Controller) RefreshToken(c *fiber.Ctx) error {
	var params apigen.RefreshTokenRequest
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID, refreshToken, err := controller.auth.ParseJWTRefreshToken(params.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	credentials, err := controller.svc.RefreshToken(c.Context(), userID, refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrRefreshTokenExpired) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(credentials)
}

func (controller *Controller) CreateCluster(c *fiber.Ctx) error {
	var params apigen.ClusterCreate
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	cluster, err := controller.svc.CreateCluster(c.Context(), params, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(cluster)
}

func (controller *Controller) DeleteCluster(c *fiber.Ctx, id int32, params apigen.DeleteClusterParams) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	err = controller.svc.DeleteCluster(c.Context(), id, utils.UnwrapOrDefault(params.Cascade, false), user.OrganizationID)
	if err != nil {
		if errors.Is(err, service.ErrClusterHasDatabaseConnections) {
			return c.Status(fiber.StatusConflict).SendString(err.Error())
		}
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (controller *Controller) GetCluster(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	cluster, err := controller.svc.GetCluster(c.Context(), id, user.OrganizationID)
	if err != nil {
		if errors.Is(err, service.ErrClusterNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(cluster)
}

func (controller *Controller) UpdateCluster(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var params apigen.ClusterCreate
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	cluster, err := controller.svc.UpdateCluster(c.Context(), id, params, user.OrganizationID)
	if err != nil {
		if errors.Is(err, service.ErrClusterNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(cluster)
}

func (controller *Controller) ListClusters(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	clusters, err := controller.svc.ListClusters(c.Context(), user.OrganizationID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(clusters)
}

func (controller *Controller) CreateDatabase(c *fiber.Ctx) error {
	var params apigen.DatabaseConnectInfo
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	database, err := controller.svc.CreateDatabase(c.Context(), params, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(database)
}

func (controller *Controller) DeleteDatabase(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	err = controller.svc.DeleteDatabase(c.Context(), id, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (controller *Controller) GetDatabase(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	database, err := controller.svc.GetDatabase(c.Context(), id, user.OrganizationID)
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

	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	database, err := controller.svc.UpdateDatabase(c.Context(), id, params, user.OrganizationID)
	if err != nil {
		if errors.Is(err, service.ErrDatabaseNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(database)
}

func (controller *Controller) ListDatabases(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	databases, err := controller.svc.ListDatabases(c.Context(), user.OrganizationID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(databases)
}

func (controller *Controller) GetDDLProgress(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	progress, err := controller.svc.GetDDLProgress(c.Context(), id, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(progress)
}

func (controller *Controller) CancelDDLProgress(c *fiber.Ctx, id int32, ddlID int64) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	err = controller.svc.CancelDDLProgress(c.Context(), id, ddlID, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (controller *Controller) TestDatabaseConnection(c *fiber.Ctx) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var params apigen.TestDatabaseConnectionPayload
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result, err := controller.svc.TestDatabaseConnection(c.Context(), params, user.OrganizationID)
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

	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	result, err := controller.svc.QueryDatabase(c.Context(), id, params, user.OrganizationID, utils.UnwrapOrDefault(params.BackgroundDDL, false))
	if err != nil {
		if errors.Is(err, service.ErrDatabaseNotFound) {
			return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("database %d not found", id))
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (controller *Controller) CreateClusterSnapshot(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var params apigen.CreateClusterSnapshotJSONRequestBody
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	snapshot, err := controller.svc.CreateClusterSnapshot(c.Context(), id, params.Name, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(snapshot)
}

func (controller *Controller) DeleteClusterSnapshot(c *fiber.Ctx, id int32, snapshotId int64) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	err = controller.svc.DeleteClusterSnapshot(c.Context(), id, snapshotId, user.OrganizationID)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (controller *Controller) RestoreClusterSnapshot(c *fiber.Ctx, id int32, snapshotId int64) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) ListClusterSnapshots(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	snapshots, err := controller.svc.ListClusterSnapshots(c.Context(), id, user.OrganizationID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(snapshots)
}

func (controller *Controller) ListClusterDiagnostics(c *fiber.Ctx, id int32, params apigen.ListClusterDiagnosticsParams) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	diagnostics, err := controller.svc.ListClusterDiagnostics(c.Context(), id, user.OrganizationID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(diagnostics)
}

func (controller *Controller) GetClusterAutoBackupConfig(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	config, err := controller.svc.GetClusterAutoBackupConfig(c.Context(), id, user.OrganizationID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(config)
}

func (controller *Controller) UpdateClusterAutoBackupConfig(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var params apigen.AutoBackupConfig
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = controller.svc.UpdateClusterAutoBackupConfig(c.Context(), id, params, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (controller *Controller) GetClusterAutoDiagnosticConfig(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	config, err := controller.svc.GetClusterAutoDiagnosticConfig(c.Context(), id, user.OrganizationID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(config)
}

func (controller *Controller) UpdateClusterAutoDiagnosticConfig(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var params apigen.AutoDiagnosticConfig
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = controller.svc.UpdateClusterAutoDiagnosticConfig(c.Context(), id, params, user.OrganizationID)
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

	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	conn, err := controller.svc.TestClusterConnection(c.Context(), params, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(conn)
}

func (controller *Controller) RunRisectlCommand(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var params apigen.RisectlCommand
	if err := c.BodyParser(&params); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result, err := controller.svc.RunRisectlCommand(c.Context(), id, params, user.OrganizationID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (controller *Controller) CreateClusterDiagnostic(c *fiber.Ctx, id int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	diagnostic, err := controller.svc.CreateClusterDiagnostic(c.Context(), id, user.OrganizationID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(diagnostic)
}

func (controller *Controller) GetClusterDiagnostic(c *fiber.Ctx, id int32, diagnosticId int32) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	diagnostic, err := controller.svc.GetClusterDiagnostic(c.Context(), id, diagnosticId, user.OrganizationID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(diagnostic)
}

func (controller *Controller) GetMaterializedViewThroughput(c *fiber.Ctx, clusterID int32) error {
	throughput, err := controller.svc.GetMaterializedViewThroughput(c.Context(), clusterID)
	if err != nil {
		if errors.Is(err, prom.ErrPrometheusEndpointNotFound) {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(throughput)
}
