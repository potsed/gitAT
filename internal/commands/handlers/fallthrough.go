package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/internal/utils"
)

// FallthroughHandler handles commands that should fall through to standard Git
type FallthroughHandler struct {
	BaseHandler
	processExecutor  *utils.ProcessExecutor
	reservedCommands map[string]bool
}

// NewFallthroughHandler creates a new fallthrough handler
func NewFallthroughHandler(cfg *config.Config, gitRepo *git.Repository) *FallthroughHandler {
	processExecutor := utils.NewProcessExecutor()
	processExecutor.SetVerbose(cfg.Fallthrough.Verbose)

	return &FallthroughHandler{
		BaseHandler:     NewBaseHandler(cfg, gitRepo),
		processExecutor: processExecutor,
		reservedCommands: map[string]bool{
			// Gitat-specific commands that should not fall through
			"semver":  true, // Renamed from "version"
			"help":    true,
			"sprout":  true, // Renamed from "branch"
			"feature": true,
			"hotfix":  true,
			"info":    true,
			"issue":   true,
			"label":   true,
			"pr":      true,
			"product": true,
			"release": true,
			"save":    true,
			"squash":  true,
			"sweep":   true,
			"dub":     true, // Renamed from "tag"
			"wip":     true,
			"work":    true,
			// Utility commands
			"changes":      true,
			"logz":         true, // Renamed from "logs"
			"shasum":       true, // Renamed from "hash"
			"id":           true,
			"path":         true,
			"main":         true,
			"master":       true,
			"root":         true,
			"ignore":       true,
			"setup-local":  true,
			"setup-remote": true,
			"security":     true,
			"_go":          true,
			"_label":       true,
			"_id":          true,
			"_path":        true,
			"_trunk":       true,
			"_security":    true,
			// Special flags that should show Gitat info, not fall through
			"-v":        true,
			"--version": true,
			"-h":        true,
			"--help":    true,
		},
	}
}

// NewFallthroughHandlerWithTestMode creates a new fallthrough handler in test mode
func NewFallthroughHandlerWithTestMode(cfg *config.Config, gitRepo *git.Repository, testResponses map[string][]string) *FallthroughHandler {
	handler := NewFallthroughHandler(cfg, gitRepo)
	processExecutor := utils.NewProcessExecutorWithTestMode(testResponses)
	processExecutor.SetVerbose(cfg.Fallthrough.Verbose)
	handler.processExecutor = processExecutor
	return handler
}

// SetTestMode enables or disables test mode for the fallthrough handler
func (f *FallthroughHandler) SetTestMode(enabled bool) {
	f.processExecutor.SetTestMode(enabled)
}

// SetTestResponses sets the automatic responses for interactive commands in test mode
func (f *FallthroughHandler) SetTestResponses(responses map[string][]string) {
	f.processExecutor.SetTestResponses(responses)
}

