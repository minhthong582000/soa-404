package random

import (
	"context"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
	"github.com/minhthong582000/soa-404/internal/entity"
	"github.com/minhthong582000/soa-404/pkg/log"
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
	s.logger.With(ctx).Info("GetRandNumber")

	randNum, err := s.RandomService.Get(ctx, request.SeedNum)
	if err != nil {
		s.logger.Errorf("invalid request: %v", err)
		return nil, err
	}

	return &pb.GetRandNumberReply{
		Number: randNum.Number,
	}, nil
}
