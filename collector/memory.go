package collector

import (
	"context"
	"log"
	"pulse-agent/models"
	"pulse-agent/types"
	"time"
)

type MemoryCollector struct {
	sys SystemInfo
}

func NewMemoryCollector(sys SystemInfo) MemoryCollector {
	return MemoryCollector{sys: sys}
}

func (m MemoryCollector) Collect(ctx context.Context, ch chan<- models.Metric) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			vMem, err := m.sys.VirtualMemory()
			if err != nil {
				log.Printf("[Memory] VirtualMemory error: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}

			var pageFileUsed uint64
			vSwap, err := m.sys.SwapMemory()
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
