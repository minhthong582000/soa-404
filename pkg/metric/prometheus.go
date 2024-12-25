package metric

import (
	"context"
	"errors"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/minhthong582000/soa-404/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusMetrics struct {
	MetricsMap map[string]prometheus.Collector
	registry   *prometheus.Registry
	mutex      *sync.RWMutex
}

func NewPrometheusMetrics() *prometheusMetrics {
	return &prometheusMetrics{
		MetricsMap: make(map[string]prometheus.Collector),
		registry:   prometheus.NewRegistry(),
		mutex:      &sync.RWMutex{},
	}
}

func (p *prometheusMetrics) RunHTTPMetricsServer(ctx context.Context, address string) {
	logger := log.GetLogger()

	router := echo.New()
	router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	if err := router.Start(address); err != nil {
		logger.Errorf("Error starting metrics server: %s", err)
	}
}

func (p *prometheusMetrics) InitMetric(metric ...*Metric) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, m := range metric {
		_, ok := p.MetricsMap[m.Name]
		if ok {
			return errors.New("metrics with the same name already registered")
		}

		var prometheusMetric prometheus.Collector

		switch m.Type {
		case Counter, Gauge:
			prometheusMetric = promauto.NewCounterVec(prometheus.CounterOpts{
				Name: m.Name,
			}, m.Labels)
		case Histogram:
			if len(m.Buckets) == 0 {
				return errors.New("empty histogram bucket")
			}
			prometheusMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
				Name:    m.Name,
				Buckets: m.Buckets,
			}, m.Labels)
		default:
			return errors.New("invalid metric type")
		}

		if err := p.registry.Register(prometheusMetric); err != nil {
			return err
		}
		p.MetricsMap[m.Name] = prometheusMetric
	}

	// Collect build info
	if err := p.registry.Register(collectors.NewBuildInfoCollector()); err != nil {
		return err
	}

	return nil
}

func (p *prometheusMetrics) Counter(metric *Metric, value float64, labelValues ...string) error {
	collector, ok := p.MetricsMap[metric.Name]
	if !ok {
		return errors.New("metric does not exist")
	}

	// Check if the metric is a counter
	counterVec, ok := collector.(*prometheus.CounterVec)
	if !ok {
		return errors.New("metric is not a counter")
	}
	counter, err := counterVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		return err
	}
	counter.Add(value)

	return nil
}

func (p *prometheusMetrics) Gauge(metric *Metric, value float64, labelValues ...string) error {
	collector, ok := p.MetricsMap[metric.Name]
	if !ok {
		return errors.New("metric does not exist")
	}

	// Check if the metric is a gauge
	gaugeVec, ok := collector.(*prometheus.GaugeVec)
	if !ok {
		return errors.New("metric is not a gauge")
	}
	gauge, err := gaugeVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		return err
	}
	gauge.Set(value)

	return nil
}

func (p *prometheusMetrics) Histogram(metric *Metric, value float64, labelValues ...string) error {
	collector, ok := p.MetricsMap[metric.Name]
	if !ok {
		return errors.New("metric does not exist")
	}

	// Check if the metric is a histogram
	histogramVec, ok := collector.(*prometheus.HistogramVec)
	if !ok {
		return errors.New("metric is not a histogram")
	}
	histogram, err := histogramVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		return err
	}
	histogram.Observe(value)

	return nil
}
