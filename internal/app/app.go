package app

import (
	"avanpost-show/internal/app/handler"
	"avanpost-show/internal/config"
	"avanpost-show/pkg/publisher"
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

type App struct {
	config   *config.Config
	db       *pgxpool.Pool
	s3Client *s3.Client
	nc       *nats.Conn
	router   *echo.Echo
	handlers []handler.Handler
	pub      *publisher.Publisher
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}
	err := app.initApp(ctx)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) initApp(ctx context.Context) (err error) {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLog,
		a.initDb,
		a.initNats,
		a.initPublisher,
		a.initRouter,
		a.initSwagger,
		a.initMiddleware,
		a.initHandlers,
		a.registerRoutes,
	}
	for _, f := range inits {
		if err = f(ctx); err != nil {
			slog.Error("init error", "err", err)
			return err
		}
	}
	return
}
