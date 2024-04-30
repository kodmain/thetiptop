package levels

// TYPE represents the log level type.
type TYPE uint8

// Log level constants.
const (
	OFF     TYPE = iota // OFF disables logging.
	PANIC               // PANIC for critical errors leading to termination.
	FATAL               // FATAL for severe errors that may cause termination.
	ERROR               // ERROR for non-fatal errors.
	WARN                // WARN for warning messages.
	SUCCESS             // SUCCESS for successful operations.
	INFO                // INFO for informational messages.
	MESSAGE             // MESSAGE for general messages.
	DEBUG               // DEBUG for debugging messages.
	TRACE               // TRACE for detailed tracing.

	DEFAULT = INFO // DEFAULT is the default log level, set to INFO.
)

// LABELS maps log level constants to their corresponding string labels.
var LABELS = []string{
	OFF:     "OFF",
	PANIC:   "PANIC",
	FATAL:   "FATAL",
	ERROR:   "ERROR",
	WARN:    "WARN",
	SUCCESS: "SUCCESS",
	INFO:    "INFO",
	MESSAGE: "MESSAGE",
	DEBUG:   "DEBUG",
	TRACE:   "TRACE",
}

// COLORS maps log level constants to their corresponding color codes.
var COLORS = []string{
	PANIC:   "9",   // Bright red for PANIC.
	FATAL:   "160", // Dark red for FATAL.
	ERROR:   "1",   // Red for ERROR.
	WARN:    "3",   // Yellow for WARN.
	SUCCESS: "2",   // Green for SUCCESS.
	INFO:    "4",   // Blue for INFO.
	MESSAGE: "7",   // White for MESSAGE.
	DEBUG:   "6",   // Cyan for DEBUG.
	TRACE:   "7",   // White for TRACE.
}

// Int returns the integer representation of the log level.
// This method is useful for comparing log levels.
//
// Returns:
// - uint8: The integer value of the log level.
func (t TYPE) Int() uint8 {
	return uint8(t)
}

// String returns the string representation of the log level.
// It provides a human-readable name for the log level.
//
// Returns:
// - string: The name of the log level, or "UNKNOWN" if the level is not defined.
func (t TYPE) String() string {
	if t >= TYPE(len(LABELS)) {
		return "UNKNOWN"
	}
	return LABELS[t]
}

// Color returns the color code associated with the log level.
// This method is useful for color-coding log messages.
//
// Returns:
// - string: The color code for the log level, or "UNKNOWN" if the level is not defined.
func (t TYPE) Color() string {
	if t >= TYPE(len(COLORS)) {
		return "UNKNOWN"
	}
	return COLORS[t]
}
