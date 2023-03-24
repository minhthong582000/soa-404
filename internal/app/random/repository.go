package random

import (
	"context"
	"math/rand"

	"github.com/minhthong582000/soa-404/internal/entity"
)

// RandomRepository
type RandomRepository interface {
	Get(ctx context.Context, seed int64) (entity.Random, error)
}

type repo struct{}

func NewRepository() RandomRepository {
	return &repo{}
}

func (r *repo) Get(ctx context.Context, seed int64) (entity.Random, error) {
	rand := rand.New(rand.NewSource(seed))
	randNum := rand.Int63()

	return entity.Random{
		Number: randNum,
	}, nil
}
