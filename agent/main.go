package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pulse-agent/collector"
	"pulse-agent/models"
	"pulse-agent/ui"
	"sync"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataPipe := make(chan models.Metric)

	p := tea.NewProgram(ui.InitialModel())

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

	go func() {
		for m := range dataPipe {
			p.Send(models.NewDataMsg(m))
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
		wg.Wait()
		close(dataPipe)
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
