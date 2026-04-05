package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pulse-agent/collector"
	"pulse-agent/models"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataPipe := make(chan models.Metric)

	fmt.Println("Pulse agent started")

	// Start collectors in separate goroutines
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		collector.CPUCollector{}.Collect(ctx, dataPipe)
	}()

	go func() {
		defer wg.Done()
		collector.MemoryCollector{}.Collect(ctx, dataPipe)
	}()

	go func() {
		defer wg.Done()
		collector.DiskCollector{}.Collect(ctx, dataPipe)
	}()

	// Listen for termination signals to gracefully shut down (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
		wg.Wait()
		close(dataPipe)
	}()

	for m := range dataPipe {
		fmt.Printf("[%s] Usage: %.2f%%\n", m.Source, m.Value)
	}
}
