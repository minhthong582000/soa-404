package main

import (
	"fmt"

	"github.com/minhthong582000/soa-404/internal/server/server"
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
	s := server.New(config)
	if err := s.Run(stopCh); err != nil {
		fmt.Println(err)
	}
}
