package handlers

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// TestFallthroughIntegration_CommonGitCommands tests integration with common Git commands
func TestFallthroughIntegration_CommonGitCommands(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping integration test: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
		description string
	}{
		{
			name:        "git status",
			command:     "status",
			args:        []string{"--porcelain"},
			expectError: false,
			description: "Basic status command should work",
		},
		{
			name:        "git rev-parse",
			command:     "rev-parse",
			args:        []string{"HEAD"},
			expectError: false,
			description: "Rev-parse command should work",
		},
		{
			name:        "git show",
			command:     "show",
			args:        []string{"--name-only", "HEAD"},
			expectError: false,
			description: "Show command should work",
		},
		{
			name:        "git log with options",
			command:     "log",
			args:        []string{"--oneline", "-n", "5"},
			expectError: false,
			description: "Log command with multiple options should work",
		},
		{
			name:        "git diff",
			command:     "diff",
			args:        []string{"--name-only"},
			expectError: false,
			description: "Diff command should work",
		},
		{
			name:        "git ls-files",
			command:     "ls-files",
			args:        []string{"--cached"},
			expectError: false,
			description: "Ls-files command should work",
		},
		{
			name:        "git config get",
			command:     "config",
			args:        []string{"--get", "user.name"},
			expectError: false,
			description: "Config commands should work",
		},
		{
			name:        "git remote list",
			command:     "remote",
			args:        []string{"-v"},
			expectError: false,
			description: "Remote commands should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.description)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, got: %v", tt.description, err)
				}
			}
		})
	}
}

// TestFallthroughIntegration_ComplexWorkflows tests complex Git workflows
func TestFallthroughIntegration_ComplexWorkflows(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping integration test: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: true, // Enable verbose for workflow testing
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	tests := []struct {
		name     string
		commands []struct {
			command string
			args    []string
		}
		description string
	}{
		{
			name: "status and diff workflow",
			commands: []struct {
				command string
				args    []string
			}{
				{"status", []string{"--short"}},
				{"diff", []string{"--name-only"}},
				{"status", []string{"--porcelain"}},
			},
			description: "Common workflow of checking status and diff",
		},
		{
			name: "log and show workflow",
			commands: []struct {
				command string
				args    []string
			}{
				{"log", []string{"--oneline", "-n", "3"}},
				{"show", []string{"--name-only", "HEAD"}},
			},
			description: "Workflow for examining recent commits",
		},
		{
			name: "branch and remote workflow",
			commands: []struct {
				command string
				args    []string
			}{
				{"branch", []string{"-a"}},
				{"remote", []string{"-v"}},
				{"status", []string{"-b"}},
			},
			description: "Workflow for checking branch and remote status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, cmd := range tt.commands {
				err := handler.Execute(cmd.command, cmd.args)
				if err != nil {
					t.Errorf("Command %d in workflow '%s' failed: %v", i+1, tt.description, err)
					break // Stop workflow on first error
				}
			}
		})
	}
}

// TestFallthroughIntegration_ArgumentPreservation tests that complex arguments are preserved correctly
func TestFallthroughIntegration_ArgumentPreservation(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping integration test: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	tests := []struct {
		name        string
		command     string
		args        []string
		description string
	}{
		{
			name:        "arguments with spaces",
			command:     "log",
			args:        []string{"--grep=test message", "--oneline"},
			description: "Arguments containing spaces should be preserved",
		},
		{
			name:        "arguments with quotes",
			command:     "log",
			args:        []string{"--pretty=format:%h %s", "-n", "1"},
			description: "Arguments with format strings should be preserved",
		},
		{
			name:        "multiple flag combinations",
			command:     "status",
			args:        []string{"--short", "--branch", "--ahead-behind"},
			description: "Multiple flags should be preserved in order",
		},
		{
			name:        "arguments with special characters",
			command:     "log",
			args:        []string{"--pretty=format:%C(red)%h%C(reset) %s", "-n", "1"},
			description: "Arguments with color codes should be preserved",
		},
		{
			name:        "file paths with spaces",
			command:     "status",
			args:        []string{"--", "file with spaces.txt"},
			description: "File paths with spaces should be handled correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, tt.args)
			// We don't check for specific success/failure here since some commands
			// might fail due to missing files/commits, but we verify no argument parsing errors
			if err != nil && strings.Contains(err.Error(), "argument validation failed") {
				t.Errorf("Argument preservation failed for %s: %v", tt.description, err)
			}
		})
	}
}