// Execute executes a Git command through fallthrough mechanism
func (f *FallthroughHandler) Execute(command string, args []string) error {
	startTime := time.Now()

	// Log debug information if verbose mode is enabled
	if f.config.Fallthrough.Verbose {
		f.logDebug("Starting fallthrough execution", map[string]interface{}{
			"command": command,
			"args":    args,
			"time":    startTime.Format(time.RFC3339),
		})
	}

	// Check if fallthrough is enabled
	if !f.config.Fallthrough.Enabled {
		return f.createUnknownCommandError(command, args)
	}

	// Check if Git executable is available
	if !f.isGitAvailable() {
		return f.createGitNotFoundError()
	}

	// Validate that the command should fall through
	if !f.shouldFallthrough(command) {
		return f.createReservedCommandError(command, args)
	}

	// Check if command is blacklisted
	if f.config.IsFallthroughBlacklisted(command) {
		if f.config.Fallthrough.Verbose {
			f.logDebug("Command is blacklisted", map[string]interface{}{
				"command": command,
			})
		}
		return f.createBlacklistedCommandError(command, args)
	}

	// Validate arguments for safety (prevent command injection)
	if err := f.ValidateArguments(args); err != nil {
		if f.config.Fallthrough.Verbose {
			f.logDebug("Argument validation failed", map[string]interface{}{
				"command": command,
				"args":    args,
				"error":   err.Error(),
			})
		}
		return f.createArgumentValidationError(command, args, err)
	}

	// Validate complex commands, aliases, and subcommands (Requirement 5.2)
	if strings.TrimSpace(command) != "" {
		if err := f.ValidateComplexCommand(command, args); err != nil {
			if f.config.Fallthrough.Verbose {
				f.logDebug("Complex command validation failed", map[string]interface{}{
					"command": command,
					"args":    args,
					"error":   err.Error(),
				})
			}
			return f.createCommandValidationError(command, args, err)
		}
	}

	// Show verbose output when fallthrough occurs
	if f.config.Fallthrough.Verbose {
		f.logVerbose("Falling through to Git command", command, args)
	}

	// Preserve complex arguments and execute the Git command
	preservedArgs := f.PreserveComplexArguments(args)

	// Execute the command and measure timing
	err := f.executeGitCommand(command, preservedArgs)

	// Log timing information
	duration := time.Since(startTime)
	if f.config.Fallthrough.Verbose {
		f.logTiming("Command execution completed", command, args, duration, err)
	}

	// Enhance error propagation from underlying Git commands
	if err != nil {
		return f.enhanceGitError(command, args, err)
	}

	return nil
}

// isGitAvailable checks if the git executable is available in PATH
func (f *FallthroughHandler) isGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// shouldFallthrough determines if a command should fall through to Git
func (f *FallthroughHandler) shouldFallthrough(command string) bool {
	// Don't fall through reserved Gitat commands
	if f.reservedCommands[command] {
		return false
	}

	// Special case: empty command should fall through to show Git help (Requirement 4.1)
	if strings.TrimSpace(command) == "" {
		return true
	}

	// All other commands can fall through
	return true
}

// executeGitCommand executes the Git command with preserved arguments
func (f *FallthroughHandler) executeGitCommand(command string, args []string) error {
	var gitArgs []string

	// Special case: empty command means just run "git" with no subcommand (shows Git help)
	if strings.TrimSpace(command) == "" {
		gitArgs = args // Just pass through any arguments (should be empty for help case)
	} else {
		// Prepare the full argument list for Git
		gitArgs = []string{command}

		// Preserve all arguments exactly as provided
		// Go's exec.Command handles argument separation correctly,
		// so we don't need to do any special escaping or parsing
		gitArgs = append(gitArgs, args...)
	}

	// Execute the Git command with preserved arguments
	return f.processExecutor.ExecuteCommand("git", gitArgs)
}

// ValidateArguments validates that arguments are safe to pass through
// This helps prevent command injection while preserving legitimate use cases
func (f *FallthroughHandler) ValidateArguments(args []string) error {
	for i, arg := range args {
		// Check for potentially dangerous patterns
		if strings.Contains(arg, "$(") || strings.Contains(arg, "`") {
			return fmt.Errorf("potentially unsafe argument at position %d: %s", i, arg)
		}

		// Allow all other arguments including those with quotes, spaces, etc.
		// Go's exec.Command handles these safely
	}
	return nil
}

// PreserveComplexArguments ensures complex arguments are handled correctly
// This method demonstrates how different argument types are preserved
func (f *FallthroughHandler) PreserveComplexArguments(args []string) []string {
	// Go's exec.Command already handles argument preservation correctly:
	// - Arguments with spaces are preserved as single arguments
	// - Quoted arguments maintain their content without the quotes being passed to the process
	// - Special characters are preserved within arguments
	// - Multiple flags and options are handled as separate arguments

	// We don't need to modify the arguments, just return them as-is
	// The key is that each element in the slice represents one argument
	return args
}

