package random

import (
	"context"
	"errors"

	"github.com/minhthong582000/soa-404/internal/entity"
	"github.com/minhthong582000/soa-404/pkg/tracing"
)

type RandomService struct {
	repo entity.IRandomRepository
}

func NewService(repo entity.IRandomRepository) entity.IRandomService {
	return &RandomService{
		repo: repo,
	}
}

func (s *RandomService) Get(ctx context.Context, seed int64) (*entity.Random, error) {
	tracer := tracing.GetCurrenTracer()
	ctx = tracer.StartSpan(ctx, "RandomService.Usecase.GetRandNumber")
	defer tracer.EndSpan(ctx)

	// Validate seed
	if seed < 2 {
		return nil, errors.New("validate: seed must be greater than 2")
	}

	randNum, err := s.repo.Get(ctx, seed)
	if err != nil {
		return nil, err
	}

	return &randNum, nil
}
