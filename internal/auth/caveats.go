package auth

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit/internal/macaroons"
	"github.com/risingwavelabs/wavekit/internal/utils"
)

var (
	ErrCaveatCheckFailed = errors.New("caveat check failed")
)

const (
	CaveatUserContext = "user_context"
	CaveatRefreshOnly = "refresh_only"
)

type CaveatParser struct {
}

func NewCaveatParser() macaroons.CaveatParser {
	return &CaveatParser{}
}

func (c *CaveatParser) Parse(s string) (macaroons.Caveat, error) {
	typ, err := utils.RetrieveFromJSON[string](s, "type")
	if err != nil {
		return nil, err
	}

	var ret macaroons.Caveat

	switch *typ {
	case CaveatUserContext:
		ret = &UserContextCaveat{}
	case CaveatRefreshOnly:
		ret = &RefreshOnlyCaveat{}
	default:
		return nil, errors.Errorf("unknown caveat type: %s", *typ)
	}

	err = ret.Decode(s)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// EncodeCaveat encodes a caveat to base64 string
func EncodeCaveat(v interface{}) (string, error) {
	json, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(json), nil
}

// DecodeCaveat decodes a base64 string to a caveat
func DecodeCaveat(s string, v interface{}) error {
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, v)
	if err != nil {
		return err
	}

	return nil
}

type UserContextCaveat struct {
	Typ    string `json:"type"`
	UserID int32  `json:"user_id"`
	OrgID  int32  `json:"org_id"`
}

func NewUserContextCaveat(userID int32, orgID int32) *UserContextCaveat {
	return &UserContextCaveat{
		Typ:    CaveatUserContext,
		UserID: userID,
		OrgID:  orgID,
	}
}

func (uc *UserContextCaveat) Encode() (string, error) {
	return EncodeCaveat(uc)
}

func (uc *UserContextCaveat) Decode(s string) error {
	return DecodeCaveat(s, uc)
}

func (uc *UserContextCaveat) Type() string {
	return uc.Typ
}

func (uc *UserContextCaveat) Validate(ctx *fiber.Ctx) error {
	ctx.Locals(ContextKeyUserID, uc.UserID)
	ctx.Locals(ContextKeyOrgID, uc.OrgID)
	return nil
}

type RefreshOnlyCaveat struct {
	Typ         string `json:"type"`
	UserID      int32  `json:"user_id"`
	AccessKeyID int64  `json:"access_key_id"`
}

func NewRefreshOnlyCaveat(userID int32, accessKeyID int64) *RefreshOnlyCaveat {
	return &RefreshOnlyCaveat{
		Typ:         CaveatRefreshOnly,
		UserID:      userID,
		AccessKeyID: accessKeyID,
	}
}

func (rc *RefreshOnlyCaveat) Encode() (string, error) {
	return EncodeCaveat(rc)
}

func (rc *RefreshOnlyCaveat) Decode(s string) error {
	return DecodeCaveat(s, rc)
}

func (rc *RefreshOnlyCaveat) Type() string {
	return rc.Typ
}

func (rc *RefreshOnlyCaveat) Validate(ctx *fiber.Ctx) error {
	if ctx.Method() == "POST" && strings.HasSuffix(ctx.Path(), "/auth/refresh") {
		return nil
	}
	return errors.Wrapf(ErrCaveatCheckFailed, "invalid request: %s %s, the token is for refresh only", ctx.Method(), ctx.Path())
}