// IsGitCommand checks if a command is a valid Git command
// This is a helper method for validation
func (f *FallthroughHandler) IsGitCommand(command string) bool {
	if !f.isGitAvailable() {
		return false
	}

	// Try to get help for the command to see if it exists
	// We use git help <command> which returns 0 for valid commands
	output, err := f.processExecutor.ExecuteCommandWithOutput("git", []string{"help", command})

	// If the command succeeded and doesn't contain "No manual entry"
	if err == nil {
		outputStr := string(output)
		return !strings.Contains(outputStr, "No manual entry") &&
			!strings.Contains(outputStr, "is not a git command")
	}

	// If there was an exit error, check the exit code
	if exitErr, ok := err.(*utils.ExitError); ok {
		// Git help returns 1 for unknown commands, 0 for valid ones
		return exitErr.ExitCode == 0
	}

	return false
}

// IsGitAlias checks if a command is a Git alias
// Requirement 5.2: Handle Git aliases correctly
func (f *FallthroughHandler) IsGitAlias(command string) bool {
	if !f.isGitAvailable() {
		return false
	}

	// Check if the command is defined as an alias in Git config
	output, err := f.processExecutor.ExecuteCommandWithOutput("git", []string{"config", "--get", "alias." + command})

	// If we got output without error, it's an alias
	return err == nil && len(strings.TrimSpace(string(output))) > 0
}

// HandleSubcommands processes commands that may have subcommands
// Requirement 5.2: Handle nested command structures correctly
func (f *FallthroughHandler) HandleSubcommands(command string, args []string) (bool, error) {
	// Some Git commands have subcommands (e.g., git remote add, git config --global)
	// We need to handle these as complete command structures

	// Commands known to have subcommands
	subcommandCommands := map[string]bool{
		"remote":        true,
		"config":        true,
		"submodule":     true,
		"worktree":      true,
		"notes":         true,
		"bundle":        true,
		"credential":    true,
		"filter-branch": true,
	}

	// If this is a subcommand-capable command, validate the full structure
	if subcommandCommands[command] {
		// For subcommand validation, we can try a dry-run or help to see if it's valid
		// But for now, we'll trust Git to handle the validation
		return true, nil
	}

	return false, nil
}

// ValidateComplexCommand validates commands with complex argument structures
func (f *FallthroughHandler) ValidateComplexCommand(command string, args []string) error {
	// Check if it's a Git alias first
	if f.IsGitAlias(command) {
		if f.config.Fallthrough.Verbose {
			f.logDebug("Command is a Git alias", map[string]interface{}{
				"command": command,
				"args":    args,
			})
		}
		return nil
	}

	// Check if it's a valid Git command
	if f.IsGitCommand(command) {
		return nil
	}

	// Check if it has subcommands
	if hasSubcommands, err := f.HandleSubcommands(command, args); err != nil {
		return err
	} else if hasSubcommands {
		return nil
	}

	// If none of the above, it might still be a valid Git command that we can't validate
	// Let Git itself handle the final validation
	return nil
}

// logVerbose logs verbose output when fallthrough occurs
func (f *FallthroughHandler) logVerbose(message, command string, args []string) {
	fmt.Fprintf(os.Stderr, "[GITAT] %s: git %s", message, command)
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, " %s", strings.Join(args, " "))
	}
	fmt.Fprintln(os.Stderr)
}

// logDebug logs debug information for command execution and argument processing
func (f *FallthroughHandler) logDebug(message string, details map[string]interface{}) {
	fmt.Fprintf(os.Stderr, "[GITAT DEBUG] %s:", message)
	for key, value := range details {
		fmt.Fprintf(os.Stderr, " %s=%v", key, value)
	}
	fmt.Fprintln(os.Stderr)
}

// logTiming logs timing information for performance monitoring
func (f *FallthroughHandler) logTiming(message, command string, args []string, duration time.Duration, err error) {
	status := "SUCCESS"
	if err != nil {
		status = "ERROR"
	}

	fmt.Fprintf(os.Stderr, "[GITAT TIMING] %s: git %s", message, command)
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, " %s", strings.Join(args, " "))
	}
	fmt.Fprintf(os.Stderr, " - Duration: %v, Status: %s", duration, status)
	if err != nil {
		fmt.Fprintf(os.Stderr, ", Error: %v", err)
	}
	fmt.Fprintln(os.Stderr)
}

