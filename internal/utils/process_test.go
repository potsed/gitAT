package utils

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestNewProcessExecutor(t *testing.T) {
	executor := NewProcessExecutor()
	if executor == nil {
		t.Fatal("NewProcessExecutor should return a non-nil ProcessExecutor")
	}
}

func TestExecuteCommand_Success(t *testing.T) {
	executor := NewProcessExecutor()
	
	// Test with a simple command that should succeed
	err := executor.ExecuteCommand("echo", []string{"hello world"})
	if err != nil {
		t.Fatalf("ExecuteCommand should succeed for valid command, got error: %v", err)
	}
}

func TestExecuteCommand_NonZeroExit(t *testing.T) {
	executor := NewProcessExecutor()
	
	// Test with a command that returns non-zero exit code
	err := executor.ExecuteCommand("sh", []string{"-c", "exit 42"})
	if err == nil {
		t.Fatal("ExecuteCommand should return error for non-zero exit code")
	}
	
	// Check that we get the correct exit code
	exitCode := GetExitCode(err)
	if exitCode != 42 {
		t.Fatalf("Expected exit code 42, got %d", exitCode)
	}
}

func TestExecuteCommand_CommandNotFound(t *testing.T) {
	executor := NewProcessExecutor()
	
	// Test with a command that doesn't exist
	err := executor.ExecuteCommand("nonexistentcommand12345", []string{})
	if err == nil {
		t.Fatal("ExecuteCommand should return error for non-existent command")
	}
	
	if !IsCommandNotFound(err) {
		t.Fatalf("Error should be recognized as command not found: %v", err)
	}
}

func TestExecuteCommandWithOutput_Success(t *testing.T) {
	executor := NewProcessExecutor()
	
	// Test command that produces output
	output, err := executor.ExecuteCommandWithOutput("echo", []string{"test output"})
	if err != nil {
		t.Fatalf("ExecuteCommandWithOutput should succeed, got error: %v", err)
	}
	
	expectedOutput := "test output\n"
	if string(output) != expectedOutput {
		t.Fatalf("Expected output %q, got %q", expectedOutput, string(output))
	}
}

func TestExecuteCommandWithOutput_NonZeroExit(t *testing.T) {
	executor := NewProcessExecutor()
	
	// Test command that returns non-zero exit and produces output
	output, err := executor.ExecuteCommandWithOutput("sh", []string{"-c", "echo 'error message' >&2; exit 1"})
	if err == nil {
		t.Fatal("ExecuteCommandWithOutput should return error for non-zero exit code")
	}
	
	// Check that we still get the output
	if !strings.Contains(string(output), "error message") {
		t.Fatalf("Expected output to contain 'error message', got %q", string(output))
	}
	
	// Check exit code
	exitCode := GetExitCode(err)
	if exitCode != 1 {
		t.Fatalf("Expected exit code 1, got %d", exitCode)
	}
}

func TestStreamOutput(t *testing.T) {
	executor := NewProcessExecutor()
	
	// Create a test input
	input := strings.NewReader("line1\nline2\nline3\n")
	
	// Capture output
	var output bytes.Buffer
	
	// Stream the input to output
	executor.streamOutput(input, &output)
	
	expected := "line1\nline2\nline3\n"
	if output.String() != expected {
		t.Fatalf("Expected output %q, got %q", expected, output.String())
	}
}

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: 0,
		},
		{
			name:     "regular error",
			err:      exec.ErrNotFound,
			expected: 0,
		},
		{
			name:     "exit error",
			err:      &ExitError{ExitCode: 42, Err: exec.ErrNotFound},
			expected: 42,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetExitCode(tt.err)
			if result != tt.expected {
				t.Fatalf("Expected exit code %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestIsCommandNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "exec.ErrNotFound",
			err:      exec.ErrNotFound,
			expected: true,
		},
		{
			name:     "PATH error",
			err:      &exec.Error{Name: "test", Err: exec.ErrNotFound},
			expected: true, // This contains "executable file not found"
		},
		{
			name:     "other error",
			err:      &ExitError{ExitCode: 1, Err: exec.ErrNotFound},
			expected: true, // This also contains "executable file not found"
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCommandNotFound(tt.err)
			if result != tt.expected {
				t.Fatalf("Expected %v, got %v for error: %v", tt.expected, result, tt.err)
			}
		})
	}
}

func TestExitError_Error(t *testing.T) {
	originalErr := exec.ErrNotFound
	exitErr := &ExitError{
		ExitCode: 42,
		Err:      originalErr,
	}
	
	if exitErr.Error() != originalErr.Error() {
		t.Fatalf("ExitError.Error() should return the wrapped error's message")
	}
}

