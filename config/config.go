package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Server config
type Server struct {
	BindAddr string `mapstructure:"bind_addr"`
	Name     string `mapstructure:"name"`
}

// Client config
type Client struct {
	BindAddr   string `mapstructure:"bind_addr"`
	ServerAddr string `mapstructure:"server_addr"`
}

// Metrics config
type Metrics struct {
	BindAddr string `mapstructure:"bind_addr"`
}

// Tracing config for OTLP
type OTLPTracing struct {
	CollectorAddr string `mapstructure:"collector_url"`
	Insecure      bool   `mapstructure:"insecure"`
}

// Tracing config
type Tracing struct {
	OLTPTracing OTLPTracing `mapstructure:"otlp"`
}

// Config struct
type Config struct {
	Server  Server  `mapstructure:"server"`
	Client  Client  `mapstructure:"client"`
	Metrics Metrics `mapstructure:"metrics"`
	Tracing Tracing `mapstructure:"tracing"`
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath("config")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file to Config struct
func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
