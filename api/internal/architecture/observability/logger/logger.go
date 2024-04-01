package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	"github.com/Code-Hex/dd"
	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger/levels"
)

// Constant for formatting log path with color.
const PATH = "| \033[38;5;2m%v\033[39;49m | \033[38;5;%sm%s\033[39;49m |"

// Logger struct holds the logging configuration and loggers for different levels.
type Logger struct {
	level   levels.TYPE // The logging level of the logger.
	success *log.Logger // Logger for success messages.
	failure *log.Logger // Logger for error messages.
}

// New creates a new Logger instance with the specified writer type and log level.
// It initializes two loggers: one for success and another for failure messages.
//
// Parameters:
// - l: levels.TYPE The logging level for the logger.
// - ws: []io.Writer The writer for success or failure messages (optional).
//
// Returns:
// - *Logger: A new Logger instance.
func New(l levels.TYPE, ws ...io.Writer) *Logger {
	logger := &Logger{level: l}

	switch len(ws) {
	case 1:
		logger.success, logger.failure = createLoggers(ws[0], ws[0])
	case 2:
		logger.success, logger.failure = createLoggers(ws[0], ws[1])
	default:
		logger.success, logger.failure = createLoggers(os.Stdout, os.Stderr)
	}
	return logger
}

// createLoggers creates two new log.Logger instances.
//
// This function assists in creating two log.Logger instances with specified parameters.
//
// Parameters:
// - successOutput: io.Writer the Writer for success logs.
// - failureOutput: io.Writer the Writer for failure logs.
//
// Returns:
// - *log.Logger: the logger for success.
// - *log.Logger: the logger for failures.
func createLoggers(successOutput, failureOutput io.Writer) (*log.Logger, *log.Logger) {
	return log.New(successOutput, "", log.Ldate|log.Ltime), log.New(failureOutput, "", log.Ldate|log.Ltime)
}

// Write writes the log message with the specified log level.
// It formats the message and decides which logger to use based on the level.
//
// Parameters:
// - level: levels.TYPE The log level for the message.
// - messages: ...any The messages or data to log.
func (l *Logger) Write(level levels.TYPE, messages ...any) {
	for _, message := range messages {
		if level <= l.level {
			var logger *log.Logger = nil
			if level <= levels.WARN {
				logger = l.failure
			} else {
				logger = l.success
			}

			if level == levels.DEBUG {
				logger.Println(fmt.Sprintf(PATH, os.Getpid(), level.Color(), level.String()), dd.Dump(message, dd.WithIndent(4)))
			} else if level <= l.level {
				logger.Println(fmt.Sprintf(PATH, os.Getpid(), level.Color(), level.String()), message)
			}
		}
	}
}

// Panic logs the error and stack trace with the PANIC level.
// It logs the error and a stack trace for debugging.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the error is not nil, false otherwise.
func (l *Logger) Panic(err error) bool {
	if err != nil {
		l.Write(levels.PANIC, err, string(debug.Stack()))
		return true
	}

	return false
}

// Fatal logs the error and stack trace with the FATAL level.
// It logs critical errors that might require the application to stop.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the error is not nil, false otherwise.
func (l *Logger) Fatal(err error) bool {
	if err != nil {
		l.Write(levels.FATAL, err, string(debug.Stack()))
		return true
	}

	return false
}

// Error logs the error with the ERROR level.
// It is used for logging general errors.
//
// Parameters:
// - err: error The error to log.
//
// Returns:
// - bool: true if the error is not nil, false otherwise.
func (l *Logger) Error(err error) bool {
	if err != nil {
		l.Write(levels.ERROR, err)
		return true
	}

	return false
}

// Success logs the success message with the SUCCESS level.
// It is used for logging successful operations.
//
// Parameters:
// - v: ...any The success messages or data to log.
func (l *Logger) Success(v ...any) {
	l.Write(levels.SUCCESS, v...)
}

// Message logs the message with the MESSAGE level.
// It is used for general-purpose logging.
//
// Parameters:
// - v: ...any The messages or data to log.
func (l *Logger) Message(v ...any) {
	l.Write(levels.MESSAGE, v...)
}

// Warn logs the warning message with the WARN level.
// It is used for logging potential issues or warnings.
//
// Parameters:
// - v: ...any The warning messages or data to log.
func (l *Logger) Warn(v ...any) {
	l.Write(levels.WARN, v...)
}

// Info logs the info message with the INFO level.
// It is used for logging informational messages.
//
// Parameters:
// - v: ...any The informational messages or data to log.
func (l *Logger) Info(v ...any) {
	l.Write(levels.INFO, v...)
}

// Debug logs the debug message with the DEBUG level.
// It provides detailed debug information for troubleshooting.
//
// Parameters:
// - v: ...any The debug messages or data to log.
func (l *Logger) Debug(v ...any) {
	l.Write(levels.DEBUG, v...)
}

// Trace logs the stack trace with the TRACE level.
// It is used for logging detailed execution traces for in-depth debugging.
func (l *Logger) Trace() {
	l.Write(levels.TRACE, string(debug.Stack()))
}

// Infof logs an informational message with formatted output.
// It is similar to Info but allows for formatted messages.
//
// Parameters:
// - format: string The format string.
// - a: ...any The arguments for formatting.
func (l *Logger) Infof(format string, a ...any) {
	l.Write(levels.INFO, fmt.Sprintf(format, a...))
}

// Warnf logs a warning message with formatted output.
// It is similar to Warn but allows for formatted messages.
//
// Parameters:
// - format: string The format string.
// - a: ...any The arguments for formatting.
func (l *Logger) Warnf(format string, a ...any) {
	l.Write(levels.WARN, fmt.Sprintf(format, a...))
}

// Successf logs a success message with formatted output.
// It is similar to Success but allows for formatted messages.
//
// Parameters:
// - format: string The format string.
// - a: ...any The arguments for formatting.
func (l *Logger) Successf(format string, a ...any) {
	l.Write(levels.SUCCESS, fmt.Sprintf(format, a...))
}

// Debugf logs a debug message with formatted output.
// It is similar to Debug but allows for formatted messages.
//
// Parameters:
// - format: string The format string.
// - a: ...any The arguments for formatting.
func (l *Logger) Debugf(format string, a ...any) {
	l.Write(levels.DEBUG, fmt.Sprintf(format, a...))
}

// Messagef logs a general message with formatted output.
// It is similar to Message but allows for formatted messages.
//
// Parameters:
// - format: string The format string.
// - a: ...any The arguments for formatting.
func (l *Logger) Messagef(format string, a ...any) {
	l.Write(levels.MESSAGE, fmt.Sprintf(format, a...))
}
