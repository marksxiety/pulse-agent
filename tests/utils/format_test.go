package utils_test

import (
	"testing"
	"time"

	"pulse-agent/utils"
)

func TestBytesToGB(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0.00 GB"},
		{1_000_000_000, "1.00 GB"},
		{1_073_741_824, "1.07 GB"},
		{500_000_000, "0.50 GB"},
	}
	for _, tt := range tests {
		got := utils.BytesToGB(tt.input)
		if got != tt.expected {
			t.Errorf("BytesToGB(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestBytesToMB(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0.0 MB"},
		{1_000_000, "1.0 MB"},
		{500_000, "0.5 MB"},
		{1_048_576, "1.0 MB"},
	}
	for _, tt := range tests {
		got := utils.BytesToMB(tt.input)
		if got != tt.expected {
			t.Errorf("BytesToMB(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFormatUptime_WithHours(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{1*time.Hour + 2*time.Minute + 3*time.Second, "1h 02m 03s"},
		{24*time.Hour + 5*time.Minute, "24h 05m 00s"},
		{1*time.Hour + 0*time.Minute + 5*time.Second, "1h 00m 05s"},
		{time.Duration(0), "0m 00s"},
	}
	for _, tt := range tests {
		got := utils.FormatUptime(tt.input)
		if got != tt.expected {
			t.Errorf("FormatUptime(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFormatUptime_WithoutHours(t *testing.T) {
	got := utils.FormatUptime(5*time.Minute + 30*time.Second)
	if got != "5m 30s" {
		t.Errorf("FormatUptime(5m30s) = %q, want %q", got, "5m 30s")
	}
}

func TestFormatUptime_RoundsToSeconds(t *testing.T) {
	got := utils.FormatUptime(5*time.Minute + 30*time.Second + 500*time.Millisecond)
	if got != "5m 31s" {
		t.Errorf("expected rounding up, got %q", got)
	}
}

func TestPadToHeight_AlreadyEnough(t *testing.T) {
	lines := []string{"a", "b", "c"}
	result := utils.PadToHeight(lines, 3)
	if len(result) != 3 {
		t.Fatalf("expected len 3, got %d", len(result))
	}
}

func TestPadToHeight_NeedsPadding(t *testing.T) {
	lines := []string{"a"}
	result := utils.PadToHeight(lines, 3)
	if len(result) != 3 {
		t.Fatalf("expected len 3, got %d", len(result))
	}
	if result[1] != "" || result[2] != "" {
		t.Fatal("padded lines should be empty strings")
	}
}

func TestPadToHeight_EmptySlice(t *testing.T) {
	result := utils.PadToHeight([]string{}, 5)
	if len(result) != 5 {
		t.Fatalf("expected len 5, got %d", len(result))
	}
}

func TestPadToHeight_LongerThanTarget(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e"}
	result := utils.PadToHeight(lines, 3)
	if len(result) != 5 {
		t.Fatalf("should not truncate, expected len 5, got %d", len(result))
	}
}
