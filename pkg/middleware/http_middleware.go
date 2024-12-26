package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/minhthong582000/soa-404/pkg/log"
)

type Middleware struct {
	logger log.Logger
}

func NewMiddleware(logger log.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			m.logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	})
}
