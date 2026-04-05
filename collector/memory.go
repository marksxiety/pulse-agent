package collector

import (
	"context"
	"log"
	"pulse-agent/models"
	"pulse-agent/types"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryCollector struct{}

func (m MemoryCollector) Collect(ctx context.Context, ch chan<- models.Metric) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			v, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("[Memory] VirtualMemory error: %v", err)
			} else {
				ch <- models.Metric{Source: types.MEM, Value: v.UsedPercent}
			}
			time.Sleep(2 * time.Second)
		}
	}
}
