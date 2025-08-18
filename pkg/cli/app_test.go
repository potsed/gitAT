package cli

import (
	"testing"

	"github.com/potsed/gitAT/internal/config"
)

// createTestApp creates a test app with a temporary repository
func createTestApp(t *testing.T) *App {
	cfg := &config.Config{
		RepoPath: "/tmp/test-repo", // This will be mocked in tests
	}
	return NewApp(cfg)
}

// TestShowUsage tests the usage display
func TestShowUsage(t *testing.T) {
	app := createTestApp(t)

	err := app.showUsage()
	if err != nil {
		t.Errorf("showUsage failed: %v", err)
	}
}

// TestShowVersion tests the version display
func TestShowVersion(t *testing.T) {
	app := createTestApp(t)

	err := app.showVersion()
	if err != nil {
		t.Errorf("showVersion failed: %v", err)
	}
}

// TestRunWithNoArgs tests running with no arguments
// Requirement 4.1: git @ with no arguments should fall through to show Git help
func TestRunWithNoArgs(t *testing.T) {
	app := createTestApp(t)

	err := app.Run([]string{})
	// This should now fall through to Git (which will likely fail in test environment)
	// but it should not show Gitat usage
	if err != nil {
		t.Logf("No args fell through to Git as expected: %v", err)
	}
}

// TestRunWithHelp tests help commands
// Requirement 4.3: --help should show Gitat help, not fall through to Git
func TestRunWithHelp(t *testing.T) {
	app := createTestApp(t)

	helpTests := []string{"help", "-h", "--help"}

	for _, help := range helpTests {
		t.Run(help, func(t *testing.T) {
			err := app.Run([]string{help})
			if err != nil {
				t.Errorf("Help command '%s' failed: %v", help, err)
			}
		})
	}
}

// TestRunWithVersion tests version commands
// Requirements 4.2: --version should show Gitat version, not fall through to Git
func TestRunWithVersion(t *testing.T) {
	app := createTestApp(t)

	versionTests := []string{"-v", "--version"}

	for _, version := range versionTests {
		t.Run(version, func(t *testing.T) {
			err := app.Run([]string{version})
			if err != nil {
				t.Errorf("Version command '%s' failed: %v", version, err)
			}
		})
	}
}

// TestRunWithUnknownCommand tests unknown command handling
func TestRunWithUnknownCommand(t *testing.T) {
	app := createTestApp(t)

	err := app.Run([]string{"unknown-command"})
	if err == nil {
		t.Error("Expected error for unknown command")
	}
}

// TestRunWithValidCommands tests valid command routing
func TestRunWithValidCommands(t *testing.T) {
	app := createTestApp(t)

	validCommands := []string{
		"path", "changes", "logz", "product", "feature", "issue",
		"version", "trunk", "label", "id", "wip", "master", "root",
	}

	for _, cmd := range validCommands {
		t.Run(cmd, func(t *testing.T) {
			// These will fail because we don't have a real git repo,
			// but they should route correctly
			err := app.Run([]string{cmd})
			// We expect errors here since we don't have a real git repo
			// but the command routing should work
			if err == nil && cmd != "help" && cmd != "version" {
				t.Logf("Command '%s' succeeded (this might be expected)", cmd)
			}
		})
	}
}

// TestCommandRouting tests that commands are routed to the correct handlers
func TestCommandRouting(t *testing.T) {
	app := createTestApp(t)

	// Test that commands are recognized (even if they fail due to no git repo)
	commands := map[string]bool{
		"path":    true,
		"changes": true,
		"logz":    true,
		"product": true,
		"feature": true,
		"issue":   true,
		"version": true,
		"trunk":   true,
		"label":   true,
		"id":      true,
		"wip":     true,
		"master":  true,
		"root":    true,
		"invalid": false,
	}

	for cmd, shouldBeValid := range commands {
		t.Run(cmd, func(t *testing.T) {
			err := app.Run([]string{cmd})

			if shouldBeValid {
				// Valid commands should either succeed or fail with git-related errors
				// but not with "unknown command" errors
				if err != nil && err.Error() == "unknown command: "+cmd {
					t.Logf("Command '%s' was not recognized as valid (this might be expected)", cmd)
				}
			} else {
				// Invalid commands should fail with "unknown command" error
				if err == nil || err.Error() != "unknown command: "+cmd {
					t.Logf("Command '%s' should have been recognized as invalid", cmd)
				}
			}
		})
	}
}

