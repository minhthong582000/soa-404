package tracing

import (
	"context"
	"errors"
	"sync"
)

// Tracer is the interface for the tracer
type Tracer interface {
	InitTracer() error
	StartSpan(ctx context.Context, name string) context.Context
	EndSpan(ctx context.Context)
}

var (
	globalTracer Tracer
	rwMutex      sync.RWMutex
)

func GetCurrenTracer() Tracer {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	return globalTracer
}

func SetglobalTracer(tracer Tracer) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	globalTracer = tracer
}

type Provider string

const (
	Jaeger  Provider = "JEAGER"
	Datadog Provider = "DATADOG"
	OTLP    Provider = "OTLP"
)

// TracerFactory returns a tracer based on the type
func TracerFactory(opts ...Option) (Tracer, error) {
	config := &TracerConfig{}
	for _, opt := range opts {
		opt(config)
	}

	var tracer Tracer
	switch config.Provider {
	case OTLP:
		tracer = NewOTLPTracer(config)
	case Jaeger:
		return nil, errors.New("jaeger not implemented yet")
	case Datadog:
		return nil, errors.New("datadog not implemented yet")
	default:
		return nil, errors.New("invalid tracer type")
	}

	err := tracer.InitTracer()
	if err != nil {
		return nil, err
	}

	SetglobalTracer(tracer)

	return GetCurrenTracer(), nil
}
