package controller

import (
	"errors"

	"github.com/cloudcarver/anchor/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/risingwavelabs/wavekit/pkg/zcore/model"
	"github.com/risingwavelabs/wavekit/pkg/zgen/apigen"
	"github.com/risingwavelabs/wavekit/pkg/zgen/querier"
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
		ID:    databaseID,
		OrgID: orgID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *Validator) AuthFunc(c *fiber.Ctx) error {
	return v.auth.Authfunc(c)
}

func (v *Validator) PreValidate(c *fiber.Ctx) error {
	return nil
}

func (v *Validator) PostValidate(c *fiber.Ctx) error {
	return nil
}

func (v *Validator) PremiumAccess(c *fiber.Ctx) error {
	return errors.New("premium access required")
}
