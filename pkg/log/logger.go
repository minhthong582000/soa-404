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
	// With returns a logger based on the root logger and decorates it with the given context and arguments.
	With(ctx context.Context, args ...interface{}) Logger
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func Init(config *config.Logs) Logger {
	logger := New(config)

	SetLogger(logger)
	return GetLogger()
}
