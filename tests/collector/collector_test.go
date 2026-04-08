package collector_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"pulse-agent/collector"
	"pulse-agent/models"
	"pulse-agent/types"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type mockSystemInfo struct {
	cpuPercentFn func(time.Duration, bool) ([]float64, error)
	cpuCountFn   func(bool) (int, error)
	virtMemFn    func() (*mem.VirtualMemoryStat, error)
	swapMemFn    func() (*mem.SwapMemoryStat, error)
	diskUsageFn  func(string) (*disk.UsageStat, error)
	diskIOFn     func() (map[string]disk.IOCountersStat, error)
}

func (m *mockSystemInfo) CPUPercent(d time.Duration, b bool) ([]float64, error) {
	return m.cpuPercentFn(d, b)
}

func (m *mockSystemInfo) CPUCount(b bool) (int, error) {
	return m.cpuCountFn(b)
}

func (m *mockSystemInfo) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return m.virtMemFn()
}

func (m *mockSystemInfo) SwapMemory() (*mem.SwapMemoryStat, error) {
	return m.swapMemFn()
}

func (m *mockSystemInfo) DiskUsage(p string) (*disk.UsageStat, error) {
	return m.diskUsageFn(p)
}

func (m *mockSystemInfo) DiskIOCounters() (map[string]disk.IOCountersStat, error) {
	return m.diskIOFn()
}

func collectWithTimeout(c collector.Collector, timeout time.Duration) []models.Metric {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan models.Metric, 100)
	go c.Collect(ctx, ch)

	var results []models.Metric
	for {
		select {
		case m := <-ch:
			results = append(results, m)
			if len(results) >= 2 {
				cancel()
				return results
			}
		case <-ctx.Done():
			return results
		}
	}
}

func TestCPUCollector_CollectsData(t *testing.T) {
	mock := &mockSystemInfo{
		cpuCountFn: func(bool) (int, error) { return 8, nil },
		cpuPercentFn: func(time.Duration, bool) ([]float64, error) {
			return []float64{42.5}, nil
		},
	}

	c := collector.NewCPUCollector(mock)
	results := collectWithTimeout(c, 5*time.Second)

	if len(results) < 1 {
		t.Fatal("expected at least 1 metric from CPU collector")
	}

	m := results[0]
	if m.Source != types.CPU {
		t.Fatalf("expected source CPU, got %s", m.Source)
	}

	payload, ok := m.Data.(models.CPUPayload)
	if !ok {
		t.Fatal("expected CPUPayload")
	}
	if payload.Percentage != 42.5 {
		t.Fatalf("expected 42.5%%, got %v%%", payload.Percentage)
	}
	if payload.CoreCount != 8 {
		t.Fatalf("expected 8 cores, got %d", payload.CoreCount)
	}
}

func TestCPUCollector_StopsOnCancel(t *testing.T) {
	callCount := 0
	var mu sync.Mutex

	mock := &mockSystemInfo{
		cpuCountFn: func(bool) (int, error) { return 4, nil },
		cpuPercentFn: func(d time.Duration, b bool) ([]float64, error) {
			mu.Lock()
			callCount++
			mu.Unlock()
			time.Sleep(50 * time.Millisecond)
			return []float64{10.0}, nil
		},
	}

	c := collector.NewCPUCollector(mock)
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan models.Metric, 10)
	go c.Collect(ctx, ch)

	<-ch
	cancel()
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	count := callCount
	mu.Unlock()

	if count > 2 {
		t.Fatalf("expected collector to stop quickly, got %d calls", count)
	}
}

func TestCPUCollector_HandlesCPUCountError(t *testing.T) {
	mock := &mockSystemInfo{
		cpuCountFn: func(bool) (int, error) { return 0, fmt.Errorf("mock error") },
		cpuPercentFn: func(time.Duration, bool) ([]float64, error) {
			return []float64{30.0}, nil
		},
	}

	c := collector.NewCPUCollector(mock)
	results := collectWithTimeout(c, 5*time.Second)

	if len(results) < 1 {
		t.Fatal("expected metric despite CPUCount error")
	}

	payload := results[0].Data.(models.CPUPayload)
	if payload.CoreCount != 0 {
		t.Fatalf("expected core count 0 on error, got %d", payload.CoreCount)
	}
}

func TestMemoryCollector_CollectsData(t *testing.T) {
	mock := &mockSystemInfo{
		virtMemFn: func() (*mem.VirtualMemoryStat, error) {
			return &mem.VirtualMemoryStat{
				Total:     16384 * 1024 * 1024,
				Used:      8192 * 1024 * 1024,
				Available: 8192 * 1024 * 1024,
			}, nil
		},
		swapMemFn: func() (*mem.SwapMemoryStat, error) {
			return &mem.SwapMemoryStat{Used: 1024 * 1024 * 1024}, nil
		},
	}

	c := collector.NewMemoryCollector(mock)
	results := collectWithTimeout(c, 5*time.Second)

	if len(results) < 1 {
		t.Fatal("expected at least 1 metric from Memory collector")
	}

	m := results[0]
	if m.Source != types.MEM {
		t.Fatalf("expected source MEM, got %s", m.Source)
	}

	payload, ok := m.Data.(models.MemoryPayload)
	if !ok {
		t.Fatal("expected MemoryPayload")
	}
	if payload.Total != 16384*1024*1024 {
		t.Fatalf("unexpected total: %d", payload.Total)
	}
	if payload.PagefileUsage != 1024*1024*1024 {
		t.Fatalf("unexpected pagefile: %d", payload.PagefileUsage)
	}
}

