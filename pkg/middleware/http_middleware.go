package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/minhthong582000/soa-404/pkg/log"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) RequestLogger() echo.MiddlewareFunc {
	logger := log.GetLogger()
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.With(
				c.Request().Context(),
				"uri", v.URI,
				"status", v.Status,
				"method", v.Method,
			).Infof("received request")

			return nil
		},
	})
}
