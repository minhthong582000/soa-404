package tracing

import (
	"context"
	"errors"
)

type Provider string

const (
	Jaeger  Provider = "JEAGER"
	Datadog Provider = "DATADOG"
	OTLP    Provider = "OTLP"
)

// Tracer is the interface for the tracer
type Tracer interface {
	InitTracer() (func(context.Context) error, error)
}

// TracerFactory returns a tracer based on the type
func TracerFactory(opts ...Option) (Tracer, error) {
	config := &TracerConfig{}
	for _, opt := range opts {
		opt(config)
	}

	switch config.Provider {
	case OTLP:
		return NewOTLPTracer(config), nil
	case Jaeger:
		return nil, errors.New("jaeger not implemented yet")
	case Datadog:
		return nil, errors.New("datadog not implemented yet")
	default:
		return nil, errors.New("invalid tracer type")
	}
}
