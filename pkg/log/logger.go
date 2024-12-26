package log

import (
	"context"
	"fmt"
	"sync"

	logger_config "github.com/minhthong582000/soa-404/pkg/config"
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

func LogFactory(config *logger_config.Logs) (Logger, error) {
	var (
		logger Logger
	)
	switch config.Provider {
	case logger_config.ZapLog:
		logger = NewZapLogger(config)
	default:
		return nil, fmt.Errorf("unsupported logger provider: %s", config.Provider)
	}

	SetLogger(logger)

	return GetLogger(), nil
}
