package metric

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Subsystem string

const (
	HTTP Subsystem = "http"
	GRPC Subsystem = "grpc"
)

type Metric struct {
	Name        string
	Type        MetricType
	Subsystem   Subsystem
	Description string
	Labels      []string
	Buckets     []float64
}

type Metrics interface {
	RunHTTPMetricsServer(ctx context.Context, address string)
	Counter(metric *Metric, value float64, labelValues ...string) error
	AddGauge(metric *Metric, value float64, labelValues ...string) error
	SetGauge(metric *Metric, value float64, labelValues ...string) error
	Histogram(metric *Metric, value float64, labelValues ...string) error
	IsMetricExist(name string) bool
}

var (
	globalMetric Metrics
	mutex        sync.Mutex
)

func GetMetric() Metrics {
	if globalMetric == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if globalMetric == nil {
			fmt.Println("Initialize with default prometheus metric")
			globalMetric = NewTmpPrometheusMetrics()
		}
	}

	return globalMetric
}

func SetMetric(metric Metrics) {
	mutex.Lock()
	defer mutex.Unlock()

	globalMetric = metric
}

type Provider string

const (
	Datadog    Provider = "DATADOG"
	Prometheus Provider = "PROMETHEUS"
)

type MetricType string

const (
	Counter   MetricType = "COUNTER"
	Gauge     MetricType = "GAUGE"
	Histogram MetricType = "HISTOGRAM"
)

func MetricFactory(opts ...Option) (Metrics, error) {
	config := &MetricsConfig{}
	for _, opt := range opts {
		opt(config)
	}

	var (
		metrics Metrics
		err     error
	)
	switch config.Provider {
	case Prometheus:
		metrics, err = NewPrometheusMetrics(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create Prometheus metrics: %w", err)
		}
	case Datadog:
		return nil, errors.New("datadog not implemented yet")
	default:
		return nil, errors.New("provider not implemented yet")
	}

	SetMetric(metrics)

	return GetMetric(), nil
}
