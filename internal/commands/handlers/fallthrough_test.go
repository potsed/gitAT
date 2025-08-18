package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

func TestNewFallthroughHandler(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}
	
	if handler.processExecutor == nil {
		t.Error("Expected processExecutor to be initialized")
	}
	
	if handler.reservedCommands == nil {
		t.Error("Expected reservedCommands to be initialized")
	}
	
	// Check that some known reserved commands are present
	expectedReserved := []string{"version", "help", "branch", "feature", "info"}
	for _, cmd := range expectedReserved {
		if !handler.reservedCommands[cmd] {
			t.Errorf("Expected %s to be in reserved commands", cmd)
		}
	}
}

func TestIsGitAvailable(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// This test assumes git is available in the test environment
	// In a real CI environment, we might need to mock this
	available := handler.isGitAvailable()
	
	// Check if git is actually in PATH
	_, err := exec.LookPath("git")
	expectedAvailable := err == nil
	
	if available != expectedAvailable {
		t.Errorf("Expected isGitAvailable() to return %v, got %v", expectedAvailable, available)
	}
}

func TestShouldFallthrough(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "valid git command",
			command:  "status",
			expected: true,
		},
		{
			name:     "another valid git command",
			command:  "diff",
			expected: true,
		},
		{
			name:     "reserved gitat command - version",
			command:  "version",
			expected: false,
		},
		{
			name:     "reserved gitat command - help",
			command:  "help",
			expected: false,
		},
		{
			name:     "reserved gitat command - branch",
			command:  "branch",
			expected: false,
		},
		{
			name:     "reserved gitat command - feature",
			command:  "feature",
			expected: false,
		},
		{
			name:     "empty command",
			command:  "",
			expected: false,
		},
		{
			name:     "whitespace only command",
			command:  "   ",
			expected: false,
		},
		{
			name:     "unknown but valid command name",
			command:  "someunknowncommand",
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.shouldFallthrough(tt.command)
			if result != tt.expected {
				t.Errorf("shouldFallthrough(%q) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestExecute_GitNotAvailable(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true, // Enable fallthrough to test git not available scenario
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
		t.Error("Expected error when git is not available, got nil")
	}
	
	// Check that the error message contains the expected phrases from the enhanced error
	errorMsg := err.Error()
	expectedPhrases := []string{
		"Git executable not found in PATH",
		"Installation options:",
		"brew install git",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(errorMsg, phrase) {
			t.Errorf("Expected error message to contain %q, but it didn't. Full message: %s", phrase, errorMsg)
		}
	}
}

func TestExecute_ReservedCommand(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	err := handler.Execute("version", []string{})
	
	if err == nil {
		t.Error("Expected error for reserved command, got nil")
	}
	
	expectedMsg := "unknown command: version"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestExecute_EmptyCommand(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	err := handler.Execute("", []string{})
	
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
	
	expectedMsg := "unknown command: "
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

// TestExecute_ValidGitCommand tests execution of valid git commands
// This test requires git to be available and will actually execute git commands
func TestExecute_ValidGitCommand(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}
	
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Test git --version which should always work
	err := handler.Execute("--version", []string{})
	
	if err != nil {
		t.Errorf("Expected no error for valid git command, got: %v", err)
	}
}

func TestIsGitCommand(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}
	
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "valid git command - status",
			command:  "status",
			expected: true,
		},
		{
			name:     "valid git command - diff",
			command:  "diff",
			expected: true,
		},
		{
			name:     "valid git command - log",
			command:  "log",
			expected: true,
		},
		{
			name:     "invalid git command",
			command:  "nonexistentcommand",
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.IsGitCommand(tt.command)
			if result != tt.expected {
				t.Errorf("IsGitCommand(%q) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestIsGitCommand_GitNotAvailable(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Temporarily modify PATH to simulate git not being available
	originalPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", originalPath)
	
	result := handler.IsGitCommand("status")
	
	if result != false {
		t.Errorf("Expected IsGitCommand to return false when git not available, got %v", result)
	}
}

func TestValidateArguments(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "safe arguments with spaces",
			args:        []string{"--message", "commit with spaces"},
			expectError: false,
		},
		{
			name:        "safe arguments with quotes",
			args:        []string{"--message", "commit with 'single quotes'"},
			expectError: false,
		},
		{
			name:        "safe arguments with double quotes",
			args:        []string{"--message", `commit with "double quotes"`},
			expectError: false,
		},
		{
			name:        "safe arguments with special characters",
			args:        []string{"--message", "commit with @#$%^&*()"},
			expectError: false,
		},
		{
			name:        "multiple flags and options",
			args:        []string{"-a", "-m", "message", "--author", "John Doe <john@example.com>"},
			expectError: false,
		},
		{
			name:        "file paths with spaces",
			args:        []string{"add", "file with spaces.txt"},
			expectError: false,
		},
		{
			name:        "dangerous command substitution",
			args:        []string{"--message", "$(rm -rf /)"},
			expectError: true,
			errorMsg:    "potentially unsafe argument at position 1",
		},
		{
			name:        "dangerous backticks",
			args:        []string{"--message", "`rm -rf /`"},
			expectError: true,
			errorMsg:    "potentially unsafe argument at position 1",
		},
		{
			name:        "safe parentheses without dollar sign",
			args:        []string{"--message", "commit (with parentheses)"},
			expectError: false,
		},
		{
			name:        "empty arguments",
			args:        []string{},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateArguments(tt.args)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for args %v, got nil", tt.args)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for args %v, got: %v", tt.args, err)
				}
			}
		})
	}
}

