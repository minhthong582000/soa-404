package random

import (
	"context"

	"github.com/minhthong582000/soa-404/internal/entity"
)

type RandomService struct {
	Repo entity.IRandomRepository
}

func NewService(repo entity.IRandomRepository) entity.IRandomService {
	return &RandomService{
		Repo: repo,
	}
}

func (s *RandomService) Get(ctx context.Context, seed int64) (*entity.Random, error) {
	randNum, err := s.Repo.Get(ctx, seed)
	if err != nil {
		return nil, err
	}

	return &randNum, nil
}
