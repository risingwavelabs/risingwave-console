package globalctx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type GlobalContext struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func New() *GlobalContext {
	ctx, cancel := context.WithCancel(context.Background())
	globalCtx := &GlobalContext{
		ctx:        ctx,
		cancelFunc: cancel,
	}

	// Setup signal handling
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signalCh:
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(signalCh)
	}()

	return globalCtx
}

func (g *GlobalContext) Context() context.Context {
	return g.ctx
}
