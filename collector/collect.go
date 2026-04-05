package collector

import (
	"context"
	"pulse-agent/models"
)

// this will serve as the interface for all collectors to implement
type Collector interface {
	Collect(ctx context.Context, ch chan<- models.Metric)
}
