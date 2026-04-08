package models_test

import (
	"testing"

	"pulse-agent/models"
	"pulse-agent/types"
)

func TestCPUPayload_SourceType(t *testing.T) {
	p := models.CPUPayload{Percentage: 50.0, CoreCount: 8}
	if p.SourceType() != types.CPU {
		t.Fatalf("expected CPU, got %s", p.SourceType())
	}
}

func TestMemoryPayload_SourceType(t *testing.T) {
	p := models.MemoryPayload{Total: 16384, Used: 8192}
	if p.SourceType() != types.MEM {
		t.Fatalf("expected MEM, got %s", p.SourceType())
	}
}

func TestDiskPayload_SourceType(t *testing.T) {
	p := models.DiskPayload{Total: 500000, Used: 200000}
	if p.SourceType() != types.DISK {
		t.Fatalf("expected DISK, got %s", p.SourceType())
	}
}