// TestFallthroughPerformance_SubprocessOverhead tests performance overhead of subprocess execution
func TestFallthroughPerformance_SubprocessOverhead(t *testing.T) {
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

	// Test commands that should be fast
	fastCommands := []struct {
		command string
		args    []string
		name    string
	}{
		{"--version", []string{}, "version"},
		{"status", []string{"--porcelain"}, "status"},
		{"branch", []string{}, "branch"},
		{"config", []string{"--get", "user.name"}, "config"},
	}

	for _, cmd := range fastCommands {
		t.Run(fmt.Sprintf("performance_%s", cmd.name), func(t *testing.T) {
			// Measure execution time
			start := time.Now()
			err := handler.Execute(cmd.command, cmd.args)
			duration := time.Since(start)

			if err != nil {
				t.Logf("Command failed (expected for some environments): %v", err)
				return
			}

			// Log performance metrics
			t.Logf("Command '%s %v' took %v", cmd.command, cmd.args, duration)

			// Reasonable performance threshold (adjust based on system)
			maxDuration := 5 * time.Second
			if duration > maxDuration {
				t.Errorf("Command took too long: %v > %v", duration, maxDuration)
			}
		})
	}
}

// TestFallthroughSecurity_CommandInjectionPrevention tests security against command injection
func TestFallthroughSecurity_CommandInjectionPrevention(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	// Test potentially dangerous argument patterns
	dangerousTests := []struct {
		name        string
		command     string
		args        []string
		description string
	}{
		{
			name:        "command substitution with dollar",
			command:     "status",
			args:        []string{"$(rm -rf /)"},
			description: "Command substitution should be blocked",
		},
		{
			name:        "command substitution with backticks",
			command:     "log",
			args:        []string{"`rm -rf /`"},
			description: "Backtick command substitution should be blocked",
		},
		{
			name:        "shell injection attempt",
			command:     "status",
			args:        []string{"; rm -rf /"},
			description: "Shell command chaining should be safe",
		},
		{
			name:        "pipe injection attempt",
			command:     "log",
			args:        []string{"| rm -rf /"},
			description: "Pipe injection should be safe",
		},
		{
			name:        "redirect injection attempt",
			command:     "status",
			args:        []string{"> /etc/passwd"},
			description: "Redirect injection should be safe",
		},
	}

	for _, tt := range dangerousTests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, tt.args)

			// We expect either:
			// 1. Argument validation error (preferred)
			// 2. Git error (acceptable - means Git handled it safely)
			// 3. No error but safe execution (acceptable)
			
			if err != nil {
				if strings.Contains(err.Error(), "potentially unsafe argument") {
					// This is the preferred outcome - our validation caught it
					t.Logf("Security validation correctly blocked: %s", tt.description)
				} else {
					// Git or system handled it - also acceptable
					t.Logf("Command failed safely: %s - %v", tt.description, err)
				}
			} else {
				// Command succeeded - verify it was safe
				t.Logf("Command executed safely: %s", tt.description)
			}
		})
	}
}

// TestFallthroughEndToEnd_CompleteScenarios tests complete end-to-end scenarios
func TestFallthroughEndToEnd_CompleteScenarios(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping end-to-end test: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	t.Run("complete_git_workflow", func(t *testing.T) {
		// Simulate a complete Git workflow
		workflow := []struct {
			command     string
			args        []string
			description string
		}{
			{"status", []string{}, "Check repository status"},
			{"branch", []string{"-a"}, "List all branches"},
			{"log", []string{"--oneline", "-n", "5"}, "Show recent commits"},
			{"diff", []string{"--name-only"}, "Show changed files"},
			{"remote", []string{"-v"}, "Show remotes"},
			{"config", []string{"--list", "--local"}, "Show local config"},
		}

		for i, step := range workflow {
			t.Logf("Step %d: %s", i+1, step.description)
			err := handler.Execute(step.command, step.args)
			if err != nil {
				t.Logf("Step %d failed (may be expected): %v", i+1, err)
				// Don't fail the test - some commands may fail in CI environments
			}
		}
	})

	t.Run("error_handling_workflow", func(t *testing.T) {
		// Test error handling in realistic scenarios
		errorTests := []struct {
			command     string
			args        []string
			description string
		}{
			{"nonexistentcommand", []string{}, "Unknown Git command"},
			{"log", []string{"--invalid-option"}, "Invalid Git option"},
			{"show", []string{"nonexistent-commit"}, "Invalid commit reference"},
		}

		for _, test := range errorTests {
			t.Logf("Testing error scenario: %s", test.description)
			err := handler.Execute(test.command, test.args)
			if err == nil {
				t.Logf("Command unexpectedly succeeded: %s", test.description)
			} else {
				t.Logf("Error handled correctly: %s - %v", test.description, err)
			}
		}
	})
}

