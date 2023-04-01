package main

import (
	"fmt"

	"github.com/minhthong582000/soa-404/internal/app/client"
	"github.com/minhthong582000/soa-404/pkg/config"
)

func main() {
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

	err = client.HttpClient(config)
	if err != nil {
		fmt.Println(err)
	}
}
