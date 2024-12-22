package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/minhthong582000/soa-404/pkg/config"
	interceptors "github.com/minhthong582000/soa-404/pkg/interceptor"
	"github.com/minhthong582000/soa-404/pkg/log"
	metric "github.com/minhthong582000/soa-404/pkg/metrics"
	"github.com/minhthong582000/soa-404/pkg/tracing"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/app/random"
)

// Server to serve the service.
type Server struct {
	config *config.Config
	ctx    context.Context
	logger log.ILogger
}

// New returns a new server.
func New(logger log.ILogger, ctx context.Context, config *config.Config) *Server {
	return &Server{
		config: config,
		ctx:    ctx,
		logger: logger,
	}
}

// Run runs server.
func (s Server) Run() error {
	lis, err := net.Listen("tcp", s.config.Server.BindAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Metrics
	metrics, err := metric.CreateMetrics(s.config.Metrics.BindAddr, s.config.Server.Name)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Infof(
		"Metrics available URL: %s, ServiceName: %s",
		s.config.Metrics.BindAddr,
		s.config.Server.Name,
	)

	// Tracing
	_, err = tracing.TracerFactory(
		tracing.WithProvider(tracing.OTLP),
		tracing.WithCollectorURL(s.config.Tracing.OLTPTracing.CollectorAddr),
		tracing.WithEnabled(s.config.Tracing.OLTPTracing.Enabled),
		tracing.WithInsecure(s.config.Tracing.OLTPTracing.Insecure),
		tracing.WithServiceName(s.config.Server.Name),
	)
	if err != nil {
		s.logger.Errorf("TracerFactory Error: %s", err)
	}

	// Register logs & metrics interceptor
	in := interceptors.NewInterceptorManager(s.logger, metrics)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			in.ExtractRequestID,
			in.Logger,
			in.Metrics,
			grpc_prometheus.UnaryServerInterceptor,
			grpc_ctxtags.UnaryServerInterceptor(),
			recovery.UnaryServerInterceptor(),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// Random service
	randomServer := random.NewServer(
		s.logger,
		random.NewService(
			random.NewRepository(),
		),
	)
	pb.RegisterRandomServiceServer(grpcServer, randomServer)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			s.logger.Info("shutting down gRPC server...")
			grpcServer.GracefulStop()
			<-s.ctx.Done()
		}
	}()

	fmt.Println("gRPC server is running on", s.config.Server.BindAddr)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
