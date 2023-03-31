package config

import (
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Server config
type Server struct {
	BindAddr string `mapstructure:"bind_addr" validate:"required"`
	Name     string `mapstructure:"name" validate:"required"`
}

// Client config
type Client struct {
	BindAddr   string `mapstructure:"bind_addr" validate:"required"`
	ServerAddr string `mapstructure:"server_addr" validate:"required"`
}

// Metrics config
type Metrics struct {
	BindAddr string `mapstructure:"bind_addr" validate:"required"`
}

// Tracing config for OTLP
type OTLPTracing struct {
	CollectorAddr string `mapstructure:"collector_url" validate:"required"`
	Insecure      bool   `mapstructure:"insecure" validate:"required"`
}

// Tracing config
type Tracing struct {
	OLTPTracing OTLPTracing `mapstructure:"otlp" validate:"required"`
}

// Config struct
type Config struct {
	Server  Server  `mapstructure:"server" validate:"required"`
	Client  Client  `mapstructure:"client" validate:"required"`
	Metrics Metrics `mapstructure:"metrics" validate:"required"`
	Tracing Tracing `mapstructure:"tracing" validate:"required"`
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")  // optionally look for config in the working directory
	v.SetEnvPrefix("soa") // set env prefix with "SOA_", e.g. SOA_SERVER_BIND_ADDR
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	// If yaml value is ${ENV}, replace the value with ENV value
	// e.g. key: ${VALUE} -> key: ABC, if `VALUE` is set to `ABC`
	for _, k := range v.AllKeys() {
		val := v.GetString(k)
		v.Set(k, os.ExpandEnv(val))
	}

	return v, nil
}

// Parse config file to Config struct
func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
