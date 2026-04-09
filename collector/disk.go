package collector

import (
	"context"
	"log"
	"pulse-agent/models"
	"pulse-agent/types"
	"runtime"
	"time"
)

type DiskCollector struct {
	sys     SystemInfo
	rootDir string
}

func diskRootPath() string {
	if runtime.GOOS == "windows" {
		return "C:\\"
	}
	return "/"
}

func NewDiskCollector(sys SystemInfo) DiskCollector {
	return DiskCollector{sys: sys, rootDir: diskRootPath()}
}

func (d DiskCollector) Collect(ctx context.Context, ch chan<- models.Metric) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			usage, err := d.sys.DiskUsage(d.rootDir)
			if err != nil {
				log.Printf("Error occurred while collecting disk usage: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			ioCounters, err := d.sys.DiskIOCounters()
			var ioStats models.DiskIOStats

			if err == nil && len(ioCounters) > 0 {
				for _, counter := range ioCounters {
					ioStats.ReadBytes += counter.ReadBytes
					ioStats.WriteBytes += counter.WriteBytes
					ioStats.ReadCount += counter.ReadCount
					ioStats.WriteCount += counter.WriteCount
				}
			}

			ch <- models.Metric{Source: types.DISK, Data: models.DiskPayload{
				Total:     usage.Total,
				Used:      usage.Used,
				Available: usage.Free,
				IOStats:   ioStats,
			}}

			time.Sleep(1 * time.Second)
		}
	}
}
