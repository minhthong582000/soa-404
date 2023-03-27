package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/minhthong582000/soa-404/internal/app/random"
	"github.com/minhthong582000/soa-404/internal/server"
	"github.com/minhthong582000/soa-404/pkg/log"
)

var Version = "1.0.0"

func main() {
	ctx := context.Background()
	_ = godotenv.Load("sample_file.env")
	logger := log.New().With(ctx, "version", Version)

	randomServer := random.NewServer(logger, random.NewService(random.NewRepository()))
	server := server.New(ctx, os.Getenv("BIND_ADDR"), randomServer)

	err := server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
