package entity

import "context"

type Random struct {
	Number int64 `json:"number"`
}

//go:generate mockery --name IRandomRepository --output ../mocks/ --case underscore
type IRandomRepository interface {
	Get(ctx context.Context, seed int64) (Random, error)
}

type IRandomService interface {
	Get(ctx context.Context, seed int64) (*Random, error)
}
