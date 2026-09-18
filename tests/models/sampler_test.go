package models_test

import (
	"testing"
	"time"

	"pulse-agent/models"
)

const testInterval = 5 * time.Second

func cpu(pct float64) *models.CPUPayload {
	return &models.CPUPayload{Percentage: pct, CoreCount: 8}
}

func memory(used, total uint64) *models.MemoryPayload {
	return &models.MemoryPayload{Used: used, Total: total}
}

func disk(used, total uint64) *models.DiskPayload {
	return &models.DiskPayload{Used: used, Total: total}
}

func TestSampler_RecordsEveryMetricOnFirstSample(t *testing.T) {
	s := models.NewSampler(testInterval)
	now := time.Now()

	s.Record(now, cpu(25), memory(50, 100), disk(75, 100))

	if s.CPU.Len() != 1 || s.Memory.Len() != 1 || s.Disk.Len() != 1 {
		t.Fatalf("expected one sample each, got cpu=%d mem=%d disk=%d",
			s.CPU.Len(), s.Memory.Len(), s.Disk.Len())
	}
	if got := s.CPU.Max(); got != 25 {
		t.Errorf("CPU sample = %v, want the reported percentage 25", got)
	}
	if got := s.Memory.Max(); got != 50 {
		t.Errorf("Memory sample = %v, want 50%% for 50/100", got)
	}
	if got := s.Disk.Max(); got != 75 {
		t.Errorf("Disk sample = %v, want 75%% for 75/100", got)
	}
}

// A nil payload means the collector has not reported yet. Recording zero for it
// would drag the trend down and misalign the three histories.
func TestSampler_SkipsMetricsThatHaveNotReported(t *testing.T) {
	s := models.NewSampler(testInterval)

	s.Record(time.Now(), nil, memory(50, 100), nil)

	if s.CPU.Len() != 0 {
		t.Errorf("CPU has %d samples, want 0 when its collector is silent", s.CPU.Len())
	}
	if s.Disk.Len() != 0 {
		t.Errorf("Disk has %d samples, want 0 when its collector is silent", s.Disk.Len())
	}
	if s.Memory.Len() != 1 {
		t.Errorf("Memory has %d samples, want 1", s.Memory.Len())
	}
}

func TestSampler_DoesNotSampleBeforeTheInterval(t *testing.T) {
	s := models.NewSampler(testInterval)
	start := time.Now()

	s.Record(start, cpu(10), memory(1, 100), disk(1, 100))
	s.Record(start.Add(time.Second), cpu(20), memory(2, 100), disk(2, 100))
	s.Record(start.Add(testInterval-time.Millisecond), cpu(30), memory(3, 100), disk(3, 100))

	if s.CPU.Len() != 1 {
		t.Errorf("CPU has %d samples, want 1: the interval had not elapsed", s.CPU.Len())
	}
	if got := s.CPU.Max(); got != 10 {
		t.Errorf("CPU sample = %v, want the first reading 10", got)
	}
}

func TestSampler_SamplesAgainOnceTheIntervalElapses(t *testing.T) {
	s := models.NewSampler(testInterval)
	start := time.Now()

	s.Record(start, cpu(10), memory(1, 100), disk(1, 100))
	s.Record(start.Add(testInterval), cpu(60), memory(2, 100), disk(2, 100))

	if s.CPU.Len() != 2 {
		t.Fatalf("CPU has %d samples, want 2", s.CPU.Len())
	}
	if got := s.CPU.Max(); got != 60 {
		t.Errorf("CPU Max = %v, want the later reading 60", got)
	}
}

func TestSampler_SkipsVolumesWithNoCapacity(t *testing.T) {
	s := models.NewSampler(testInterval)

	s.Record(time.Now(), cpu(10), memory(0, 0), disk(0, 0))

	if s.Memory.Len() != 0 {
		t.Errorf("Memory has %d samples, want 0 for a zero total", s.Memory.Len())
	}
	if s.Disk.Len() != 0 {
		t.Errorf("Disk has %d samples, want 0 for a zero total", s.Disk.Len())
	}
	if s.CPU.Len() != 1 {
		t.Errorf("CPU has %d samples, want 1", s.CPU.Len())
	}
}

// The widget advertises a three-hour trend. That only holds while the interval
// and the history capacity stay in step, so pin the relationship.
func TestSampler_FullHistoryCoversThreeHours(t *testing.T) {
	if got := testInterval * time.Duration(models.HistoryCapacity); got != 3*time.Hour {
		t.Errorf("HistoryCapacity (%d) at %v covers %v, want 3h",
			models.HistoryCapacity, testInterval, got)
	}
}
