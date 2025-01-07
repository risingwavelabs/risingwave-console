//go:build !ut
// +build !ut

package e2e

import (
	"log"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/risingwavelabs/wavekit/internal/apigen"
	"github.com/risingwavelabs/wavekit/internal/apps/server"
	"github.com/risingwavelabs/wavekit/wire"
)

func newHttpExpect(handler http.Handler, t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewRequireReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

var (
	apiServer *server.Server
	token     string
	mu        sync.Mutex
)

var (
	globalUsername = "root"
	globalPassword = "123456"
)

func TestMain(m *testing.M) {
	server, err := wire.InitializeServer()
	if err != nil {
		log.Fatal(err)
	}
	apiServer = server

	os.Exit(m.Run())
}

func getTestEngine(t *testing.T) *httpexpect.Expect {
	endpoint := os.Getenv("TEST_ENDPOINT")
	if endpoint != "" {
		return httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  endpoint,
			Reporter: httpexpect.NewRequireReporter(t),
			Printers: []httpexpect.Printer{
				httpexpect.NewDebugPrinter(t, true),
			},
		})

	}
	return newHttpExpect(adaptor.FiberApp(apiServer.GetApp()), t)
}

type AutenticatedTestEngine struct {
	*httpexpect.Expect
	authInfo apigen.Credentials
}

func getAuthenticatedTestEngine(t *testing.T) *AutenticatedTestEngine {
	te := getTestEngine(t)

	authInfo := loginAccount(t, globalUsername, globalPassword)

	return &AutenticatedTestEngine{
		Expect: te.Builder(func(req *httpexpect.Request) {
			req.WithHeader("Authorization", "Bearer "+token)
		}),
		authInfo: authInfo,
	}
}