func TestNewProcessExecutorWithTestMode(t *testing.T) {
	responses := map[string][]string{
		"git add -p": {"y", "q"},
	}
	
	executor := NewProcessExecutorWithTestMode(responses)
	if executor == nil {
		t.Fatal("NewProcessExecutorWithTestMode should return a non-nil ProcessExecutor")
	}
	
	if !executor.testMode {
		t.Fatal("ProcessExecutor should be in test mode")
	}
	
	if len(executor.testResponses) != 1 {
		t.Fatalf("Expected 1 test response, got %d", len(executor.testResponses))
	}
}

func TestSetTestMode(t *testing.T) {
	executor := NewProcessExecutor()
	
	// Initially should not be in test mode
	if executor.testMode {
		t.Fatal("ProcessExecutor should not be in test mode initially")
	}
	
	// Enable test mode
	executor.SetTestMode(true)
	if !executor.testMode {
		t.Fatal("ProcessExecutor should be in test mode after SetTestMode(true)")
	}
	
	// Disable test mode
	executor.SetTestMode(false)
	if executor.testMode {
		t.Fatal("ProcessExecutor should not be in test mode after SetTestMode(false)")
	}
}

func TestSetTestResponses(t *testing.T) {
	executor := NewProcessExecutor()
	
	responses := map[string][]string{
		"git add -p": {"y", "n", "q"},
		"git rebase -i": {":wq"},
	}
	
	executor.SetTestResponses(responses)
	
	if len(executor.testResponses) != 2 {
		t.Fatalf("Expected 2 test responses, got %d", len(executor.testResponses))
	}
	
	if len(executor.testResponses["git add -p"]) != 3 {
		t.Fatalf("Expected 3 responses for 'git add -p', got %d", len(executor.testResponses["git add -p"]))
	}
}

func TestIsInteractiveCommand(t *testing.T) {
	executor := NewProcessExecutor()
	
	tests := []struct {
		name     string
		command  string
		args     []string
		expected bool
	}{
		{
			name:     "git add -p",
			command:  "git",
			args:     []string{"add", "-p"},
			expected: true,
		},
		{
			name:     "git add --patch",
			command:  "git",
			args:     []string{"add", "--patch"},
			expected: true,
		},
		{
			name:     "git add -i",
			command:  "git",
			args:     []string{"add", "-i"},
			expected: true,
		},
		{
			name:     "git add file.txt",
			command:  "git",
			args:     []string{"add", "file.txt"},
			expected: false,
		},
		{
			name:     "git rebase -i",
			command:  "git",
			args:     []string{"rebase", "-i", "HEAD~3"},
			expected: true,
		},
		{
			name:     "git rebase main",
			command:  "git",
			args:     []string{"rebase", "main"},
			expected: false,
		},
		{
			name:     "git commit -v",
			command:  "git",
			args:     []string{"commit", "-v"},
			expected: true,
		},
		{
			name:     "git commit -m",
			command:  "git",
			args:     []string{"commit", "-m", "message"},
			expected: false,
		},
		{
			name:     "git merge",
			command:  "git",
			args:     []string{"merge", "feature-branch"},
			expected: true,
		},
		{
			name:     "git status",
			command:  "git",
			args:     []string{"status"},
			expected: false,
		},
		{
			name:     "non-git command",
			command:  "echo",
			args:     []string{"hello"},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.IsInteractiveCommand(tt.command, tt.args)
			if result != tt.expected {
				t.Fatalf("Expected %v for command '%s %v', got %v", tt.expected, tt.command, tt.args, result)
			}
		})
	}
}

func TestContainsFlag(t *testing.T) {
	executor := NewProcessExecutor()
	
	tests := []struct {
		name     string
		args     []string
		flags    []string
		expected bool
	}{
		{
			name:     "single flag present",
			args:     []string{"add", "-p", "file.txt"},
			flags:    []string{"-p"},
			expected: true,
		},
		{
			name:     "multiple flags, one present",
			args:     []string{"add", "--patch", "file.txt"},
			flags:    []string{"-p", "--patch"},
			expected: true,
		},
		{
			name:     "flag not present",
			args:     []string{"add", "file.txt"},
			flags:    []string{"-p", "--patch"},
			expected: false,
		},
		{
			name:     "empty args",
			args:     []string{},
			flags:    []string{"-p"},
			expected: false,
		},
		{
			name:     "empty flags",
			args:     []string{"add", "-p"},
			flags:    []string{},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.containsFlag(tt.args, tt.flags...)
			if result != tt.expected {
				t.Fatalf("Expected %v for args %v and flags %v, got %v", tt.expected, tt.args, tt.flags, result)
			}
		})
	}
}

