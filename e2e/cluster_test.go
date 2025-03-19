//go:build !ut
// +build !ut

package e2e

import "testing"

func TestCluster(t *testing.T) {
	te := getAuthenticatedTestEngine(t)

	te.GET("/api/v1/clusters").
		Expect().
		Status(200)
}
