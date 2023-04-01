package main

import (
	"context"
	"fmt"

	"github.com/minhthong582000/soa-404/internal/app/random"
	"github.com/minhthong582000/soa-404/internal/server"
	"github.com/minhthong582000/soa-404/pkg/config"
	"github.com/minhthong582000/soa-404/pkg/log"
	"github.com/minhthong582000/soa-404/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()

	// Logging
	logger := log.New().With(ctx)

	// Read config
	v, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	config, err := config.ParseConfig(v)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Tracing
	tracer, err := tracing.TracerFactory(tracing.OLTP, tracing.TracerConfig{
		ServiceName:  config.Server.Name,
		CollectorURL: config.Tracing.OLTPTracing.CollectorAddr,
		Insecure:     config.Tracing.OLTPTracing.Insecure,
	})
	if err != nil {
		logger.Errorf("TracerFactory Error: %s", err)
	}
	cleanup, err := tracer.InitTracer()
	if err != nil {
		logger.Errorf("InitTracer Error: %s", err)
	}
	defer cleanup(ctx)
	tp := otel.Tracer(config.Server.Name)

	randomServer := random.NewServer(
		logger,
		tp,
		random.NewService(
			tp,
			random.NewRepository(tp),
		),
	)
	server := server.New(logger, ctx, config, randomServer)

	err = server.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
