package middleware

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LogMiddleware() echo.MiddlewareFunc {
	logger := slog.Default()
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogMethod:   true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, v.Method,
					slog.Int("status", v.Status), slog.String("remoteIP", v.RemoteIP),
					slog.String("uri", v.URI),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "ERROR",
					slog.String("remoteIP", v.RemoteIP), slog.Int("status", v.Status),
					slog.String("method", v.Method), slog.String("uri", v.URI),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	})
}
