package context_keys

type ContextKey int

const (
	RequestIDKey ContextKey = iota
	CorrelationIDKey
)
