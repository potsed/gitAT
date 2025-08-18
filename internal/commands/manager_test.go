package commands

import (
	"os"
	"os/exec"
	"testing"

	"github.com/potsed/gitAT/internal/config"
)

func TestManager_Execute_WithFallthrough(t *testing.T) {
	// Skip test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git executable not found, skipping fallthrough tests")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gitat-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize a git repository in the temp directory
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user for the test repository
	configCmds := [][]string{
		{"config", "user.name", "Test User"},
		{"config", "user.email", "test@example.com"},
	}
	for _, args := range configCmds {
		cmd := exec.Command("git", args...)
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to configure git: %v", err)
		}
	}

	// Create test config
	cfg := &config.Config{
		RepoPath: tempDir,
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}

	// Create manager
	manager := NewManager(cfg)

	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
		description string
	}{
		{
			name:        "fallthrough_git_status",
			command:     "status",
			args:        []string{},
			expectError: false,
			description: "Should execute git status through fallthrough",
		},
		{
			name:        "fallthrough_git_log",
			command:     "log",
			args:        []string{"--oneline"},
			expectError: false,
			description: "Should execute git log --oneline through fallthrough",
		},
		{
			name:        "fallthrough_git_diff",
			command:     "diff",
			args:        []string{"--cached"},
			expectError: false,
			description: "Should execute git diff --cached through fallthrough",
		},
		{
			name:        "fallthrough_invalid_git_command",
			command:     "invalidcommand",
			args:        []string{},
			expectError: true,
			description: "Should fail for invalid git commands",
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.Execute(tt.command, tt.args)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error for command '%s' but got none", tt.command)
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for command '%s': %v", tt.command, err)
			}
		})
	}
}

func TestManager_Execute_FallthroughIntegration(t *testing.T) {
	// Skip test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git executable not found, skipping fallthrough integration tests")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gitat-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize a git repository in the temp directory
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user for the test repository
	configCmds := [][]string{
		{"config", "user.name", "Test User"},
		{"config", "user.email", "test@example.com"},
	}
	for _, args := range configCmds {
		cmd := exec.Command("git", args...)
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to configure git: %v", err)
		}
	}

	// Create a test file and commit it
	testFile := tempDir + "/test.txt"
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit the file
	addCmd := exec.Command("git", "add", "test.txt")
	addCmd.Dir = tempDir
	if err := addCmd.Run(); err != nil {
		t.Fatalf("Failed to add test file: %v", err)
	}

	commitCmd := exec.Command("git", "commit", "-m", "Initial commit")
	commitCmd.Dir = tempDir
	if err := commitCmd.Run(); err != nil {
		t.Fatalf("Failed to commit test file: %v", err)
	}

	// Create test config
	cfg := &config.Config{
		RepoPath: tempDir,
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}

	// Create manager
	manager := NewManager(cfg)

	// Test complex Git commands through fallthrough
	complexTests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
		description string
	}{
		{
			name:        "git_log_with_format",
			command:     "log",
			args:        []string{"--pretty=format:%h %s", "-1"},
			expectError: false,
			description: "Should handle git log with custom format",
		},
		{
			name:        "git_show_with_options",
			command:     "show",
			args:        []string{"--name-only", "HEAD"},
			expectError: false,
			description: "Should handle git show with options",
		},
		{
			name:        "git_config_list",
			command:     "config",
			args:        []string{"--list"},
			expectError: false,
			description: "Should handle git config listing through fallthrough",
		},
	}

	for _, tt := range complexTests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to the test directory for git commands
			originalDir, _ := os.Getwd()
			os.Chdir(tempDir)
			defer os.Chdir(originalDir)

			err := manager.Execute(tt.command, tt.args)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error for command '%s %v' but got none", tt.command, tt.args)
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for command '%s %v': %v", tt.command, tt.args, err)
			}
		})
	}
}

func TestManager_NewManager_FallthroughInitialization(t *testing.T) {
	cfg := &config.Config{
		RepoPath: "/tmp/test",
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}

	manager := NewManager(cfg)

	// Verify that the fallthrough handler is properly initialized
	if manager.fallthroughHandler == nil {
		t.Error("Fallthrough handler should be initialized in NewManager")
	}

	// Verify that all other handlers are still initialized
	if manager.work == nil {
		t.Error("Work handler should be initialized")
	}
	if manager.version == nil {
		t.Error("Version handler should be initialized")
	}
	if manager.save == nil {
		t.Error("Save handler should be initialized")
	}

	// Verify config and git repository are set
	if manager.config != cfg {
		t.Error("Config should be set correctly")
	}
	if manager.git == nil {
		t.Error("Git repository should be initialized")
	}
}

func TestManager_Execute_ReservedCommands(t *testing.T) {
	cfg := &config.Config{
		RepoPath: "/tmp/test",
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}

	manager := NewManager(cfg)

	// Test a few key reserved Gitat commands to ensure they don't fall through
	// Using commands that are less likely to have interactive forms or complex dependencies
	reservedCommands := []string{
		"info", "sweep", "save",
	}

	for _, cmd := range reservedCommands {
		t.Run("reserved_command_"+cmd, func(t *testing.T) {
			// These should execute the Gitat handlers, not fall through
			// We expect them to execute (may error due to test environment, but shouldn't fall through)
			err := manager.Execute(cmd, []string{})
			
			// The key test is that it doesn't return "unknown command" error
			// which would indicate fallthrough failure
			if err != nil && err.Error() == "unknown command: "+cmd {
				t.Errorf("Command '%s' should not fall through to unknown command error", cmd)
			}
		})
	}
}

func TestManager_Execute_FallthroughHandlerIntegration(t *testing.T) {
	// Skip test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git executable not found, skipping fallthrough handler integration test")
	}

	cfg := &config.Config{
		RepoPath: "/tmp/test",
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}

	manager := NewManager(cfg)

	// Test that unknown commands fall through to the fallthrough handler
	// Use a command that's not reserved by Gitat
	err := manager.Execute("remote", []string{})
	
	// The command should execute through fallthrough, not return "unknown command"
	if err != nil && err.Error() == "unknown command: remote" {
		t.Error("Command 'remote' should fall through to git, not return unknown command error")
	}
	
	// Test with a clearly non-existent command
	err = manager.Execute("nonexistentcommand", []string{})
	
	// This should still go through fallthrough (and git will handle the error)
	if err != nil && err.Error() == "unknown command: nonexistentcommand" {
		t.Error("Even non-existent commands should fall through to git, not return unknown command error")
	}
}

func TestManager_Execute_FallthroughDisabled(t *testing.T) {
	// Test behavior when fallthrough is disabled (if we add that config option)
	cfg := &config.Config{
		RepoPath: "/tmp/test",
		Fallthrough: config.FallthroughConfig{
			Enabled: false, // Disabled
			Verbose: false,
		},
	}

	manager := NewManager(cfg)

	// Even with fallthrough disabled, the handler should still be initialized
	// The actual disable logic would be in the fallthrough handler itself
	if manager.fallthroughHandler == nil {
		t.Error("Fallthrough handler should be initialized even when disabled")
	}
}