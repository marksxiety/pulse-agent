package collector

import (
	"log"
	"pulse-agent/models"
	"pulse-agent/types"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryCollector struct{}

func (m MemoryCollector) Collect(ch chan<- models.Metric) {
	for {
		v, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("[Memory] VirtualMemory error: %v", err)
		} else {
			ch <- models.Metric{Source: types.MEM, Value: v.UsedPercent}
		}
		time.Sleep(2 * time.Second)
	}
}
