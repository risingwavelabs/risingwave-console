package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/auth"
	"github.com/risingwavelabs/wavekit/internal/model"
	"github.com/risingwavelabs/wavekit/internal/model/querier"
)

type Validator struct {
	model model.ModelInterface
}

func NewValidator(model model.ModelInterface) apigen.Validator {
	return &Validator{model: model}
}

func (v *Validator) GetOrgID(c *fiber.Ctx) int32 {
	return c.Locals(auth.ContextKeyOrgID).(int32)
}

func (v *Validator) OwnDatabase(c *fiber.Ctx, orgID int32, databaseID int32) error {
	_, err := v.model.GetOrgDatabaseByID(c.Context(), querier.GetOrgDatabaseByIDParams{
		ID:             databaseID,
		OrganizationID: orgID,
	})
	if err != nil {
		return err
	}
	return nil
}
