package cmd

import (
	"testing"
)

func TestConfig(t *testing.T) {
	initializeConfig()

	expectedDataShards := 5
	if cfg.GetInt("dataShards") != expectedDataShards {
		t.Errorf("Expected dataShards to be %d, but got %d", expectedDataShards, cfg.GetInt("dataShards"))
	}

	expectedParityShards := 5
	if cfg.GetInt("dataShards") != expectedParityShards {
		t.Errorf("Expected dataShards to be %d, but got %d", expectedParityShards, cfg.GetInt("dataShards"))
	}
}

func BenchmarkConfig(b *testing.B) {
	initializeConfig()
	for i := 0; i < b.N; i++ {
		cfg.GetInt("dataShards")
	}
}