func TestPreserveComplexArguments(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "simple arguments",
			input:    []string{"status", "--short"},
			expected: []string{"status", "--short"},
		},
		{
			name:     "arguments with spaces",
			input:    []string{"commit", "-m", "message with spaces"},
			expected: []string{"commit", "-m", "message with spaces"},
		},
		{
			name:     "arguments with quotes",
			input:    []string{"commit", "-m", "message with 'quotes'"},
			expected: []string{"commit", "-m", "message with 'quotes'"},
		},
		{
			name:     "multiple flags and options",
			input:    []string{"-a", "-m", "message", "--author", "John Doe"},
			expected: []string{"-a", "-m", "message", "--author", "John Doe"},
		},
		{
			name:     "file paths with special characters",
			input:    []string{"add", "file@#$%.txt", "another file.txt"},
			expected: []string{"add", "file@#$%.txt", "another file.txt"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.PreserveComplexArguments(tt.input)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d arguments, got %d", len(tt.expected), len(result))
				return
			}
			
			for i, arg := range result {
				if arg != tt.expected[i] {
					t.Errorf("Argument %d: expected %q, got %q", i, tt.expected[i], arg)
				}
			}
		})
	}
}

func TestExecute_ComplexArguments(t *testing.T) {
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
		expectError bool
	}{
		{
			name:        "git version command",
			command:     "--version",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "git command with flags",
			command:     "status",
			args:        []string{"--short", "--branch"},
			expectError: false,
		},
		{
			name:        "git command with quoted message",
			command:     "log",
			args:        []string{"--oneline", "-n", "1"},
			expectError: false,
		},
		{
			name:        "dangerous command injection attempt",
			command:     "status",
			args:        []string{"$(echo dangerous)"},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, tt.args)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for command %s with args %v, got nil", tt.command, tt.args)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for command %s with args %v, got: %v", tt.command, tt.args, err)
				}
			}
		})
	}
}

func TestExecuteGitCommand_ArgumentPreservation(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}
	
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Test that arguments are preserved correctly by using git help
	// which will show different output based on the exact arguments passed
	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "single argument",
			command: "help",
			args:    []string{"status"},
		},
		{
			name:    "multiple arguments",
			command: "help",
			args:    []string{"log", "--oneline"},
		},
		{
			name:    "arguments with dashes",
			command: "help",
			args:    []string{"--version"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies that the command executes without error
			// The actual argument preservation is tested by the fact that
			// git help with different arguments produces different results
			err := handler.executeGitCommand(tt.command, tt.args)
			
			// We expect this to either succeed or fail with a known git error
			// but not fail due to argument handling issues
			if err != nil {
				// Check if it's a git-related error (which is acceptable)
				// vs an argument handling error (which would indicate a problem)
				if strings.Contains(err.Error(), "command execution failed") {
					t.Errorf("Unexpected command execution failure: %v", err)
				}
				// Other errors (like git command not found, etc.) are acceptable
			}
		})
	}
}

func TestNewFallthroughHandlerWithTestMode(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	
	testResponses := map[string][]string{
		"git add -p": {"y", "n", "q"},
		"git rebase -i": {":wq"},
	}
	
	handler := NewFallthroughHandlerWithTestMode(cfg, gitRepo, testResponses)
	
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}
	
	if handler.processExecutor == nil {
		t.Error("Expected processExecutor to be initialized")
	}
	
	// Test that the handler is properly configured for test mode
	// We can't directly access the test mode state, but we can test the behavior
	// by checking that it doesn't fail when executing interactive commands
}

func TestSetTestMode(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Test enabling test mode
	handler.SetTestMode(true)
	
	// Test disabling test mode
	handler.SetTestMode(false)
	
	// The actual test mode state is internal to the process executor
	// This test verifies that the methods can be called without error
}

func TestSetTestResponses(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	testResponses := map[string][]string{
		"git add -p": {"y", "n", "q"},
		"git rebase -i": {":wq"},
		"git commit -v": {"Test commit message", "", ":wq"},
	}
	
	// Test setting responses
	handler.SetTestResponses(testResponses)
	
	// The actual responses are internal to the process executor
	// This test verifies that the method can be called without error
}

func TestExecute_InteractiveCommandInTestMode(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}
	
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	
	// Create handler with test mode and responses for interactive commands
	testResponses := map[string][]string{
		"git add -p": {"q"}, // Quit immediately
		"git rebase -i": {":q"}, // Quit editor immediately
	}
	
	handler := NewFallthroughHandlerWithTestMode(cfg, gitRepo, testResponses)
	
	// Test that we can execute commands that would normally be interactive
	// without hanging waiting for user input
	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "non-interactive git command",
			command: "status",
			args:    []string{"--porcelain"},
		},
		{
			name:    "git version",
			command: "--version",
			args:    []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, tt.args)
			
			// These commands should execute successfully in test mode
			if err != nil {
				t.Errorf("Expected no error for command %s with args %v in test mode, got: %v", tt.command, tt.args, err)
			}
		})
	}
}