// Error creation methods for comprehensive error handling

// createGitNotFoundError creates a clear error message when Git executable is not found
func (f *FallthroughHandler) createGitNotFoundError() error {
	return fmt.Errorf(`Git executable not found in PATH.

To use Gitat's fallthrough functionality, you need to have Git installed.

Installation options:
  • macOS: brew install git
  • Ubuntu/Debian: sudo apt-get install git
  • CentOS/RHEL: sudo yum install git
  • Windows: Download from https://git-scm.com/download/win

After installation, make sure 'git' is available in your PATH.
You can test this by running: git --version`)
}

// createUnknownCommandError creates a helpful error message for unknown commands
func (f *FallthroughHandler) createUnknownCommandError(command string, args []string) error {
	var suggestions []string

	// Suggest similar Gitat commands
	gitAtCommands := []string{
		"branch", "feature", "hotfix", "info", "issue", "label", "pr",
		"product", "release", "save", "squash", "sweep", "tag", "wip", "work",
		"changes", "logz", "shasum", "id", "path", "main", "master", "root", "ignore",
	}

	// Find similar commands using simple string matching
	for _, cmd := range gitAtCommands {
		if strings.Contains(cmd, command) || strings.Contains(command, cmd) {
			suggestions = append(suggestions, cmd)
		}
	}

	errorMsg := fmt.Sprintf("Unknown command: %s", command)

	if len(suggestions) > 0 {
		errorMsg += fmt.Sprintf("\n\nDid you mean one of these Gitat commands?\n")
		for _, suggestion := range suggestions {
			errorMsg += fmt.Sprintf("  • git @ %s\n", suggestion)
		}
	}

	errorMsg += fmt.Sprintf(`
Fallthrough is currently disabled. To enable Git command fallthrough:
  git config at.fallthrough.enabled true

Available Gitat commands:
  • git @ help        - Show Gitat help
  • git @ version     - Show Gitat version
  • git @ info        - Show repository information
  • git @ branch      - Branch management
  • git @ feature     - Feature branch operations
  • git @ work        - Work with current changes

For a complete list of commands, run: git @ help`)

	return fmt.Errorf("%s", errorMsg)
}

// createReservedCommandError creates an error for reserved Gitat commands
func (f *FallthroughHandler) createReservedCommandError(command string, args []string) error {
	return fmt.Errorf(`'%s' is a reserved Gitat command.

Use 'git @ %s' to access Gitat functionality, or 'git %s' to use the standard Git command directly.

For help with Gitat commands, run: git @ help`, command, command, command)
}

// createBlacklistedCommandError creates an error for blacklisted commands
func (f *FallthroughHandler) createBlacklistedCommandError(command string, args []string) error {
	return fmt.Errorf(`Command '%s' is blacklisted from fallthrough.

This command has been explicitly disabled for security or compatibility reasons.
To use this command, run it directly with Git: git %s

To remove this command from the blacklist:
  git config --unset at.fallthrough.blacklist.%s

Current blacklist: %v`, command, command, command, f.config.Fallthrough.Blacklist)
}

// createArgumentValidationError creates an error for argument validation failures
func (f *FallthroughHandler) createArgumentValidationError(command string, args []string, validationErr error) error {
	return fmt.Errorf(`Argument validation failed for command '%s'.

Error: %v

This error occurred while validating arguments for security.
If you believe this is a false positive, please report it as an issue.

To bypass this validation, run the command directly with Git:
  git %s %s`, command, validationErr, command, strings.Join(args, " "))
}

// createCommandValidationError creates an error for command validation failures
func (f *FallthroughHandler) createCommandValidationError(command string, args []string, validationErr error) error {
	return fmt.Errorf(`Command validation failed for '%s'.

Error: %v

This might be because:
  • The command doesn't exist in your Git installation
  • The command syntax is incorrect
  • The command requires additional setup

To verify the command exists, try running it directly:
  git %s %s

For Git command help, run: git help %s`, command, validationErr, command, strings.Join(args, " "), command)
}

