package handlers

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/internal/utils"
)

func TestCreateGitNotFoundError(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	err := handler.createGitNotFoundError()
	
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	errorMsg := err.Error()
	
	// Check that the error message contains helpful information
	expectedPhrases := []string{
		"Git executable not found in PATH",
		"Installation options:",
		"brew install git",
		"sudo apt-get install git",
		"git --version",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}

func TestCreateUnknownCommandError(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: false, // Fallthrough disabled
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name            string
		command         string
		args            []string
		expectedPhrases []string
	}{
		{
			name:    "command with no suggestions",
			command: "xyz123",
			args:    []string{},
			expectedPhrases: []string{
				"Unknown command: xyz123",
				"Fallthrough is currently disabled",
				"git config at.fallthrough.enabled true",
				"Available Gitat commands:",
				"git @ help",
			},
		},
		{
			name:    "command that might match feature",
			command: "feat",
			args:    []string{},
			expectedPhrases: []string{
				"Unknown command: feat",
				"Did you mean one of these Gitat commands?",
				"git @ feature",
			},
		},
		{
			name:    "command that might match branch",
			command: "br",
			args:    []string{},
			expectedPhrases: []string{
				"Unknown command: br",
				"Did you mean one of these Gitat commands?",
				"git @ branch",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.createUnknownCommandError(tt.command, tt.args)
			
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			
			errorMsg := err.Error()
			
			for _, phrase := range tt.expectedPhrases {
				if !strings.Contains(errorMsg, phrase) {
					t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
				}
			}
		})
	}
}

func TestCreateReservedCommandError(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "version command",
			command: "version",
			args:    []string{},
		},
		{
			name:    "help command",
			command: "help",
			args:    []string{},
		},
		{
			name:    "branch command with args",
			command: "branch",
			args:    []string{"--list"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.createReservedCommandError(tt.command, tt.args)
			
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			
			errorMsg := err.Error()
			
			expectedPhrases := []string{
				fmt.Sprintf("'%s' is a reserved Gitat command", tt.command),
				fmt.Sprintf("Use 'git @ %s'", tt.command),
				fmt.Sprintf("or 'git %s'", tt.command),
				"git @ help",
			}
			
			for _, phrase := range expectedPhrases {
				if !strings.Contains(errorMsg, phrase) {
					t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
				}
			}
		})
	}
}

func TestCreateBlacklistedCommandError(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Blacklist: []string{"status", "diff", "log"},
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "blacklisted status command",
			command: "status",
			args:    []string{},
		},
		{
			name:    "blacklisted diff command with args",
			command: "diff",
			args:    []string{"--cached"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.createBlacklistedCommandError(tt.command, tt.args)
			
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			
			errorMsg := err.Error()
			
			expectedPhrases := []string{
				fmt.Sprintf("Command '%s' is blacklisted from fallthrough", tt.command),
				"explicitly disabled for security or compatibility reasons",
				fmt.Sprintf("git %s", tt.command),
				"git config --unset at.fallthrough.blacklist",
				"Current blacklist:",
			}
			
			for _, phrase := range expectedPhrases {
				if !strings.Contains(errorMsg, phrase) {
					t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
				}
			}
		})
	}
}

func TestCreateArgumentValidationError(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	validationErr := fmt.Errorf("potentially unsafe argument at position 1: $(rm -rf /)")
	command := "commit"
	args := []string{"-m", "$(rm -rf /)"}
	
	err := handler.createArgumentValidationError(command, args, validationErr)
	
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	errorMsg := err.Error()
	
	expectedPhrases := []string{
		"Argument validation failed for command 'commit'",
		"potentially unsafe argument at position 1",
		"validating arguments for security",
		"git commit -m $(rm -rf /)",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}

func TestCreateCommandValidationError(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	validationErr := fmt.Errorf("command not found")
	command := "nonexistent"
	args := []string{"--flag"}
	
	err := handler.createCommandValidationError(command, args, validationErr)
	
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	errorMsg := err.Error()
	
	expectedPhrases := []string{
		"Command validation failed for 'nonexistent'",
		"command not found",
		"The command doesn't exist",
		"git nonexistent --flag",
		"git help nonexistent",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}

func TestEnhanceGitError(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name            string
		command         string
		args            []string
		inputError      error
		expectedPhrases []string
		expectedExitCode int
	}{
		{
			name:    "exit code 1 - general git error",
			command: "status",
			args:    []string{},
			inputError: &utils.ExitError{
				ExitCode: 1,
				Err:      fmt.Errorf("fatal: not a git repository"),
			},
			expectedPhrases: []string{
				"Git command failed: git status",
				"Original error:",
				"fatal: not a git repository",
				"git help status",
			},
			expectedExitCode: 1,
		},
		{
			name:    "exit code 127 - command not found",
			command: "nonexistent",
			args:    []string{},
			inputError: &utils.ExitError{
				ExitCode: 127,
				Err:      fmt.Errorf("git: 'nonexistent' is not a git command"),
			},
			expectedPhrases: []string{
				"Git command not found: 'nonexistent'",
				"doesn't exist in your Git version",
				"typo in the command name",
				"git help -a",
			},
			expectedExitCode: 0, // This error type doesn't preserve exit code
		},
		{
			name:    "exit code 128 - repository error",
			command: "diff",
			args:    []string{"--cached"},
			inputError: &utils.ExitError{
				ExitCode: 128,
				Err:      fmt.Errorf("fatal: not a git repository"),
			},
			expectedPhrases: []string{
				"Git repository error:",
				"not in a Git repository",
				"git init",
			},
			expectedExitCode: 128,
		},
		{
			name:    "other exit code",
			command: "merge",
			args:    []string{"feature-branch"},
			inputError: &utils.ExitError{
				ExitCode: 2,
				Err:      fmt.Errorf("merge conflict"),
			},
			expectedPhrases: []string{
				"Git command failed with exit code 2",
				"git merge feature-branch",
				"merge conflict",
			},
			expectedExitCode: 2,
		},
		{
			name:    "command not found error",
			command: "status",
			args:    []string{},
			inputError: fmt.Errorf("executable file not found in PATH: git"),
			expectedPhrases: []string{
				"Git executable not found in PATH",
				"Installation options:",
				"brew install git",
			},
			expectedExitCode: 0,
		},
		{
			name:    "generic error",
			command: "log",
			args:    []string{"--oneline"},
			inputError: fmt.Errorf("some generic error"),
			expectedPhrases: []string{
				"Git command execution failed: git log --oneline",
				"some generic error",
				"git help log",
			},
			expectedExitCode: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.enhanceGitError(tt.command, tt.args, tt.inputError)
			
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			
			errorMsg := err.Error()
			
			for _, phrase := range tt.expectedPhrases {
				if !strings.Contains(errorMsg, phrase) {
					t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
				}
			}
			
			// Check exit code preservation for ExitError types
			if tt.expectedExitCode > 0 {
				if exitErr, ok := err.(*utils.ExitError); ok {
					if exitErr.ExitCode != tt.expectedExitCode {
						t.Errorf("Expected exit code %d, got %d", tt.expectedExitCode, exitErr.ExitCode)
					}
				} else {
					t.Errorf("Expected ExitError with exit code %d, got different error type", tt.expectedExitCode)
				}
			}
		})
	}
}

