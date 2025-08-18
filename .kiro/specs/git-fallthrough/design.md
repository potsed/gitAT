# Design Document: Git Fallthrough Feature

## Overview

The Git fallthrough feature enhances Gitat by adding seamless integration with standard Git commands. When a user executes `git @ <command>` and no Gitat-specific handler exists for that command, the system will automatically execute the equivalent `git <command>` with all arguments preserved. This design maintains backward compatibility while making Gitat a complete replacement for the `git` command.

## Architecture

### Current Command Flow
```
User Input: git @ <command> [args]
    ↓
CLI App (pkg/cli/app.go)
    ↓
Commands Manager (internal/commands/manager.go)
    ↓
Switch Statement → Specific Handler OR Error
```

### Enhanced Command Flow with Fallthrough
```
User Input: git @ <command> [args]
    ↓
CLI App (pkg/cli/app.go)
    ↓
Commands Manager (internal/commands/manager.go)
    ↓
Switch Statement → Specific Handler OR Fallthrough Handler
    ↓
Fallthrough Handler → Execute: git <command> [args]
```

## Components and Interfaces

### 1. Fallthrough Handler

**Location:** `internal/commands/handlers/fallthrough.go`

```go
type FallthroughHandler struct {
    config *config.Config
    git    *git.Repository
}

func NewFallthroughHandler(cfg *config.Config, gitRepo *git.Repository) *FallthroughHandler

func (f *FallthroughHandler) Execute(command string, args []string) error
func (f *FallthroughHandler) isGitCommand(command string) bool
func (f *FallthroughHandler) shouldFallthrough(command string) bool
```

**Responsibilities:**
- Execute standard Git commands with preserved arguments
- Validate that the command should fall through (not reserved Gitat commands)
- Handle process execution and output streaming
- Preserve exit codes and error handling

### 2. Enhanced Commands Manager

**Location:** `internal/commands/manager.go` (modified)

**Changes:**
- Add fallthrough handler instance
- Modify `Execute()` method to use fallthrough as default case
- Add configuration for fallthrough behavior

```go
type Manager struct {
    // ... existing fields ...
    fallthrough *handlers.FallthroughHandler
}

func (m *Manager) Execute(command string, args []string) error {
    switch command {
    // ... existing cases ...
    default:
        return m.fallthrough.Execute(command, args)
    }
}
```

### 3. Process Execution Utility

**Location:** `internal/utils/process.go` (new)

```go
type ProcessExecutor struct{}

func NewProcessExecutor() *ProcessExecutor
func (p *ProcessExecutor) ExecuteCommand(command string, args []string) error
func (p *ProcessExecutor) ExecuteCommandWithOutput(command string, args []string) ([]byte, error)
```

**Responsibilities:**
- Handle subprocess execution
- Stream stdout/stderr in real-time
- Preserve exit codes
- Handle interactive commands

### 4. Configuration Enhancement

**Location:** `internal/config/config.go` (modified)

Add fallthrough configuration options:

```go
type Config struct {
    // ... existing fields ...
    Fallthrough FallthroughConfig `json:"fallthrough"`
}

type FallthroughConfig struct {
    Enabled     bool     `json:"enabled"`
    Verbose     bool     `json:"verbose"`
    Blacklist   []string `json:"blacklist"`
}
```

## Data Models

### Command Execution Context

```go
type ExecutionContext struct {
    Command     string
    Args        []string
    WorkingDir  string
    Environment map[string]string
}
```

### Fallthrough Result

```go
type FallthroughResult struct {
    Command    string
    Args       []string
    ExitCode   int
    Output     []byte
    Error      error
    Duration   time.Duration
}
```

## Error Handling

### Error Categories

1. **Git Not Found Error**
   - When `git` executable is not in PATH
   - Return clear error message with installation guidance

2. **Command Execution Error**
   - When the underlying Git command fails
   - Preserve original error message and exit code
   - No modification of Git's error output

3. **Argument Processing Error**
   - When arguments cannot be properly parsed or passed
   - Provide debugging information in verbose mode

