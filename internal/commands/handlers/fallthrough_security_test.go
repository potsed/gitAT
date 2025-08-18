package handlers

import (
	"strings"
	"testing"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// TestFallthroughSecurity_CommandInjection tests various command injection attack vectors
func TestFallthroughSecurity_CommandInjection(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	injectionTests := []struct {
		name        string
		command     string
		args        []string
		expectBlock bool
		description string
	}{
		{
			name:        "dollar_command_substitution",
			command:     "status",
			args:        []string{"$(whoami)"},
			expectBlock: true,
			description: "Dollar sign command substitution should be blocked",
		},
		{
			name:        "backtick_command_substitution",
			command:     "log",
			args:        []string{"`id`"},
			expectBlock: true,
			description: "Backtick command substitution should be blocked",
		},
		{
			name:        "nested_command_substitution",
			command:     "diff",
			args:        []string{"$(echo $(whoami))"},
			expectBlock: true,
			description: "Nested command substitution should be blocked",
		},
		{
			name:        "semicolon_command_chaining",
			command:     "status",
			args:        []string{"; rm -rf /tmp/test"},
			expectBlock: false, // This is safe as it's passed as an argument
			description: "Semicolon in arguments should be safe",
		},
		{
			name:        "pipe_redirection",
			command:     "log",
			args:        []string{"| cat /etc/passwd"},
			expectBlock: false, // This is safe as it's passed as an argument
			description: "Pipe in arguments should be safe",
		},
		{
			name:        "output_redirection",
			command:     "status",
			args:        []string{"> /etc/passwd"},
			expectBlock: false, // This is safe as it's passed as an argument
			description: "Output redirection in arguments should be safe",
		},
		{
			name:        "input_redirection",
			command:     "log",
			args:        []string{"< /etc/passwd"},
			expectBlock: false, // This is safe as it's passed as an argument
			description: "Input redirection in arguments should be safe",
		},
		{
			name:        "ampersand_background",
			command:     "status",
			args:        []string{"& sleep 10"},
			expectBlock: false, // This is safe as it's passed as an argument
			description: "Background execution in arguments should be safe",
		},
		{
			name:        "double_ampersand_and",
			command:     "log",
			args:        []string{"&& rm -rf /"},
			expectBlock: false, // This is safe as it's passed as an argument
			description: "Logical AND in arguments should be safe",
		},
		{
			name:        "double_pipe_or",
			command:     "status",
			args:        []string{"|| rm -rf /"},
			expectBlock: false, // This is safe as it's passed as an argument
			description: "Logical OR in arguments should be safe",
		},
	}

	for _, tt := range injectionTests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, tt.args)

			if tt.expectBlock {
				// We expect this to be blocked by argument validation
				if err == nil {
					t.Errorf("Expected %s to be blocked, but it wasn't", tt.description)
				} else if !strings.Contains(err.Error(), "potentially unsafe argument") {
					t.Errorf("Expected security validation error for %s, got: %v", tt.description, err)
				} else {
					t.Logf("Correctly blocked: %s", tt.description)
				}
			} else {
				// These should be safe (passed as arguments to Git)
				// We don't check for success/failure since Git might reject them for other reasons
				if err != nil && strings.Contains(err.Error(), "potentially unsafe argument") {
					t.Errorf("Safe pattern incorrectly blocked for %s: %v", tt.description, err)
				} else {
					t.Logf("Safely handled: %s", tt.description)
				}
			}
		})
	}
}

// TestFallthroughSecurity_PathTraversal tests path traversal attack prevention
func TestFallthroughSecurity_PathTraversal(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	pathTraversalTests := []struct {
		name        string
		command     string
		args        []string
		description string
	}{
		{
			name:        "relative_path_traversal",
			command:     "log",
			args:        []string{"../../../etc/passwd"},
			description: "Relative path traversal should be handled safely",
		},
		{
			name:        "absolute_path_access",
			command:     "show",
			args:        []string{"/etc/passwd"},
			description: "Absolute path access should be handled safely",
		},
		{
			name:        "encoded_path_traversal",
			command:     "diff",
			args:        []string{"..%2F..%2F..%2Fetc%2Fpasswd"},
			description: "URL-encoded path traversal should be handled safely",
		},
		{
			name:        "null_byte_injection",
			command:     "status",
			args:        []string{"file.txt\x00/etc/passwd"},
			description: "Null byte injection should be handled safely",
		},
	}

	for _, tt := range pathTraversalTests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, tt.args)
			// We don't expect these to be blocked by our validation since they're
			// legitimate file paths that Git should handle safely
			// We just verify they don't cause system-level issues
			t.Logf("Path traversal test '%s' result: %v", tt.description, err)
		})
	}
}

