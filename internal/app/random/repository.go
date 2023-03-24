package random

import (
	"context"
	"math/rand"

	"github.com/minhthong582000/soa-404/internal/entity"
)

type RandomRepo struct{}

func NewRepository() entity.IRandomRepository {
	return &RandomRepo{}
}

func (r *RandomRepo) Get(ctx context.Context, seed int64) (entity.Random, error) {
	rand := rand.New(rand.NewSource(seed))
	randNum := rand.Int63()

	return entity.Random{
		Number: randNum,
	}, nil
}
