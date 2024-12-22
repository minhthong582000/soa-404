package client

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/pkg/config"
	"github.com/minhthong582000/soa-404/pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Client is a simple client for the Random service.
type Client struct {
	randClient pb.RandomServiceClient
}

// NewClient creates a new client.
func NewClient(randClient pb.RandomServiceClient) *Client {
	return &Client{
		randClient: randClient,
	}
}

// GetRandNumber gets a random number from the server.
func (c Client) GetRandNumber(ctx context.Context, seed int64) (int64, error) {
	reply, err := c.randClient.GetRandNumber(ctx, &pb.GetRandNumberRequest{
		SeedNum: seed,
	})
	if err != nil {
		return -1, err
	}

	return reply.Number, nil
}

// HttpClient runs the client.
func HttpClient(config *config.Config) error {
	// Tracing
	tracer, err := tracing.TracerFactory(
		tracing.WithProvider(tracing.OTLP),
		tracing.WithCollectorURL(config.Tracing.OLTPTracing.CollectorAddr),
		tracing.WithEnabled(config.Tracing.OLTPTracing.Enabled),
		tracing.WithInsecure(config.Tracing.OLTPTracing.Insecure),
		tracing.WithServiceName(config.Client.Name),
	)
	if err != nil {
		log.Fatalf("TracerFactory Error: %s", err)
	}
	cleanup, err := tracer.InitTracer()
	if err != nil {
		log.Fatalf("InitTracer Error: %s", err)
	}
	defer cleanup(context.Background())

	kacp := keepalive.ClientParameters{
		Timeout: 10 * time.Second,
		Time:    1 * time.Minute,
	}
	// Set up a connection to the server
	conn, err := grpc.NewClient(
		config.Client.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	randClient := pb.NewRandomServiceClient(conn)
	client := NewClient(randClient)

	router := echo.New()
	router.Use(middleware.RequestID())
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

		// Call the server
		randNum, err := client.GetRandNumber(c.Request().Context(), seed)
		if err != nil {
			return c.String(500, "failed to get random number")
		}

		// Return the random number in JSON
		return c.JSON(200, map[string]int64{
			"number": randNum,
		})
	})
	if err := router.Start(config.Client.BindAddr); err != nil {
		log.Fatal(err)
	}

	return nil
}