func TestMemoryCollector_HandlesSwapError(t *testing.T) {
	mock := &mockSystemInfo{
		virtMemFn: func() (*mem.VirtualMemoryStat, error) {
			return &mem.VirtualMemoryStat{
				Total: 8192, Used: 4096, Available: 4096,
			}, nil
		},
		swapMemFn: func() (*mem.SwapMemoryStat, error) {
			return nil, fmt.Errorf("swap error")
		},
	}

	c := collector.NewMemoryCollector(mock)
	results := collectWithTimeout(c, 5*time.Second)

	if len(results) < 1 {
		t.Fatal("expected metric despite swap error")
	}

	payload := results[0].Data.(models.MemoryPayload)
	if payload.PagefileUsage != 0 {
		t.Fatalf("expected 0 pagefile on error, got %d", payload.PagefileUsage)
	}
}

func TestDiskCollector_CollectsData(t *testing.T) {
	mock := &mockSystemInfo{
		diskUsageFn: func(string) (*disk.UsageStat, error) {
			return &disk.UsageStat{
				Total: 500 * 1024 * 1024 * 1024,
				Used:  200 * 1024 * 1024 * 1024,
				Free:  300 * 1024 * 1024 * 1024,
			}, nil
		},
		diskIOFn: func() (map[string]disk.IOCountersStat, error) {
			return map[string]disk.IOCountersStat{
				"sda": {
					ReadBytes: 1000, WriteBytes: 2000,
					ReadCount: 10, WriteCount: 20,
				},
			}, nil
		},
	}

	c := collector.NewDiskCollector(mock)
	results := collectWithTimeout(c, 5*time.Second)

	if len(results) < 1 {
		t.Fatal("expected at least 1 metric from Disk collector")
	}

	m := results[0]
	if m.Source != types.DISK {
		t.Fatalf("expected source DISK, got %s", m.Source)
	}

	payload, ok := m.Data.(models.DiskPayload)
	if !ok {
		t.Fatal("expected DiskPayload")
	}
	if payload.Available != 300*1024*1024*1024 {
		t.Fatalf("unexpected available: %d", payload.Available)
	}
	if payload.IOStats.ReadBytes != 1000 {
		t.Fatalf("unexpected read bytes: %d", payload.IOStats.ReadBytes)
	}
}

func TestDiskCollector_AggregatesMultipleIOCounters(t *testing.T) {
	mock := &mockSystemInfo{
		diskUsageFn: func(string) (*disk.UsageStat, error) {
			return &disk.UsageStat{Total: 1000, Used: 500, Free: 500}, nil
		},
		diskIOFn: func() (map[string]disk.IOCountersStat, error) {
			return map[string]disk.IOCountersStat{
				"sda": {ReadBytes: 100, WriteBytes: 200, ReadCount: 1, WriteCount: 2},
				"sdb": {ReadBytes: 300, WriteBytes: 400, ReadCount: 3, WriteCount: 4},
			}, nil
		},
	}

	c := collector.NewDiskCollector(mock)
	results := collectWithTimeout(c, 5*time.Second)

	if len(results) < 1 {
		t.Fatal("expected metric")
	}

	payload := results[0].Data.(models.DiskPayload)
	if payload.IOStats.ReadBytes != 400 {
		t.Fatalf("expected aggregated read 400, got %d", payload.IOStats.ReadBytes)
	}
	if payload.IOStats.WriteBytes != 600 {
		t.Fatalf("expected aggregated write 600, got %d", payload.IOStats.WriteBytes)
	}
	if payload.IOStats.ReadCount != 4 {
		t.Fatalf("expected aggregated read count 4, got %d", payload.IOStats.ReadCount)
	}
}

func TestDiskCollector_HandlesDiskIOError(t *testing.T) {
	mock := &mockSystemInfo{
		diskUsageFn: func(string) (*disk.UsageStat, error) {
			return &disk.UsageStat{Total: 1000, Used: 500, Free: 500}, nil
		},
		diskIOFn: func() (map[string]disk.IOCountersStat, error) {
			return nil, fmt.Errorf("io error")
		},
	}

	c := collector.NewDiskCollector(mock)
	results := collectWithTimeout(c, 5*time.Second)

	if len(results) < 1 {
		t.Fatal("expected metric despite IO error")
	}

	payload := results[0].Data.(models.DiskPayload)
	if payload.IOStats.ReadBytes != 0 || payload.IOStats.WriteBytes != 0 {
		t.Fatal("expected zero IO stats on error")
	}
}
