package random

import (
	"context"
	"errors"

	"github.com/minhthong582000/soa-404/internal/entity"
	"go.opentelemetry.io/otel/trace"
)

type RandomService struct {
	repo entity.IRandomRepository

	tracer trace.Tracer
}

func NewService(tracer trace.Tracer, repo entity.IRandomRepository) entity.IRandomService {
	return &RandomService{
		repo:   repo,
		tracer: tracer,
	}
}

func (s *RandomService) Get(ctx context.Context, seed int64) (*entity.Random, error) {
	ctx, span := s.tracer.Start(ctx, "RandomService.Usecase.GetRandNumber")
	defer span.End()

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