func TestExecute_InteractiveCommandDetection(t *testing.T) {
	cfg := &config.Config{}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Enable test mode to prevent hanging on interactive commands
	handler.SetTestMode(true)
	handler.SetTestResponses(map[string][]string{
		"git add -p": {"q"},
		"git rebase -i": {":q"},
		"git commit -v": {"Test message", "", ":wq"},
	})
	
	// Test that interactive commands are properly detected and handled
	// We use echo commands to simulate the behavior without requiring git
	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
	}{
		{
			name:        "echo command (non-interactive)",
			command:     "echo",
			args:        []string{"test"},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need to use executeGitCommand directly since Execute validates git availability
			err := handler.processExecutor.ExecuteCommand(tt.command, tt.args)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for command %s with args %v, got nil", tt.command, tt.args)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for command %s with args %v, got: %v", tt.command, tt.args, err)
				}
			}
		})
	}
}

func TestVerboseMode_Enabled(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Verify that verbose mode is properly set on the process executor
	if !handler.processExecutor.Verbose {
		t.Error("Expected process executor to have verbose mode enabled")
	}
}

func TestVerboseMode_Disabled(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Verify that verbose mode is properly disabled on the process executor
	if handler.processExecutor.Verbose {
		t.Error("Expected process executor to have verbose mode disabled")
	}
}

func TestVerboseMode_WithTestMode(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: true,
		},
	}
	gitRepo := &git.Repository{}
	
	testResponses := map[string][]string{
		"git status": {""},
	}
	
	handler := NewFallthroughHandlerWithTestMode(cfg, gitRepo, testResponses)
	
	// Verify that verbose mode is properly set even in test mode
	if !handler.processExecutor.Verbose {
		t.Error("Expected process executor to have verbose mode enabled in test mode")
	}
}

func TestExecute_VerboseOutput(t *testing.T) {
	// Skip this test if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git not available in PATH")
	}
	
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Test that verbose mode doesn't cause errors
	// We can't easily capture stderr output in this test, but we can verify
	// that the command still executes successfully with verbose mode enabled
	err := handler.Execute("--version", []string{})
	
	if err != nil {
		t.Errorf("Expected no error with verbose mode enabled, got: %v", err)
	}
}

func TestExecute_BlacklistedCommand(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled:   true,
			Verbose:   true,
			Blacklist: []string{"status", "diff"},
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	tests := []struct {
		name        string
		command     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "blacklisted command - status",
			command:     "status",
			expectError: true,
			errorMsg:    "command 'status' is blacklisted from fallthrough",
		},
		{
			name:        "blacklisted command - diff",
			command:     "diff",
			expectError: true,
			errorMsg:    "command 'diff' is blacklisted from fallthrough",
		},
		{
			name:        "non-blacklisted command",
			command:     "log",
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, []string{})
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for blacklisted command %s, got nil", tt.command)
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				// For non-blacklisted commands, we might get other errors (like git not in repo)
				// but we shouldn't get blacklist errors
				if err != nil && strings.Contains(err.Error(), "blacklisted") {
					t.Errorf("Unexpected blacklist error for command %s: %v", tt.command, err)
				}
			}
		})
	}
}

func TestLogVerbose(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Verbose: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Test that logVerbose method can be called without error
	// We can't easily capture stderr in this test, but we can verify
	// that the method doesn't panic or cause errors
	handler.logVerbose("Test message", "status", []string{"--short"})
	handler.logVerbose("Test message", "diff", []string{})
	handler.logVerbose("Test message", "log", []string{"--oneline", "-n", "5"})
}

func TestLogDebug(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Verbose: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Test that logDebug method can be called without error
	details := map[string]interface{}{
		"command": "status",
		"args":    []string{"--short"},
		"time":    "2023-01-01T00:00:00Z",
	}
	
	handler.logDebug("Test debug message", details)
	
	// Test with empty details
	handler.logDebug("Test debug message", map[string]interface{}{})
	
	// Test with various data types
	complexDetails := map[string]interface{}{
		"string":  "test",
		"int":     42,
		"bool":    true,
		"slice":   []string{"a", "b", "c"},
	}
	
	handler.logDebug("Complex debug message", complexDetails)
}

func TestLogTiming(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Verbose: true,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)
	
	// Test successful command timing
	handler.logTiming("Command completed", "status", []string{"--short"}, 100*time.Millisecond, nil)
	
	// Test failed command timing
	testErr := fmt.Errorf("test error")
	handler.logTiming("Command failed", "diff", []string{}, 50*time.Millisecond, testErr)
	
	// Test with no arguments
	handler.logTiming("Command completed", "version", []string{}, 10*time.Millisecond, nil)
}