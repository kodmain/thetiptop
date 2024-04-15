package levels_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger/levels"
	"github.com/stretchr/testify/assert"
)

// TestLevelString tests the String method of the levels.TYPE type.
// It verifies that each log level is correctly converted to its string representation.
func TestLevelString(t *testing.T) {
	tests := []struct {
		level    levels.TYPE
		expected string
	}{
		{levels.OFF, "OFF"},
		{levels.PANIC, "PANIC"},
		{levels.FATAL, "FATAL"},
		{levels.ERROR, "ERROR"},
		{levels.SUCCESS, "SUCCESS"},
		{levels.MESSAGE, "MESSAGE"},
		{levels.WARN, "WARN"},
		{levels.INFO, "INFO"},
		{levels.DEBUG, "DEBUG"},
		{levels.TRACE, "TRACE"},
		{levels.TYPE(99), "UNKNOWN"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.level.String(), "La chaîne de caractères devrait correspondre au niveau")
	}
}

// TestLevelColor tests the Color method of the levels.TYPE type.
// It verifies that each log level is associated with the correct color code.
// This test also includes a case for an unknown log level.
func TestLevelColor(t *testing.T) {
	tests := []struct {
		level    levels.TYPE
		expected string
	}{
		{levels.PANIC, "9"},
		{levels.FATAL, "160"},
		{levels.ERROR, "1"},
		{levels.SUCCESS, "2"},
		{levels.MESSAGE, "7"},
		{levels.WARN, "3"},
		{levels.INFO, "4"},
		{levels.DEBUG, "6"},
		{levels.TRACE, "7"},
		{levels.TYPE(99), "UNKNOWN"}, // Test pour un niveau inconnu
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.level.Color(), "La couleur devrait correspondre au niveau")
	}
}

// TestLevelInt tests the Int method of the levels.TYPE type.
// It verifies that each log level is correctly converted to its corresponding integer value.
func TestLevelInt(t *testing.T) {
	tests := []struct {
		name     string
		level    levels.TYPE
		expected uint8
	}{
		{"OFF", levels.OFF, 0},
		{"PANIC", levels.PANIC, 1},
		{"FATAL", levels.FATAL, 2},
		{"ERROR", levels.ERROR, 3},
		{"WARN", levels.WARN, 4},
		{"SUCCESS", levels.SUCCESS, 5},
		{"INFO", levels.INFO, 6},
		{"MESSAGE", levels.MESSAGE, 7},
		{"DEBUG", levels.DEBUG, 8},
		{"TRACE", levels.TRACE, 9},
		{"TYPE99", levels.TYPE(99), 99}, // Test for an unknown log level
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.level.Int(), "The integer value should match the log level %d we got %d for %s", test.level.Int(), test.level.Int(), test.name)
	}
}
