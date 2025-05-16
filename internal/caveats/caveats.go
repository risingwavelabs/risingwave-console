package caveats

import (
	"github.com/cloudcarver/anchor/pkg/macaroons"
	"github.com/gofiber/fiber/v2"
)

const ContextKeyOrgID = "wavekit_org_id"

type OrgContextCaveat struct {
	Typ   string `json:"type"`
	OrgID int32  `json:"org_id"`
}

func NewOrgContextCaveat(orgID int32) macaroons.Caveat {
	return &OrgContextCaveat{Typ: ContextKeyOrgID, OrgID: orgID}
}

func (caveat *OrgContextCaveat) Validate(c *fiber.Ctx) error {
	c.Locals(ContextKeyOrgID, caveat.OrgID)
	return nil
}

func (caveat *OrgContextCaveat) Type() string {
	return caveat.Typ
}
