package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/minhthong582000/soa-404/pkg/log"
	"github.com/minhthong582000/soa-404/pkg/metric"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Logger() echo.MiddlewareFunc {
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

func (m *Middleware) Metrics() echo.MiddlewareFunc {
	metr := metric.GetMetric()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()
			path := c.Path() // contains route path ala `/users/:id`

			if metr.IsMetricExist(metric.Http_request_inflight.Name) {
				metr.Counter(metric.Grpc_request_inflight, 1, path)
				defer func() {
					metr.Counter(metric.Grpc_request_inflight, -1, path)
				}()
			}

			reqSz := computeApproximateRequestSize(c.Request())

			err := next(c)

			status := c.Response().Status
			if err != nil {
				var httpError *echo.HTTPError
				if errors.As(err, &httpError) {
					status = httpError.Code
				}
				if status == 0 || status == http.StatusOK {
					status = http.StatusInternalServerError
				}
			}
			statusStr := strconv.Itoa(status)
			resSz := c.Response().Size

			if metr.IsMetricExist(metric.Http_request_total.Name) {
				metr.Counter(metric.Http_request_total, 1, path, statusStr)
			}
			if metr.IsMetricExist(metric.Http_request_duration_seconds.Name) {
				metr.Histogram(metric.Http_request_duration_seconds, time.Since(startTime).Seconds(), path, statusStr)
			}
			if metr.IsMetricExist(metric.Http_request_size_bytes.Name) {
				metr.Histogram(metric.Http_request_size_bytes, float64(reqSz), path)
			}
			if metr.IsMetricExist(metric.Http_response_size_bytes.Name) {
				metr.Histogram(metric.Http_response_size_bytes, float64(resSz), path)
			}

			return nil
		}
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
