package mock

import (
	_ "go.uber.org/mock/mockgen/model"
)

//go:generate mockgen -destination=mock_tracing.go -package=mock github.com/minhthong582000/soa-404/pkg/tracing Tracer
