package random

import (
	"context"

	"github.com/bufbuild/protovalidate-go"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/entity"
	"github.com/minhthong582000/soa-404/pkg/tracing"
)

type RandomServer struct {
	pb.UnimplementedRandomServiceServer
	RandomService entity.IRandomService
}

func NewServer(randomService entity.IRandomService) *RandomServer {
	return &RandomServer{
		RandomService: randomService,
	}
}

func (s RandomServer) GetRandNumber(ctx context.Context, request *pb.GetRandNumberRequest) (*pb.GetRandNumberReply, error) {
	tracer := tracing.GetTracer()
	ctx = tracer.StartSpan(ctx, "RandomService.Handler.GetRandNumber")
	defer tracer.EndSpan(ctx)

	if err := protovalidate.Validate(request); err != nil {
		return nil, err
	}

	randNum, err := s.RandomService.Get(ctx, request.SeedNum)
	if err != nil {
		return nil, err
	}

	return &pb.GetRandNumberReply{
		Number: randNum.Number,
	}, nil
}
