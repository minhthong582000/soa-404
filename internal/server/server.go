package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/minhthong582000/soa-404/pkg/config"
	interceptors "github.com/minhthong582000/soa-404/pkg/interceptor"
	"github.com/minhthong582000/soa-404/pkg/log"
	metric "github.com/minhthong582000/soa-404/pkg/metrics"
	"github.com/minhthong582000/soa-404/pkg/tracing"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"

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
	tracer, err := tracing.TracerFactory(tracing.OLTP, tracing.TracerConfig{
		ServiceName:  s.config.Server.Name,
		CollectorURL: s.config.Tracing.OLTPTracing.CollectorAddr,
		Insecure:     s.config.Tracing.OLTPTracing.Insecure,
	})
	if err != nil {
		s.logger.Errorf("TracerFactory Error: %s", err)
	}
	cleanup, err := tracer.InitTracer()
	if err != nil {
		s.logger.Errorf("InitTracer Error: %s", err)
	}
	tp := otel.Tracer(s.config.Server.Name)

	// Register logs & metrics interceptor
	in := interceptors.NewInterceptorManager(s.logger, metrics)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(in.Logger),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.ChainUnaryInterceptor(
			in.Metrics,
			grpc_prometheus.UnaryServerInterceptor,
			otelgrpc.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			recovery.UnaryServerInterceptor(),
		),
	)

	// Random service
	randomServer := random.NewServer(
		s.logger,
		tp,
		random.NewService(
			tp,
			random.NewRepository(tp),
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
			err = cleanup(s.ctx)
			if err != nil {
				s.logger.Errorf("trace cleanup error: %s", err)
			}
			<-s.ctx.Done()
		}
	}()

	fmt.Println("gRPC server is running on", s.config.Server.BindAddr)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
