package models

import (
	"encoding/json"
	"pulse-agent/types"
)

type Metric struct {
	Source types.Source `json:"source"`
	Data   Payload      `json:"-"`
}

type NewDataMsg Metric

func (m Metric) MarshalJSON() ([]byte, error) {
	switch p := m.Data.(type) {
	case CPUPayload:
		return json.Marshal(struct {
			Source types.Source `json:"source"`
			CPU    CPUPayload   `json:"cpu"`
		}{m.Source, p})
	case MemoryPayload:
		return json.Marshal(struct {
			Source types.Source  `json:"source"`
			Memory MemoryPayload `json:"memory"`
		}{m.Source, p})
	case DiskPayload:
		return json.Marshal(struct {
			Source types.Source `json:"source"`
			Disk   DiskPayload  `json:"disk"`
		}{m.Source, p})
	default:
		return json.Marshal(struct {
			Source types.Source `json:"source"`
		}{m.Source})
	}
}
