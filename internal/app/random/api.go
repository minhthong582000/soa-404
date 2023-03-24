package random

import (
	"context"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
)

type RandomServer struct {
	pb.UnimplementedRandomServiceServer

	randomService RandomService
}

func NewServer(randomService RandomService) *RandomServer {
	return &RandomServer{
		randomService: randomService,
	}
}

func (s RandomServer) GetRandNumber(ctx context.Context, request *pb.GetRandNumberRequest) (*pb.GetRandNumberReply, error) {
	randNum, err := s.randomService.Get(ctx, request.SeedNum)
	if err != nil {
		return nil, err
	}

	return &pb.GetRandNumberReply{
		Number: randNum.Number,
	}, nil
}
