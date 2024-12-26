package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	grpcUtils "github.com/minhthong582000/soa-404/pkg/grpc"
	"github.com/minhthong582000/soa-404/pkg/grpc_errors"
	"github.com/minhthong582000/soa-404/pkg/log"
	"github.com/minhthong582000/soa-404/pkg/metric"
)

// Interceptor
type Interceptor struct {
}

// InterceptorManager constructor
func NewInterceptor() *Interceptor {
	return &Interceptor{}
}

// Logger Interceptor
func (im *Interceptor) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	logger := log.GetLogger()

	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	if err != nil {
		logger.With(
			ctx,
			"Method", info.FullMethod,
			"Time", time.Since(start),
			"Metadata", md,
		).Error(err)
	} else {
		logger.With(
			ctx,
			"Method", info.FullMethod,
			"Time", time.Since(start),
			"Metadata", md,
		).Infof("Success")
	}

	return reply, err
}
func (im *Interceptor) Metrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metr := metric.GetMetric()

	startTime := time.Now()
	serviceName, methodName := grpcUtils.SplitMethodName(info.FullMethod)

	// Post Message Received
	if metr.IsMetricExist(metric.Grpc_server_msg_received_total.Name) {
		metr.Counter(metric.Grpc_server_msg_received_total, 1, string(metric.Unary), serviceName, methodName)
	}

	// Call
	resp, err := handler(ctx, req)

	// Post Message Sent
	if metr.IsMetricExist(metric.Grpc_server_msg_sent_total.Name) {
		metr.Counter(metric.Grpc_server_msg_sent_total, 1, string(metric.Unary), serviceName, methodName)
	}

	// Post Call
	status := http.StatusOK
	if err != nil {
		status = grpc_errors.MapGRPCErrCodeToHttpStatus(grpc_errors.ParseGRPCErrStatusCode(err))
	}
	statusStr := strconv.Itoa(status)
	if metr.IsMetricExist(metric.Grpc_server_handled_total.Name) {
		metr.Counter(metric.Grpc_server_handled_total, 1, string(metric.Unary), serviceName, methodName, statusStr)
	}
	if metr.IsMetricExist(metric.Grpc_server_handling_seconds.Name) {
		metr.Histogram(metric.Grpc_server_handling_seconds, time.Since(startTime).Seconds(), string(metric.Unary), serviceName, methodName)
	}

	return resp, err
}
