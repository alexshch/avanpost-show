package app

import (
	"context"
	"log/slog"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	stackSize      = 1 << 10 // 1 KB
	bodyLimit      = "2M"
	readTimeout    = 15 * time.Second
	writeTimeout   = 15 * time.Second
	gzipLevel      = 5
)

func (a *App) Run(ctx context.Context) (err error) {
	go func() {
		err := a.runHttpServer()
		if err != nil {
			slog.Error("run server error", "err", err)
		}
	}()
	return
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.router.Shutdown(ctx); err != nil {
		return err
	}
	a.db.Close()
	return nil
}

func (a *App) runHttpServer() error {
	a.router.Server.ReadTimeout = readTimeout
	a.router.Server.WriteTimeout = writeTimeout
	a.router.Server.MaxHeaderBytes = maxHeaderBytes

	return a.router.Start(a.config.Addr())
}