// TestFallthroughSecurity_ArgumentValidation tests comprehensive argument validation
func TestFallthroughSecurity_ArgumentValidation(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	validationTests := []struct {
		name        string
		args        []string
		expectError bool
		description string
	}{
		{
			name:        "safe_arguments",
			args:        []string{"--oneline", "-n", "10"},
			expectError: false,
			description: "Normal Git arguments should pass validation",
		},
		{
			name:        "arguments_with_spaces",
			args:        []string{"--grep=test message", "--author=John Doe"},
			expectError: false,
			description: "Arguments with spaces should pass validation",
		},
		{
			name:        "arguments_with_quotes",
			args:        []string{"--pretty=format:'%h %s'", "--since='2023-01-01'"},
			expectError: false,
			description: "Arguments with quotes should pass validation",
		},
		{
			name:        "command_substitution_dollar",
			args:        []string{"--message=$(whoami)"},
			expectError: true,
			description: "Command substitution with $ should be blocked",
		},
		{
			name:        "command_substitution_backtick",
			args:        []string{"--message=`id`"},
			expectError: true,
			description: "Command substitution with backticks should be blocked",
		},
		{
			name:        "nested_substitution",
			args:        []string{"--grep=$(echo $(whoami))"},
			expectError: true,
			description: "Nested command substitution should be blocked",
		},
		{
			name:        "safe_parentheses",
			args:        []string{"--grep=fix(auth): update validation"},
			expectError: false,
			description: "Safe parentheses usage should pass validation",
		},
		{
			name:        "safe_dollar_without_parentheses",
			args:        []string{"--grep=cost $100"},
			expectError: false,
			description: "Dollar sign without parentheses should be safe",
		},
		{
			name:        "empty_arguments",
			args:        []string{},
			expectError: false,
			description: "Empty argument list should pass validation",
		},
		{
			name:        "very_long_argument",
			args:        []string{strings.Repeat("a", 10000)},
			expectError: false,
			description: "Very long arguments should pass validation",
		},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.ValidateArguments(tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected validation error for %s, but got none", tt.description)
				} else {
					t.Logf("Correctly blocked unsafe arguments: %s - %v", tt.description, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected safe arguments to pass validation for %s, but got error: %v", tt.description, err)
				} else {
					t.Logf("Correctly allowed safe arguments: %s", tt.description)
				}
			}
		})
	}
}

// TestFallthroughSecurity_ReservedCommandProtection tests protection of reserved commands
func TestFallthroughSecurity_ReservedCommandProtection(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	reservedCommands := []string{
		"version", "help", "branch", "feature", "hotfix", "info",
		"issue", "label", "pr", "product", "release", "save",
		"squash", "sweep", "tag", "wip", "work", "changes",
		"logz", "shasum", "id", "path", "master", "root", "ignore",
	}

	for _, cmd := range reservedCommands {
		t.Run("reserved_"+cmd, func(t *testing.T) {
			err := handler.Execute(cmd, []string{})
			if err == nil {
				t.Errorf("Reserved command '%s' should not fall through to Git", cmd)
			} else if !strings.Contains(err.Error(), "reserved") {
				t.Errorf("Expected reserved command error for '%s', got: %v", cmd, err)
			} else {
				t.Logf("Correctly protected reserved command: %s", cmd)
			}
		})
	}
}

