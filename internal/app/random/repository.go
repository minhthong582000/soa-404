package random

import (
	"context"
	"math/rand"

	"github.com/minhthong582000/soa-404/internal/entity"
	"github.com/minhthong582000/soa-404/pkg/tracing"
)

type RandomRepo struct {
}

func NewRepository() *RandomRepo {
	return &RandomRepo{}
}

func (r *RandomRepo) Get(ctx context.Context, seed int64) (entity.Random, error) {
	tracer := tracing.GetTracer()
	ctx = tracer.StartSpan(ctx, "RandomService.Repository.GetRandNumber")
	defer tracer.EndSpan(ctx)

	rand := rand.New(rand.NewSource(seed))
	randNum := rand.Int63()

	return entity.Random{
		Number: randNum,
	}, nil
}
