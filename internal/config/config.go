package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FallthroughConfig holds configuration for Git fallthrough functionality
type FallthroughConfig struct {
	Enabled   bool     `json:"enabled"`
	Verbose   bool     `json:"verbose"`
	Blacklist []string `json:"blacklist"`
}

// Config holds the GitAT configuration
type Config struct {
	// Git repository configuration
	RepoPath string
	Trunk    string
	Product  string
	Feature  string
	Task     string
	Branch   string
	Version  string
	WIP      string

	// Application settings
	Verbose     bool
	DryRun      bool
	Fallthrough FallthroughConfig `json:"fallthrough"`
}

// Load loads the GitAT configuration from Git config
func Load() (*Config, error) {
	cfg := &Config{}

	// Get repository path
	repoPath, err := getGitRepoPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository path: %w", err)
	}
	cfg.RepoPath = repoPath

	// Load GitAT configuration from Git config
	cfg.Trunk = getGitConfig("at.trunk")
	cfg.Product = getGitConfig("at.product")
	cfg.Feature = getGitConfig("at.feature")
	cfg.Task = getGitConfig("at.task")
	cfg.Branch = getGitConfig("at.branch")
	cfg.Version = getGitConfig("at.version")
	cfg.WIP = getGitConfig("at.wip")

	// Load fallthrough configuration with defaults
	cfg.Fallthrough = loadFallthroughConfig()

	return cfg, nil
}

// getGitRepoPath returns the path to the Git repository root
func getGitRepoPath() (string, error) {
	// Start from current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree to find .git
	for {
		gitDir := filepath.Join(currentDir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return currentDir, nil
		}

		// Move up one directory
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached root directory
			return "", fmt.Errorf("not in a Git repository")
		}
		currentDir = parent
	}
}

// getGitConfig retrieves a Git configuration value
func getGitConfig(key string) string {
	// This is a placeholder - we'll implement actual Git config reading
	// For now, return empty string
	return ""
}

// Save saves the GitAT configuration to Git config
func (c *Config) Save() error {
	// This is a placeholder - we'll implement actual Git config writing
	return nil
}

// SetProduct sets the product name
func (c *Config) SetProduct(name string) error {
	c.Product = name
	return c.Save()
}

// SetTrunk sets the trunk branch
func (c *Config) SetTrunk(branch string) error {
	c.Trunk = branch
	return c.Save()
}

// SetFeature sets the current feature
func (c *Config) SetFeature(feature string) error {
	c.Feature = feature
	return c.Save()
}

// SetTask sets the current task/issue ID
func (c *Config) SetTask(task string) error {
	c.Task = task
	return c.Save()
}

// SetBranch sets the working branch
func (c *Config) SetBranch(branch string) error {
	c.Branch = branch
	return c.Save()
}

// SetVersion sets the current version
func (c *Config) SetVersion(version string) error {
	c.Version = version
	return c.Save()
}

// SetWIP sets the WIP branch
func (c *Config) SetWIP(branch string) error {
	c.WIP = branch
	return c.Save()
}

// loadFallthroughConfig loads fallthrough configuration with default values
func loadFallthroughConfig() FallthroughConfig {
	cfg := FallthroughConfig{
		Enabled: true,  // Enable fallthrough by default
		Verbose: false, // Disable verbose mode by default
		Blacklist: []string{
			"--version", // Prevent fallthrough for version command
			"--help",    // Prevent fallthrough for help command
		},
	}

	// Load from Git config if available
	if enabled := getGitConfig("at.fallthrough.enabled"); enabled != "" {
		cfg.Enabled = enabled == "true"
	}

	if verbose := getGitConfig("at.fallthrough.verbose"); verbose != "" {
		cfg.Verbose = verbose == "true"
	}

	// TODO: Load blacklist from Git config (comma-separated values)
	// For now, use defaults

	return cfg
}

// SetFallthroughEnabled enables or disables fallthrough functionality
func (c *Config) SetFallthroughEnabled(enabled bool) error {
	c.Fallthrough.Enabled = enabled
	return c.Save()
}

// SetFallthroughVerbose enables or disables verbose mode for fallthrough
func (c *Config) SetFallthroughVerbose(verbose bool) error {
	c.Fallthrough.Verbose = verbose
	return c.Save()
}

// AddToFallthroughBlacklist adds a command to the fallthrough blacklist
func (c *Config) AddToFallthroughBlacklist(command string) error {
	// Check if command is already in blacklist
	for _, cmd := range c.Fallthrough.Blacklist {
		if cmd == command {
			return nil // Already in blacklist
		}
	}
	
	c.Fallthrough.Blacklist = append(c.Fallthrough.Blacklist, command)
	return c.Save()
}

// RemoveFromFallthroughBlacklist removes a command from the fallthrough blacklist
func (c *Config) RemoveFromFallthroughBlacklist(command string) error {
	for i, cmd := range c.Fallthrough.Blacklist {
		if cmd == command {
			// Remove command from slice
			c.Fallthrough.Blacklist = append(c.Fallthrough.Blacklist[:i], c.Fallthrough.Blacklist[i+1:]...)
			return c.Save()
		}
	}
	return nil // Command not found in blacklist
}

// IsFallthroughBlacklisted checks if a command is blacklisted from fallthrough
func (c *Config) IsFallthroughBlacklisted(command string) bool {
	for _, cmd := range c.Fallthrough.Blacklist {
		if cmd == command {
			return true
		}
	}
	return false
}

// GetFallthroughBlacklist returns a copy of the current blacklist
func (c *Config) GetFallthroughBlacklist() []string {
	blacklist := make([]string, len(c.Fallthrough.Blacklist))
	copy(blacklist, c.Fallthrough.Blacklist)
	return blacklist
}

// ClearFallthroughBlacklist removes all commands from the blacklist
func (c *Config) ClearFallthroughBlacklist() error {
	c.Fallthrough.Blacklist = []string{}
	return c.Save()
}

// SetFallthroughBlacklist replaces the entire blacklist with the provided commands
func (c *Config) SetFallthroughBlacklist(commands []string) error {
	c.Fallthrough.Blacklist = make([]string, len(commands))
	copy(c.Fallthrough.Blacklist, commands)
	return c.Save()
}

// ValidateFallthroughConfig validates the fallthrough configuration
func (c *Config) ValidateFallthroughConfig() error {
	// Check for duplicate entries in blacklist
	seen := make(map[string]bool)
	for _, cmd := range c.Fallthrough.Blacklist {
		if seen[cmd] {
			return fmt.Errorf("duplicate command in fallthrough blacklist: %s", cmd)
		}
		seen[cmd] = true
	}
	
	// Validate that blacklisted commands are reasonable
	for _, cmd := range c.Fallthrough.Blacklist {
		if strings.TrimSpace(cmd) == "" {
			return fmt.Errorf("empty command in fallthrough blacklist")
		}
		
		// Warn about potentially problematic blacklist entries
		if cmd == "help" || cmd == "--help" {
			// This is fine, but might be confusing
		}
	}
	
	return nil
} 