func TestGetAvailableGitAtCommands(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	commands := handler.GetAvailableGitAtCommands()
	
	if len(commands) == 0 {
		t.Error("Expected non-empty list of available commands")
	}
	
	// Check that some expected commands are present
	expectedCommands := []string{"branch", "feature", "info", "help", "version"}
	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range commands {
			if cmd == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command %q to be in available commands list", expected)
		}
	}
}

func TestSuggestSimilarCommands(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name            string
		input           string
		expectedContains []string
		expectedEmpty   bool
	}{
		{
			name:            "partial match - feat",
			input:           "feat",
			expectedContains: []string{"feature"},
		},
		{
			name:            "partial match - br",
			input:           "br",
			expectedContains: []string{"branch"},
		},
		{
			name:            "partial match - inf",
			input:           "inf",
			expectedContains: []string{"info"},
		},
		{
			name:            "common abbreviation - st",
			input:           "st",
			expectedContains: []string{"status"},
		},
		{
			name:            "common abbreviation - co",
			input:           "co",
			expectedContains: []string{"checkout"},
		},
		{
			name:          "no matches",
			input:         "xyz123",
			expectedEmpty: true,
		},
		{
			name:            "exact match",
			input:           "branch",
			expectedContains: []string{"branch"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := handler.SuggestSimilarCommands(tt.input)
			
			if tt.expectedEmpty {
				if len(suggestions) != 0 {
					t.Errorf("Expected empty suggestions for %q, got %v", tt.input, suggestions)
				}
				return
			}
			
			if len(suggestions) == 0 {
				t.Errorf("Expected non-empty suggestions for %q", tt.input)
				return
			}
			
			for _, expected := range tt.expectedContains {
				found := false
				for _, suggestion := range suggestions {
					if suggestion == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected suggestion %q for input %q, got %v", expected, tt.input, suggestions)
				}
			}
		})
	}
}

func TestExecute_FallthroughDisabled(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	err := handler.Execute("status", []string{})
	
	if err == nil {
		t.Fatal("Expected error when fallthrough is disabled, got nil")
	}
	
	errorMsg := err.Error()
	
	expectedPhrases := []string{
		"Unknown command: status",
		"Fallthrough is currently disabled",
		"git config at.fallthrough.enabled true",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}

func TestExecute_GitNotAvailable_EnhancedError(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Temporarily modify PATH to simulate git not being available
	originalPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", originalPath)
	
	err := handler.Execute("status", []string{})
	
	if err == nil {
		t.Fatal("Expected error when git is not available, got nil")
	}
	
	errorMsg := err.Error()
	
	expectedPhrases := []string{
		"Git executable not found in PATH",
		"Installation options:",
		"brew install git",
		"sudo apt-get install git",
		"git --version",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}

func TestExecute_ReservedCommand_EnhancedError(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	err := handler.Execute("version", []string{})
	
	if err == nil {
		t.Fatal("Expected error for reserved command, got nil")
	}
	
	errorMsg := err.Error()
	
	expectedPhrases := []string{
		"'version' is a reserved Gitat command",
		"Use 'git @ version'",
		"or 'git version'",
		"git @ help",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}

func TestExecute_BlacklistedCommand_EnhancedError(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled:   true,
			Blacklist: []string{"status", "diff"},
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	err := handler.Execute("status", []string{})
	
	if err == nil {
		t.Fatal("Expected error for blacklisted command, got nil")
	}
	
	errorMsg := err.Error()
	
	expectedPhrases := []string{
		"Command 'status' is blacklisted from fallthrough",
		"explicitly disabled for security or compatibility reasons",
		"git status",
		"git config --unset at.fallthrough.blacklist",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}

func TestExecute_ArgumentValidation_EnhancedError(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Use arguments that will trigger validation error
	err := handler.Execute("commit", []string{"-m", "$(rm -rf /)"})
	
	if err == nil {
		t.Fatal("Expected error for dangerous arguments, got nil")
	}
	
	errorMsg := err.Error()
	
	expectedPhrases := []string{
		"Argument validation failed for command 'commit'",
		"potentially unsafe argument",
		"validating arguments for security",
		"git commit",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}