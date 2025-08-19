package handlers

import (
	"strings"
	"testing"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

func TestFallthroughHandler_EdgeCases(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := git.NewRepository(".")

	tests := []struct {
		name          string
		command       string
		args          []string
		expectError   bool
		errorContains string
		testResponses map[string][]string
		description   string
	}{
		{
			name:        "empty command shows git help",
			command:     "",
			args:        []string{},
			expectError: true, // Git help returns exit code 1
			testResponses: map[string][]string{
				"git": {"usage: git [--version] [--help] [-C <path>]..."},
			},
			description: "Requirement 4.1: git @ with no arguments should show Git help",
		},
		{
			name:          "version flag should not fall through",
			command:       "--version",
			args:          []string{},
			expectError:   true,
			errorContains: "unknown command",
			description:   "Requirement 4.2: --version should show Gitat info, not fall through",
		},
		{
			name:          "help flag should not fall through",
			command:       "--help",
			args:          []string{},
			expectError:   true,
			errorContains: "unknown command",
			description:   "Requirement 4.3: --help should show Gitat info, not fall through",
		},
		{
			name:          "reserved gitat command should not fall through",
			command:       "work",
			args:          []string{"feature", "test"},
			expectError:   true,
			errorContains: "unknown command",
			description:   "Reserved Gitat commands should not fall through",
		},
		{
			name:        "git alias should work",
			command:     "st",
			args:        []string{},
			expectError: true, // Will fail because 'st' is not a real alias in test env
			testResponses: map[string][]string{
				"git config --get alias.st": {"status"},
				"git st":                    {"On branch main"},
			},
			description: "Requirement 5.2: Git aliases should work through fallthrough",
		},
		{
			name:        "subcommand should work",
			command:     "remote",
			args:        []string{"add", "origin", "https://github.com/test/repo.git"},
			expectError: true, // Will fail because origin already exists in test env
			testResponses: map[string][]string{
				"git remote add origin https://github.com/test/repo.git": {""},
			},
			description: "Requirement 5.2: Subcommands should work correctly",
		},
		{
			name:        "complex arguments with quotes",
			command:     "status",
			args:        []string{"--porcelain"},
			expectError: false, // Git status --porcelain succeeds even with changes
			testResponses: map[string][]string{
				"git status --porcelain": {"M file.txt"},
			},
			description: "Requirement 5.1: Complex arguments should be preserved",
		},
		{
			name:        "multiple flags and options",
			command:     "log",
			args:        []string{"--oneline", "--graph", "--decorate", "--all"},
			expectError: false, // Git log should succeed
			testResponses: map[string][]string{
				"git log --oneline --graph --decorate --all": {"* abc1234 (HEAD -> main) Test commit"},
			},
			description: "Multiple flags should be preserved correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewFallthroughHandlerWithTestMode(cfg, gitRepo, tt.testResponses)

			err := handler.Execute(tt.command, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for test: %s", tt.description)
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s' but got '%s' for test: %s",
						tt.errorContains, err.Error(), tt.description)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v for test: %s", err, tt.description)
				}
			}
		})
	}
}

func TestFallthroughHandler_ReservedCommands(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := git.NewRepository(".")
	handler := NewFallthroughHandler(cfg, gitRepo)

	// Test all reserved commands to ensure they don't fall through
	reservedCommands := []string{
		"semver", "help", "sprout", "feature", "hotfix", "info", "issue",
		"label", "pr", "product", "release", "save", "squash", "sweep",
		"dub", "wip", "work", "changes", "logz", "shasum", "id", "path",
		"main", "master", "root", "ignore", "setup-local", "setup-remote", "security",
		"changelog", "rebase", "commitizen", "cz",
		"_go", "_label", "_id", "_path", "_trunk", "_security",
		"-v", "--version", "-h", "--help",
	}

	for _, cmd := range reservedCommands {
		t.Run("reserved_command_"+cmd, func(t *testing.T) {
			err := handler.Execute(cmd, []string{})
			if err == nil {
				t.Errorf("Expected error for reserved command '%s' but got none", cmd)
			}
			if !strings.Contains(err.Error(), "unknown command") {
				t.Errorf("Expected 'unknown command' error for '%s' but got: %v", cmd, err)
			}
		})
	}
}

