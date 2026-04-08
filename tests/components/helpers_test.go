package components_test

import (
	"testing"

	"pulse-agent/components"
)

func TestSeparator_NotEmpty(t *testing.T) {
	result := components.Separator(components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty separator")
	}
}

func TestOverflowGuard_WithinLimit(t *testing.T) {
	lines := []string{"a", "b", "c"}
	result := components.OverflowGuard(lines, 5)
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
}

func TestOverflowGuard_ExactLimit(t *testing.T) {
	lines := []string{"a", "b", "c"}
	result := components.OverflowGuard(lines, 3)
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
}

func TestOverflowGuard_ExceedsLimit(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e"}
	result := components.OverflowGuard(lines, 3)
	if len(result) != 3 {
		t.Fatalf("expected 3 lines (2 clipped + 1 hidden msg), got %d", len(result))
	}
	if result[len(result)-1] == "" {
		t.Fatal("last line should be the 'hidden' message, not empty")
	}
}

func TestOverflowGuard_Empty(t *testing.T) {
	result := components.OverflowGuard([]string{}, 5)
	if len(result) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(result))
	}
}

func TestRow_NotEmpty(t *testing.T) {
	result := components.Row("Label", "Value")
	if result == "" {
		t.Fatal("expected non-empty row")
	}
}

func TestRow_LongLabel(t *testing.T) {
	result := components.Row("Very Long Label Name", "val")
	if result == "" {
		t.Fatal("expected non-empty row for long label")
	}
}

func TestRow_LongValue(t *testing.T) {
	result := components.Row("L", "Very Long Value String")
	if result == "" {
		t.Fatal("expected non-empty row for long value")
	}
}
