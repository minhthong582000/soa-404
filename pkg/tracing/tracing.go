package tracing

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Tracer is the interface for the tracer
type Tracer interface {
	StartSpan(ctx context.Context, name string) context.Context
	EndSpan(ctx context.Context)
	GetTraceID(ctx context.Context) string
	GetSpanID(ctx context.Context) string
}

var (
	globalTracer Tracer
	mutex        sync.Mutex
)

func GetTracer() Tracer {
	if globalTracer == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if globalTracer == nil {
			fmt.Println("Initialize with default empty tracer")
			globalTracer = NewTmpOLTPTracer()
		}
	}

	return globalTracer
}

func SetTracer(tracer Tracer) {
	mutex.Lock()
	defer mutex.Unlock()

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

	var (
		tracer Tracer
		err    error
	)
	switch config.Provider {
	case OTLP:
		tracer, err = NewOTLPTracer(config)
		if err != nil {
			return nil, err
		}
	case Jaeger:
		return nil, errors.New("jaeger not implemented yet")
	case Datadog:
		return nil, errors.New("datadog not implemented yet")
	default:
		return nil, errors.New("invalid tracer type")
	}

	SetTracer(tracer)

	return GetTracer(), nil
}
