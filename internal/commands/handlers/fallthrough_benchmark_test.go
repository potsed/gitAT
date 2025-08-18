package handlers

import (
	"os/exec"
	"testing"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// BenchmarkFallthrough_BasicCommands benchmarks basic Git command execution
func BenchmarkFallthrough_BasicCommands(b *testing.B) {
	// Skip benchmark if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		b.Skip("Skipping benchmark: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false, // Disable verbose for benchmarking
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	b.Run("git_version", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.Execute("--version", []string{})
		}
	})

	b.Run("git_status", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.Execute("status", []string{"--porcelain"})
		}
	})

	b.Run("git_branch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.Execute("branch", []string{})
		}
	})
}

// BenchmarkFallthrough_ArgumentProcessing benchmarks argument processing overhead
func BenchmarkFallthrough_ArgumentProcessing(b *testing.B) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	b.Run("simple_args", func(b *testing.B) {
		args := []string{"--porcelain"}
		for i := 0; i < b.N; i++ {
			_ = handler.ValidateArguments(args)
			_ = handler.PreserveComplexArguments(args)
		}
	})

	b.Run("complex_args", func(b *testing.B) {
		args := []string{"--pretty=format:%h %s", "--grep=test message", "--author=John Doe"}
		for i := 0; i < b.N; i++ {
			_ = handler.ValidateArguments(args)
			_ = handler.PreserveComplexArguments(args)
		}
	})

	b.Run("many_args", func(b *testing.B) {
		args := []string{
			"--oneline", "--graph", "--decorate", "--all",
			"--since=2023-01-01", "--until=2023-12-31",
			"--author=test", "--grep=feature", "-n", "100",
		}
		for i := 0; i < b.N; i++ {
			_ = handler.ValidateArguments(args)
			_ = handler.PreserveComplexArguments(args)
		}
	})
}

// BenchmarkFallthrough_CommandValidation benchmarks command validation performance
func BenchmarkFallthrough_CommandValidation(b *testing.B) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	b.Run("should_fallthrough", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.shouldFallthrough("status")
			_ = handler.shouldFallthrough("log")
			_ = handler.shouldFallthrough("diff")
		}
	})

	b.Run("reserved_command_check", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.shouldFallthrough("version")
			_ = handler.shouldFallthrough("help")
			_ = handler.shouldFallthrough("branch")
		}
	})

	b.Run("blacklist_check", func(b *testing.B) {
		cfg.Fallthrough.Blacklist = []string{"push", "pull", "fetch", "merge", "rebase"}
		for i := 0; i < b.N; i++ {
			_ = cfg.IsFallthroughBlacklisted("push")
			_ = cfg.IsFallthroughBlacklisted("status")
			_ = cfg.IsFallthroughBlacklisted("log")
		}
	})
}

// BenchmarkFallthrough_ErrorHandling benchmarks error handling performance
func BenchmarkFallthrough_ErrorHandling(b *testing.B) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	b.Run("create_error_messages", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.createGitNotFoundError()
			_ = handler.createUnknownCommandError("unknown", []string{})
			_ = handler.createReservedCommandError("version", []string{})
		}
	})

	b.Run("suggest_similar_commands", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.SuggestSimilarCommands("st")
			_ = handler.SuggestSimilarCommands("co")
			_ = handler.SuggestSimilarCommands("feat")
		}
	})
}

// BenchmarkFallthrough_ConfigOperations benchmarks configuration operations
func BenchmarkFallthrough_ConfigOperations(b *testing.B) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled:   true,
			Verbose:   false,
			Blacklist: []string{"push", "pull", "fetch"},
		},
	}

	b.Run("blacklist_operations", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cfg.IsFallthroughBlacklisted("push")
			_ = cfg.AddToFallthroughBlacklist("merge")
			_ = cfg.RemoveFromFallthroughBlacklist("merge")
		}
	})

	b.Run("config_validation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cfg.ValidateFallthroughConfig()
		}
	})

	b.Run("blacklist_copy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cfg.GetFallthroughBlacklist()
		}
	})
}

// BenchmarkFallthrough_ProcessExecution benchmarks process execution overhead
func BenchmarkFallthrough_ProcessExecution(b *testing.B) {
	// Skip benchmark if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		b.Skip("Skipping benchmark: git not available in PATH")
	}

	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}
	handler := NewFallthroughHandler(cfg, gitRepo)

	b.Run("direct_git_execution", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.executeGitCommand("--version", []string{})
		}
	})

	b.Run("full_fallthrough_pipeline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = handler.Execute("--version", []string{})
		}
	})
}

// BenchmarkFallthrough_MemoryUsage benchmarks memory usage patterns
func BenchmarkFallthrough_MemoryUsage(b *testing.B) {
	cfg := &config.Config{
		Fallthrough: config.FallthroughConfig{
			Enabled: true,
			Verbose: false,
		},
	}
	gitRepo := &git.Repository{}

	b.Run("handler_creation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewFallthroughHandler(cfg, gitRepo)
		}
	})

	b.Run("test_mode_handler_creation", func(b *testing.B) {
		testResponses := map[string][]string{
			"git add -p": {"y", "n", "q"},
		}
		for i := 0; i < b.N; i++ {
			_ = NewFallthroughHandlerWithTestMode(cfg, gitRepo, testResponses)
		}
	})

	b.Run("large_argument_lists", func(b *testing.B) {
		handler := NewFallthroughHandler(cfg, gitRepo)
		largeArgs := make([]string, 100)
		for i := range largeArgs {
			largeArgs[i] = "arg" + string(rune(i))
		}

		for i := 0; i < b.N; i++ {
			_ = handler.ValidateArguments(largeArgs)
			_ = handler.PreserveComplexArguments(largeArgs)
		}
	})
}