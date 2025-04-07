package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/globalctx"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logger.NewLogAgent("metrics")

var WorkerGoroutines = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "wavekit_worker_goroutines",
		Help: "The number of goroutines that are running",
	},
)

type MetricsServer struct {
	port      int
	server    *http.Server
	globalCtx *globalctx.GlobalContext
}

func (m *MetricsServer) Start() {
	go func() {
		log.Infof("metrics server is listening on port %d", m.port)
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("metrics server exited", zap.Error(err))
		}
	}()

	// Shutdown the server when the global context is done
	go func() {
		<-m.globalCtx.Context().Done()
		if err := m.server.Shutdown(context.Background()); err != nil {
			log.Error("metrics server shutdown error", zap.Error(err))
		} else {
			log.Info("metrics server shutdown gracefully")
		}
	}()

	ready := make(chan struct{})

	go func() {
		for range 5 {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics", m.port))
			if err == nil {
				resp.Body.Close()
				close(ready)
				return
			}
			time.Sleep(time.Second)
		}
	}()

	// Wait for the server to be ready or timeout
	select {
	case <-ready:
		log.Info("metrics server started successfully")
	case <-time.After(5 * time.Second):
		panic("timed out waiting for metrics server to start")
	}
}

func NewMetricsServer(cfg *config.Config, globalCtx *globalctx.GlobalContext) *MetricsServer {
	port := 9020
	if cfg.MetricsPort != 0 {
		port = cfg.MetricsPort
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &MetricsServer{
		port:      port,
		server:    server,
		globalCtx: globalCtx,
	}
}
