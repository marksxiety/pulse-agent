package collector

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo interface {
	CPUPercent(interval time.Duration, percpu bool) ([]float64, error)
	CPUCount(logical bool) (int, error)
	VirtualMemory() (*mem.VirtualMemoryStat, error)
	SwapMemory() (*mem.SwapMemoryStat, error)
	DiskUsage(path string) (*disk.UsageStat, error)
	DiskIOCounters() (map[string]disk.IOCountersStat, error)
}

type LiveSystemInfo struct{}

func (LiveSystemInfo) CPUPercent(interval time.Duration, percpu bool) ([]float64, error) {
	return cpu.Percent(interval, percpu)
}

func (LiveSystemInfo) CPUCount(logical bool) (int, error) {
	return cpu.Counts(logical)
}

func (LiveSystemInfo) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

func (LiveSystemInfo) SwapMemory() (*mem.SwapMemoryStat, error) {
	return mem.SwapMemory()
}

func (LiveSystemInfo) DiskUsage(path string) (*disk.UsageStat, error) {
	return disk.Usage(path)
}

func (LiveSystemInfo) DiskIOCounters() (map[string]disk.IOCountersStat, error) {
	return disk.IOCounters()
}
