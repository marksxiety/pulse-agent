package models_test

import (
	"math"
	"testing"

	"pulse-agent/models"
)

func TestNewMetricHistory(t *testing.T) {
	h := models.NewMetricHistory()
	if h == nil {
		t.Fatal("expected non-nil history")
	}
	if h.Len() != 0 {
		t.Fatalf("expected empty history, got len=%d", h.Len())
	}
}

func TestPush_IncrementsLength(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(1.0)
	h.Push(2.0)
	h.Push(3.0)

	if h.Len() != 3 {
		t.Fatalf("expected len=3, got %d", h.Len())
	}
}

func TestPush_WrapsAround(t *testing.T) {
	h := models.NewMetricHistory()

	for i := 0; i < models.HistoryCapacity+10; i++ {
		h.Push(float64(i))
	}

	if h.Len() != models.HistoryCapacity {
		t.Fatalf("expected len=%d, got %d", models.HistoryCapacity, h.Len())
	}

	first := h.Min()
	last := h.Max()

	if first != 10 {
		t.Fatalf("expected min=10 (oldest after wrap), got %v", first)
	}
	if last != float64(models.HistoryCapacity+9) {
		t.Fatalf("expected max=%d, got %v", models.HistoryCapacity+9, last)
	}
}

func TestMin_Empty(t *testing.T) {
	h := models.NewMetricHistory()
	if h.Min() != 0 {
		t.Fatalf("expected 0 for empty history, got %v", h.Min())
	}
}

func TestMin_SingleValue(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(42.0)
	if h.Min() != 42.0 {
		t.Fatalf("expected 42, got %v", h.Min())
	}
}

func TestMin_MultipleValues(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(10.0)
	h.Push(3.0)
	h.Push(7.0)

	if h.Min() != 3.0 {
		t.Fatalf("expected 3, got %v", h.Min())
	}
}

func TestMin_AfterWrap(t *testing.T) {
	h := models.NewMetricHistory()
	for i := 0; i < models.HistoryCapacity; i++ {
		h.Push(100.0)
	}
	h.Push(1.0)
	h.Push(50.0)

	if h.Min() != 1.0 {
		t.Fatalf("expected 1, got %v", h.Min())
	}
}

func TestMax_Empty(t *testing.T) {
	h := models.NewMetricHistory()
	if h.Max() != 0 {
		t.Fatalf("expected 0 for empty history, got %v", h.Max())
	}
}

func TestMax_SingleValue(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(42.0)
	if h.Max() != 42.0 {
		t.Fatalf("expected 42, got %v", h.Max())
	}
}

func TestMax_MultipleValues(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(5.0)
	h.Push(99.0)
	h.Push(12.0)

	if h.Max() != 99.0 {
		t.Fatalf("expected 99, got %v", h.Max())
	}
}

func TestAvg_Empty(t *testing.T) {
	h := models.NewMetricHistory()
	if h.Avg() != 0 {
		t.Fatalf("expected 0 for empty history, got %v", h.Avg())
	}
}

func TestAvg_SingleValue(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(42.0)
	if h.Avg() != 42.0 {
		t.Fatalf("expected 42, got %v", h.Avg())
	}
}

func TestAvg_MultipleValues(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(10.0)
	h.Push(20.0)
	h.Push(30.0)

	expected := 20.0
	if math.Abs(h.Avg()-expected) > 1e-9 {
		t.Fatalf("expected %v, got %v", expected, h.Avg())
	}
}

func TestAvg_AfterWrap(t *testing.T) {
	h := models.NewMetricHistory()
	for i := 0; i < models.HistoryCapacity; i++ {
		h.Push(0.0)
	}
	h.Push(100.0)

	expected := 100.0 / float64(models.HistoryCapacity)
	if math.Abs(h.Avg()-expected) > 1e-9 {
		t.Fatalf("expected %v, got %v", expected, h.Avg())
	}
}

func TestDownsample_Empty(t *testing.T) {
	h := models.NewMetricHistory()
	result := h.Downsample(10)
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestDownsample_SmallerThanCount(t *testing.T) {
	h := models.NewMetricHistory()
	for i := 1; i <= 9; i++ {
		h.Push(float64(i))
	}

	result := h.Downsample(3)
	if len(result) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(result))
	}

	if result[0] != 2.0 {
		t.Fatalf("bucket 0: expected avg of [1,2,3] = 2, got %v", result[0])
	}
	if result[1] != 5.0 {
		t.Fatalf("bucket 1: expected avg of [4,5,6] = 5, got %v", result[1])
	}
	if result[2] != 8.0 {
		t.Fatalf("bucket 2: expected avg of [7,8,9] = 8, got %v", result[2])
	}
}

func TestDownsample_LargerThanCount(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(10.0)
	h.Push(20.0)

	result := h.Downsample(10)
	if len(result) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(result))
	}
	if result[0] != 10.0 || result[1] != 20.0 {
		t.Fatalf("expected [10, 20], got %v", result)
	}
}

func TestDownsample_ExactlyMatchesCount(t *testing.T) {
	h := models.NewMetricHistory()
	h.Push(5.0)
	h.Push(10.0)
	h.Push(15.0)

	result := h.Downsample(3)
	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}
	for i, exp := range []float64{5.0, 10.0, 15.0} {
		if result[i] != exp {
			t.Fatalf("element %d: expected %v, got %v", i, exp, result[i])
		}
	}
}
