package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/minhthong582000/soa-404/internal/app/random"
	"github.com/minhthong582000/soa-404/internal/server"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load("sample_file.env")
	fmt.Println(os.Getenv("BIND_ADDR"))

	randomServer := random.NewServer(random.NewService(random.NewRepository()))
	server := server.New(ctx, os.Getenv("BIND_ADDR"), randomServer)

	err := server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
