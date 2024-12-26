package metric

// List of metric labels
const (
	Method string = "method"
	Status string = "status"
	Path   string = "path"
)

// GrpcType is the type of RPC call.
type grpcType string

const (
	Unary        grpcType = "unary"
	ClientStream grpcType = "client_stream"
	ServerStream grpcType = "server_stream"
	BidiStream   grpcType = "bidi_stream"
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
	Labels:      []string{Path},
	Buckets:     []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
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
	Type:        Gauge,
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

// grpc_server_handling_seconds is a histogram metric that measures the duration of the request in seconds.
var Grpc_server_handling_seconds *Metric = &Metric{
	Name:        "server_handling_seconds",
	Description: "Histogram metric that measures the duration of the request in seconds.",
	Subsystem:   GRPC,
	Type:        Histogram,
	Labels:      []string{"grpc_type", "grpc_service", "grpc_method"},
	Buckets:     []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
}

// grpc_server_handled_total is a counter metric that measures the total number of requests.
var Grpc_server_handled_total *Metric = &Metric{
	Name:        "server_handled_total",
	Description: "Total number of RPCs completed on the server, regardless of success or failure.",
	Subsystem:   GRPC,
	Type:        Counter,
	Labels:      []string{"grpc_type", "grpc_service", "grpc_method", Status},
}

// grpc_server_started_total is a gauge metric that measures the number of requests currently in progress.
var Grpc_server_started_total *Metric = &Metric{
	Name:        "server_started_total",
	Description: "Total number of RPCs started on the server.",
	Subsystem:   GRPC,
	Type:        Counter,
	Labels:      []string{"grpc_type", "grpc_service", "grpc_method"},
}

// grpc_server_msg_received_total is a counter metric that measures the total number of RPC stream messages received on the server.
var Grpc_server_msg_received_total *Metric = &Metric{
	Name:        "server_msg_received_total",
	Description: "Total number of RPC stream messages received on the server.",
	Subsystem:   GRPC,
	Type:        Counter,
	Labels:      []string{"grpc_type", "grpc_service", "grpc_method"},
}

// grpc_server_msg_sent_total is a counter metric that measures the total number of gRPC stream messages sent by the server.
var Grpc_server_msg_sent_total *Metric = &Metric{
	Name:        "server_msg_sent_total",
	Description: "Total number of gRPC stream messages sent by the server.",
	Subsystem:   GRPC,
	Type:        Counter,
	Labels:      []string{"grpc_type", "grpc_service", "grpc_method"},
}
