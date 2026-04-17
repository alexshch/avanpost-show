package main

import (
	_ "avanpost-show/docs"
	"avanpost-show/internal/app"
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	server, err := app.NewApp(ctx)
	if err != nil {
		slog.Error("failed to init", "err", err)
	}
	if err = server.Run(ctx); err != nil {
		slog.Error("run server error", "err", err)
		return
	}
	<-ctx.Done()
	ctx, stop = context.WithTimeout(context.Background(), 3*time.Second)
	defer stop()
	_ = server.Stop(ctx)
}
