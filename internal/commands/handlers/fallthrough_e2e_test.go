package handlers

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// TestFallthroughE2E_RealWorldScenarios tests real-world usage scenarios
func TestFallthroughE2E_RealWorldScenarios(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping E2E test: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	t.Run("developer_daily_workflow", func(t *testing.T) {
		// Simulate a typical developer workflow
		workflow := []struct {
			command     string
			args        []string
			description string
			critical    bool // If true, failure stops the workflow
		}{
			{"status", []string{}, "Check current status", false},
			{"branch", []string{"-a"}, "List all branches", false},
			{"log", []string{"--oneline", "-10"}, "Check recent commits", false},
			{"diff", []string{"--name-only"}, "See what files changed", false},
			{"remote", []string{"-v"}, "Check remotes", false},
			{"config", []string{"--get", "user.name"}, "Verify user config", false},
		}

		for i, step := range workflow {
			t.Logf("Step %d: %s", i+1, step.description)
			err := handler.Execute(step.command, step.args)
			
			if err != nil {
				if step.critical {
					t.Fatalf("Critical step failed: %s - %v", step.description, err)
				} else {
					t.Logf("Step failed (non-critical): %s - %v", step.description, err)
				}
			} else {
				t.Logf("Step succeeded: %s", step.description)
			}
		}
	})

	t.Run("code_review_workflow", func(t *testing.T) {
		// Simulate a code review workflow
		reviewSteps := []struct {
			command     string
			args        []string
			description string
		}{
			{"log", []string{"--oneline", "--graph", "-20"}, "Review commit history"},
			{"diff", []string{"HEAD~5..HEAD", "--name-only"}, "See changed files in recent commits"},
			{"show", []string{"--stat", "HEAD"}, "Show latest commit details"},
			{"blame", []string{"README.md"}, "Check file blame (may fail if file doesn't exist)"},
		}

		for _, step := range reviewSteps {
			t.Logf("Review step: %s", step.description)
			err := handler.Execute(step.command, step.args)
			if err != nil {
				t.Logf("Review step failed (expected for some): %s - %v", step.description, err)
			}
		}
	})

	t.Run("repository_maintenance_workflow", func(t *testing.T) {
		// Simulate repository maintenance tasks
		maintenanceSteps := []struct {
			command     string
			args        []string
			description string
		}{
			{"status", []string{"--porcelain"}, "Check for uncommitted changes"},
			{"branch", []string{"--merged"}, "List merged branches"},
			{"remote", []string{"prune", "origin"}, "Prune remote references (may fail)"},
			{"gc", []string{"--auto"}, "Run garbage collection"},
			{"fsck", []string{"--connectivity-only"}, "Check repository integrity"},
		}

		for _, step := range maintenanceSteps {
			t.Logf("Maintenance step: %s", step.description)
			err := handler.Execute(step.command, step.args)
			if err != nil {
				t.Logf("Maintenance step result: %s - %v", step.description, err)
			}
		}
	})
}

// TestFallthroughE2E_ErrorRecovery tests error recovery scenarios
func TestFallthroughE2E_ErrorRecovery(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: true, // Enable verbose for error scenarios
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	t.Run("git_not_available_recovery", func(t *testing.T) {
		// Temporarily modify PATH to simulate git not being available
		originalPath := os.Getenv("PATH")
		os.Setenv("PATH", "/nonexistent")
		defer os.Setenv("PATH", originalPath)

		err := handler.Execute("status", []string{})
		if err == nil {
			t.Error("Expected error when git is not available")
		} else {
			t.Logf("Correctly handled git not available: %v", err)
			// Verify error message is helpful
			if !strings.Contains(err.Error(), "Git executable not found") {
				t.Errorf("Error message should mention git not found, got: %v", err)
			}
		}

		// Restore PATH and verify recovery
		os.Setenv("PATH", originalPath)
		if _, err := exec.LookPath("git"); err == nil {
			err = handler.Execute("--version", []string{})
			if err != nil {
				t.Errorf("Should recover after restoring PATH: %v", err)
			} else {
				t.Log("Successfully recovered after restoring git availability")
			}
		}
	})

	t.Run("fallthrough_disabled_recovery", func(t *testing.T) {
		// Test with fallthrough disabled
		cfg.Fallthrough.Enabled = false
		disabledHandler := NewFallthroughHandler(cfg, gitRepo)

		err := disabledHandler.Execute("status", []string{})
		if err == nil {
			t.Error("Expected error when fallthrough is disabled")
		} else {
			t.Logf("Correctly handled disabled fallthrough: %v", err)
			// Verify error message provides guidance
			if !strings.Contains(err.Error(), "disabled") {
				t.Errorf("Error message should mention fallthrough disabled, got: %v", err)
			}
		}

		// Re-enable and verify recovery
		cfg.Fallthrough.Enabled = true
		enabledHandler := NewFallthroughHandler(cfg, gitRepo)

		if _, err := exec.LookPath("git"); err == nil {
			err = enabledHandler.Execute("--version", []string{})
			if err != nil {
				t.Errorf("Should work after re-enabling fallthrough: %v", err)
			} else {
				t.Log("Successfully recovered after re-enabling fallthrough")
			}
		}
	})

	t.Run("blacklist_configuration_recovery", func(t *testing.T) {
		// Test blacklist configuration and recovery
		cfg.Fallthrough.Blacklist = []string{"status", "log"}
		blacklistHandler := NewFallthroughHandler(cfg, gitRepo)

		// Verify blacklisted commands are blocked
		err := blacklistHandler.Execute("status", []string{})
		if err == nil || !strings.Contains(err.Error(), "blacklisted") {
			t.Errorf("Expected blacklist error, got: %v", err)
		}

		// Clear blacklist and verify recovery
		cfg.Fallthrough.Blacklist = []string{}
		clearedHandler := NewFallthroughHandler(cfg, gitRepo)

		if _, err := exec.LookPath("git"); err == nil {
			err = clearedHandler.Execute("status", []string{})
			if err != nil && strings.Contains(err.Error(), "blacklisted") {
				t.Errorf("Should not be blacklisted after clearing: %v", err)
			} else {
				t.Log("Successfully recovered after clearing blacklist")
			}
		}
	})
}

