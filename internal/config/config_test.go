package config

import (
	"testing"
)

func TestConfig_Load(t *testing.T) {
	// This is a basic test to ensure the package compiles
	// TODO: Add proper tests with mock Git repository
	t.Run("basic load test", func(t *testing.T) {
		// For now, just test that we can create a config
		cfg := &Config{
			RepoPath: "/tmp/test",
			Trunk:    "main",
			Product:  "TestProduct",
		}

		if cfg.RepoPath != "/tmp/test" {
			t.Errorf("Expected RepoPath to be /tmp/test, got %s", cfg.RepoPath)
		}

		if cfg.Trunk != "main" {
			t.Errorf("Expected Trunk to be main, got %s", cfg.Trunk)
		}

		if cfg.Product != "TestProduct" {
			t.Errorf("Expected Product to be TestProduct, got %s", cfg.Product)
		}
	})
}

func TestConfig_SetProduct(t *testing.T) {
	cfg := &Config{}

	err := cfg.SetProduct("TestProduct")
	if err != nil {
		t.Errorf("SetProduct failed: %v", err)
	}

	if cfg.Product != "TestProduct" {
		t.Errorf("Expected Product to be TestProduct, got %s", cfg.Product)
	}
}

func TestConfig_SetTrunk(t *testing.T) {
	cfg := &Config{}

	err := cfg.SetTrunk("main")
	if err != nil {
		t.Errorf("SetTrunk failed: %v", err)
	}

	if cfg.Trunk != "main" {
		t.Errorf("Expected Trunk to be main, got %s", cfg.Trunk)
	}
}
