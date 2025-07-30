package config

import (
	"fmt"
	"os"
	"path/filepath"
)

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
	Verbose bool
	DryRun  bool
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