4. **Blacklisted Command Error**
   - When a command is explicitly blacklisted from fallthrough
   - Suggest using native Git command directly

### Error Handling Strategy

```go
func (f *FallthroughHandler) Execute(command string, args []string) error {
    // Validate git executable exists
    if !f.isGitAvailable() {
        return fmt.Errorf("git executable not found in PATH")
    }
    
    // Check if command should fall through
    if !f.shouldFallthrough(command) {
        return fmt.Errorf("unknown command: %s", command)
    }
    
    // Execute with proper error propagation
    return f.executeGitCommand(command, args)
}
```

## Testing Strategy

### Unit Tests

1. **Fallthrough Handler Tests**
   - Test command validation logic
   - Test argument preservation
   - Test error handling scenarios
   - Mock Git executable for controlled testing

2. **Process Executor Tests**
   - Test subprocess execution
   - Test output streaming
   - Test exit code preservation
   - Test interactive command handling

3. **Integration Tests**
   - Test common Git commands through fallthrough
   - Test complex argument scenarios
   - Test error propagation
   - Test configuration options

### Test Structure

```
tests/
├── unit/
│   ├── fallthrough_handler_test.go
│   ├── process_executor_test.go
│   └── manager_integration_test.go
├── integration/
│   ├── git_commands_test.go
│   └── fallthrough_scenarios_test.go
└── fixtures/
    ├── mock_git_repo/
    └── test_configs/
```

### Test Cases

1. **Basic Fallthrough**
   - `git @ status` → `git status`
   - `git @ diff` → `git diff`
   - `git @ log --oneline` → `git log --oneline`

2. **Complex Arguments**
   - Commands with quotes and spaces
   - Commands with multiple flags
   - Commands with file paths

3. **Interactive Commands**
   - `git @ add -p`
   - `git @ rebase -i`
   - `git @ merge` (with conflicts)

4. **Error Scenarios**
   - Invalid Git commands
   - Git not installed
   - Repository not initialized

## Implementation Phases

### Phase 1: Core Fallthrough Mechanism
- Implement `FallthroughHandler`
- Implement `ProcessExecutor`
- Modify `Manager.Execute()` method
- Basic unit tests

### Phase 2: Configuration and Customization
- Add configuration options
- Implement verbose mode
- Add command blacklisting
- Configuration tests

### Phase 3: Enhanced Error Handling
- Improve error messages
- Add debugging capabilities
- Handle edge cases
- Comprehensive error tests

### Phase 4: Interactive Command Support
- Enhance process execution for interactive commands
- Test with complex Git workflows
- Performance optimization
- Integration tests

## Security Considerations

### Command Injection Prevention
- Validate command names against known Git commands
- Sanitize arguments to prevent shell injection
- Use `exec.Command()` with separate arguments (not shell execution)

### Blacklist Management
- Maintain list of commands that should not fall through
- Allow configuration of additional blacklisted commands
- Prevent execution of potentially dangerous commands

### Privilege Escalation
- Ensure fallthrough doesn't elevate privileges
- Maintain same security context as original Git commands
- Validate working directory and file access

## Performance Considerations

### Process Overhead
- Minimize subprocess creation overhead
- Stream output efficiently for large Git operations
- Handle memory usage for large Git outputs

### Caching Strategy
- Cache Git executable location
- Cache command validation results
- Optimize repeated command executions

## Backward Compatibility

### Existing Command Preservation
- All existing Gitat commands remain unchanged
- No modification to existing command behavior
- Fallthrough only activates for unknown commands

### Configuration Migration
- Default configuration enables fallthrough
- Existing configurations remain valid
- Graceful handling of missing configuration options

## Future Enhancements

### Smart Command Suggestions
- Suggest Gitat alternatives for common Git commands
- Provide hints about available Gitat features
- Educational mode for learning Gitat commands

### Command Analytics
- Track which commands fall through most frequently
- Identify candidates for native Gitat implementation
- Usage statistics for feature prioritization

### Advanced Integration
- Git alias support through fallthrough
- Custom Git command integration
- Plugin system for extending fallthrough behavior