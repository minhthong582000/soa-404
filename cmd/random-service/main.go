package main

import (
	"context"
	"fmt"

	"github.com/minhthong582000/soa-404/internal/server"
	"github.com/minhthong582000/soa-404/pkg/config"
)

func main() {
	ctx := context.Background()

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

	server := server.New(ctx, config)
	if err := server.Run(); err != nil {
		fmt.Println(err)
	}
}