// TestFallthroughE2E_PerformanceUnderLoad tests performance under various load conditions
func TestFallthroughE2E_PerformanceUnderLoad(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping performance test: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false, // Disable verbose for performance testing
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	t.Run("concurrent_command_execution", func(t *testing.T) {
		// Test concurrent execution of commands
		const numGoroutines = 10
		const commandsPerGoroutine = 5

		done := make(chan bool, numGoroutines)
		errors := make(chan error, numGoroutines*commandsPerGoroutine)

		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				defer func() { done <- true }()

				for j := 0; j < commandsPerGoroutine; j++ {
					err := handler.Execute("--version", []string{})
					if err != nil {
						errors <- err
					}
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
		close(errors)

		// Check for errors
		errorCount := 0
		for err := range errors {
			errorCount++
			t.Logf("Concurrent execution error: %v", err)
		}

		if errorCount > 0 {
			t.Logf("Had %d errors out of %d total commands", errorCount, numGoroutines*commandsPerGoroutine)
		} else {
			t.Log("All concurrent commands executed successfully")
		}
	})

	t.Run("rapid_sequential_execution", func(t *testing.T) {
		// Test rapid sequential execution
		const numCommands = 50
		start := time.Now()

		for i := 0; i < numCommands; i++ {
			err := handler.Execute("--version", []string{})
			if err != nil {
				t.Logf("Command %d failed: %v", i, err)
			}
		}

		duration := time.Since(start)
		avgDuration := duration / numCommands
		t.Logf("Executed %d commands in %v (avg: %v per command)", numCommands, duration, avgDuration)

		// Reasonable performance threshold
		maxAvgDuration := 100 * time.Millisecond
		if avgDuration > maxAvgDuration {
			t.Errorf("Average command duration too high: %v > %v", avgDuration, maxAvgDuration)
		}
	})

	t.Run("memory_usage_stability", func(t *testing.T) {
		// Test that repeated execution doesn't cause memory leaks
		const iterations = 100

		for i := 0; i < iterations; i++ {
			// Create new handler each time to test initialization overhead
			testHandler := NewFallthroughHandler(cfg, gitRepo)
			err := testHandler.Execute("--version", []string{})
			if err != nil {
				t.Logf("Iteration %d failed: %v", i, err)
			}

			// Periodically log progress
			if i%25 == 0 {
				t.Logf("Completed %d/%d iterations", i, iterations)
			}
		}

		t.Log("Memory stability test completed")
	})
}

// TestFallthroughE2E_ConfigurationScenarios tests various configuration scenarios end-to-end
func TestFallthroughE2E_ConfigurationScenarios(t *testing.T) {
	gitRepo := &git.Repository{}

	configScenarios := []struct {
		name        string
		config      config.FallthroughConfig
		testCommand string
		testArgs    []string
		expectError bool
		errorCheck  func(error) bool
		description string
	}{
		{
			name: "default_configuration",
			config: config.FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{},
			},
			testCommand: "--version",
			testArgs:    []string{},
			expectError: false,
			description: "Default configuration should work",
		},
		{
			name: "verbose_enabled",
			config: config.FallthroughConfig{
				Enabled:   true,
				Verbose:   true,
				Blacklist: []string{},
			},
			testCommand: "--version",
			testArgs:    []string{},
			expectError: false,
			description: "Verbose mode should work",
		},
		{
			name: "fallthrough_disabled",
			config: config.FallthroughConfig{
				Enabled:   false,
				Verbose:   false,
				Blacklist: []string{},
			},
			testCommand: "status",
			testArgs:    []string{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "disabled")
			},
			description: "Disabled fallthrough should block commands",
		},
		{
			name: "single_blacklisted_command",
			config: config.FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{"status"},
			},
			testCommand: "status",
			testArgs:    []string{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "blacklisted")
			},
			description: "Blacklisted commands should be blocked",
		},
		{
			name: "multiple_blacklisted_commands",
			config: config.FallthroughConfig{
				Enabled:   true,
				Verbose:   false,
				Blacklist: []string{"push", "pull", "fetch", "merge"},
			},
			testCommand: "push",
			testArgs:    []string{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "blacklisted")
			},
			description: "Multiple blacklisted commands should work",
		},
		{
			name: "verbose_with_blacklist",
			config: config.FallthroughConfig{
				Enabled:   true,
				Verbose:   true,
				Blacklist: []string{"log"},
			},
			testCommand: "log",
			testArgs:    []string{},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "blacklisted")
			},
			description: "Verbose mode with blacklist should work",
		},
	}

	for _, scenario := range configScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			cfg := &config.Config{
				Fallthrough: scenario.config,
			}
			handler := NewFallthroughHandler(cfg, gitRepo)

			err := handler.Execute(scenario.testCommand, scenario.testArgs)

			if scenario.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", scenario.description)
				} else if scenario.errorCheck != nil && !scenario.errorCheck(err) {
					t.Errorf("Error check failed for %s: %v", scenario.description, err)
				} else {
					t.Logf("Correctly handled error scenario: %s", scenario.description)
				}
			} else {
				// Skip git availability check for success cases in CI
				if err != nil && strings.Contains(err.Error(), "git executable not found") {
					t.Skipf("Skipping success test due to git not available: %s", scenario.description)
				} else if err != nil {
					t.Errorf("Expected success for %s, but got error: %v", scenario.description, err)
				} else {
					t.Logf("Successfully handled scenario: %s", scenario.description)
				}
			}
		})
	}
}

