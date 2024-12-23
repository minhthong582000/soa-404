package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/app/random"
	"github.com/minhthong582000/soa-404/pkg/config"
	"github.com/minhthong582000/soa-404/pkg/log"
	metric "github.com/minhthong582000/soa-404/pkg/metrics"
	"github.com/minhthong582000/soa-404/pkg/middleware"
	"github.com/minhthong582000/soa-404/pkg/tracing"
)

// Server to serve the service.
type Server struct {
	config *config.Config
	ctx    context.Context
}

// New returns a new server.
func New(ctx context.Context, config *config.Config) *Server {
	return &Server{
		config: config,
		ctx:    ctx,
	}
}

// Run runs server.
func (s Server) Run() error {
	// gRPC listener
	lis, err := net.Listen("tcp", s.config.Server.BindAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Logger
	logger := log.Init(&s.config.Logs)

	// Metrics
	metrics, err := metric.CreateMetrics(s.config.Metrics.BindAddr, s.config.Server.Name)
	if err != nil {
		logger.Errorf("CreateMetrics Error: %s", err)
	}
	logger.Infof(
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
		logger.Errorf("TracerFactory Error: %s", err)
	}

	// Register logs & metrics & trace interceptor
	in := middleware.NewInterceptor(metrics)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			in.Logger,
			in.Metrics,
			grpc_prometheus.UnaryServerInterceptor,
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

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Info("shutting down gRPC server...")
			grpcServer.GracefulStop()
			<-s.ctx.Done()
		}
	}()

	logger.Infof("gRPC server is running on %s", s.config.Server.BindAddr)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
