package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/app/random"
	"google.golang.org/grpc"
)

// Server to serve the service.
type Server struct {
	grpcServer *grpc.Server
	bindAddr   string
	ctx        context.Context

	randomServer *random.RandomServer
}

// New returns a new server.
func New(ctx context.Context, bindAddr string, randomServer *random.RandomServer) *Server {
	grpcServer := grpc.NewServer()

	return &Server{
		grpcServer:   grpcServer,
		bindAddr:     bindAddr,
		ctx:          ctx,
		randomServer: randomServer,
	}
}

// Run runs server.
func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.bindAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Random service
	pb.RegisterRandomServiceServer(s.grpcServer, s.randomServer)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			s.grpcServer.GracefulStop()

			<-s.ctx.Done()
		}
	}()

	fmt.Println("gRPC server is running on", s.bindAddr)
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
