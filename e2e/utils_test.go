//go:build !ut
// +build !ut

package e2e

import (
	"testing"

	"github.com/risingwavelabs/wavekit/internal/apigen"
)

func loginAccount(t *testing.T, name, password string) apigen.Credentials {
	t.Helper()

	te := getTestEngine(t)

	var credentials apigen.Credentials
	te.POST("/api/v1/auth/sign-in").
		WithJSON(apigen.SignInRequest{
			Name:     name,
			Password: password,
		}).
		Expect().
		Status(200).
		JSON().
		Decode(&credentials)

	// te.GET("/api/v1/auth/ping").
	// 	WithHeader("Authorization", "Bearer "+credentials.AccessToken).
	// 	Expect().
	// 	Status(200)

	mu.Lock()
	defer mu.Unlock()
	token = credentials.AccessToken

	return credentials
}
