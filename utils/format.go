package utils

import (
	"fmt"
	"time"
)

func BytesToGB(b uint64) string { return fmt.Sprintf("%.2f GB", float64(b)/1e9) }

func BytesToMB(b uint64) string { return fmt.Sprintf("%.1f MB", float64(b)/1e6) }

func FormatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", h, m, s)
	}
	return fmt.Sprintf("%dm %02ds", m, s)
}

func PadToHeight(lines []string, n int) []string {
	for len(lines) < n {
		lines = append(lines, "")
	}
	return lines
}
