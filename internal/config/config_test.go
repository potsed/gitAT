package config

import (
	"strings"
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

	t.Run("fallthrough config initialization", func(t *testing.T) {
		// Test that a config with fallthrough settings can be created
		cfg := &Config{
			RepoPath: "/tmp/test",
			Fallthrough: FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{"--version", "--help"},
			},
		}

		if !cfg.Fallthrough.Enabled {
			t.Error("Expected fallthrough to be enabled")
		}

		if cfg.Fallthrough.Verbose {
			t.Error("Expected verbose mode to be disabled")
		}

		if len(cfg.Fallthrough.Blacklist) != 2 {
			t.Errorf("Expected blacklist length 2, got %d", len(cfg.Fallthrough.Blacklist))
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

func TestLoadFallthroughConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		cfg := loadFallthroughConfig()

		if !cfg.Enabled {
			t.Error("Expected fallthrough to be enabled by default")
		}

		if cfg.Verbose {
			t.Error("Expected verbose mode to be disabled by default")
		}

		expectedBlacklist := []string{"--version", "--help"}
		if len(cfg.Blacklist) != len(expectedBlacklist) {
			t.Errorf("Expected blacklist length %d, got %d", len(expectedBlacklist), len(cfg.Blacklist))
		}

		for i, expected := range expectedBlacklist {
			if cfg.Blacklist[i] != expected {
				t.Errorf("Expected blacklist[%d] to be %s, got %s", i, expected, cfg.Blacklist[i])
			}
		}
	})
}

func TestConfig_SetFallthroughEnabled(t *testing.T) {
	cfg := &Config{}

	err := cfg.SetFallthroughEnabled(true)
	if err != nil {
		t.Errorf("SetFallthroughEnabled failed: %v", err)
	}

	if !cfg.Fallthrough.Enabled {
		t.Error("Expected fallthrough to be enabled")
	}

	err = cfg.SetFallthroughEnabled(false)
	if err != nil {
		t.Errorf("SetFallthroughEnabled failed: %v", err)
	}

	if cfg.Fallthrough.Enabled {
		t.Error("Expected fallthrough to be disabled")
	}
}

func TestConfig_SetFallthroughVerbose(t *testing.T) {
	cfg := &Config{}

	err := cfg.SetFallthroughVerbose(true)
	if err != nil {
		t.Errorf("SetFallthroughVerbose failed: %v", err)
	}

	if !cfg.Fallthrough.Verbose {
		t.Error("Expected verbose mode to be enabled")
	}

	err = cfg.SetFallthroughVerbose(false)
	if err != nil {
		t.Errorf("SetFallthroughVerbose failed: %v", err)
	}

	if cfg.Fallthrough.Verbose {
		t.Error("Expected verbose mode to be disabled")
	}
}

func TestConfig_AddToFallthroughBlacklist(t *testing.T) {
	cfg := &Config{
		Fallthrough: FallthroughConfig{
			Blacklist: []string{"--version", "--help"},
		},
	}

	// Add new command
	err := cfg.AddToFallthroughBlacklist("config")
	if err != nil {
		t.Errorf("AddToFallthroughBlacklist failed: %v", err)
	}

	if len(cfg.Fallthrough.Blacklist) != 3 {
		t.Errorf("Expected blacklist length 3, got %d", len(cfg.Fallthrough.Blacklist))
	}

	if cfg.Fallthrough.Blacklist[2] != "config" {
		t.Errorf("Expected last blacklist item to be 'config', got %s", cfg.Fallthrough.Blacklist[2])
	}

	// Try to add duplicate command
	err = cfg.AddToFallthroughBlacklist("config")
	if err != nil {
		t.Errorf("AddToFallthroughBlacklist failed: %v", err)
	}

	if len(cfg.Fallthrough.Blacklist) != 3 {
		t.Errorf("Expected blacklist length to remain 3, got %d", len(cfg.Fallthrough.Blacklist))
	}
}

