package loggerEntity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{INFO, "INFO"},
		{WARNING, "WARNING"},
		{ERROR, "ERROR"},
		{SUCCESS, "SUCCESS"},
		{PANIC, "PANIC"},
		{LogLevel(100), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			actual := Name(tt.level)
			assert.Equal(t, tt.expected, actual, "they should be equal")
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{INFO, "INFO"},
		{WARNING, "WARNING"},
		{ERROR, "ERROR"},
		{SUCCESS, "SUCCESS"},
		{PANIC, "PANIC"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			actual := tt.level.String()
			assert.Equal(t, tt.expected, actual, "they should be equal")
		})
	}
}
