package random

import (
	"context"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/entity"
	"github.com/minhthong582000/soa-404/pkg/log"
	"go.opentelemetry.io/otel/trace"
)

type RandomServer struct {
	pb.UnimplementedRandomServiceServer
	logger log.ILogger
	tracer trace.Tracer

	RandomService entity.IRandomService
}

func NewServer(logger log.ILogger, tracer trace.Tracer, randomService entity.IRandomService) *RandomServer {
	return &RandomServer{
		RandomService: randomService,
		logger:        logger,
		tracer:        tracer,
	}
}

func (s RandomServer) GetRandNumber(ctx context.Context, request *pb.GetRandNumberRequest) (*pb.GetRandNumberReply, error) {
	ctx, span := s.tracer.Start(ctx, "RandomService.Usecase.GetRandNumber")
	defer span.End()

	randNum, err := s.RandomService.Get(ctx, request.SeedNum)
	if err != nil {
		return nil, err
	}

	return &pb.GetRandNumberReply{
		Number: randNum.Number,
	}, nil
}
