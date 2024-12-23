package middleware

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/minhthong582000/soa-404/pkg/grpc_errors"
	"github.com/minhthong582000/soa-404/pkg/log"
	"github.com/minhthong582000/soa-404/pkg/metrics"
)

// Interceptor
type Interceptor struct {
	metr metrics.Metrics
}

// InterceptorManager constructor
func NewInterceptor(metr metrics.Metrics) *Interceptor {
	return &Interceptor{metr: metr}
}

// Logger Interceptor
func (im *Interceptor) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	logger := log.GetLogger()

	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	if err != nil {
		logger.With(ctx).Errorf("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)
	} else {
		logger.With(ctx).Infof("Method: %s, Time: %v, Metadata: %v", info.FullMethod, time.Since(start), md)
	}

	return reply, err
}

func (im *Interceptor) Metrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	var status = http.StatusOK
	if err != nil {
		status = grpc_errors.MapGRPCErrCodeToHttpStatus(grpc_errors.ParseGRPCErrStatusCode(err))
	}
	im.metr.ObserveResponseTime(status, info.FullMethod, info.FullMethod, time.Since(start).Seconds())
	im.metr.IncHits(status, info.FullMethod, info.FullMethod)

	return resp, err
}
