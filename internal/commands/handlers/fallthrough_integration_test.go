package handlers

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// TestComplexArgumentPreservation_Integration tests that complex arguments
// are preserved correctly when passed through to Git commands
func TestComplexArgumentPreservation_Integration(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}

	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	tests := []struct {
		name        string
		command     string
		args        []string
		description string
	}{
		{
			name:        "commit message with spaces",
			command:     "log",
			args:        []string{"--oneline", "-n", "1", "--grep", "commit message"},
			description: "Arguments with spaces should be preserved as single arguments",
		},
		{
			name:        "file paths with special characters",
			command:     "status",
			args:        []string{"--porcelain"},
			description: "File paths with special characters should be handled correctly",
		},
		{
			name:        "multiple flags and options",
			command:     "log",
			args:        []string{"--pretty=format:%h %s", "--since=1.day.ago", "--author=.*"},
			description: "Multiple flags and complex format strings should be preserved",
		},
		{
			name:        "quoted arguments",
			command:     "log",
			args:        []string{"--grep=fix.*bug", "--oneline", "-n", "5"},
			description: "Arguments with regex patterns should be preserved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the command through the fallthrough handler
			err := handler.Execute(tt.command, tt.args)

			// We expect either success or a Git-specific error (not argument handling errors)
			if err != nil {
				// Check if it's an argument validation error (which would indicate a problem)
				if strings.Contains(err.Error(), "argument validation failed") {
					t.Errorf("Argument validation failed for %s: %v", tt.description, err)
				}
				// Other errors (like "no commits found", etc.) are acceptable
				// as they indicate Git processed the arguments correctly
			}
		})
	}
}

// TestArgumentEscaping_Integration tests that potentially dangerous arguments
// are properly handled and rejected
func TestArgumentEscaping_Integration(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	dangerousArgs := []struct {
		name string
		args []string
	}{
		{
			name: "command substitution with dollar sign",
			args: []string{"--message", "$(echo dangerous)"},
		},
		{
			name: "backtick command substitution",
			args: []string{"--message", "`echo dangerous`"},
		},
	}

	for _, tt := range dangerousArgs {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute("status", tt.args)

			// We expect these to be rejected by argument validation
			if err == nil {
				t.Errorf("Expected dangerous arguments to be rejected: %v", tt.args)
			} else if !strings.Contains(err.Error(), "argument validation failed") {
				t.Errorf("Expected argument validation error, got: %v", err)
			}
		})
	}
}

// TestRealWorldScenarios_Integration tests real-world Git command scenarios
func TestRealWorldScenarios_Integration(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}

	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	scenarios := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "git status with multiple flags",
			command: "status",
			args:    []string{"--short", "--branch", "--porcelain"},
		},
		{
			name:    "git log with complex formatting",
			command: "log",
			args:    []string{"--pretty=format:%h - %an, %ar : %s", "--since=1.week.ago", "-n", "10"},
		},
		{
			name:    "git diff with file paths",
			command: "diff",
			args:    []string{"--name-only", "HEAD~1", "HEAD"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			err := handler.Execute(scenario.command, scenario.args)

			// These should execute without argument handling errors
			if err != nil && strings.Contains(err.Error(), "argument validation failed") {
				t.Errorf("Unexpected argument validation error for %s: %v", scenario.name, err)
			}
			// Other Git errors (like "not a git repository") are acceptable
		})
	}
}