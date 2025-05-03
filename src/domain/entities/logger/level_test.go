package loggerEntity

import (
	"testing"
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
			if actual != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, actual)
			}
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
			if actual != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}
