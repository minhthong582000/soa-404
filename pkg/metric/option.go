package metric

// TracerConfig is the configuration for the tracer
type MetricsConfig struct {
	Provider                  Provider
	Metrics                   []*Metric
	PrometheusNativeHistogram bool
}

type Option func(f *MetricsConfig)

func WithProvider(provider Provider) Option {
	return func(f *MetricsConfig) {
		f.Provider = provider
	}
}

func WithMetrics(metrics ...*Metric) Option {
	return func(f *MetricsConfig) {
		f.Metrics = metrics
	}
}

func WithPrometheusNativeHistogram(prometheusNativeHistogram bool) Option {
	return func(f *MetricsConfig) {
		f.PrometheusNativeHistogram = prometheusNativeHistogram
	}
}