// TestFallthroughSecurity_BlacklistBypass tests that blacklist cannot be bypassed
func TestFallthroughSecurity_BlacklistBypass(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled:   true,
			Verbose:   false,
			Blacklist: []string{"push", "pull", "fetch"},
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	bypassTests := []struct {
		name        string
		command     string
		args        []string
		description string
	}{
		{
			name:        "direct_blacklisted_command",
			command:     "push",
			args:        []string{},
			description: "Direct blacklisted command should be blocked",
		},
		{
			name:        "blacklisted_with_args",
			command:     "pull",
			args:        []string{"origin", "main"},
			description: "Blacklisted command with arguments should be blocked",
		},
		{
			name:        "case_variation",
			command:     "PUSH",
			args:        []string{},
			description: "Case variation should not bypass blacklist (if implemented)",
		},
		{
			name:        "whitespace_variation",
			command:     " push ",
			args:        []string{},
			description: "Whitespace should not bypass blacklist (if trimmed)",
		},
	}

	for _, tt := range bypassTests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute(tt.command, tt.args)
			
			// For direct blacklisted commands, we expect them to be blocked
			if tt.command == "push" || tt.command == "pull" || tt.command == "fetch" {
				if err == nil {
					t.Errorf("Blacklisted command '%s' should be blocked", tt.command)
				} else if !strings.Contains(err.Error(), "blacklisted") {
					t.Errorf("Expected blacklist error for '%s', got: %v", tt.command, err)
				} else {
					t.Logf("Correctly blocked blacklisted command: %s", tt.description)
				}
			} else {
				// For variations, we expect them to either be blocked or treated as unknown
				t.Logf("Bypass test result for %s: %v", tt.description, err)
			}
		})
	}
}

// TestFallthroughSecurity_ConfigurationTampering tests protection against configuration tampering
func TestFallthroughSecurity_ConfigurationTampering(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled:   true,
			Verbose:   false,
			Blacklist: []string{"push", "pull"},
		},
	}

	t.Run("blacklist_validation", func(t *testing.T) {
		// Test that configuration validation catches issues
		originalBlacklist := cfg.Fallthrough.Blacklist

		// Test duplicate entries
		cfg.Fallthrough.Blacklist = []string{"push", "pull", "push"}
		err := cfg.ValidateFallthroughConfig()
		if err == nil {
			t.Error("Expected validation error for duplicate blacklist entries")
		}

		// Test empty entries
		cfg.Fallthrough.Blacklist = []string{"push", "", "pull"}
		err = cfg.ValidateFallthroughConfig()
		if err == nil {
			t.Error("Expected validation error for empty blacklist entries")
		}

		// Restore original
		cfg.Fallthrough.Blacklist = originalBlacklist
	})

	t.Run("blacklist_immutability", func(t *testing.T) {
		// Test that getting blacklist returns a copy
		blacklist1 := cfg.GetFallthroughBlacklist()
		blacklist2 := cfg.GetFallthroughBlacklist()

		// Modify one copy
		if len(blacklist1) > 0 {
			blacklist1[0] = "modified"
		}

		// Verify the other copy is unchanged
		if len(blacklist2) > 0 && blacklist2[0] == "modified" {
			t.Error("Blacklist copy was not properly isolated")
		}

		// Verify original config is unchanged
		if len(cfg.Fallthrough.Blacklist) > 0 && cfg.Fallthrough.Blacklist[0] == "modified" {
			t.Error("Original blacklist was modified through copy")
		}
	})
}

// TestFallthroughSecurity_ResourceExhaustion tests protection against resource exhaustion attacks
func TestFallthroughSecurity_ResourceExhaustion(t *testing.T) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	t.Run("large_argument_count", func(t *testing.T) {
		// Test with a large number of arguments
		largeArgs := make([]string, 10000)
		for i := range largeArgs {
			largeArgs[i] = "arg" + string(rune(i%26+97)) // a-z repeated
		}

		err := handler.ValidateArguments(largeArgs)
		if err != nil {
			t.Logf("Large argument validation result: %v", err)
		}
		// We don't expect this to fail, just verify it doesn't hang or crash
	})

	t.Run("very_long_arguments", func(t *testing.T) {
		// Test with very long individual arguments
		veryLongArg := strings.Repeat("a", 100000)
		args := []string{veryLongArg}

		err := handler.ValidateArguments(args)
		if err != nil {
			t.Logf("Long argument validation result: %v", err)
		}
		// We don't expect this to fail, just verify it doesn't hang or crash
	})

	t.Run("repeated_command_execution", func(t *testing.T) {
		// Test rapid repeated command execution
		for i := 0; i < 100; i++ {
			err := handler.ValidateArguments([]string{"--version"})
			if err != nil {
				t.Errorf("Validation failed on iteration %d: %v", i, err)
				break
			}
		}
	})
}