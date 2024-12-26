package metric

import (
	"context"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/minhthong582000/soa-404/pkg/log"
)

type prometheusMetrics struct {
	metricsMap map[string]prometheus.Collector
	registry   *prometheus.Registry
	config     *MetricsConfig
}

func NewPrometheusMetrics(config *MetricsConfig) (*prometheusMetrics, error) {
	p := &prometheusMetrics{
		metricsMap: make(map[string]prometheus.Collector),
		registry:   prometheus.NewRegistry(),
		config:     config,
	}

	for _, m := range config.Metrics {
		_, ok := p.metricsMap[m.Name]
		if ok {
			return nil, fmt.Errorf("duplicate metric name: \"%s\"", m.Name)
		}

		var prometheusMetric prometheus.Collector

		switch m.Type {
		case Counter:
			prometheusMetric = promauto.NewCounterVec(prometheus.CounterOpts{
				Name:      m.Name,
				Help:      m.Description,
				Subsystem: string(m.Subsystem),
			}, m.Labels)
		case Gauge:
			prometheusMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name:      m.Name,
				Help:      m.Description,
				Subsystem: string(m.Subsystem),
			}, m.Labels)
		case Histogram:
			if len(m.Buckets) == 0 {
				return nil, errors.New("empty histogram bucket")
			}
			prometheusMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
				Name:      m.Name,
				Subsystem: string(m.Subsystem),
				Help:      m.Description,
				Buckets:   m.Buckets,
			}, m.Labels)
		default:
			return nil, errors.New("invalid metric type")
		}

		if err := p.registry.Register(prometheusMetric); err != nil {
			return nil, err
		}
		p.metricsMap[m.Name] = prometheusMetric
	}

	// Collect build info
	if err := p.registry.Register(collectors.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	return p, nil
}

func NewTmpPrometheusMetrics() *prometheusMetrics {
	p := &prometheusMetrics{
		metricsMap: make(map[string]prometheus.Collector),
		registry:   prometheus.NewRegistry(),
		config:     &MetricsConfig{},
	}
	p.registry.Register(collectors.NewBuildInfoCollector())
	return p
}

func (p *prometheusMetrics) RunHTTPMetricsServer(ctx context.Context, address string) {
	logger := log.GetLogger()

	router := echo.New()
	router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	if err := router.Start(address); err != nil {
		logger.Errorf("Error starting metrics server: %s", err)
	}
}

func (p *prometheusMetrics) IsMetricExist(name string) bool {
	_, ok := p.metricsMap[name]
	return ok
}

func (p *prometheusMetrics) Counter(metric *Metric, value float64, labelValues ...string) error {
	collector, ok := p.metricsMap[metric.Name]
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

func (p *prometheusMetrics) gauge(metric *Metric, labelValues ...string) (prometheus.Gauge, error) {
	collector, ok := p.metricsMap[metric.Name]
	if !ok {
		return nil, errors.New("metric does not exist")
	}

	// Check if the metric is a gauge
	gaugeVec, ok := collector.(*prometheus.GaugeVec)
	if !ok {
		return nil, errors.New("metric is not a gauge")
	}
	gauge, err := gaugeVec.GetMetricWithLabelValues(labelValues...)
	if err != nil {
		return nil, err
	}

	return gauge, nil
}

func (p *prometheusMetrics) AddGauge(metric *Metric, value float64, labelValues ...string) error {
	gauge, err := p.gauge(metric, labelValues...)
	if err != nil {
		return err
	}
	gauge.Add(value)

	return nil
}

func (p *prometheusMetrics) SetGauge(metric *Metric, value float64, labelValues ...string) error {
	gauge, err := p.gauge(metric, labelValues...)
	if err != nil {
		return err
	}
	gauge.Set(value)

	return nil
}

func (p *prometheusMetrics) Histogram(metric *Metric, value float64, labelValues ...string) error {
	collector, ok := p.metricsMap[metric.Name]
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
