package models

import "pulse-agent/types"

type Payload interface {
	SourceType() types.Source
}

type CPUPayload struct {
	Percentage float64 `json:"percentage"`
	CoreCount  int     `json:"core_count"`
}

type MemoryPayload struct {
	Total         uint64 `json:"total"`
	Used          uint64 `json:"used"`
	Available     uint64 `json:"available"`
	PagefileUsage uint64 `json:"pagefile_usage"`
}

type DiskPayload struct {
	Total     uint64      `json:"total"`
	Used      uint64      `json:"used"`
	Available uint64      `json:"available"`
	IOStats   DiskIOStats `json:"io_stats"`
}

type DiskIOStats struct {
	ReadBytes  uint64 `json:"read_bytes"`
	WriteBytes uint64 `json:"write_bytes"`
	ReadCount  uint64 `json:"read_count"`
	WriteCount uint64 `json:"write_count"`
}

func (c CPUPayload) SourceType() types.Source {
	return types.CPU
}

func (m MemoryPayload) SourceType() types.Source {
	return types.MEM
}

func (d DiskPayload) SourceType() types.Source {
	return types.DISK
}
