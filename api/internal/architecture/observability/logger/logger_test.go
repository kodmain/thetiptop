package logger_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.SetLevel(levels.TRACE)
}

// newTestLogger creates and returns a new logger instance for testing.
// This logger is configured with a file writer and trace level logging.
func newTestLogger(writers ...io.Writer) *logger.Logger {
	return logger.New(levels.TRACE, writers...)
}

func TestNewTestLogger(t *testing.T) {
	var success bytes.Buffer
	var fail bytes.Buffer

	// Définition des cas de test
	testCases := []struct {
		name       string
		parameters []io.Writer
	}{
		{
			name:       "No parameters",
			parameters: nil,
		},
		{
			name:       "With success parameter",
			parameters: []io.Writer{&success},
		},
		{
			name:       "With success and fail parameters",
			parameters: []io.Writer{&success, &fail},
		},
	}

	// Boucle de test sur les cas de test
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := newTestLogger(tc.parameters...)
			assert.NotNil(t, logger, "Logger should not be nil")
		})
	}
}

// TestLoggerErrorMethods tests the error handling methods of the logger.
// It verifies that Error, Fatal, and Panic methods return true when provided with an error.
func TestLoggerErrorMethods(t *testing.T) {
	logger := newTestLogger()

	err := errors.New("test error")
	assert.True(t, logger.Error(err), "Error should return true if an error is provided")
	assert.True(t, logger.Fatal(err), "Fatal should return true if an error is provided")
	assert.True(t, logger.Panic(err), "Panic should return true if an error is provided")
}

// TestLoggerLevels tests the logger's ability to handle different logging levels.
// It writes messages at various levels and checks if they are correctly written to the respective output files.
func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	tests := []struct {
		level    levels.TYPE
		message  string
		contains string
	}{
		{levels.DEBUG, "debug message", "debug message"},
		{levels.INFO, "info message", "info message"},
		{levels.WARN, "warn message", "warn message"},
		{levels.ERROR, "error message", "error message"},
		{levels.FATAL, "fatal message", "fatal message"},
		{levels.PANIC, "panic message", "panic message"},
		{levels.SUCCESS, "success message", "success message"},
		{levels.MESSAGE, "simple message", "simple message"},
		{levels.TRACE, "trace message", "trace message"},
	}

	for _, test := range tests {
		logger.Write(test.level, test.message)
		assert.True(t, bytes.Contains(buf.Bytes(), []byte(test.contains)), "File should contain the message")

		/*
			if test.level.Int() > levels.WARN.Int() {
				exists, err := fs.ExistsFile(writers.FILE_STDOUT)
				assert.NoError(t, err, "Error should be nil")
				assert.True(t, exists, "Stdout file should exist")
				ok, err := fs.Contains(writers.FILE_STDOUT, test.contains)
				assert.Nil(t, err, "Error should be nil")
				assert.True(t, ok, "File should contain the message")
			} else {
				exists, err := fs.ExistsFile(writers.FILE_STDERR)
				assert.NoError(t, err, "Error should be nil")
				assert.True(t, exists, "Stderr file should exist")
				ok, err := fs.Contains(writers.FILE_STDERR, test.contains)
				assert.Nil(t, err, "Error should be nil")
				assert.True(t, ok, "File should contain the message", writers.FILE_STDERR)
			}
		*/
	}
}
