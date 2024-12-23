package log

import (
	"context"
	"fmt"
	"sync"

	"github.com/minhthong582000/soa-404/pkg/config"
)

var (
	globalLogger Logger
	mutex        sync.Mutex
)

func GetLogger() Logger {
	if globalLogger == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if globalLogger == nil {
			fmt.Println("Initialize with default zap logger")
			globalLogger = NewTmpLogger()
		}
	}

	return globalLogger
}

func SetLogger(tracer Logger) {
	mutex.Lock()
	defer mutex.Unlock()

	globalLogger = tracer
}

type Provider string

const (
	Zap Provider = "Zap"
)

// Logger is a logger that supports log levels, context and structured logging.
type Logger interface {
	// With returns a logger based off the root logger and decorates it with the given context and arguments.
	With(ctx context.Context, args ...interface{}) Logger
	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...interface{})
	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
}

func Init(config *config.Logs) Logger {
	logger := New(config)

	SetLogger(logger)
	return GetLogger()
}
