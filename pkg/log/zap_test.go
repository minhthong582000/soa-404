package log

import (
	"context"
	"os"
	"testing"

	logger_config "github.com/minhthong582000/soa-404/pkg/config"
	"github.com/minhthong582000/soa-404/pkg/grpc"
	"github.com/minhthong582000/soa-404/pkg/tracing"
	mock_tracer "github.com/minhthong582000/soa-404/pkg/tracing/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

// TestNewZapLogger tests various configurations of the Zap logger.
func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name           string
		config         *logger_config.Logs
		expectFile     bool
		expectLogFile  string
		expectLogLevel string
	}{
		{
			name: "Default Config",
			config: &logger_config.Logs{
				Level:       "debug",
				Development: true,
				Path:        "",
			},
			expectFile:     false,
			expectLogLevel: "debug",
		},
		{
			name: "Invalid Level Config",
			config: &logger_config.Logs{
				Level: "invalid",
				Path:  "",
			},
			expectFile:     false,
			expectLogLevel: "debug",
		},
		{
			name: "File Output Config",
			config: &logger_config.Logs{
				Level:       "info",
				Development: false,
			},
			expectFile:     true,
			expectLogLevel: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle temporary directory for file output
			if tt.expectFile {
				tempDir := t.TempDir()
				tt.config.Path = tempDir
				tt.expectLogFile = tempDir + "/app.log"
			}

			logger := NewZapLogger(tt.config)
			assert.NotNil(t, logger, "Logger should not be nil")
			assert.Equal(t, tt.expectLogLevel, logger.Desugar().Level().String(), "Expected log level to match")

			logger.Info("Test Log Message")

			if tt.expectFile {
				// Check if log file was created
				_, err := os.Stat(tt.expectLogFile)
				assert.False(t, os.IsNotExist(err), "Expected log file to exist at %s", tt.expectLogFile)
			}
		})
	}
}

// TestNewForTest tests the NewForTest logger factory function.
func TestNewForTest(t *testing.T) {
	logger, recordedLogs := NewForTest(nil)

	assert.NotNil(t, logger, "Expected non-nil logger")
	assert.NotNil(t, recordedLogs, "Expected non-nil recorded logs")

	expectedLogMsg := "Test Log Entry"
	logger.Info(expectedLogMsg)

	assert.Equal(t, 1, recordedLogs.Len(), "Expected one log entry")
	entry := recordedLogs.All()[0]
	assert.Equal(t, expectedLogMsg, entry.Message, "Expected log message to match")
	assert.Equal(t, "info", entry.Level.String(), "Expected log level to be Info")
}

// TestZapLogger_With tests the With method for contextual logging.
func TestZapLogger_With(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTracer := mock_tracer.NewMockTracer(ctrl)
	tracing.SetTracer(mockTracer)

	tests := []struct {
		name      string
		config    *logger_config.Logs
		setupMock func()

		mdContextKey string
		contextKey   interface{}

		expectedLogKey string
		expectedValue  string
	}{
		{
			name: "Custom Context Key",
			config: &logger_config.Logs{
				AdditionalFields: []logger_config.AdditionalField{
					{
						FieldName: "custom_key",
						ValueFrom: "customKey",
					},
				},
			},
			setupMock: func() {
				mockTracer.EXPECT().GetTraceID(gomock.Any()).Return("").Times(1)
				mockTracer.EXPECT().GetSpanID(gomock.Any()).Return("").Times(1)
			},
			mdContextKey:   "customKey",
			expectedLogKey: "custom_key",
			expectedValue:  "customValue",
		},
		{
			name: "Trace ID Context",
			setupMock: func() {
				mockTracer.EXPECT().GetTraceID(gomock.Any()).Return("trace-12345").Times(1)
				mockTracer.EXPECT().GetSpanID(gomock.Any()).Return("").Times(1)
			},
			expectedLogKey: "trace_id",
			expectedValue:  "trace-12345",
		},
		{
			name: "Span ID Context",
			setupMock: func() {
				mockTracer.EXPECT().GetTraceID(gomock.Any()).Return("").Times(1)
				mockTracer.EXPECT().GetSpanID(gomock.Any()).Return("span-abcdef").Times(1)
			},
			expectedLogKey: "span_id",
			expectedValue:  "span-abcdef",
		},
		{
			name: "Request ID Context",
			setupMock: func() {
				mockTracer.EXPECT().GetTraceID(gomock.Any()).Return("").Times(1)
				mockTracer.EXPECT().GetSpanID(gomock.Any()).Return("").Times(1)
			},
			mdContextKey:   grpc.RequestIDHeader,
			expectedLogKey: "request_id",
			expectedValue:  "123456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create a fresh logger and observed logs for each test case
			logger, recordedLogs := NewForTest(tt.config)

			// Create context with the given key-value pair
			ctx := context.Background()

			if tt.mdContextKey != "" {
				md := metadata.Pairs(tt.mdContextKey, tt.expectedValue)
				ctx = metadata.NewIncomingContext(ctx, md)
			}
			if tt.contextKey != nil {
				// Otherwise, set regular context value
				ctx = context.WithValue(ctx, tt.contextKey, tt.expectedValue)
			}

			// Log using the With method
			expectedLogMsg := "Test Log Entry"
			logger.With(ctx).Info("Test Log Entry")

			// Verify log entry was created
			assert.Equal(t, 1, recordedLogs.Len(), "Expected one log entry")
			entry := recordedLogs.All()[0]
			assert.Equal(t, expectedLogMsg, entry.Message, "Expected log message to match")

			// Verify the log contains the expected context value
			assert.Equal(t, tt.expectedValue, entry.ContextMap()[tt.expectedLogKey], "Expected log context to match")
		})
	}
}

// TestNewTmpLogger tests the temporary logger creation.
func TestNewTmpLogger(t *testing.T) {
	logger := NewTmpLogger()
	assert.NotNil(t, logger, "Expected non-nil temporary logger")

	logger.Info("Temporary Logger Test")
}
