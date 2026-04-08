package models_test

import (
	"encoding/json"
	"testing"

	"pulse-agent/models"
	"pulse-agent/types"
)

func TestMetric_MarshalJSON_CPU(t *testing.T) {
	m := models.Metric{
		Source: types.CPU,
		Data:   models.CPUPayload{Percentage: 45.5, CoreCount: 8},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := raw["source"]; !ok {
		t.Fatal("missing 'source' field")
	}
	if _, ok := raw["cpu"]; !ok {
		t.Fatal("missing 'cpu' field")
	}
	if _, ok := raw["memory"]; ok {
		t.Fatal("unexpected 'memory' field")
	}
}

func TestMetric_MarshalJSON_Memory(t *testing.T) {
	m := models.Metric{
		Source: types.MEM,
		Data: models.MemoryPayload{
			Total:     16384,
			Used:      8192,
			Available: 8192,
		},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := raw["memory"]; !ok {
		t.Fatal("missing 'memory' field")
	}
	if _, ok := raw["cpu"]; ok {
		t.Fatal("unexpected 'cpu' field")
	}
}

func TestMetric_MarshalJSON_Disk(t *testing.T) {
	m := models.Metric{
		Source: types.DISK,
		Data: models.DiskPayload{
			Total: 500000000000,
			Used:  200000000000,
		},
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := raw["disk"]; !ok {
		t.Fatal("missing 'disk' field")
	}
}

func TestMetric_MarshalJSON_UnknownPayload(t *testing.T) {
	m := models.Metric{
		Source: types.CPU,
		Data:   nil,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := raw["source"]; !ok {
		t.Fatal("missing 'source' field")
	}
	if _, ok := raw["cpu"]; ok {
		t.Fatal("unexpected 'cpu' field for nil payload")
	}
	if _, ok := raw["memory"]; ok {
		t.Fatal("unexpected 'memory' field for nil payload")
	}
	if _, ok := raw["disk"]; ok {
		t.Fatal("unexpected 'disk' field for nil payload")
	}
}
