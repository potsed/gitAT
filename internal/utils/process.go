package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// ProcessExecutor handles subprocess execution with real-time output streaming
type ProcessExecutor struct {
	testMode bool
	testResponses map[string][]string
	Verbose bool
}

// NewProcessExecutor creates a new ProcessExecutor instance
func NewProcessExecutor() *ProcessExecutor {
	return &ProcessExecutor{
		testMode: false,
		testResponses: make(map[string][]string),
		Verbose: false,
	}
}

// NewProcessExecutorWithTestMode creates a ProcessExecutor in test mode
func NewProcessExecutorWithTestMode(responses map[string][]string) *ProcessExecutor {
	if responses == nil {
		responses = make(map[string][]string)
	}
	return &ProcessExecutor{
		testMode: true,
		testResponses: responses,
		Verbose: false,
	}
}

// SetTestMode enables or disables test mode
func (p *ProcessExecutor) SetTestMode(enabled bool) {
	p.testMode = enabled
}

// SetTestResponses sets the automatic responses for test mode
func (p *ProcessExecutor) SetTestResponses(responses map[string][]string) {
	p.testResponses = responses
}

// SetVerbose enables or disables verbose logging
func (p *ProcessExecutor) SetVerbose(verbose bool) {
	p.Verbose = verbose
}

// ExecuteCommand executes a command with real-time output streaming
// Returns the exit code of the command or an error if execution fails
func (p *ProcessExecutor) ExecuteCommand(command string, args []string) error {
	startTime := time.Now()
	
	if p.Verbose {
		p.logDebug("Starting command execution", map[string]interface{}{
			"command": command,
			"args":    args,
			"time":    startTime.Format(time.RFC3339),
		})
	}
	
	// Check if this is likely an interactive command
	if p.IsInteractiveCommand(command, args) {
		err := p.ExecuteInteractiveCommand(command, args)
		if p.Verbose {
			p.logTiming("Interactive command completed", command, args, time.Since(startTime), err)
		}
		return err
	}
	
	cmd := exec.Command(command, args...)
	
	// Set up pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	
	// Connect stdin to allow interactive commands
	cmd.Stdin = os.Stdin
	
	// Start the command
	if err := cmd.Start(); err != nil {
		if p.Verbose {
			p.logDebug("Failed to start command", map[string]interface{}{
				"command": command,
				"args":    args,
				"error":   err.Error(),
			})
		}
		return fmt.Errorf("failed to start command: %w", err)
	}
	
	if p.Verbose {
		p.logDebug("Command started successfully", map[string]interface{}{
			"command": command,
			"args":    args,
			"pid":     cmd.Process.Pid,
		})
	}
	
	// Stream output in real-time
	go p.streamOutput(stdout, os.Stdout)
	go p.streamOutput(stderr, os.Stderr)
	
	// Wait for command to complete
	err = cmd.Wait()
	
	// Log timing information
	duration := time.Since(startTime)
	if p.Verbose {
		p.logTiming("Command execution completed", command, args, duration, err)
	}
	
	// Handle exit codes properly
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Get the exit code from the process state
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				// Create a custom error that preserves the exit code
				return &ExitError{
					ExitCode: status.ExitStatus(),
					Err:      exitError,
				}
			}
		}
		return fmt.Errorf("command execution failed: %w", err)
	}
	
	return nil
}

// ExecuteCommandWithOutput executes a command and returns its output
// This method captures output instead of streaming it
func (p *ProcessExecutor) ExecuteCommandWithOutput(command string, args []string) ([]byte, error) {
	// Handle test mode
	if p.testMode {
		cmdKey := p.GenerateCommandKey(command, args)
		if responses, exists := p.testResponses[cmdKey]; exists && len(responses) > 0 {
			// Return the first response as output
			return []byte(responses[0]), nil
		}
		// If no test response is defined, return empty output
		return []byte{}, nil
	}
	
	cmd := exec.Command(command, args...)
	
	// Connect stdin for interactive commands
	cmd.Stdin = os.Stdin
	
	// Execute and capture output
	output, err := cmd.CombinedOutput()
	
	// Handle exit codes properly
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				return output, &ExitError{
					ExitCode: status.ExitStatus(),
					Err:      exitError,
				}
			}
		}
		return output, fmt.Errorf("command execution failed: %w", err)
	}
	
	return output, nil
}

// ExecuteInteractiveCommand executes a command that requires interactive input
// Handles TTY setup and test mode automatic responses
func (p *ProcessExecutor) ExecuteInteractiveCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	
	if p.testMode {
		return p.executeInteractiveCommandInTestMode(cmd, command, args)
	}
	
	// For interactive commands, we need to connect all stdio directly
	// This allows proper TTY handling for commands like git add -p, git rebase -i
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// Run the command and wait for completion
	err := cmd.Run()
	
	// Handle exit codes properly
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				return &ExitError{
					ExitCode: status.ExitStatus(),
					Err:      exitError,
				}
			}
		}
		return fmt.Errorf("interactive command execution failed: %w", err)
	}
	
	return nil
}

// executeInteractiveCommandInTestMode handles interactive commands in test mode
func (p *ProcessExecutor) executeInteractiveCommandInTestMode(cmd *exec.Cmd, command string, args []string) error {
	// Create pipes for stdin to provide automatic responses
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	
	// Set up stdout and stderr pipes for capturing output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	
	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start interactive command: %w", err)
	}
	
	// Handle automatic responses in a goroutine
	go p.handleTestModeResponses(stdin, command, args)
	
	// Stream output
	go p.streamOutput(stdout, os.Stdout)
	go p.streamOutput(stderr, os.Stderr)
	
	// Wait for command to complete
	err = cmd.Wait()
	
	// Handle exit codes properly
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				return &ExitError{
					ExitCode: status.ExitStatus(),
					Err:      exitError,
				}
			}
		}
		return fmt.Errorf("interactive command execution failed: %w", err)
	}
	
	return nil
}

