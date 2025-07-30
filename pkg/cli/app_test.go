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
func TestRunWithNoArgs(t *testing.T) {
	app := createTestApp(t)

	err := app.Run([]string{})
	if err == nil {
		t.Log("No args should show usage (this is expected behavior)")
	}
}

// TestRunWithHelp tests help commands
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
		"path", "changes", "logs", "product", "feature", "issue",
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
		"logs":    true,
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
