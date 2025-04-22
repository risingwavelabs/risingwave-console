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
	auth  auth.AuthInterface
}

func NewValidator(model model.ModelInterface, auth auth.AuthInterface) apigen.Validator {
	return &Validator{model: model, auth: auth}
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

func (v *Validator) PreValidate(c *fiber.Ctx) error {
	return v.auth.Authfunc(c)
}

func (v *Validator) PostValidate(c *fiber.Ctx) error {
	return nil
}
