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
	live := collector.LiveSystemInfo{}

	// tea.WithAltScreen() gives bubbletea full ownership of the terminal
	// viewport — it clears the screen on every frame, which prevents the
	// "ghost content" multiplication you get when resizing without it.
	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		collector.NewCPUCollector(live).Collect(ctx, dataPipe)
	}()

	go func() {
		defer wg.Done()
		collector.NewMemoryCollector(live).Collect(ctx, dataPipe)
	}()

	go func() {
		defer wg.Done()
		collector.NewDiskCollector(live).Collect(ctx, dataPipe)
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
