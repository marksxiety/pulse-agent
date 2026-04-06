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
			vMem, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("[Memory] VirtualMemory error: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}

			var pageFileUsed uint64
			vSwap, err := mem.SwapMemory()
			if err != nil {
				log.Printf("[Memory] SwapMemory error: %v", err)
			} else {
				pageFileUsed = vSwap.Used
			}

			ch <- models.Metric{Source: types.MEM, Data: models.MemoryPayload{
				Total:         vMem.Total,
				Used:          vMem.Used,
				Available:     vMem.Available,
				PagefileUsage: pageFileUsed,
			}}
			time.Sleep(2 * time.Second)
		}
	}
}
