package main

import (
	"fmt"
	"pulse-agent/collector"
	"pulse-agent/models"
)

func main() {
	dataPipe := make(chan models.Metric)

	fmt.Println("Pulse agent started")

	// start collectors in separate `goroutines`
	go collector.CPUCollector{}.Collect(dataPipe)
	go collector.MemoryCollector{}.Collect(dataPipe)
	go collector.DiskCollector{}.Collect(dataPipe)

	// receive data from the channel
	for m := range dataPipe {
		fmt.Printf("[%s] Usage: %.2f%%\n", m.Source, m.Value)
	}
}
