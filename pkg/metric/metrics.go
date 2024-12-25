package metric

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Metric struct {
	Name    string
	Type    MetricType
	Labels  []string
	Buckets []float64
}

type Metrics interface {
	InitMetric(metric ...*Metric) error
	RunHTTPMetricsServer(ctx context.Context, address string)
	Counter(metric *Metric, value float64, labelValues ...string) error
	Gauge(metric *Metric, value float64, labelValues ...string) error
	Histogram(metric *Metric, value float64, labelValues ...string) error
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

	var metrics Metrics
	switch config.Provider {
	case Prometheus:
		metrics = NewPrometheusMetrics()
	case Datadog:
		return nil, errors.New("datadog not implemented yet")
	default:
		return nil, errors.New("provider not implemented yet")
	}

	if err := metrics.InitMetric(config.Metrics...); err != nil {
		return nil, err
	}

	SetMetric(metrics)

	return GetMetric(), nil
}
