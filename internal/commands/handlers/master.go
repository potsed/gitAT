package handlers

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// MainHandler handles main/master/root branch switching commands
type MainHandler struct {
	BaseHandler
}

// NewMainHandler creates a new main handler
func NewMainHandler(cfg *config.Config, gitRepo *git.Repository) *MainHandler {
	return &MainHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the main/master/root command
func (m *MainHandler) Execute(args []string) error {
	// Check for help flags
	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showUsage()
		}
	}

	// Get trunk branch configuration
	trunkBranch := m.config.Trunk
	if trunkBranch == "" {
		// Try to get from git config
		trunkBranch, _ = m.git.GetConfig("at.trunk")
	}
	
	// Default fallbacks
	if trunkBranch == "" {
		trunkBranch = "main"
		// Check if main exists
		_, err := m.git.Run("rev-parse", "--verify", "main")
		if err != nil {
			trunkBranch = "master"
			// Check if master exists
			_, err := m.git.Run("rev-parse", "--verify", "master")
			if err != nil {
				return fmt.Errorf("could not determine trunk branch - neither 'main' nor 'master' exist. Please configure with: git @ _trunk <branch>")
			}
		}
	}

	// Validate trunk branch exists
	_, err := m.git.Run("rev-parse", "--verify", trunkBranch)
	if err != nil {
		return fmt.Errorf("trunk branch '%s' does not exist. Please ensure the trunk branch exists or configure it with: git @ _trunk <branch>", trunkBranch)
	}

	// Get current branch
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// If already on trunk branch, just show info
	if currentBranch == trunkBranch {
		output.Info("Already on trunk branch: %s", trunkBranch)
		return nil
	}

	// Check for uncommitted changes
	hasUncommittedChanges := m.hasUncommittedChanges()
	if hasUncommittedChanges {
		output.Warning("You have uncommitted changes on branch: %s", currentBranch)
		
		var saveWIP bool
		err := huh.NewConfirm().
			Title("Save changes to WIP?").
			Description("Do you want to save your changes as WIP before switching branches?").
			Value(&saveWIP).
			WithTheme(huh.ThemeBase()).
			Run()
		
		if err != nil {
			return fmt.Errorf("failed to get user input: %w", err)
		}
		
		if saveWIP {
			// Save changes to WIP
			err = m.saveToWIP()
			if err != nil {
				return fmt.Errorf("failed to save changes to WIP: %w", err)
			}
			output.Success("Changes saved to WIP")
		} else {
			output.Info("Continuing without saving changes")
		}
	}

	// Switch to trunk branch
	output.Info("Switching to trunk branch: %s", trunkBranch)
	_, err = m.git.Run("checkout", trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to switch to trunk branch '%s': %w", trunkBranch, err)
	}

	// Update trunk branch
	output.Info("Updating trunk branch...")
	_, err = m.git.Run("remote", "get-url", "origin")
	if err == nil {
		_, err = m.git.Run("pull", "origin", trunkBranch)
		if err != nil {
			output.Warning("Failed to pull latest changes from remote, but continuing...")
		}
	}

	output.Success("Successfully switched to trunk branch: %s", trunkBranch)
	return nil
}

// hasUncommittedChanges checks if there are uncommitted changes
func (m *MainHandler) hasUncommittedChanges() bool {
	// Check for uncommitted staged changes
	_, err := m.git.Run("diff", "--cached", "--quiet")
	hasStaged := err != nil
	
	// Check for uncommitted unstaged changes
	_, err = m.git.Run("diff", "--quiet")
	hasUnstaged := err != nil
	
	return hasStaged || hasUnstaged
}

// saveToWIP saves current changes to WIP
func (m *MainHandler) saveToWIP() error {
	// Get current branch
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Stage all changes
	_, err = m.git.Run("add", "-A")
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Create WIP commit
	_, err = m.git.Run("commit", "-m", fmt.Sprintf("[WIP] %s", currentBranch))
	if err != nil {
		return fmt.Errorf("failed to create WIP commit: %w", err)
	}

	return nil
}

// showUsage displays the main/master/root command usage
func (m *MainHandler) showUsage() error {
	return output.Markdown(`# Main/Master/Root Command

Switch to the trunk branch (main/master).

## Usage

` + "```" + `bash
git @ main
git @ master
git @ root
` + "```" + `

## Description

The main command switches to the trunk branch (typically 'main' or 'master') and pulls the latest changes.
The master and root commands are aliases for the main command.

If you have uncommitted changes, the command will ask if you want to save them as WIP before switching.

## Behavior

1. Determines the trunk branch:
   - Checks GitAT configuration (at.trunk)
   - Falls back to 'main' if it exists
   - Falls back to 'master' if 'main' doesn't exist
   - Shows error if neither exists

2. Checks for uncommitted changes and offers to save to WIP

3. Switches to the trunk branch if not already on it

4. Pulls latest changes from remote (if origin exists)

## Configuration

To configure a specific trunk branch:

` + "```" + `bash
git @ _trunk develop
` + "```" + `

## Examples

` + "```" + `bash
# Switch to trunk branch
git @ main

# Switch to trunk branch (alternative)
git @ master
git @ root
` + "```" + `
`)
}