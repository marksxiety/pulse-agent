package models_test

import (
	"testing"

	"pulse-agent/models"
)

func TestNewSeries_SummarisesHistory(t *testing.T) {
	h := models.NewMetricHistory()
	for _, v := range []float64{10, 20, 30, 40} {
		h.Push(v)
	}

	got := models.NewSeries(h, 42.5, 2)

	if got.Percent != 42.5 {
		t.Errorf("Percent = %v, want 42.5", got.Percent)
	}
	if got.Max != 40 {
		t.Errorf("Max = %v, want 40", got.Max)
	}
	if got.Avg != 25 {
		t.Errorf("Avg = %v, want 25", got.Avg)
	}
	if len(got.Trend) != 2 {
		t.Fatalf("Trend has %d samples, want 2 (downsampled)", len(got.Trend))
	}
}

func TestNewSeries_KeepsShortHistoryIntact(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(1)
	h.Push(2)

	got := models.NewSeries(h, 2, 120)
	if len(got.Trend) != 2 {
		t.Errorf("Trend has %d samples, want the 2 that exist", len(got.Trend))
	}
}

func TestNewSeries_EmptyHistory(t *testing.T) {
	got := models.NewSeries(models.NewMetricHistory(), 0, 120)

	if got.Trend != nil {
		t.Errorf("Trend = %v, want nil for an empty history", got.Trend)
	}
	if got.Max != 0 || got.Avg != 0 {
		t.Errorf("Max/Avg = %v/%v, want 0/0", got.Max, got.Avg)
	}
}
