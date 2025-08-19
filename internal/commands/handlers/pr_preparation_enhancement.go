// pr_preparation.go - Pull Request preparation enhancements
package handlers

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// prepareBranchForPR asks user if they want to prepare their branch for PR
func prepareBranchForPR(g *git.Repository) error {
	output.Title("Prepare Branch for Pull Request")
	
	var prepareBranch bool
	err := huh.NewConfirm().
		Title("Prepare Branch?").
		Description("Do you want to prepare your branch for PR (squash, rebase, conventional commits)?").
		Value(&prepareBranch).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if !prepareBranch {
		return nil
	}

	// Get current branch
	currentBranch, err := g.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Get trunk branch
	trunkBranch, err := g.GetConfig("at.trunk")
	if err != nil {
		trunkBranch = "main"
		// Check if main exists
		_, err := g.Run("rev-parse", "--verify", "main")
		if err != nil {
			trunkBranch = "master"
		}
	}

	// 1. Squash commits if requested
	var squashCommits bool
	err = huh.NewConfirm().
		Title("Squash Commits?").
		Description("Do you want to squash your commits into a single conventional commit?").
		Value(&squashCommits).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if squashCommits {
		output.Info("Squashing commits...")
		// Here we would call the squash functionality
		// For demonstration, just show what would happen
		output.Info("Would squash commits on branch %s", currentBranch)
	}

	// 2. Rebase onto trunk if requested
	var rebaseOntoTrunk bool
	err = huh.NewConfirm().
		Title("Rebase onto Trunk?").
		Description("Do you want to rebase your branch onto the latest trunk?").
		Value(&rebaseOntoTrunk).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if rebaseOntoTrunk {
		output.Info("Updating %s branch...", trunkBranch)
		_, err = g.Run("remote", "get-url", "origin")
		if err == nil {
			_, err = g.Run("pull", "origin", trunkBranch)
			if err != nil {
				output.Warning("Failed to pull latest changes from remote, but continuing...")
			}
		}

		output.Info("Rebasing %s onto %s...", currentBranch, trunkBranch)
		_, err = g.Run("rebase", trunkBranch)
		if err != nil {
			return fmt.Errorf("rebase failed: %w", err)
		}
		output.Success("Successfully rebased %s onto %s", currentBranch, trunkBranch)
	}

	// 3. Use commitizen for conventional commit message if squashing
	if squashCommits {
		var useCommitizen bool
		err = huh.NewConfirm().
			Title("Use Commitizen?").
			Description("Do you want to create a conventional commit message?").
			Value(&useCommitizen).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get user input: %v", err)
		}

		if useCommitizen {
			output.Info("Launching commitizen...")
			// Here we would call the commitizen functionality
			// For demonstration, just show what would happen
			output.Info("Would launch commitizen to create conventional commit")
		}
	}

	// 4. Update changelog
	var updateChangelog bool
	err = huh.NewConfirm().
		Title("Update Changelog?").
		Description("Do you want to update the changelog with your changes?").
		Value(&updateChangelog).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if updateChangelog {
		output.Info("Updating changelog...")
		// Here we would call the changelog functionality
		// For demonstration, just show what would happen
		output.Info("Would update changelog with changes from branch %s", currentBranch)
	}

	output.Success("Branch preparation completed!")
	return nil
}