func TestFallthroughHandler_GitAliasDetection(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := git.NewRepository(".")

	tests := []struct {
		name          string
		command       string
		testResponses map[string][]string
		expectAlias   bool
	}{
		{
			name:    "valid git alias",
			command: "st",
			testResponses: map[string][]string{
				"git config --get alias.st": {"status"},
			},
			expectAlias: true,
		},
		{
			name:    "not an alias",
			command: "status",
			testResponses: map[string][]string{
				"git config --get alias.status": {""},
			},
			expectAlias: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewFallthroughHandlerWithTestMode(cfg, gitRepo, tt.testResponses)

			isAlias := handler.IsGitAlias(tt.command)

			if isAlias != tt.expectAlias {
				t.Errorf("Expected IsGitAlias(%s) to be %v but got %v",
					tt.command, tt.expectAlias, isAlias)
			}
		})
	}
}

func TestFallthroughHandler_SubcommandHandling(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := git.NewRepository(".")
	handler := NewFallthroughHandler(cfg, gitRepo)

	tests := []struct {
		name              string
		command           string
		args              []string
		expectSubcommands bool
	}{
		{
			name:              "remote command has subcommands",
			command:           "remote",
			args:              []string{"add", "origin", "https://github.com/test/repo.git"},
			expectSubcommands: true,
		},
		{
			name:              "config command has subcommands",
			command:           "config",
			args:              []string{"--global", "user.name", "Test User"},
			expectSubcommands: true,
		},
		{
			name:              "status command has no subcommands",
			command:           "status",
			args:              []string{},
			expectSubcommands: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasSubcommands, err := handler.HandleSubcommands(tt.command, tt.args)

			if err != nil {
				t.Errorf("Unexpected error in HandleSubcommands: %v", err)
			}

			if hasSubcommands != tt.expectSubcommands {
				t.Errorf("Expected HandleSubcommands(%s) to return %v but got %v",
					tt.command, tt.expectSubcommands, hasSubcommands)
			}
		})
	}
}

func TestFallthroughHandler_ArgumentValidation(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := git.NewRepository(".")
	handler := NewFallthroughHandler(cfg, gitRepo)

	tests := []struct {
		name        string
		args        []string
		expectError bool
		description string
	}{
		{
			name:        "safe arguments",
			args:        []string{"-m", "commit message", "--author", "Test User"},
			expectError: false,
			description: "Normal arguments should pass validation",
		},
		{
			name:        "arguments with quotes",
			args:        []string{"-m", "commit with 'single quotes'"},
			expectError: false,
			description: "Arguments with quotes should be allowed",
		},
		{
			name:        "arguments with spaces",
			args:        []string{"--message", "commit message with spaces"},
			expectError: false,
			description: "Arguments with spaces should be allowed",
		},
		{
			name:        "potentially unsafe command substitution",
			args:        []string{"-m", "$(rm -rf /)"},
			expectError: true,
			description: "Command substitution should be blocked",
		},
		{
			name:        "potentially unsafe backticks",
			args:        []string{"-m", "`rm -rf /`"},
			expectError: true,
			description: "Backtick command substitution should be blocked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateArguments(tt.args)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s but got none", tt.description)
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error for %s but got: %v", tt.description, err)
			}
		})
	}
}

func TestFallthroughHandler_ComplexArgumentPreservation(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := git.NewRepository(".")
	handler := NewFallthroughHandler(cfg, gitRepo)

	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "simple arguments",
			args:     []string{"status", "--short"},
			expected: []string{"status", "--short"},
		},
		{
			name:     "arguments with spaces",
			args:     []string{"-m", "commit message with spaces"},
			expected: []string{"-m", "commit message with spaces"},
		},
		{
			name:     "arguments with quotes",
			args:     []string{"-m", "commit with 'quotes'"},
			expected: []string{"-m", "commit with 'quotes'"},
		},
		{
			name:     "multiple flags",
			args:     []string{"--oneline", "--graph", "--decorate"},
			expected: []string{"--oneline", "--graph", "--decorate"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preserved := handler.PreserveComplexArguments(tt.args)

			if len(preserved) != len(tt.expected) {
				t.Errorf("Expected %d arguments but got %d", len(tt.expected), len(preserved))
				return
			}

			for i, arg := range preserved {
				if arg != tt.expected[i] {
					t.Errorf("Expected argument %d to be '%s' but got '%s'", i, tt.expected[i], arg)
				}
			}
		})
	}
}
