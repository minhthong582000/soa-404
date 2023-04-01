package tracing

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Enum of tracing types
const (
	OLTP = iota
	Jaeger
)

// Tracer is the interface for the tracer
type Tracer interface {
	InitTracer() (func(context.Context) error, error)
}

// TracerConfig is the configuration for the tracer
type TracerConfig struct {
	ServiceName  string
	CollectorURL string
	Insecure     bool
}

// OTLPTracer is a tracer that uses OpenTelemetry to send traces to an OTLP collector
type OLTPTracer struct {
	config TracerConfig
}

func NewOLTPTracer(config TracerConfig) Tracer {
	return &OLTPTracer{
		config: config,
	}
}

func (t OLTPTracer) InitTracer() (func(context.Context) error, error) {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if t.config.Insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(t.config.CollectorURL),
		),
	)
	if err != nil {
		return nil, err
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", t.config.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return exporter.Shutdown, nil
}

// TracerFactory returns a tracer based on the type
func TracerFactory(TracerType int, config TracerConfig) (Tracer, error) {
	switch TracerType {
	case OLTP:
		return NewOLTPTracer(config), nil
	case Jaeger:
		return nil, errors.New("jaeger not implemented yet")
	default:
		return nil, errors.New("invalid tracer type")
	}
}
