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
	interceptors "github.com/minhthong582000/soa-404/pkg/interceptor"
	"github.com/minhthong582000/soa-404/pkg/log"
	metric "github.com/minhthong582000/soa-404/pkg/metrics"
	"github.com/minhthong582000/soa-404/pkg/tracing"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/app/random"
)

// Server to serve the service.
type Server struct {
	bindAddr        string
	metricsBindAddr string
	ctx             context.Context
	logger          log.ILogger

	randomServer *random.RandomServer
}

// New returns a new server.
func New(logger log.ILogger, ctx context.Context, bindAddr string, metricsBindAddr string, randomServer *random.RandomServer) *Server {
	return &Server{
		bindAddr:        bindAddr,
		ctx:             ctx,
		randomServer:    randomServer,
		logger:          logger,
		metricsBindAddr: metricsBindAddr,
	}
}

// Run runs server.
func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.bindAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Metrics
	metrics, err := metric.CreateMetrics(s.metricsBindAddr, "random")
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Infof(
		"Metrics available URL: %s, ServiceName: %s",
		s.metricsBindAddr,
		"random",
	)

	// Tracing
	tracer, err := tracing.TracerFactory(tracing.OLTP, tracing.TracerConfig{
		ServiceName:  "random",
		CollectorURL: "localhost:8069",
		Insecure:     true,
	})
	if err != nil {
		s.logger.Errorf("TracerFactory Error: %s", err)
	}
	cleanup, err := tracer.InitTracer()
	if err != nil {
		s.logger.Errorf("InitTracer Error: %s", err)
	}

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
	pb.RegisterRandomServiceServer(grpcServer, s.randomServer)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			s.logger.Info("shutting down gRPC server...")
			grpcServer.GracefulStop()
			cleanup(s.ctx)
			<-s.ctx.Done()
		}
	}()

	fmt.Println("gRPC server is running on", s.bindAddr)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
