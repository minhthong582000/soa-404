package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

var (
	RequestIDHeader = "x-request-id"
	ClientIPHeader  = "x-client-ip"
)

func GetRequestIDFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	requestIds := md.Get(RequestIDHeader)
	if len(requestIds) == 0 {
		return ""
	}

	return requestIds[0]
}
