package main

import (
	"fmt"

	"github.com/minhthong582000/soa-404/internal/server/client"
	"github.com/minhthong582000/soa-404/pkg/config"
	"github.com/minhthong582000/soa-404/pkg/signals"
)

func main() {
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

	stopCh := signals.SetupSignalHandler()
	c := client.New(config)
	if err := c.Run(stopCh); err != nil {
		fmt.Println(err)
	}
}
