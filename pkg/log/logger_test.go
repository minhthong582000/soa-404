package log

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	logger_config "github.com/minhthong582000/soa-404/pkg/config"
)

func TestGetLogger_DefaultInitialization(t *testing.T) {
	ResetGlobalLogger()
	logger := GetLogger()

	assert.NotNil(t, logger, "Expected a non-nil logger from GetLogger")
}

func TestSetLogger(t *testing.T) {
	ResetGlobalLogger()
	logger := NewTmpLogger()
	SetLogger(logger)

	assert.Equal(t, logger, GetLogger(), "Expected global logger to be the tmp logger set by SetLogger")
}

func TestLogFactory(t *testing.T) {
	tests := []struct {
		name             string
		config           *logger_config.Logs
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "Valid ZapLog Provider",
			config: &logger_config.Logs{
				Provider: logger_config.ZapLog,
			},
			expectError: false,
		},
		{
			name: "Invalid Provider",
			config: &logger_config.Logs{
				Provider: "invalidProvider",
			},
			expectError:      true,
			expectedErrorMsg: "unsupported logger provider: invalidProvider",
		},
		{
			name: "Empty Provider",
			config: &logger_config.Logs{
				Provider: "",
			},
			expectError:      true,
			expectedErrorMsg: "unsupported logger provider: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetGlobalLogger()

			logger, err := LogFactory(tt.config)

			if tt.expectError {
				assert.Error(t, err, "Expected an error for invalid provider")
				assert.Nil(t, logger, "Logger should be nil on error")
				assert.EqualError(t, err, tt.expectedErrorMsg, "Error message should match expected")
			} else {
				assert.NoError(t, err, "Expected no error for valid provider")
				assert.NotNil(t, logger, "Logger should not be nil for valid provider")
				assert.IsType(t, &zapLogger{}, logger, "Expected logger to be of type *zapLogger")
			}
		})
	}
}

// Ensure thread safety when accessing logger
func TestGetLogger_ThreadSafety(t *testing.T) {
	ResetGlobalLogger()
	var wg sync.WaitGroup
	routines := 10
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func() {
			defer wg.Done()
			GetLogger()
		}()
	}

	wg.Wait()
	assert.NotNil(t, GetLogger(), "Expected non-nil logger even with concurrent calls")
}
