// rebase.go - Rebase command handler
package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// RebaseHandler handles rebase-related commands
type RebaseHandler struct {
	BaseHandler
}

// NewRebaseHandler creates a new rebase handler
func NewRebaseHandler(cfg *config.Config, gitRepo *git.Repository) *RebaseHandler {
	return &RebaseHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the rebase command
func (r *RebaseHandler) Execute(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "interactive", "i":
			return r.interactiveRebase(args[1:])
		case "onto":
			if len(args) < 2 {
				return fmt.Errorf("target branch required for rebase onto")
			}
			return r.rebaseOnto(args[1], args[2:])
		case "abort":
			return r.abortRebase()
		case "continue":
			return r.continueRebase()
		case "skip":
			return r.skipRebase()
		case "-h", "--help", "help":
			return r.showUsage()
		default:
			// Default to rebase onto specified branch
			return r.rebaseOnto(args[0], args[1:])
		}
	}
	
	// Default action: rebase onto trunk
	return r.rebaseOntoTrunk()
}

// rebaseOntoTrunk rebases current branch onto trunk branch
func (r *RebaseHandler) rebaseOntoTrunk() error {
	// Get trunk branch
	trunkBranch := r.getTrunkBranch()

	// Get current branch
	currentBranch, err := r.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Check if already on trunk
	if currentBranch == trunkBranch {
		output.Info("Already on trunk branch: %s", trunkBranch)
		return nil
	}

	// Show what will happen
	output.Title("Rebasing onto Trunk")
	output.Info("Current branch: %s", currentBranch)
	output.Info("Target branch: %s", trunkBranch)

	// Get user confirmation
	var proceed bool
	err = huh.NewConfirm().
		Title("Proceed with Rebase?").
		Description(fmt.Sprintf("Rebase %s onto %s?", currentBranch, trunkBranch)).
		Value(&proceed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %w", err)
	}

	if !proceed {
		output.Info("Rebase cancelled")
		return nil
	}

	// Update trunk branch first
	output.Info("Updating %s branch...", trunkBranch)
	_, err = r.git.Run("remote", "get-url", "origin")
	if err == nil {
		_, err = r.git.Run("pull", "origin", trunkBranch)
		if err != nil {
			output.Warning("Failed to pull latest changes from remote, but continuing...")
		}
	}

	// Perform rebase
	output.Info("Rebasing %s onto %s...", currentBranch, trunkBranch)
	_, err = r.git.Run("rebase", trunkBranch)
	if err != nil {
		return fmt.Errorf("rebase failed: %w", err)
	}

	output.Success("Successfully rebased %s onto %s", currentBranch, trunkBranch)
	return nil
}

// rebaseOnto rebases current branch onto specified branch
func (r *RebaseHandler) rebaseOnto(targetBranch string, options []string) error {
	// Validate target branch exists
	_, err := r.git.Run("rev-parse", "--verify", targetBranch)
	if err != nil {
		return fmt.Errorf("target branch '%s' does not exist", targetBranch)
	}

	// Get current branch
	currentBranch, err := r.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Check if already on target
	if currentBranch == targetBranch {
		output.Info("Already on target branch: %s", targetBranch)
		return nil
	}

	// Show what will happen
	output.Title("Rebasing onto Branch")
	output.Info("Current branch: %s", currentBranch)
	output.Info("Target branch: %s", targetBranch)

	// Build rebase command
	args := []string{"rebase", targetBranch}
	args = append(args, options...)

	// Get user confirmation for non-interactive rebase
	if !r.containsFlag(options, "-i", "--interactive") {
		var proceed bool
		err = huh.NewConfirm().
			Title("Proceed with Rebase?").
			Description(fmt.Sprintf("Rebase %s onto %s?", currentBranch, targetBranch)).
			Value(&proceed).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get user input: %w", err)
		}

		if !proceed {
			output.Info("Rebase cancelled")
			return nil
		}
	}

	// Perform rebase
	output.Info("Executing: git %s", strings.Join(args, " "))
	_, err = r.git.Run(args...)
	if err != nil {
		return fmt.Errorf("rebase failed: %w", err)
	}

	output.Success("Successfully rebased %s onto %s", currentBranch, targetBranch)
	return nil
}