func TestGenerateCommandKey(t *testing.T) {
	executor := NewProcessExecutor()
	
	tests := []struct {
		name     string
		command  string
		args     []string
		expected string
	}{
		{
			name:     "git add -p",
			command:  "git",
			args:     []string{"add", "-p"},
			expected: "git add -p",
		},
		{
			name:     "git add file",
			command:  "git",
			args:     []string{"add", "file.txt"},
			expected: "git add",
		},
		{
			name:     "git rebase -i",
			command:  "git",
			args:     []string{"rebase", "-i", "HEAD~3"},
			expected: "git rebase -i",
		},
		{
			name:     "simple command",
			command:  "echo",
			args:     []string{"hello"},
			expected: "echo hello",
		},
		{
			name:     "command without args",
			command:  "pwd",
			args:     []string{},
			expected: "pwd",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.GenerateCommandKey(tt.command, tt.args)
			if result != tt.expected {
				t.Fatalf("Expected %q for command '%s %v', got %q", tt.expected, tt.command, tt.args, result)
			}
		})
	}
}

func TestGetDefaultTestResponses(t *testing.T) {
	executor := NewProcessExecutor()
	
	tests := []struct {
		name     string
		command  string
		args     []string
		expected []string
	}{
		{
			name:     "git add -p",
			command:  "git",
			args:     []string{"add", "-p"},
			expected: []string{"y", "q"},
		},
		{
			name:     "git rebase -i",
			command:  "git",
			args:     []string{"rebase", "-i", "HEAD~3"},
			expected: []string{":wq"},
		},
		{
			name:     "git merge",
			command:  "git",
			args:     []string{"merge", "feature"},
			expected: []string{"", ""},
		},
		{
			name:     "git commit -v",
			command:  "git",
			args:     []string{"commit", "-v"},
			expected: []string{"Test commit message", "", ":wq"},
		},
		{
			name:     "unknown command",
			command:  "unknown",
			args:     []string{"arg"},
			expected: []string{""},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.GetDefaultTestResponses(tt.command, tt.args)
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d responses, got %d", len(tt.expected), len(result))
			}
			
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Fatalf("Expected response %d to be %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestExecuteInteractiveCommand_TestMode(t *testing.T) {
	// Create executor in test mode with custom responses
	responses := map[string][]string{
		"git add -p": {"y", "n", "q"},
	}
	executor := NewProcessExecutorWithTestMode(responses)
	
	// Test with echo command to simulate interactive behavior
	// We can't easily test real git commands in unit tests, so we use echo
	err := executor.ExecuteInteractiveCommand("echo", []string{"test"})
	if err != nil {
		t.Fatalf("ExecuteInteractiveCommand should succeed in test mode, got error: %v", err)
	}
}

func TestExecuteCommand_InteractiveDetection(t *testing.T) {
	executor := NewProcessExecutor()
	executor.SetTestMode(true)
	
	// Test that interactive commands are detected and routed correctly
	// Using echo as a safe command for testing
	err := executor.ExecuteCommand("echo", []string{"test"})
	if err != nil {
		t.Fatalf("ExecuteCommand should handle non-interactive commands correctly, got error: %v", err)
	}
}

func TestSetVerbose(t *testing.T) {
	executor := NewProcessExecutor()
	
	// Initially should not be in verbose mode
	if executor.Verbose {
		t.Fatal("ProcessExecutor should not be in verbose mode initially")
	}
	
	// Enable verbose mode
	executor.SetVerbose(true)
	if !executor.Verbose {
		t.Fatal("ProcessExecutor should be in verbose mode after SetVerbose(true)")
	}
	
	// Disable verbose mode
	executor.SetVerbose(false)
	if executor.Verbose {
		t.Fatal("ProcessExecutor should not be in verbose mode after SetVerbose(false)")
	}
}

func TestVerboseMode_Enabled(t *testing.T) {
	executor := NewProcessExecutor()
	executor.SetVerbose(true)
	
	// Test that verbose mode doesn't cause errors during command execution
	err := executor.ExecuteCommand("echo", []string{"test verbose"})
	if err != nil {
		t.Fatalf("ExecuteCommand should succeed with verbose mode enabled, got error: %v", err)
	}
}

func TestVerboseMode_Disabled(t *testing.T) {
	executor := NewProcessExecutor()
	executor.SetVerbose(false)
	
	// Test that commands work normally with verbose mode disabled
	err := executor.ExecuteCommand("echo", []string{"test non-verbose"})
	if err != nil {
		t.Fatalf("ExecuteCommand should succeed with verbose mode disabled, got error: %v", err)
	}
}

func TestVerboseMode_WithTestMode(t *testing.T) {
	responses := map[string][]string{
		"echo test": {""},
	}
	executor := NewProcessExecutorWithTestMode(responses)
	executor.SetVerbose(true)
	
	// Test that verbose mode works with test mode
	err := executor.ExecuteCommand("echo", []string{"test"})
	if err != nil {
		t.Fatalf("ExecuteCommand should succeed with both verbose and test mode enabled, got error: %v", err)
	}
}

func TestVerboseMode_WithInteractiveCommand(t *testing.T) {
	executor := NewProcessExecutor()
	executor.SetVerbose(true)
	executor.SetTestMode(true)
	
	// Set up test responses for interactive command
	responses := map[string][]string{
		"echo interactive": {"test response"},
	}
	executor.SetTestResponses(responses)
	
	// Test that verbose mode works with interactive commands
	err := executor.ExecuteInteractiveCommand("echo", []string{"interactive"})
	if err != nil {
		t.Fatalf("ExecuteInteractiveCommand should succeed with verbose mode enabled, got error: %v", err)
	}
}

func TestLogDebug(t *testing.T) {
	executor := NewProcessExecutor()
	executor.SetVerbose(true)
	
	// Test that logDebug method can be called without error
	details := map[string]interface{}{
		"command": "echo",
		"args":    []string{"test"},
		"time":    "2023-01-01T00:00:00Z",
	}
	
	executor.logDebug("Test debug message", details)
	
	// Test with empty details
	executor.logDebug("Test debug message", map[string]interface{}{})
	
	// Test with various data types
	complexDetails := map[string]interface{}{
		"string":  "test",
		"int":     42,
		"bool":    true,
		"slice":   []string{"a", "b", "c"},
	}
	
	executor.logDebug("Complex debug message", complexDetails)
}

func TestLogTiming(t *testing.T) {
	executor := NewProcessExecutor()
	executor.SetVerbose(true)
	
	// Test successful command timing
	executor.logTiming("Command completed", "echo", []string{"test"}, 100*time.Millisecond, nil)
	
	// Test failed command timing
	testErr := &ExitError{ExitCode: 1, Err: exec.ErrNotFound}
	executor.logTiming("Command failed", "false", []string{}, 50*time.Millisecond, testErr)
	
	// Test with no arguments
	executor.logTiming("Command completed", "true", []string{}, 10*time.Millisecond, nil)
	
	// Test with regular error (not ExitError)
	regularErr := exec.ErrNotFound
	executor.logTiming("Command error", "nonexistent", []string{}, 5*time.Millisecond, regularErr)
}

func TestExecuteCommand_VerboseTiming(t *testing.T) {
	executor := NewProcessExecutor()
	executor.SetVerbose(true)
	
	// Test that timing information is logged for successful commands
	err := executor.ExecuteCommand("echo", []string{"timing test"})
	if err != nil {
		t.Fatalf("ExecuteCommand should succeed, got error: %v", err)
	}
	
	// Test that timing information is logged for failed commands
	err = executor.ExecuteCommand("sh", []string{"-c", "exit 1"})
	if err == nil {
		t.Fatal("ExecuteCommand should fail for exit 1")
	}
	
	// Verify it's an ExitError with correct code
	exitCode := GetExitCode(err)
	if exitCode != 1 {
		t.Fatalf("Expected exit code 1, got %d", exitCode)
	}
}

func TestExecuteCommandWithOutput_VerboseMode(t *testing.T) {
	executor := NewProcessExecutor()
	executor.SetVerbose(true)
	
	// Test that ExecuteCommandWithOutput works with verbose mode
	// Note: ExecuteCommandWithOutput doesn't currently have verbose logging,
	// but we test that it still works when verbose mode is enabled
	output, err := executor.ExecuteCommandWithOutput("echo", []string{"verbose output test"})
	if err != nil {
		t.Fatalf("ExecuteCommandWithOutput should succeed with verbose mode, got error: %v", err)
	}
	
	expected := "verbose output test\n"
	if string(output) != expected {
		t.Fatalf("Expected output %q, got %q", expected, string(output))
	}
}