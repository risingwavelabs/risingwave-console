//go:build !ut
// +build !ut

package e2e

import (
	"testing"

	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/stretchr/testify/assert"
)

func registerAccount(t *testing.T, phone, username, password string) {
	t.Helper()

	te := getTestEngine(t)

	te.POST("/api/v1/auth/code").
		WithJSON(apigen.PostAuthCodeJSONBody{
			Phone: phone,
			Typ:   apigen.Register,
		}).
		Expect().
		Status(202)

	var param = apigen.PostAuthRegisterJSONBody{
		Phone:    phone,
		Username: username,
		Password: password,
	}

	te.POST("/api/v1/auth/register").
		WithJSON(param).
		Expect().
		Status(200)
}

func loginAccount(t *testing.T, _, phone, password string) apigen.AuthInfo {
	t.Helper()

	te := getTestEngine(t)

	var authInfo apigen.AuthInfo
	te.POST("/api/v1/auth/login").
		WithJSON(apigen.PostAuthLoginJSONBody{
			Phone:    phone,
			Password: password,
		}).
		Expect().
		Status(200).
		JSON().
		Decode(&authInfo)

	te.GET("/api/v1/auth/ping").
		WithHeader("Authorization", "Bearer "+authInfo.Token).
		Expect().
		Status(200)

	var orgs apigen.OrgsList
	te.GET("/api/v1/orgs").
		WithHeader("Authorization", "Bearer "+authInfo.Token).
		Expect().
		Status(200).
		JSON().
		Decode(&orgs)
	assert.Len(t, orgs, 1)
	assert.Equal(t, "sage的小组", orgs[0].Name)

	mu.Lock()
	defer mu.Unlock()
	token = authInfo.Token

	return authInfo
}
