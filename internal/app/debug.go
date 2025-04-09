package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/risingwavelabs/wavekit/internal/config"
	"github.com/risingwavelabs/wavekit/internal/globalctx"
	"github.com/risingwavelabs/wavekit/internal/logger"
	"go.uber.org/zap"
)

var log = logger.NewLogAgent("debug-server")

type DebugServer struct {
	globalCtx *globalctx.GlobalContext
	port      int
	enable    bool
}

func NewDebugServer(cfg *config.Config, globalCtx *globalctx.GlobalContext) *DebugServer {
	return &DebugServer{
		globalCtx: globalCtx,
		port:      cfg.Debug.Port,
		enable:    cfg.Debug.Enable,
	}
}

func (d *DebugServer) Start() error {
	if !d.enable {
		return nil
	}
	if d.port == 0 {
		d.port = 8777
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", d.port),
		Handler: mux,
	}

	go func() {
		log.Info("debug server is listening", zap.Int("port", d.port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("debug server exited", zap.Error(err))
		}
	}()

	<-d.globalCtx.Context().Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("debug server shutdown error", zap.Error(err))
	} else {
		log.Info("debug server shutdown gracefully")
	}

	return nil
}
