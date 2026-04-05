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
			v, err := disk.Usage("/")
			if err != nil {
				log.Printf("Error occurred while collecting disk usage: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			ch <- models.Metric{
				Source: types.DISK,
				Value:  v.UsedPercent,
			}
			time.Sleep(1 * time.Second)
		}
	}
}
