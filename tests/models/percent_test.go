package models_test

import (
	"math"
	"testing"

	"pulse-agent/models"
)

func TestPercent_ZeroTotalReadsAsZero(t *testing.T) {
	if got := models.Percent(500, 0); got != 0 {
		t.Errorf("Percent(500, 0) = %v, want 0", got)
	}
	if got := models.Percent(0, 0); got != 0 {
		t.Errorf("Percent(0, 0) = %v, want 0", got)
	}
}

func TestPercent_Normal(t *testing.T) {
	tests := []struct {
		used, total uint64
		want        float64
	}{
		{50, 100, 50},
		{0, 100, 0},
		{100, 100, 100},
		{1, 4, 25},
	}
	for _, tt := range tests {
		if got := models.Percent(tt.used, tt.total); got != tt.want {
			t.Errorf("Percent(%d, %d) = %v, want %v", tt.used, tt.total, got, tt.want)
		}
	}
}

func TestPercent_DoesNotOverflowOnLargeVolumes(t *testing.T) {
	// A 4 TB volume: the naive uint64 multiply used by percentage helpers
	// elsewhere overflows at this size, so this pins the float64 path.
	const total = uint64(4_000_000_000_000)
	got := models.Percent(total/2, total)
	if math.Abs(got-50) > 0.0001 {
		t.Errorf("Percent(half of %d) = %v, want ~50", total, got)
	}
}
