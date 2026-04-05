package collector

import "pulse-agent/models"

// this will serve as the interface for all collectors to implement
type Collector interface {
	Collect(ch chan<- models.Metric)
}
