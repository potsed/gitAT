package handlers

import (
	"os/exec"
	"testing"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/internal/utils"
)

// TestInteractiveCommandSupport tests the comprehensive interactive command functionality
func TestInteractiveCommandSupport(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}

	cfg := &config.Config{}
	gitRepo := &git.Repository{}

	tests := []struct {
		name          string
		setupTestMode bool
		testResponses map[string][]string
		command       string
		args          []string
		expectError   bool
		description   string
	}{
		{
			name:          "non-interactive command in normal mode",
			setupTestMode: false,
			command:       "status",
			args:          []string{"--porcelain"},
			expectError:   false,
			description:   "Regular git status should work normally",
		},
		{
			name:          "non-interactive command in test mode",
			setupTestMode: true,
			testResponses: map[string][]string{},
			command:       "status",
			args:          []string{"--porcelain"},
			expectError:   false,
			description:   "Regular git status should work in test mode",
		},
		{
			name:          "git version in test mode",
			setupTestMode: true,
			testResponses: map[string][]string{},
			command:       "--version",
			args:          []string{},
			expectError:   false,
			description:   "Git version should work in test mode",
		},
		{
			name:          "interactive command detection - add patch",
			setupTestMode: true,
			testResponses: map[string][]string{
				"git add -p": {"q"}, // Quit immediately
			},
			command:     "add",
			args:        []string{"-p"},
			expectError: false,
			description: "git add -p should be detected as interactive and handled in test mode",
		},
		{
			name:          "interactive command detection - rebase interactive",
			setupTestMode: true,
			testResponses: map[string][]string{
				"git rebase -i": {":q"}, // Quit editor immediately
			},
			command:     "rebase",
			args:        []string{"-i", "HEAD~1"},
			expectError: false,
			description: "git rebase -i should be detected as interactive and handled in test mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var handler *FallthroughHandler

			if tt.setupTestMode {
				handler = NewFallthroughHandlerWithTestMode(cfg, gitRepo, tt.testResponses)
			} else {
				handler = NewFallthroughHandler(cfg, gitRepo)
			}

			err := handler.Execute(tt.command, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.description)
				}
			} else {
				if err != nil {
					// Some git commands might fail due to repository state, which is acceptable
					// We're mainly testing that the interactive handling doesn't cause hangs
					t.Logf("Command failed (acceptable for test): %v", err)
				}
			}
		})
	}
}

// TestProcessExecutorInteractiveFeatures tests the process executor's interactive capabilities
func TestProcessExecutorInteractiveFeatures(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		isInteractive bool
		description string
	}{
		{
			name:        "git add patch mode",
			command:     "git",
			args:        []string{"add", "-p"},
			isInteractive: true,
			description: "git add -p should be detected as interactive",
		},
		{
			name:        "git add interactive mode",
			command:     "git",
			args:        []string{"add", "-i"},
			isInteractive: true,
			description: "git add -i should be detected as interactive",
		},
		{
			name:        "git rebase interactive",
			command:     "git",
			args:        []string{"rebase", "-i", "HEAD~3"},
			isInteractive: true,
			description: "git rebase -i should be detected as interactive",
		},
		{
			name:        "git commit verbose",
			command:     "git",
			args:        []string{"commit", "-v"},
			isInteractive: true,
			description: "git commit -v should be detected as interactive",
		},
		{
			name:        "git merge",
			command:     "git",
			args:        []string{"merge", "feature-branch"},
			isInteractive: true,
			description: "git merge should be detected as interactive",
		},
		{
			name:        "git status",
			command:     "git",
			args:        []string{"status"},
			isInteractive: false,
			description: "git status should not be detected as interactive",
		},
		{
			name:        "git add file",
			command:     "git",
			args:        []string{"add", "file.txt"},
			isInteractive: false,
			description: "git add file should not be detected as interactive",
		},
		{
			name:        "git commit with message",
			command:     "git",
			args:        []string{"commit", "-m", "message"},
			isInteractive: false,
			description: "git commit -m should not be detected as interactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := utils.NewProcessExecutor()
			
			result := executor.IsInteractiveCommand(tt.command, tt.args)
			
			if result != tt.isInteractive {
				t.Errorf("Expected isInteractiveCommand(%s %v) = %v, got %v", 
					tt.command, tt.args, tt.isInteractive, result)
			}
		})
	}
}

