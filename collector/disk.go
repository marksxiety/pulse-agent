package collector

import (
	"context"
	"log"
	"pulse-agent/models"
	"pulse-agent/types"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

type DiskCollector struct{}

func (d DiskCollector) Collect(ctx context.Context, ch chan<- models.Metric) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			usage, err := disk.Usage("/")
			if err != nil {
				log.Printf("Error occurred while collecting disk usage: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			ioCounters, err := disk.IOCounters()
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