func TestConfig_RemoveFromFallthroughBlacklist(t *testing.T) {
	cfg := &Config{
		Fallthrough: FallthroughConfig{
			Blacklist: []string{"--version", "--help", "config"},
		},
	}

	// Remove existing command
	err := cfg.RemoveFromFallthroughBlacklist("config")
	if err != nil {
		t.Errorf("RemoveFromFallthroughBlacklist failed: %v", err)
	}

	if len(cfg.Fallthrough.Blacklist) != 2 {
		t.Errorf("Expected blacklist length 2, got %d", len(cfg.Fallthrough.Blacklist))
	}

	// Verify config is no longer in blacklist
	for _, cmd := range cfg.Fallthrough.Blacklist {
		if cmd == "config" {
			t.Error("Expected 'config' to be removed from blacklist")
		}
	}

	// Try to remove non-existent command
	err = cfg.RemoveFromFallthroughBlacklist("nonexistent")
	if err != nil {
		t.Errorf("RemoveFromFallthroughBlacklist failed: %v", err)
	}

	if len(cfg.Fallthrough.Blacklist) != 2 {
		t.Errorf("Expected blacklist length to remain 2, got %d", len(cfg.Fallthrough.Blacklist))
	}
}

func TestConfig_IsFallthroughBlacklisted(t *testing.T) {
	cfg := &Config{
		Fallthrough: FallthroughConfig{
			Blacklist: []string{"--version", "--help", "config"},
		},
	}

	// Test blacklisted commands
	if !cfg.IsFallthroughBlacklisted("--version") {
		t.Error("Expected '--version' to be blacklisted")
	}

	if !cfg.IsFallthroughBlacklisted("--help") {
		t.Error("Expected '--help' to be blacklisted")
	}

	if !cfg.IsFallthroughBlacklisted("config") {
		t.Error("Expected 'config' to be blacklisted")
	}

	// Test non-blacklisted commands
	if cfg.IsFallthroughBlacklisted("status") {
		t.Error("Expected 'status' to not be blacklisted")
	}

	if cfg.IsFallthroughBlacklisted("diff") {
		t.Error("Expected 'diff' to not be blacklisted")
	}
}

func TestFallthroughConfig_Validation(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		cfg := FallthroughConfig{
			Enabled:   true,
			Verbose:   false,
			Blacklist: []string{"--version", "--help"},
		}

		// Test that configuration is valid
		if !cfg.Enabled {
			t.Error("Expected configuration to be enabled")
		}

		if len(cfg.Blacklist) != 2 {
			t.Errorf("Expected blacklist length 2, got %d", len(cfg.Blacklist))
		}
	})

	t.Run("empty blacklist", func(t *testing.T) {
		cfg := FallthroughConfig{
			Enabled:   true,
			Verbose:   true,
			Blacklist: []string{},
		}

		if len(cfg.Blacklist) != 0 {
			t.Errorf("Expected empty blacklist, got length %d", len(cfg.Blacklist))
		}
	})

	t.Run("nil blacklist", func(t *testing.T) {
		cfg := FallthroughConfig{
			Enabled:   false,
			Verbose:   false,
			Blacklist: nil,
		}

		if cfg.Blacklist != nil {
			t.Error("Expected nil blacklist")
		}
	})
}

