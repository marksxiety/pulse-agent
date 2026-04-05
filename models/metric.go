package models

import "pulse-agent/types"

type Metric struct {
	Source types.Source
	Value  float64
}
