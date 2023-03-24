package random

import (
	"context"

	"github.com/minhthong582000/soa-404/internal/entity"
)

type RandomService interface {
	Get(ctx context.Context, seed int64) (*Random, error)
}

type Random struct {
	entity.Random
}

type service struct {
	repo RandomRepository
}

func NewService(repo RandomRepository) RandomService {
	return &service{
		repo: repo,
	}
}

func (s *service) Get(ctx context.Context, seed int64) (*Random, error) {
	randNum, err := s.repo.Get(ctx, seed)
	if err != nil {
		return nil, err
	}

	return &Random{randNum}, nil
}
