package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/middleware"
	"github.com/risingwavelabs/wavekit/internal/service"
)

type Controller struct {
	mid *middleware.Middleware
	svc service.ServiceInterface
}

var _ apigen.ServerInterface = &Controller{}

func NewController(
	s service.ServiceInterface,
	mid *middleware.Middleware,
) *Controller {
	return &Controller{
		svc: s,
		mid: mid,
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

	user, err := middleware.GetUser(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	credentials, err := controller.svc.RefreshToken(c.Context(), user.ID, params.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrRefreshTokenExpired) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(credentials)
}

func (controller *Controller) CreateCluster(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) DeleteCluster(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) GetCluster(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) UpdateCluster(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) ListClusters(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) CreateDatabase(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) DeleteDatabase(c *fiber.Ctx, id int32) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) GetDatabase(c *fiber.Ctx, id int32) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) UpdateDatabase(c *fiber.Ctx, id int32) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) ListDatabases(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) CreateClusterSnapshot(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) DeleteClusterSnapshot(c *fiber.Ctx, id string, snapshotId string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) RestoreClusterSnapshot(c *fiber.Ctx, id string, snapshotId string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) ListClusterSnapshots(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) GetClusterSnapshotConfig(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) UpdateClusterSnapshotConfig(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) ListClusterDiagnostics(c *fiber.Ctx, id string, params apigen.ListClusterDiagnosticsParams) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) GetClusterDiagnosticConfig(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}

func (controller *Controller) UpdateClusterDiagnosticConfig(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusOK).SendString("Hello, World!")
}