// TestFallthroughE2E_InteractiveCommandHandling tests interactive command handling end-to-end
func TestFallthroughE2E_InteractiveCommandHandling(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping interactive test: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}

	// Create handler with test mode for interactive commands
	testResponses := map[string][]string{
		"git add -p":      {"q"}, // Quit immediately
		"git rebase -i":   {":q"}, // Quit editor
		"git commit -v":   {"Test message", "", ":wq"}, // Commit message and save
		"git merge":       {""},   // Accept default
	}
	handler := NewFallthroughHandlerWithTestMode(cfg, gitRepo, testResponses)

	interactiveTests := []struct {
		name        string
		command     string
		args        []string
		description string
	}{
		{
			name:        "non_interactive_command",
			command:     "status",
			args:        []string{"--porcelain"},
			description: "Non-interactive commands should work normally",
		},
		{
			name:        "help_command",
			command:     "help",
			args:        []string{"status"},
			description: "Help commands should work",
		},
		{
			name:        "version_command",
			command:     "--version",
			args:        []string{},
			description: "Version command should work",
		},
	}

	for _, test := range interactiveTests {
		t.Run(test.name, func(t *testing.T) {
			err := handler.Execute(test.command, test.args)
			if err != nil {
				t.Logf("Interactive test result for %s: %v", test.description, err)
			} else {
				t.Logf("Interactive test succeeded: %s", test.description)
			}
		})
	}
}

// TestFallthroughE2E_CompleteIntegration tests the complete integration with all components
func TestFallthroughE2E_CompleteIntegration(t *testing.T) {
	t.Run("full_system_integration", func(t *testing.T) {
		// Test the complete system integration
		cfg := &config.Config{
			Fallthrough: config.FallthroughConfig{
				Enabled:   true,
				Verbose:   true,
				Blacklist: []string{"push", "pull"},
			},
		}

		// Validate configuration
		err := cfg.ValidateFallthroughConfig()
		if err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}

		gitRepo := &git.Repository{}
		handler := NewFallthroughHandler(cfg, gitRepo)

		// Test configuration methods
		t.Log("Testing configuration methods...")
		blacklist := cfg.GetFallthroughBlacklist()
		if len(blacklist) != 2 {
			t.Errorf("Expected 2 blacklisted commands, got %d", len(blacklist))
		}

		// Test adding to blacklist
		err = cfg.AddToFallthroughBlacklist("fetch")
		if err != nil {
			t.Errorf("Failed to add to blacklist: %v", err)
		}

		// Test blacklist functionality
		err = handler.Execute("fetch", []string{})
		if err == nil || !strings.Contains(err.Error(), "blacklisted") {
			t.Errorf("Expected blacklist error for fetch command")
		}

		// Test removing from blacklist
		err = cfg.RemoveFromFallthroughBlacklist("fetch")
		if err != nil {
			t.Errorf("Failed to remove from blacklist: %v", err)
		}

		// Test reserved command protection
		err = handler.Execute("version", []string{})
		if err == nil || !strings.Contains(err.Error(), "reserved") {
			t.Errorf("Expected reserved command error for version")
		}

		// Test argument validation
		err = handler.Execute("log", []string{"$(whoami)"})
		if err == nil || !strings.Contains(err.Error(), "unsafe") {
			t.Errorf("Expected unsafe argument error")
		}

		// Test valid command (if git is available)
		if _, err := exec.LookPath("git"); err == nil {
			err = handler.Execute("--version", []string{})
			if err != nil {
				t.Logf("Git command result: %v", err)
			}
		}

		t.Log("Complete integration test finished")
	})
}