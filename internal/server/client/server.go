package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/app/client"
	"github.com/minhthong582000/soa-404/pkg/config"
	"github.com/minhthong582000/soa-404/pkg/log"
	"github.com/minhthong582000/soa-404/pkg/metric"
	http_middleware "github.com/minhthong582000/soa-404/pkg/middleware"
	"github.com/minhthong582000/soa-404/pkg/tracing"
)

// Server to serve the service.
type Server struct {
	config *config.Config
}

// New returns a new server.
func New(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

// HttpClient runs the client.
func (s Server) Run(stopCh <-chan struct{}) error {
	// Logs
	s.config.Logs.Provider = config.ZapLog
	logger, err := log.LogFactory(&s.config.Logs)
	if err != nil {
		return fmt.Errorf("error initializing logger: %v", err)
	}
	httpMiddleware := http_middleware.NewMiddleware()

	// Tracing
	_, err = tracing.TracerFactory(
		tracing.WithProvider(tracing.OTLP),
		tracing.WithCollectorURL(s.config.Tracing.OLTPTracing.CollectorAddr),
		tracing.WithEnabled(s.config.Tracing.OLTPTracing.Enabled),
		tracing.WithInsecure(s.config.Tracing.OLTPTracing.Insecure),
		tracing.WithServiceName(s.config.Client.Name),
	)
	if err != nil {
		return fmt.Errorf("error initializing tracer: %v", err)
	}

	// Metrics
	metrics, err := metric.MetricFactory(
		metric.WithProvider(metric.Prometheus),
		metric.WithMetrics(
			metric.Http_request_inflight,
			metric.Http_request_total,
			metric.Http_request_duration_seconds,
			metric.Http_response_size_bytes,
			metric.Http_request_size_bytes,
		),
	)
	if err != nil {
		return fmt.Errorf("error initializing metrics: %v", err)
	}

	kacp := keepalive.ClientParameters{
		Timeout: 10 * time.Second,
		Time:    1 * time.Minute,
	}
	// Set up a connection to the server
	conn, err := grpc.NewClient(
		s.config.Client.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return fmt.Errorf("unable to connect to \"%s\": %v", s.config.Client.ServerAddr, err)
	}
	defer conn.Close()

	randClient := pb.NewRandomServiceClient(conn)
	client := client.NewClient(randClient)

	router := echo.New()
	router.Use(middleware.RequestID())
	router.Use(httpMiddleware.Logger())
	router.Use(httpMiddleware.Metrics())
	router.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	router.GET("/random", func(c echo.Context) error {
		seedStr := c.QueryParam("seed")

		// Check if seed is empty
		if seedStr == "" {
			return c.String(400, "seed is required")
		}

		// Convert seed to int64
		seed, err := strconv.ParseInt(seedStr, 10, 64)
		if err != nil {
			return c.String(400, "seed must be an integer")
		}

		// Extract Client IP
		clientIP := c.RealIP()

		// Add client IP to gRPC metadata
		ctx := metadata.AppendToOutgoingContext(c.Request().Context(), "x-client-ip", clientIP)

		// Call the server
		randNum, err := client.GetRandNumber(ctx, seed)
		if err != nil {
			return c.String(500, "failed to get random number")
		}

		// Return the random number in JSON
		return c.JSON(200, map[string]int64{
			"number": randNum,
		})
	})

	errCh := make(chan error, 1)
	defer func() {
		logger.Info("Shutting down HTTP server...")
		close(errCh)
		logger.Info("Bye!")
	}()

	// Run HTTP server
	go func() {
		if err := router.Start(s.config.Client.BindAddr); err != nil {
			errCh <- err
		}
	}()

	// Run metrics server
	go func() {
		metrics.RunHTTPMetricsServer(context.Background(), s.config.Metrics.BindAddr)
	}()

	select {
	case <-stopCh:
		logger.Infof("Received shutdown signal")
	case err := <-errCh:
		return err
	}

	if err := router.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}
