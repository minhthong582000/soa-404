package tracing

// TracerConfig is the configuration for the tracer
type TracerConfig struct {
	Provider     Provider
	ServiceName  string
	CollectorURL string
	Insecure     bool
	Enabled      bool
}

type Option func(f *TracerConfig)

func WithProvider(provider Provider) Option {
	return func(f *TracerConfig) {
		f.Provider = provider
	}
}

func WithServiceName(serviceName string) Option {
	return func(f *TracerConfig) {
		f.ServiceName = serviceName
	}
}

func WithCollectorURL(CollectorURL string) Option {
	return func(f *TracerConfig) {
		f.CollectorURL = CollectorURL
	}
}

func WithInsecure(isInsecure bool) Option {
	return func(f *TracerConfig) {
		f.Insecure = isInsecure
	}
}

func WithEnabled(isEnabled bool) Option {
	return func(f *TracerConfig) {
		f.Enabled = isEnabled
	}
}
