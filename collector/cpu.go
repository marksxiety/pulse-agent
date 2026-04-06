package collector

import (
	"context"
	"log"
	"pulse-agent/models"
	"pulse-agent/types"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUCollector struct{}

func (c CPUCollector) Collect(ctx context.Context, ch chan<- models.Metric) {

	cores, err := cpu.Counts(true)
	if err != nil {
		cores = 0
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			v, err := cpu.Percent(time.Second, false)
			if err != nil {
				log.Printf("[CPU] Percent error: %v", err)
			} else if len(v) > 0 {
				ch <- models.Metric{
					Source: types.CPU,
					Data:   models.CPUPayload{Percentage: v[0], CoreCount: cores},
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}
