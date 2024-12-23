package config

import (
	"errors"

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
	Name       string `mapstructure:"name" validate:"required"`
}

// Logger config
type Logs struct {
	Development      bool              `mapstructure:"development"`
	Level            string            `mapstructure:"level" validate:"required,oneof=debug info warn error dpanic panic fatal"`
	Path             string            `mapstructure:"path"`
	AdditionalFields []AdditionalField `mapstructure:"additional_fields" validate:"dive"`
}

type AdditionalField struct {
	FieldName string `mapstructure:"field_name" validate:"required"`
	ValueFrom string `mapstructure:"value_from" validate:"required"`
}

// Metrics config
type Metrics struct {
	BindAddr string `mapstructure:"bind_addr" validate:"required"`
}

// Tracing config for OTLP
type OTLPTracing struct {
	Enabled       bool   `mapstructure:"enabled"`
	CollectorAddr string `mapstructure:"collector_url" validate:"required"`
	Insecure      bool   `mapstructure:"insecure"`
}

// Tracing config
type Tracing struct {
	OLTPTracing OTLPTracing `mapstructure:"otlp" validate:"required"`
}

// Config struct
type Config struct {
	Server  Server  `mapstructure:"server" validate:"required"`
	Client  Client  `mapstructure:"client" validate:"required"`
	Logs    Logs    `mapstructure:"logs" validate:"required"`
	Metrics Metrics `mapstructure:"metrics" validate:"required"`
	Tracing Tracing `mapstructure:"tracing" validate:"required"`
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")  // optionally look for config in the working directory
	v.SetEnvPrefix("SOA") // set env prefix with "SOA_", e.g. SOA_SERVER_BIND_ADDR
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	// TODO: If yaml value is ${ENV}, replace the value with ENV value
	// e.g. key: ${VALUE} -> key: ABC, if `VALUE` is set to `ABC`

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