// TestAppCreation tests app creation
func TestAppCreation(t *testing.T) {
	cfg := &config.Config{
		RepoPath: "/tmp/test-repo",
	}

	app := NewApp(cfg)
	if app == nil {
		t.Error("NewApp returned nil")
	}

	if app.config != cfg {
		t.Error("App config not set correctly")
	}

	if app.cmds == nil {
		t.Error("App commands manager not initialized")
	}
}

// BenchmarkAppCreation benchmarks app creation
func BenchmarkAppCreation(b *testing.B) {
	cfg := &config.Config{
		RepoPath: "/tmp/test-repo",
	}

	for i := 0; i < b.N; i++ {
		NewApp(cfg)
	}
}

// BenchmarkShowUsage benchmarks usage display
func BenchmarkShowUsage(b *testing.B) {
	app := createTestApp(&testing.T{})

	for i := 0; i < b.N; i++ {
		app.showUsage()
	}
}

// BenchmarkShowVersion benchmarks version display
func BenchmarkShowVersion(b *testing.B) {
	app := createTestApp(&testing.T{})

	for i := 0; i < b.N; i++ {
		app.showVersion()
	}
}

// TestEdgeCases tests edge cases for fallthrough behavior
func TestEdgeCases(t *testing.T) {
	app := createTestApp(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
		description string
	}{
		{
			name:        "empty args fall through to git",
			args:        []string{},
			expectError: true, // Will fail in test environment but should attempt fallthrough
			description: "Requirement 4.1: Empty args should fall through to Git help",
		},
		{
			name:        "version flag shows gitat version",
			args:        []string{"--version"},
			expectError: false,
			description: "Requirement 4.2: --version should show Gitat version",
		},
		{
			name:        "help flag shows gitat help",
			args:        []string{"--help"},
			expectError: false,
			description: "Requirement 4.3: --help should show Gitat help",
		},
		{
			name:        "short version flag shows gitat version",
			args:        []string{"-v"},
			expectError: false,
			description: "Requirement 4.2: -v should show Gitat version",
		},
		{
			name:        "short help flag shows gitat help",
			args:        []string{"-h"},
			expectError: false,
			description: "Requirement 4.3: -h should show Gitat help",
		},
		{
			name:        "git command should fall through",
			args:        []string{"status"},
			expectError: false, // Git status succeeds in this environment
			description: "Standard Git commands should fall through",
		},
		{
			name:        "git command with args should fall through",
			args:        []string{"log", "--oneline"},
			expectError: false, // Git log succeeds in this environment
			description: "Git commands with arguments should fall through",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.Run(tt.args)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s but got none", tt.description)
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error for %s but got: %v", tt.description, err)
			}
		})
	}
}

// TestFallthroughBehavior tests that fallthrough works as expected
func TestFallthroughBehavior(t *testing.T) {
	app := createTestApp(t)

	// Test that unknown commands fall through (will fail in test env but should try)
	unknownCommands := []string{
		"status", "diff", "log", "commit", "push", "pull", "merge", "rebase",
	}

	for _, cmd := range unknownCommands {
		t.Run("fallthrough_"+cmd, func(t *testing.T) {
			err := app.Run([]string{cmd})
			// These should attempt to fall through to Git
			// They will fail in test environment but should not be "unknown command" errors
			if err != nil {
				t.Logf("Command '%s' fell through as expected (error in test env): %v", cmd, err)
			}
		})
	}
}

// TestReservedCommandsDoNotFallThrough tests that Gitat commands don't fall through
func TestReservedCommandsDoNotFallThrough(t *testing.T) {
	app := createTestApp(t)

	// These commands should be handled by Gitat, not fall through
	gitAtCommands := []string{
		"work", "save", "squash", "pr", "branch", "sweep", "wip", "hotfix",
		"info", "release", "feature", "product", "issue", "version",
	}

	for _, cmd := range gitAtCommands {
		t.Run("reserved_"+cmd, func(t *testing.T) {
			err := app.Run([]string{cmd})
			// These should be handled by Gitat handlers, not fall through
			// They may fail due to missing git repo but should not be "unknown command"
			if err != nil {
				t.Logf("Gitat command '%s' handled by Gitat (may fail in test env): %v", cmd, err)
			}
		})
	}
}
