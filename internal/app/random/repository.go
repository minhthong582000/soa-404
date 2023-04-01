package random

import (
	"context"
	"math/rand"

	"github.com/minhthong582000/soa-404/internal/entity"
	"go.opentelemetry.io/otel/trace"
)

type RandomRepo struct {
	tracer trace.Tracer
}

func NewRepository(tracer trace.Tracer) entity.IRandomRepository {
	return &RandomRepo{
		tracer: tracer,
	}
}

func (r *RandomRepo) Get(ctx context.Context, seed int64) (entity.Random, error) {
	_, span := r.tracer.Start(ctx, "RandomService.Repository.GetRandNumber")
	defer span.End()

	rand := rand.New(rand.NewSource(seed))
	randNum := rand.Int63()

	return entity.Random{
		Number: randNum,
	}, nil
}