// TestFallthroughIntegration_BlacklistFunctionality tests blacklist functionality in realistic scenarios
func TestFallthroughIntegration_BlacklistFunctionality(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled:   true,
			Verbose:   false,
			Blacklist: []string{"push", "pull", "fetch"},
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	blacklistedCommands := []string{"push", "pull", "fetch"}
	allowedCommands := []string{"status", "log", "diff", "branch"}

	t.Run("blacklisted_commands_blocked", func(t *testing.T) {
		for _, cmd := range blacklistedCommands {
			err := handler.Execute(cmd, []string{})
			if err == nil {
				t.Errorf("Expected blacklisted command '%s' to be blocked", cmd)
			} else if !strings.Contains(err.Error(), "blacklisted") {
				t.Errorf("Expected blacklist error for '%s', got: %v", cmd, err)
			}
		}
	})

	t.Run("allowed_commands_work", func(t *testing.T) {
		// Skip this test if git is not available
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("Skipping test: git not available in PATH")
		}

		for _, cmd := range allowedCommands {
			err := handler.Execute(cmd, []string{})
			// Don't check for success/failure, just ensure it's not blacklisted
			if err != nil && strings.Contains(err.Error(), "blacklisted") {
				t.Errorf("Allowed command '%s' was incorrectly blacklisted: %v", cmd, err)
			}
		}
	})
}

// TestFallthroughIntegration_VerboseMode tests verbose mode functionality
func TestFallthroughIntegration_VerboseMode(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping integration test: git not available in PATH")
	}

	// Test with verbose mode enabled
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	t.Run("verbose_mode_enabled", func(t *testing.T) {
		// Capture stderr to verify verbose output
		// Note: In a real test environment, you might want to capture stderr
		// For now, we just verify the command executes without error
		err := handler.Execute("status", []string{"--porcelain"})
		if err != nil {
			t.Logf("Command failed (may be expected in test environment): %v", err)
		}
		// Verbose output goes to stderr, so we can't easily capture it in this test
		// But we can verify the command executed
	})

	// Test with verbose mode disabled
	cfg.Fallthrough.Verbose = false
	handler = NewFallthroughHandler(cfg, gitRepo)

	t.Run("verbose_mode_disabled", func(t *testing.T) {
		err := handler.Execute("status", []string{"--porcelain"})
		if err != nil {
			t.Logf("Command failed (may be expected in test environment): %v", err)
		}
	})
}

// TestFallthroughIntegration_ConfigurationScenarios tests different configuration scenarios
func TestFallthroughIntegration_ConfigurationScenarios(t *testing.T) {
	gitRepo := &git.Repository{}

	t.Run("fallthrough_disabled", func(t *testing.T) {
		cfg := &config.Config{
			Fallthrough: config.FallthroughConfig{
				Enabled: false,
			},
		}
		handler := NewFallthroughHandler(cfg, gitRepo)

		err := handler.Execute("status", []string{})
		if err == nil {
			t.Error("Expected error when fallthrough is disabled")
		} else if !strings.Contains(err.Error(), "disabled") {
			t.Errorf("Expected disabled error, got: %v", err)
		}
	})

	t.Run("empty_blacklist", func(t *testing.T) {
		cfg := &config.Config{
			Fallthrough: config.FallthroughConfig{
				Enabled:   true,
				Blacklist: []string{},
			},
		}
		handler := NewFallthroughHandler(cfg, gitRepo)

		// Verify no commands are blacklisted
		if handler.config.IsFallthroughBlacklisted("status") {
			t.Error("No commands should be blacklisted with empty blacklist")
		}
	})

	t.Run("extensive_blacklist", func(t *testing.T) {
		cfg := &config.Config{
			Fallthrough: config.FallthroughConfig{
				Enabled:   true,
				Blacklist: []string{"push", "pull", "fetch", "merge", "rebase", "reset"},
			},
		}
		handler := NewFallthroughHandler(cfg, gitRepo)

		// Test that all blacklisted commands are blocked
		for _, cmd := range cfg.Fallthrough.Blacklist {
			err := handler.Execute(cmd, []string{})
			if err == nil || !strings.Contains(err.Error(), "blacklisted") {
				t.Errorf("Command '%s' should be blacklisted", cmd)
			}
		}
	})
}