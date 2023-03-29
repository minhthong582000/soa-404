package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/minhthong582000/soa-404/internal/app/client"
)

func main() {
	_ = godotenv.Load("sample_file.env")

	err := client.HttpClient(os.Getenv("CLIENT_ADDR"), os.Getenv("BIND_ADDR"))
	if err != nil {
		fmt.Println(err)
	}
}