// interactiveRebase performs interactive rebase
func (r *RebaseHandler) interactiveRebase(args []string) error {
	var commits string
	
	if len(args) == 0 {
		// Default to last 5 commits
		commits = "HEAD~5"
	} else {
		commits = args[0]
	}

	// Get current branch
	currentBranch, err := r.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	output.Title("Interactive Rebase")
	output.Info("Current branch: %s", currentBranch)
	output.Info("Commits to rebase: %s", commits)

	// Perform interactive rebase
	output.Info("Starting interactive rebase...")
	_, err = r.git.Run("rebase", "-i", commits)
	if err != nil {
		return fmt.Errorf("interactive rebase failed: %w", err)
	}

	output.Success("Interactive rebase completed")
	return nil
}

// abortRebase aborts current rebase
func (r *RebaseHandler) abortRebase() error {
	output.Info("Aborting current rebase...")
	_, err := r.git.Run("rebase", "--abort")
	if err != nil {
		return fmt.Errorf("failed to abort rebase: %w", err)
	}

	output.Success("Rebase aborted successfully")
	return nil
}

// continueRebase continues current rebase
func (r *RebaseHandler) continueRebase() error {
	output.Info("Continuing rebase...")
	_, err := r.git.Run("rebase", "--continue")
	if err != nil {
		return fmt.Errorf("failed to continue rebase: %w", err)
	}

	output.Success("Rebase continued successfully")
	return nil
}

// skipRebase skips current commit in rebase
func (r *RebaseHandler) skipRebase() error {
	output.Info("Skipping current commit in rebase...")
	_, err := r.git.Run("rebase", "--skip")
	if err != nil {
		return fmt.Errorf("failed to skip commit in rebase: %w", err)
	}

	output.Success("Commit skipped successfully")
	return nil
}

// Helper methods
func (r *RebaseHandler) getTrunkBranch() string {
	trunkBranch := r.config.Trunk
	if trunkBranch == "" {
		trunkBranch, _ = r.git.GetConfig("at.trunk")
	}
	if trunkBranch == "" {
		trunkBranch = "main"
		// Check if main exists
		_, err := r.git.Run("rev-parse", "--verify", "main")
		if err != nil {
			trunkBranch = "master"
			// Check if master exists
			_, err := r.git.Run("rev-parse", "--verify", "master")
			if err != nil {
				// Neither exists, return main as default
				return "main"
			}
		}
	}
	return trunkBranch
}

// containsFlag checks if args contain any of the specified flags
func (r *RebaseHandler) containsFlag(args []string, flags ...string) bool {
	for _, arg := range args {
		for _, flag := range flags {
			if arg == flag {
				return true
			}
		}
	}
	return false
}

// showUsage displays the rebase command usage
func (r *RebaseHandler) showUsage() error {
	return output.Markdown(`# Rebase Command

Rebase current branch onto another branch with enhanced options.

## Usage

` + "```" + `bash
git @ rebase [command] [options]
git @ rebase onto <branch> [git-rebase-options]
git @ rebase interactive [commit-range]
git @ rebase abort
git @ rebase continue
git @ rebase skip
` + "```" + `

## Commands

• **onto <branch>**: Rebase current branch onto specified branch
• **interactive, i [commit-range]**: Start interactive rebase (default: last 5 commits)
• **abort**: Abort current rebase operation
• **continue**: Continue current rebase operation
• **skip**: Skip current commit in rebase

## Options

• **-h, --help**: Show this help message
• **All standard git rebase options** are supported

## Examples

` + "```" + `bash
# Rebase onto trunk branch
git @ rebase

# Rebase onto specific branch
git @ rebase onto develop

# Interactive rebase of last 3 commits
git @ rebase interactive HEAD~3

# Rebase with git options
git @ rebase onto main --autosquash

# Abort current rebase
git @ rebase abort
` + "```" + `

## Features

• **Trunk Integration**: Automatically detects and uses trunk branch
• **Safety Checks**: Validates branches exist before rebasing
• **Interactive Support**: Full support for interactive rebasing
• **Conflict Resolution**: Commands for aborting, continuing, and skipping
• **Remote Updates**: Automatically pulls latest changes before rebasing
• **User Confirmation**: Asks for confirmation before destructive operations

## Workflow

1. **Detect Target**: Automatically determine target branch or use specified one
2. **Update Target**: Pull latest changes from remote (if available)
3. **Confirm Action**: Ask user for confirmation
4. **Execute Rebase**: Run git rebase with appropriate options
5. **Show Results**: Display success or error messages

## Best Practices

• Always rebase onto the latest trunk branch before creating PRs
• Use interactive rebase to clean up commit history
• Abort rebase if conflicts are too complex to resolve
• Continue rebase after resolving conflicts
• Skip commits that cause persistent conflicts
`)
}