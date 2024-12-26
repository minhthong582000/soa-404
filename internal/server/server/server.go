package server

import (
	"context"
	"fmt"
	"net"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/app/random"
	"github.com/minhthong582000/soa-404/pkg/config"
	"github.com/minhthong582000/soa-404/pkg/log"
	"github.com/minhthong582000/soa-404/pkg/metric"
	"github.com/minhthong582000/soa-404/pkg/middleware"
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

// Run runs server.
func (s Server) Run(stopCh <-chan struct{}) error {
	// gRPC listener
	lis, err := net.Listen("tcp", s.config.Server.BindAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Logger
	s.config.Logs.Provider = config.ZapLog
	logger, err := log.LogFactory(&s.config.Logs)
	if err != nil {
		return fmt.Errorf("error initializing logger: %v", err)
	}

	// Metrics
	metrics, err := metric.MetricFactory(
		metric.WithProvider(metric.Prometheus),
		metric.WithMetrics(
			metric.Grpc_server_handled_total,
			metric.Grpc_server_msg_received_total,
			metric.Grpc_server_msg_sent_total,
			metric.Grpc_server_handling_seconds,
		),
	)
	if err != nil {
		return fmt.Errorf("error initializing metrics: %v", err)
	}

	// Tracing
	_, err = tracing.TracerFactory(
		tracing.WithProvider(tracing.OTLP),
		tracing.WithCollectorURL(s.config.Tracing.OLTPTracing.CollectorAddr),
		tracing.WithEnabled(s.config.Tracing.OLTPTracing.Enabled),
		tracing.WithInsecure(s.config.Tracing.OLTPTracing.Insecure),
		tracing.WithServiceName(s.config.Server.Name),
	)
	if err != nil {
		return fmt.Errorf("error initializing tracer: %v", err)
	}

	// Register logs & metrics & trace interceptor
	in := middleware.NewInterceptor()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			in.Logger,
			in.Metrics,
			grpc_ctxtags.UnaryServerInterceptor(),
			recovery.UnaryServerInterceptor(),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	randomServer := random.NewServer(
		random.NewService(
			random.NewRepository(),
		),
	)
	pb.RegisterRandomServiceServer(grpcServer, randomServer)

	errCh := make(chan error, 1)
	defer func() {
		logger.Infof("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		close(errCh)
		logger.Info("Bye!")
	}()

	// Run gRPC server
	go func() {
		logger.Infof("gRPC server is running on %s", s.config.Server.BindAddr)
		if err := grpcServer.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	// Run metrics server
	go func() {
		logger.Infof(
			"Metrics available URL: %s, ServiceName: %s",
			s.config.Metrics.BindAddr,
			s.config.Server.Name,
		)
		metrics.RunHTTPMetricsServer(context.Background(), s.config.Metrics.BindAddr)
	}()

	// Wait for shutdown signal
	select {
	case <-stopCh:
		logger.Infof("Received shutdown signal")
	case err := <-errCh:
		return err
	}

	return nil
}
