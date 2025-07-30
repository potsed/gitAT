# GitAT Go Project Structure

This document describes the Go project structure for GitAT, following modern Go best practices.

## Project Layout

```
gitAT/
├── cmd/
│   └── gitat/                 # Main application entry point
│       └── main.go           # Application entry point
├── internal/                  # Private application code
│   ├── commands/             # Command implementations
│   │   └── manager.go        # Command manager
│   ├── config/               # Configuration management
│   │   ├── config.go         # Configuration struct and methods
│   │   └── config_test.go    # Configuration tests
│   ├── git/                  # Git operations
│   │   └── repository.go     # Git repository wrapper
│   └── utils/                # Internal utilities
├── pkg/                       # Public libraries
│   ├── cli/                  # CLI framework
│   │   └── app.go            # CLI application
│   └── workflow/             # Workflow management
├── docs/                      # Documentation
│   └── GO_STRUCTURE.md       # This file
├── scripts/                   # Build and deployment scripts
├── tests/                     # Test files (from shell version)
├── build/                     # Build artifacts (generated)
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── Makefile                   # Build and development tasks
├── .gitignore                 # Git ignore patterns
└── README.md                  # Project documentation
```

## Package Organization

### `cmd/gitat/`

Contains the main application entry point. This is where the `main()` function lives and where the application starts.

### `internal/`

Contains private application code that should not be imported by other projects.

- **`commands/`**: All command implementations (work, save, squash, etc.)
- **`config/`**: Configuration management and Git config integration
- **`git/`**: Git operations wrapper and repository management
- **`utils/`**: Internal utility functions

### `pkg/`

Contains public libraries that could be imported by other projects.

- **`cli/`**: CLI framework and command parsing
- **`workflow/`**: Workflow management and business logic

## Key Design Principles

### 1. **Separation of Concerns**

- Configuration management is separate from Git operations
- CLI parsing is separate from command implementation
- Each command has its own implementation

### 2. **Dependency Injection**

- Configuration is injected into command managers
- Git repository is injected into commands that need it
- Easy to test and mock dependencies

### 3. **Error Handling**

- All functions return errors where appropriate
- Errors are wrapped with context using `fmt.Errorf` and `%w`
- Graceful error handling throughout the application

### 4. **Testing**

- Each package has corresponding test files
- Tests are organized alongside the code they test
- Mock interfaces for external dependencies

## Building and Development

### Prerequisites

- Go 1.24+ (latest stable)
- Git (for version information)
- Make (for build automation)

### Build Commands

```bash
# Build for current platform
make build-local

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Install dependencies
make deps

# Development setup
make dev-setup
```

### Development Workflow

1. **Setup**: `make dev-setup`
2. **Build**: `make build-local`
3. **Test**: `make test`
4. **Format**: `make fmt`
5. **Lint**: `make lint`

## Migration from Shell Scripts

The Go version maintains the same command interface as the shell script version:

```bash
# Shell version
git @ work feature add-auth
git @ save "Add authentication"

# Go version (same interface)
git @ work feature add-auth
git @ save "Add authentication"
```

### Key Differences

1. **Performance**: Compiled binary is faster than interpreted scripts
2. **Distribution**: Single binary vs. multiple shell scripts
3. **Cross-platform**: Native binaries for each platform
4. **Dependencies**: No external shell dependencies
5. **Testing**: Native Go testing framework

## Configuration

The Go version uses the same Git configuration keys as the shell version:

- `at.product` - Product name
- `at.feature` - Current feature
- `at.task` - Current task/issue ID
- `at.branch` - Working branch
- `at.trunk` - Trunk branch
- `at.version` - Current version
- `at.wip` - Work in progress branch

## Testing Strategy

### Unit Tests

- Each package has corresponding test files
- Mock external dependencies (Git, filesystem)
- Test configuration loading and validation
- Test command parsing and validation

### Integration Tests

- Test with real Git repositories
- Test command workflows end-to-end
- Test error conditions and edge cases

### Test Organization

```
internal/
├── config/
│   ├── config.go
│   └── config_test.go
├── git/
│   ├── repository.go
│   └── repository_test.go
└── commands/
    ├── manager.go
    └── manager_test.go
```

## Future Enhancements

1. **Plugin System**: Allow custom commands via plugins
2. **Configuration UI**: Interactive configuration setup
3. **Remote Integration**: Direct integration with GitHub/GitLab APIs
4. **Workflow Templates**: Predefined workflow templates
5. **Performance Monitoring**: Built-in performance metrics
6. **Logging**: Structured logging for debugging

## Contributing

When adding new features:

1. **Add tests first** (TDD approach)
2. **Follow Go conventions** (gofmt, golint)
3. **Update documentation** (README, this file)
4. **Add to Makefile** if needed
5. **Update version** in main.go

## References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
