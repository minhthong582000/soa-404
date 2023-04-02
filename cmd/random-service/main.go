package main

import (
	"context"
	"fmt"

	"github.com/minhthong582000/soa-404/internal/server"
	"github.com/minhthong582000/soa-404/pkg/config"
	"github.com/minhthong582000/soa-404/pkg/log"
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

	server := server.New(logger, ctx, config)

	err = server.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
