package components_test

import (
	"testing"

	"pulse-agent/components"
)

func TestProgressBar_ZeroPercent(t *testing.T) {
	result := components.ProgressBar(0, 20, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestProgressBar_Full(t *testing.T) {
	result := components.ProgressBar(100, 10, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestProgressBar_Half(t *testing.T) {
	result := components.ProgressBar(50, 10, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestProgressBar_NegativeClampsToZero(t *testing.T) {
	result := components.ProgressBar(-10, 10, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result for negative input")
	}
}

func TestProgressBar_Over100ClampsTo100(t *testing.T) {
	result := components.ProgressBar(150, 10, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestProgressBar_ZeroWidth(t *testing.T) {
	result := components.ProgressBar(50, 0, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result even with 0 width")
	}
}

func TestSparkline_Empty(t *testing.T) {
	result := components.Sparkline(nil, components.ColorCPUAccent)
	if result != "" {
		t.Fatalf("expected empty string for nil input, got %q", result)
	}

	result = components.Sparkline([]float64{}, components.ColorCPUAccent)
	if result != "" {
		t.Fatalf("expected empty string for empty slice, got %q", result)
	}
}

func TestSparkline_SingleValue(t *testing.T) {
	result := components.Sparkline([]float64{50.0}, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result for single value")
	}
}

func TestSparkline_AllSameValues(t *testing.T) {
	result := components.Sparkline([]float64{50, 50, 50, 50, 50}, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result for flat values")
	}
}

func TestSparkline_VaryingValues(t *testing.T) {
	result := components.Sparkline([]float64{10, 20, 50, 80, 100}, components.ColorCPUAccent)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestCentreBlock_NarrowTerminal(t *testing.T) {
	result := components.CentreBlock("test", 4)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestCentreBlock_WideTerminal(t *testing.T) {
	result := components.CentreBlock("test", 100)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestCentreBlock_ZeroWidth(t *testing.T) {
	result := components.CentreBlock("test", 0)
	if result == "" {
		t.Fatal("expected non-empty result even with 0 width")
	}
}
