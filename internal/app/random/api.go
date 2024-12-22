package random

import (
	"context"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/entity"
	"github.com/minhthong582000/soa-404/pkg/log"
	"github.com/minhthong582000/soa-404/pkg/tracing"
)

type RandomServer struct {
	pb.UnimplementedRandomServiceServer
	logger log.ILogger

	RandomService entity.IRandomService
}

func NewServer(logger log.ILogger, randomService entity.IRandomService) *RandomServer {
	return &RandomServer{
		RandomService: randomService,
		logger:        logger,
	}
}

func (s RandomServer) GetRandNumber(ctx context.Context, request *pb.GetRandNumberRequest) (*pb.GetRandNumberReply, error) {
	tracer := tracing.GetCurrenTracer()
	ctx = tracer.StartSpan(ctx, "RandomService.Handler.GetRandNumber")
	defer tracer.EndSpan(ctx)

	randNum, err := s.RandomService.Get(ctx, request.SeedNum)
	if err != nil {
		return nil, err
	}

	return &pb.GetRandNumberReply{
		Number: randNum.Number,
	}, nil
}