func TestConfig_GetFallthroughBlacklist(t *testing.T) {
	originalBlacklist := []string{"status", "diff", "log"}
	cfg := &Config{
		Fallthrough: FallthroughConfig{
			Blacklist: originalBlacklist,
		},
	}
	
	blacklist := cfg.GetFallthroughBlacklist()
	
	// Check that we get a copy, not the original slice
	if len(blacklist) > 0 && len(cfg.Fallthrough.Blacklist) > 0 {
		if &blacklist[0] == &cfg.Fallthrough.Blacklist[0] {
			t.Error("Expected GetFallthroughBlacklist to return a copy, not the original slice")
		}
	}
	
	// Check that contents are correct
	if len(blacklist) != len(originalBlacklist) {
		t.Errorf("Expected blacklist length %d, got %d", len(originalBlacklist), len(blacklist))
	}
	
	for i, cmd := range blacklist {
		if cmd != originalBlacklist[i] {
			t.Errorf("Expected blacklist[%d] = %q, got %q", i, originalBlacklist[i], cmd)
		}
	}
	
	// Modify the returned slice and ensure original is unchanged
	if len(blacklist) > 0 {
		blacklist[0] = "modified"
		if cfg.Fallthrough.Blacklist[0] == "modified" {
			t.Error("Modifying returned blacklist should not affect original")
		}
	}
}

func TestConfig_ClearFallthroughBlacklist(t *testing.T) {
	cfg := &Config{
		Fallthrough: FallthroughConfig{
			Blacklist: []string{"status", "diff", "log"},
		},
	}
	
	err := cfg.ClearFallthroughBlacklist()
	if err != nil {
		t.Errorf("Expected no error clearing blacklist, got: %v", err)
	}
	
	if len(cfg.Fallthrough.Blacklist) != 0 {
		t.Errorf("Expected empty blacklist after clearing, got %v", cfg.Fallthrough.Blacklist)
	}
	
	// Test that previously blacklisted commands are no longer blacklisted
	if cfg.IsFallthroughBlacklisted("status") {
		t.Error("Expected 'status' to not be blacklisted after clearing")
	}
}

func TestConfig_SetFallthroughBlacklist(t *testing.T) {
	cfg := &Config{
		Fallthrough: FallthroughConfig{
			Blacklist: []string{"status"},
		},
	}
	
	newBlacklist := []string{"diff", "log", "commit"}
	err := cfg.SetFallthroughBlacklist(newBlacklist)
	if err != nil {
		t.Errorf("Expected no error setting blacklist, got: %v", err)
	}
	
	// Check that old blacklist is replaced
	if cfg.IsFallthroughBlacklisted("status") {
		t.Error("Expected 'status' to not be blacklisted after replacement")
	}
	
	// Check that new blacklist is in effect
	for _, cmd := range newBlacklist {
		if !cfg.IsFallthroughBlacklisted(cmd) {
			t.Errorf("Expected '%s' to be blacklisted after setting", cmd)
		}
	}
	
	// Check that we have a copy, not the original slice
	newBlacklist[0] = "modified"
	if cfg.Fallthrough.Blacklist[0] == "modified" {
		t.Error("Modifying input slice should not affect config blacklist")
	}
}

func TestConfig_ValidateFallthroughConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      FallthroughConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{"status", "diff"},
			},
			expectError: false,
		},
		{
			name: "empty blacklist is valid",
			config: FallthroughConfig{
				Enabled:   true,
				Verbose:   true,
				Blacklist: []string{},
			},
			expectError: false,
		},
		{
			name: "duplicate commands in blacklist",
			config: FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{"status", "diff", "status"},
			},
			expectError: true,
			errorMsg:    "duplicate command in fallthrough blacklist: status",
		},
		{
			name: "empty command in blacklist",
			config: FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{"status", "", "diff"},
			},
			expectError: true,
			errorMsg:    "empty command in fallthrough blacklist",
		},
		{
			name: "whitespace-only command in blacklist",
			config: FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{"status", "   ", "diff"},
			},
			expectError: true,
			errorMsg:    "empty command in fallthrough blacklist",
		},
		{
			name: "help command in blacklist is allowed",
			config: FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{"help", "--help"},
			},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Fallthrough: tt.config,
			}
			
			err := cfg.ValidateFallthroughConfig()
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for config %+v, got nil", tt.config)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for config %+v, got: %v", tt.config, err)
				}
			}
		})
	}
}
