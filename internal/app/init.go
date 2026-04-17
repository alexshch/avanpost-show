package app

import (
	"avanpost-show/docs"
	"avanpost-show/internal/app/handler"
	appconfig "avanpost-show/internal/config"
	"avanpost-show/pkg/config"
	middleware2 "avanpost-show/pkg/middleware"
	"avanpost-show/pkg/postgres"
	"avanpost-show/pkg/publisher"
	"avanpost-show/pkg/router"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"
	echoswagger "github.com/swaggo/echo-swagger"
)

func (a *App) initConfig(_ context.Context) (err error) {
	a.config, err = appconfig.NewConfigLoader().Load()
	return
}

func (a *App) initLog(_ context.Context) (err error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.GetLogLevel(a.config.LogLevel),
	}))
	slog.SetDefault(logger)
	slog.Info("config", "val", a.config)
	return
}

func (a *App) initDb(_ context.Context) (err error) {
	if strings.TrimSpace(a.config.Db.ConnectionString) == "" {
		err = errors.New("db: an empty connection string")
		return
	}
	db, err := postgres.NewPgxConn(&a.config.Db)
	a.db = db
	return
}

//func (a *App) initVault(ctx context.Context) (err error) {
//	vaultClient, err := vault.InitVault(a.config.HashiCorpVault, a.config.VaultGRPC)
//	if err != nil {
//		panic(err)
//	}
//	a.vaultClient = vaultClient
//	return nil
//}

func (a *App) initNats(ctx context.Context) (err error) {
	natsConfig := a.config.Nats
	a.nc, err = nats.Connect(natsConfig.Url, nats.UserInfo(natsConfig.Login, natsConfig.Password))
	if err != nil {
		return err
	}
	return err
}

func (a *App) initPublisher(ctx context.Context) error {
	a.pub = publisher.NewPublisher(a.nc)
	return nil
}

func (a *App) initRouter(_ context.Context) (err error) {
	a.router = router.New()
	a.router.HideBanner = true
	a.router.Validator = router.NewValidator()
	return
}

func (a *App) initSwagger(_ context.Context) error {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "API Presentation Service"
	docs.SwaggerInfo.Description = "API service"
	docs.SwaggerInfo.BasePath = a.config.Http.BaseAPI

	a.router.GET("/swagger/*", echoswagger.WrapHandler)

	a.router.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	return nil
}

func (a *App) initMiddleware(_ context.Context) error {
	a.router.Use(middleware.RequestID())

	a.router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	a.router.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Limit: bodyLimit,
	}))

	a.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		//AllowOrigins:     []string{"http://localhost:3300"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAccessControlAllowOrigin,
			echo.HeaderAuthorization,
			echo.HeaderContentDisposition,
			echo.HeaderAccessControlExposeHeaders,
		},
		ExposeHeaders: []string{
			echo.HeaderContentDisposition,
			echo.HeaderContentLength,
			echo.HeaderSetCookie,
		},
	}))

	a.router.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		echoError := &echo.HTTPError{}
		if errors.As(err, &echoError) {
			code = echoError.Code
		}
		c.Logger().Error(err)
		_ = c.NoContent(code)
	}

	a.router.Use(middleware2.LogMiddleware())

	return nil
}

func (a *App) initHandlers(ctx context.Context) (err error) {
	a.handlers = []handler.Handler{
		handler.NewUserHandler(a.db, a.pub),
	}
	return
}

func (a *App) registerRoutes(ctx context.Context) (err error) {
	api := a.router.Group(a.config.Http.BaseAPI)
	for i := range a.handlers {
		a.handlers[i].Register(api)
	}
	return
}
