package metric

// List of metric labels
const (
	Method string = "method"
	Status string = "status"
	Path   string = "path"
)

// Size
const (
	_           = iota // ignore first value by assigning to blank identifier
	bKB float64 = 1 << (10 * iota)
	bMB
)

// sizeBuckets is the buckets for request/response size. Here we define a spectrum from 1KB through 1NB up to 10MB.
var sizeBuckets = []float64{1.0 * bKB, 2.0 * bKB, 5.0 * bKB, 10.0 * bKB, 100 * bKB, 500 * bKB, 1.0 * bMB, 2.5 * bMB, 5.0 * bMB, 10.0 * bMB}

//
// List of default Http metrics
//

// http_request_duration_seconds is a histogram metric that measures the duration of the request in seconds.
var Http_request_duration_seconds *Metric = &Metric{
	Name:        "request_duration_seconds",
	Subsystem:   HTTP,
	Description: "Histogram metric that measures the duration of the request in seconds.",
	Type:        Histogram,
	Labels:      []string{Path, Status},
	Buckets:     []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
}

// http_request_total is a counter metric that measures the total number of requests.
var Http_request_total *Metric = &Metric{
	Name:        "request_total",
	Description: "Counter metric that measures the total number of requests.",
	Subsystem:   HTTP,
	Type:        Counter,
	Labels:      []string{Path, Status},
}

// http_request_inflight is a gauge metric that measures the number of requests currently in progress.
var Http_request_inflight *Metric = &Metric{
	Name:        "request_inflight",
	Description: "Gauge metric that measures the number of requests currently in progress.",
	Subsystem:   HTTP,
	Type:        Counter,
	Labels:      []string{Path},
}

// http_response_size_bytes is a histogram metric that measures the size of the response in bytes.
var Http_response_size_bytes *Metric = &Metric{
	Name:        "response_size_bytes",
	Description: "Histogram metric that measures the size of the response in bytes.",
	Subsystem:   HTTP,
	Type:        Histogram,
	Labels:      []string{Path},
	Buckets:     sizeBuckets,
}

// http_request_size_bytes is a histogram metric that measures the size of the request in bytes.
var Http_request_size_bytes *Metric = &Metric{
	Name:        "request_size_bytes",
	Description: "Histogram metric that measures the size of the request in bytes.",
	Subsystem:   HTTP,
	Type:        Histogram,
	Labels:      []string{Path},
	Buckets:     sizeBuckets,
}

//
// List of default GRPC metrics
//

// grpc_request_duration_seconds is a histogram metric that measures the duration of the request in seconds.
var Grpc_request_duration_seconds *Metric = &Metric{
	Name:        "request_duration_seconds",
	Description: "Histogram metric that measures the duration of the request in seconds.",
	Subsystem:   GRPC,
	Type:        Histogram,
	Labels:      []string{Method, Status},
	Buckets:     []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
}

// grpc_request_total is a counter metric that measures the total number of requests.
var Grpc_request_total *Metric = &Metric{
	Name:        "request_total",
	Description: "Counter metric that measures the total number of requests.",
	Subsystem:   GRPC,
	Type:        Counter,
	Labels:      []string{Method, Status},
}

// grpc_request_inflight is a gauge metric that measures the number of requests currently in progress.
var Grpc_request_inflight *Metric = &Metric{
	Name:        "request_inflight",
	Description: "Gauge metric that measures the number of requests currently in progress.",
	Subsystem:   GRPC,
	Type:        Gauge,
	Labels:      []string{Method},
}
