package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

// OTLPTracer is a tracer that uses OpenTelemetry to send traces to an OTLP collector
type OTLPTracer struct {
	tracer trace.Tracer
	config *TracerConfig
}

func NewOTLPTracer(config *TracerConfig) *OTLPTracer {
	return &OTLPTracer{
		config: config,
	}
}

func (t *OTLPTracer) InitTracer() error {
	ctx := context.Background()

	if !t.config.Enabled {
		return nil
	}

	// Setup GRPC connection to collector
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if t.config.Insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(t.config.CollectorURL),
		),
	)
	if err != nil {
		return err
	}

	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", t.config.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	t.tracer = otel.Tracer(t.config.ServiceName)

	return nil
}

func (t *OTLPTracer) StartSpan(ctx context.Context, name string) context.Context {
	if !t.config.Enabled {
		return ctx
	}

	ctx, _ = t.tracer.Start(ctx, name)
	return ctx
}

func (t *OTLPTracer) EndSpan(ctx context.Context) {
	if !t.config.Enabled {
		return
	}

	span := trace.SpanFromContext(ctx)

	if span != nil {
		span.End()
	}
}

func (t *OTLPTracer) GetTraceID(ctx context.Context) string {
	if !t.config.Enabled {
		return ""
	}

	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sc := span.SpanContext()
	return sc.TraceID().String()
}

func (t *OTLPTracer) GetSpanID(ctx context.Context) string {
	if !t.config.Enabled {
		return ""
	}

	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sc := span.SpanContext()
	return sc.SpanID().String()
}