// enhanceGitError enhances error messages from underlying Git commands
func (f *FallthroughHandler) enhanceGitError(command string, args []string, gitErr error) error {
	// Check if this is an ExitError with exit code
	if exitErr, ok := gitErr.(*utils.ExitError); ok {
		// Preserve the original exit code but enhance the error message
		switch exitErr.ExitCode {
		case 1:
			// Common Git error - usually means command failed but is valid
			return &utils.ExitError{
				ExitCode: exitErr.ExitCode,
				Err: fmt.Errorf("Git command failed: git %s %s\n\nOriginal error: %v\n\nThis error came from Git itself. For help with this Git command, run:\n  git help %s",
					command, strings.Join(args, " "), exitErr.Err, command),
			}
		case 127:
			// Command not found
			return fmt.Errorf(`Git command not found: '%s'

This usually means:
  • The command doesn't exist in your Git version
  • There's a typo in the command name
  • The command requires a Git plugin or extension

To see available Git commands, run: git help -a
For help with a specific command, run: git help <command>`, command)
		case 128:
			// Git repository error
			return &utils.ExitError{
				ExitCode: exitErr.ExitCode,
				Err:      fmt.Errorf("Git repository error: %v\n\nThis usually means you're not in a Git repository or the repository is corrupted.\nTo initialize a new repository, run: git init", exitErr.Err),
			}
		default:
			// Other exit codes - preserve original error but add context
			return &utils.ExitError{
				ExitCode: exitErr.ExitCode,
				Err: fmt.Errorf("Git command failed with exit code %d: git %s %s\n\nOriginal error: %v",
					exitErr.ExitCode, command, strings.Join(args, " "), exitErr.Err),
			}
		}
	}

	// Check for command not found errors
	if utils.IsCommandNotFound(gitErr) {
		return f.createGitNotFoundError()
	}

	// For other errors, add helpful context
	return fmt.Errorf("Git command execution failed: git %s %s\n\nError: %v\n\nFor help with this command, run: git help %s",
		command, strings.Join(args, " "), gitErr, command)
}

// GetAvailableGitAtCommands returns a list of available Gitat commands for suggestions
func (f *FallthroughHandler) GetAvailableGitAtCommands() []string {
	return []string{
		"branch", "feature", "hotfix", "info", "issue", "label", "pr",
		"product", "release", "save", "squash", "sweep", "tag", "wip", "work",
		"changes", "logz", "shasum", "id", "path", "main", "master", "root", "ignore",
		"initlocal", "initremote", "security", "version", "help",
	}
}

// SuggestSimilarCommands suggests similar commands based on string similarity
func (f *FallthroughHandler) SuggestSimilarCommands(command string) []string {
	var suggestions []string
	availableCommands := f.GetAvailableGitAtCommands()

	// Simple similarity check - contains or is contained
	for _, cmd := range availableCommands {
		if strings.Contains(cmd, command) || strings.Contains(command, cmd) {
			suggestions = append(suggestions, cmd)
		}
	}

	// Try common typos and abbreviations (even if we have direct matches)
	commonMappings := map[string][]string{
		"st":    {"status"}, // Common Git alias, suggest fallthrough
		"co":    {"checkout"},
		"br":    {"branch"},
		"ci":    {"commit"},
		"stat":  {"status"},
		"check": {"checkout"},
		"comm":  {"commit"},
		"feat":  {"feature"},
		"hot":   {"hotfix"},
		"rel":   {"release"},
		"inf":   {"info"},
		"iss":   {"issue"},
		"lab":   {"label"},
	}

	if mapped, exists := commonMappings[command]; exists {
		// Add mapped suggestions, avoiding duplicates
		for _, mappedCmd := range mapped {
			found := false
			for _, existing := range suggestions {
				if existing == mappedCmd {
					found = true
					break
				}
			}
			if !found {
				suggestions = append(suggestions, mappedCmd)
			}
		}
	}

	return suggestions
}
