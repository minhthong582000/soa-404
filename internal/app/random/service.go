package random

import (
	"context"
	"errors"

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
	// Validate seed
	if seed < 2 {
		return nil, errors.New("validate: seed must be greater than 2")
	}

	randNum, err := s.Repo.Get(ctx, seed)
	if err != nil {
		return nil, err
	}

	return &randNum, nil
}