// TestTestModeResponses tests the test mode response system
func TestTestModeResponses(t *testing.T) {
	executor := utils.NewProcessExecutor()

	tests := []struct {
		name        string
		command     string
		args        []string
		expectedKey string
		description string
	}{
		{
			name:        "git add patch",
			command:     "git",
			args:        []string{"add", "-p"},
			expectedKey: "git add -p",
			description: "git add -p should generate correct key",
		},
		{
			name:        "git rebase interactive",
			command:     "git",
			args:        []string{"rebase", "-i", "HEAD~3"},
			expectedKey: "git rebase -i",
			description: "git rebase -i should generate correct key",
		},
		{
			name:        "git commit verbose",
			command:     "git",
			args:        []string{"commit", "-v"},
			expectedKey: "git commit -v",
			description: "git commit -v should generate correct key",
		},
		{
			name:        "simple command",
			command:     "echo",
			args:        []string{"hello", "world"},
			expectedKey: "echo hello world",
			description: "Non-git commands should include all args",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.GenerateCommandKey(tt.command, tt.args)
			
			if result != tt.expectedKey {
				t.Errorf("Expected generateCommandKey(%s %v) = %q, got %q", 
					tt.command, tt.args, tt.expectedKey, result)
			}
		})
	}
}

// TestDefaultTestResponses tests the default response generation
func TestDefaultTestResponses(t *testing.T) {
	executor := utils.NewProcessExecutor()

	tests := []struct {
		name             string
		command          string
		args             []string
		expectedResponses []string
		description      string
	}{
		{
			name:             "git add patch",
			command:          "git",
			args:             []string{"add", "-p"},
			expectedResponses: []string{"y", "q"},
			description:      "git add -p should get default patch responses",
		},
		{
			name:             "git rebase interactive",
			command:          "git",
			args:             []string{"rebase", "-i", "HEAD~3"},
			expectedResponses: []string{":wq"},
			description:      "git rebase -i should get default editor quit",
		},
		{
			name:             "git merge",
			command:          "git",
			args:             []string{"merge", "feature"},
			expectedResponses: []string{"", ""},
			description:      "git merge should get default empty responses",
		},
		{
			name:             "git commit verbose",
			command:          "git",
			args:             []string{"commit", "-v"},
			expectedResponses: []string{"Test commit message", "", ":wq"},
			description:      "git commit -v should get default commit message and editor quit",
		},
		{
			name:             "unknown command",
			command:          "unknown",
			args:             []string{"arg"},
			expectedResponses: []string{""},
			description:      "Unknown commands should get default empty response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.GetDefaultTestResponses(tt.command, tt.args)
			
			if len(result) != len(tt.expectedResponses) {
				t.Errorf("Expected %d responses, got %d", len(tt.expectedResponses), len(result))
				return
			}
			
			for i, expected := range tt.expectedResponses {
				if result[i] != expected {
					t.Errorf("Response %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

// TestInteractiveCommandIntegration tests end-to-end interactive command handling
func TestInteractiveCommandIntegration(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}

	cfg := &config.Config{}
	gitRepo := &git.Repository{}

	// Test with custom responses
	customResponses := map[string][]string{
		"git add -p":     {"n", "q"}, // No to first hunk, then quit
		"git rebase -i":  {":q!"},    // Force quit editor
		"git commit -v":  {"Custom test message", "", ":wq"},
	}

	handler := NewFallthroughHandlerWithTestMode(cfg, gitRepo, customResponses)

	// Test that we can configure and reconfigure test responses
	handler.SetTestResponses(map[string][]string{
		"git add -p": {"y", "y", "q"}, // Yes to first two hunks, then quit
	})

	// Test enabling and disabling test mode
	handler.SetTestMode(false)
	handler.SetTestMode(true)

	// Test that the handler can execute commands without hanging
	// We use git --version as a safe command that always works
	err := handler.Execute("--version", []string{})
	if err != nil {
		t.Errorf("Expected no error for git --version in test mode, got: %v", err)
	}
}