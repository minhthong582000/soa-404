package log

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/metadata"
	"gopkg.in/natefinch/lumberjack.v2"

	logger_config "github.com/minhthong582000/soa-404/pkg/config"
	grpcHeader "github.com/minhthong582000/soa-404/pkg/grpc"
	"github.com/minhthong582000/soa-404/pkg/tracing"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

type zapLogger struct {
	*zap.SugaredLogger
	config *logger_config.Logs
}

// New creates a new logger using the default configuration.
func NewZapLogger(config *logger_config.Logs) *zapLogger {
	// Get the log level from the config
	level, exist := loggerLevelMap[config.Level]
	if !exist {
		level = zapcore.DebugLevel
	}

	var encoderCfg zapcore.EncoderConfig
	if !config.Development {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}
	encoderCfg.TimeKey = "ts"
	encoderCfg.LevelKey = "level"
	encoderCfg.NameKey = "logger"
	encoderCfg.CallerKey = "caller"
	encoderCfg.FunctionKey = zapcore.OmitKey
	encoderCfg.MessageKey = "msg"
	encoderCfg.StacktraceKey = "stacktrace"
	encoderCfg.LineEnding = zapcore.DefaultLineEnding
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderCfg.EncodeTime = zapcore.EpochTimeEncoder
	encoderCfg.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	var writeSyncer zapcore.WriteSyncer
	if config.Path != "" {
		writeSyncer = fileWriteSyncer(config.Path)
	} else {
		writeSyncer = zapcore.AddSync(os.Stderr)
	}

	encoder := zapcore.NewJSONEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, writeSyncer, zap.NewAtomicLevelAt(level))
	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))

	zap.ReplaceGlobals(l)

	return &zapLogger{
		l.Sugar(),
		config,
	}
}

func fileWriteSyncer(path string) zapcore.WriteSyncer {
	// Create the directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(fmt.Sprintf("Failed to create log directory: %v", err))
		}
	}
	file := fmt.Sprintf("%s/%s.log", path, "app")

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   file,
		MaxSize:    5, // MB
		MaxBackups: 3,
	})

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stderr),
		zapcore.AddSync(w),
	)
}

// NewTmpLogger creates a temporary logger if the logger is not yet initialized.
func NewTmpLogger() *zapLogger {
	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	l, err := c.Build()
	if err != nil {
		panic(err)
	}

	return &zapLogger{
		l.Sugar(),
		&logger_config.Logs{},
	}
}

// NewForTest returns a new logger and the corresponding observed logs which can be used in unit tests to verify log entries.
func NewForTest(config *logger_config.Logs) (*zapLogger, *observer.ObservedLogs) {
	if config == nil {
		config = &logger_config.Logs{}
	}

	core, recorded := observer.New(zapcore.InfoLevel)
	l := zap.New(core)

	return &zapLogger{
		l.Sugar(),
		config,
	}, recorded
}

func (l *zapLogger) With(ctx context.Context, args ...interface{}) Logger {
	if ctx == nil {
		return l
	}

	tracer := tracing.GetTracer()

	requestID := grpcHeader.GetRequestIDFromContext(ctx)
	if requestID != "" {
		args = append(args, zap.String("request_id", requestID))
	}
	traceID := tracer.GetTraceID(ctx)
	if traceID != "" {
		args = append(args, zap.String("trace_id", traceID))
	}
	spanID := tracer.GetSpanID(ctx)
	if spanID != "" {
		args = append(args, zap.String("span_id", spanID))
	}
	if len(l.config.AdditionalFields) > 0 {
		if headers, ok := metadata.FromIncomingContext(ctx); ok {
			for _, header := range l.config.AdditionalFields {
				if value := headers.Get(header.ValueFrom); len(value) > 0 {
					args = append(args, zap.String(header.FieldName, value[0]))
				}
			}
		}
	}

	if len(args) > 0 {
		return &zapLogger{l.SugaredLogger.With(args...), l.config}
	}

	return l
}
