package logger

import (
	"errors"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

// TestLoggerMethods tests various methods of the logger package.
// It checks the behavior of different logging methods, particularly focusing on error handling and ensuring no panics.
func TestLoggerMethods(t *testing.T) {
	// Error simulation
	testError := errors.New("test error")

	// Test for Panic, Fatal, Error methods
	// Ensures that Panic, Fatal, and Error methods return true when an error is provided.
	assert.True(t, Panic(testError), "Panic should return true if an error is provided")
	assert.True(t, Fatal(testError), "Fatal should return true if an error is provided")
	assert.True(t, Error(testError), "Error should return true if an error is provided")

	assert.False(t, Panic(nil), "Panic should return false if an error is provided")
	assert.False(t, Fatal(nil), "Fatal should return false if an error is provided")
	assert.False(t, Error(nil), "Error should return false if an error is provided")

	// Test for other methods, ensuring they do not produce an error
	// Verifies that methods such as Success, Message, Warn, Info, Debug, and Trace do not cause panics.
	assert.NotPanics(t, func() { Success("test success") }, "Success should not panic")
	assert.NotPanics(t, func() { Message("test message") }, "Message should not panic")
	assert.NotPanics(t, func() { Warn("test warn") }, "Warn should not panic")
	assert.NotPanics(t, func() { Info("test info") }, "Info should not panic")
	assert.NotPanics(t, func() { Debug("test debug") }, "Debug should not panic")
	assert.NotPanics(t, func() { Trace() }, "Trace should not panic")

	assert.NotPanics(t, func() { Successf("test success") }, "Success should not panic")
	assert.NotPanics(t, func() { Messagef("test message") }, "Message should not panic")
	assert.NotPanics(t, func() { Warnf("test warn") }, "Warn should not panic")
	assert.NotPanics(t, func() { Infof("test info") }, "Info should not panic")
	assert.NotPanics(t, func() { Debugf("test debug") }, "Debug should not panic")
}

// TestSetLevel tests the SetLevel function of the logger package.
func TestSetLevel(t *testing.T) {
	// Test setting different log levels
	SetLevel(levels.DEBUG)
	assert.Equal(t, levels.DEBUG, standard().level, "SetLevel should set the log level to levels.DEBUG")

	SetLevel(levels.INFO)
	assert.Equal(t, levels.INFO, standard().level, "SetLevel should set the log level to LevelInfo")

	SetLevel(levels.WARN)
	assert.Equal(t, levels.WARN, standard().level, "SetLevel should set the log level to LevelWarn")

	SetLevel(levels.ERROR)
	assert.Equal(t, levels.ERROR, standard().level, "SetLevel should set the log level to LevelError")

	SetLevel(levels.FATAL)
	assert.Equal(t, levels.FATAL, standard().level, "SetLevel should set the log level to LevelFatal")

	SetLevel(levels.PANIC)
	assert.Equal(t, levels.PANIC, standard().level, "SetLevel should set the log level to LevelPanic")
}
