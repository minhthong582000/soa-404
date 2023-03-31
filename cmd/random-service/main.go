package main

import (
	"context"
	"fmt"

	"github.com/minhthong582000/soa-404/config"
	"github.com/minhthong582000/soa-404/internal/app/random"
	"github.com/minhthong582000/soa-404/internal/server"
	"github.com/minhthong582000/soa-404/pkg/log"
)

func main() {
	ctx := context.Background()
	logger := log.New().With(ctx)

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

	randomServer := random.NewServer(logger, random.NewService(random.NewRepository()))
	server := server.New(logger, ctx, config, randomServer)

	err = server.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
