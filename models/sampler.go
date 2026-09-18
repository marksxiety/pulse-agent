package models

import "time"

// Sampler folds a stream of metrics into three rolling histories, recording at
// most one sample per interval.
//
// With HistoryCapacity samples and a 5s interval the histories cover three
// hours of trend, which is the window both frontends present.
type Sampler struct {
	interval time.Duration
	last     time.Time

	CPU    *MetricHistory
	Memory *MetricHistory
	Disk   *MetricHistory
}

func NewSampler(interval time.Duration) *Sampler {
	return &Sampler{
		interval: interval,
		CPU:      NewMetricHistory(),
		Memory:   NewMetricHistory(),
		Disk:     NewMetricHistory(),
	}
}

// Record stores the current readings once at least one interval has elapsed
// since the previous sample.
//
// A nil payload means that collector has not reported yet, so only metrics with
// data are recorded. That is what keeps the three histories aligned instead of
// pushing zeroes for a collector that is merely lagging.
func (s *Sampler) Record(now time.Time, cpu *CPUPayload, memory *MemoryPayload, disk *DiskPayload) {
	if now.Sub(s.last) < s.interval {
		return
	}
	s.last = now

	if cpu != nil {
		s.CPU.Push(cpu.Percentage)
	}
	if memory != nil && memory.Total > 0 {
		s.Memory.Push(Percent(memory.Used, memory.Total))
	}
	if disk != nil && disk.Total > 0 {
		s.Disk.Push(Percent(disk.Used, disk.Total))
	}
}