// handleTestModeResponses provides automatic responses for interactive commands in test mode
func (p *ProcessExecutor) handleTestModeResponses(stdin io.WriteCloser, command string, args []string) {
	defer stdin.Close()
	
	// Generate a key for this command
	cmdKey := p.GenerateCommandKey(command, args)
	
	// Get responses for this command
	responses, exists := p.testResponses[cmdKey]
	if !exists {
		// Use default responses based on command type
		responses = p.GetDefaultTestResponses(command, args)
	}
	
	// Send each response
	for _, response := range responses {
		fmt.Fprintln(stdin, response)
	}
}

// GenerateCommandKey creates a key for looking up test responses
func (p *ProcessExecutor) GenerateCommandKey(command string, args []string) string {
	if len(args) == 0 {
		return command
	}
	
	// For most commands, just join command and args
	return fmt.Sprintf("%s %s", command, strings.Join(args, " "))
}

// GetDefaultTestResponses returns default responses for common interactive commands
func (p *ProcessExecutor) GetDefaultTestResponses(command string, args []string) []string {
	if command == "git" && len(args) > 0 {
		subcommand := args[0]
		
		switch subcommand {
		case "add":
			// For git add -p, provide default responses
			if p.containsFlag(args, "-p", "--patch") {
				return []string{"y", "q"} // Yes to first hunk, then quit
			}
		case "rebase":
			// For git rebase -i, just continue with default editor behavior
			if p.containsFlag(args, "-i", "--interactive") {
				return []string{":wq"} // Save and quit in editor
			}
		case "merge":
			// For merge conflicts, accept default merge
			return []string{"", ""} // Just press enter for defaults
		case "commit":
			// For interactive commit, provide a default message
			if p.containsFlag(args, "-v", "--verbose") {
				return []string{"Test commit message", "", ":wq"}
			}
		}
	}
	
	// Default: just press enter
	return []string{""}
}

// containsFlag checks if args contains any of the specified flags
func (p *ProcessExecutor) containsFlag(args []string, flags ...string) bool {
	for _, arg := range args {
		for _, flag := range flags {
			if arg == flag {
				return true
			}
		}
	}
	return false
}

// IsInteractiveCommand determines if a command is likely to be interactive
func (p *ProcessExecutor) IsInteractiveCommand(command string, args []string) bool {
	if command == "git" && len(args) > 0 {
		subcommand := args[0]
		
		// Known interactive Git commands
		interactiveCommands := map[string]bool{
			"add":    p.containsFlag(args, "-p", "--patch", "-i", "--interactive"),
			"rebase": p.containsFlag(args, "-i", "--interactive"),
			"commit": p.containsFlag(args, "-v", "--verbose") || p.containsFlag(args, "-e", "--edit"),
			"merge":  true, // Merge can be interactive when there are conflicts
			"cherry-pick": p.containsFlag(args, "-e", "--edit"),
			"revert": p.containsFlag(args, "-e", "--edit"),
		}
		
		if interactive, exists := interactiveCommands[subcommand]; exists {
			return interactive
		}
	}
	
	// Default to non-interactive
	return false
}

// streamOutput streams data from reader to writer line by line
func (p *ProcessExecutor) streamOutput(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Fprintln(writer, scanner.Text())
	}
}

// ExitError represents an error with an exit code
type ExitError struct {
	ExitCode int
	Err      error
}

func (e *ExitError) Error() string {
	return e.Err.Error()
}

// GetExitCode returns the exit code from an error, or 0 if not an ExitError
func GetExitCode(err error) int {
	if exitErr, ok := err.(*ExitError); ok {
		return exitErr.ExitCode
	}
	return 0
}

// IsCommandNotFound checks if an error indicates the command was not found
func IsCommandNotFound(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for exec.ErrNotFound or "executable file not found" error
	if err == exec.ErrNotFound {
		return true
	}
	
	// Check error message for common "not found" patterns
	errMsg := err.Error()
	return strings.Contains(errMsg, "executable file not found") ||
		   strings.Contains(errMsg, "no such file or directory") ||
		   strings.Contains(errMsg, "command not found")
}

// logDebug logs debug information for process execution
func (p *ProcessExecutor) logDebug(message string, details map[string]interface{}) {
	fmt.Fprintf(os.Stderr, "[PROCESS DEBUG] %s:", message)
	for key, value := range details {
		fmt.Fprintf(os.Stderr, " %s=%v", key, value)
	}
	fmt.Fprintln(os.Stderr)
}

// logTiming logs timing information for performance monitoring
func (p *ProcessExecutor) logTiming(message, command string, args []string, duration time.Duration, err error) {
	status := "SUCCESS"
	exitCode := 0
	if err != nil {
		status = "ERROR"
		if exitErr, ok := err.(*ExitError); ok {
			exitCode = exitErr.ExitCode
		}
	}
	
	fmt.Fprintf(os.Stderr, "[PROCESS TIMING] %s: %s", message, command)
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, " %s", strings.Join(args, " "))
	}
	fmt.Fprintf(os.Stderr, " - Duration: %v, Status: %s, ExitCode: %d", duration, status, exitCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, ", Error: %v", err)
	}
	fmt.Fprintln(os.Stderr)
}