package cmd

import (
	"log/slog"
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  slog.Level
	}{
		{name: "default", want: slog.LevelDebug},
		{name: "debug", value: "debug", want: slog.LevelDebug},
		{name: "info", value: "info", want: slog.LevelInfo},
		{name: "warn", value: "warn", want: slog.LevelWarn},
		{name: "error", value: "error", want: slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogLevel(tt.value)
			if err != nil {
				t.Fatalf("parseLogLevel() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("parseLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLogLevelRejectsInvalidValue(t *testing.T) {
	level, err := parseLogLevel("verbose")
	if err == nil {
		t.Fatal("parseLogLevel() error = nil, want an error")
	}
	if level != slog.LevelDebug {
		t.Errorf("parseLogLevel() fallback = %v, want %v", level, slog.LevelDebug)
	}
}
