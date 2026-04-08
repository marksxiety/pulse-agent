package components_test

import (
	"testing"

	"pulse-agent/components"
)

func TestInfoModal_DefaultScroll(t *testing.T) {
	result := components.InfoModal(0)
	if result == "" {
		t.Fatal("expected non-empty modal output")
	}
}

func TestInfoModal_ScrollBeyondContent(t *testing.T) {
	result := components.InfoModal(9999)
	if result == "" {
		t.Fatal("expected non-empty modal output even with excessive scroll")
	}
}

func TestInfoModal_NegativeScroll(t *testing.T) {
	result := components.InfoModal(-5)
	if result == "" {
		t.Fatal("expected non-empty modal output even with negative scroll")
	}
}

func TestInfoModal_DifferentScrollPositions(t *testing.T) {
	for scroll := 0; scroll < 20; scroll++ {
		result := components.InfoModal(scroll)
		if result == "" {
			t.Fatalf("expected non-empty modal at scroll=%d", scroll)
		}
	}
}
