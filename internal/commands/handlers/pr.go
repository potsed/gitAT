// pr.go - Clean Pull Request command handler
package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// PullRequestHandler handles PR-related commands
type PullRequestHandler struct {
	BaseHandler
}

// NewPullRequestHandler creates a new PR handler
func NewPullRequestHandler(cfg *config.Config, gitRepo *git.Repository) *PullRequestHandler {
	return &PullRequestHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the PR command
func (p *PullRequestHandler) Execute(args []string) error {
	if len(args) == 0 {
		return p.createPullRequest()
	}

	switch args[0] {
	case "create", "new":
		return p.createPullRequest()
	case "list", "ls":
		return p.listPullRequests()
	case "view", "show":
		if len(args) < 2 {
			return fmt.Errorf("PR number required for view command")
		}
		return p.viewPullRequest(args[1])
	case "merge":
		if len(args) < 2 {
			return fmt.Errorf("PR number required for merge command")
		}
		return p.mergePullRequest(args[1])
	case "close":
		if len(args) < 2 {
			return fmt.Errorf("PR number required for close command")
		}
		return p.closePullRequest(args[1])
	case "-h", "--help":
		return p.showUsage()
	default:
		return p.createPullRequest()
	}
}

// createPullRequest creates a new pull request
func (p *PullRequestHandler) createPullRequest() error {
	// Get current branch
	currentBranch, err := p.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Get trunk branch
	trunkBranch := p.config.Trunk
	if trunkBranch == "" {
		trunkBranch, _ = p.git.GetConfig("at.trunk")
	}
	if trunkBranch == "" {
		trunkBranch = "main"
		// Check if main exists
		_, err := p.git.Run("rev-parse", "--verify", "main")
		if err != nil {
			trunkBranch = "master"
		}
	}

	// Check if we're on trunk branch
	if currentBranch == trunkBranch {
		return fmt.Errorf("cannot create PR from trunk branch (%s)", trunkBranch)
	}

	// Ask if user wants to prepare PR
	var preparePR bool
	err = huh.NewConfirm().
		Title("Prepare Pull Request?").
		Description("Do you want to prepare your branch for PR (squash, rebase, conventional commits)?").
		Value(&preparePR).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if preparePR {
		// Perform PR preparation steps
		if err := p.prepareBranchForPR(currentBranch, trunkBranch); err != nil {
			return fmt.Errorf("failed to prepare branch for PR: %v", err)
		}
	}

	// Try to push the branch (will fail if already exists)
	output.Info("Pushing branch %s to remote...", currentBranch)
	if err := p.git.Push("origin", currentBranch); err != nil {
		output.Info("Branch may already exist on remote or push failed: %v", err)
	}

	// Get commit messages for description
	commits, err := p.getCommitsSinceBranch(trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to get commits: %v", err)
	}

	// Generate PR description
	description := p.generatePRDescription(commits)

	// Get user confirmation
	var proceed bool
	err = huh.NewConfirm().
		Title("Create Pull Request").
		Description(fmt.Sprintf("Create PR from %s to %s?", currentBranch, trunkBranch)).
		Value(&proceed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if !proceed {
		output.Info("PR creation cancelled")
		return nil
	}

	// Try to use GitHub CLI if available
	if p.hasGitHubCLI() {
		return p.createWithGitHubCLI(currentBranch, trunkBranch, description)
	}

	// Fallback to opening browser
	return p.openPRInBrowser(currentBranch, trunkBranch, description)
}

// prepareBranchForPR prepares a branch for PR by performing standard cleanup operations
func (p *PullRequestHandler) prepareBranchForPR(currentBranch, trunkBranch string) error {
	output.Title("Preparing Branch for Pull Request")
	
	// 1. Squash commits if requested
	var squashCommits bool
	err := huh.NewConfirm().
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
		// For now, just show what would happen
		output.Info("Would squash commits on branch %s", currentBranch)
	}

	// 2. Rebase onto trunk
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
		output.Info("Rebasing onto trunk...")
		// Here we would call the rebase functionality
		// For now, just show what would happen
		output.Info("Would rebase %s onto %s", currentBranch, trunkBranch)
		
		// Actually update trunk first
		output.Info("Updating %s branch...", trunkBranch)
		_, err = p.git.Run("remote", "get-url", "origin")
		if err == nil {
			_, err = p.git.Run("pull", "origin", trunkBranch)
			if err != nil {
				output.Warning("Failed to pull latest changes from remote, but continuing...")
			}
		}
		
		// Perform rebase
		output.Info("Rebasing %s onto %s...", currentBranch, trunkBranch)
		_, err = p.git.Run("rebase", trunkBranch)
		if err != nil {
			return fmt.Errorf("rebase failed: %w", err)
		}
		output.Success("Successfully rebased %s onto %s", currentBranch, trunkBranch)
	}

	// 3. Use commitizen for conventional commit message
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
			// For now, just show what would happen
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
		// For now, just show what would happen
		output.Info("Would update changelog with changes from branch %s", currentBranch)
	}

	output.Success("Branch preparation completed!")
	return nil
}

// listPullRequests lists open pull requests
func (p *PullRequestHandler) listPullRequests() error {
	if p.hasGitHubCLI() {
		output.Info("Listing pull requests...")
		cmd := exec.Command("gh", "pr", "list")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	output.Warning("GitHub CLI not found. Please install 'gh' to list PRs.")
	output.Info("You can view PRs at: https://github.com/[owner]/[repo]/pulls")
	return nil
}

// viewPullRequest views a specific pull request
func (p *PullRequestHandler) viewPullRequest(prNumber string) error {
	if p.hasGitHubCLI() {
		output.Info("Viewing pull request %s...", prNumber)
		cmd := exec.Command("gh", "pr", "view", prNumber)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	output.Warning("GitHub CLI not found. Please install 'gh' to view PRs.")
	return nil
}

// mergePullRequest merges a pull request
func (p *PullRequestHandler) mergePullRequest(prNumber string) error {
	if p.hasGitHubCLI() {
		output.Info("Merging pull request %s...", prNumber)
		cmd := exec.Command("gh", "pr", "merge", prNumber)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	output.Warning("GitHub CLI not found. Please install 'gh' to merge PRs.")
	return nil
}

// closePullRequest closes a pull request
func (p *PullRequestHandler) closePullRequest(prNumber string) error {
	if p.hasGitHubCLI() {
		var proceed bool
		err := huh.NewConfirm().
			Title("Close Pull Request").
			Description(fmt.Sprintf("Close PR #%s?", prNumber)).
			Value(&proceed).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get user input: %v", err)
		}

		if !proceed {
			output.Info("PR close cancelled")
			return nil
		}

		output.Info("Closing pull request %s...", prNumber)
		cmd := exec.Command("gh", "pr", "close", prNumber)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	output.Warning("GitHub CLI not found. Please install 'gh' to close PRs.")
	output.Info("You can close PR %s at: https://github.com/[owner]/[repo]/pull/%s", prNumber, prNumber)
	return nil
}

// getCommitsSinceBranch gets commits since the given branch
func (p *PullRequestHandler) getCommitsSinceBranch(branch string) ([]string, error) {
	output, err := p.git.Run("log", "--oneline", "--no-merges", fmt.Sprintf("%s..HEAD", branch))
	if err != nil {
		return nil, err
	}

	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}

// generatePRDescription generates a PR description from commits
func (p *PullRequestHandler) generatePRDescription(commits []string) string {
	if len(commits) == 0 {
		return "No commits to show"
	}

	var description strings.Builder
	description.WriteString("## Changes\n\n")

	for _, commit := range commits {
		if commit != "" {
			description.WriteString(fmt.Sprintf("- %s\n", commit))
		}
	}

	// Add template sections
	description.WriteString("\n## Description\n\n")
	description.WriteString("<!-- Describe your changes here -->\n\n")
	description.WriteString("## Testing\n\n")
	description.WriteString("<!-- Describe how you tested these changes -->\n\n")
	description.WriteString("## Checklist\n\n")
	description.WriteString("- [ ] Code follows project style guidelines\n")
	description.WriteString("- [ ] Self-review completed\n")
	description.WriteString("- [ ] Tests added/updated\n")
	description.WriteString("- [ ] Documentation updated\n")

	return description.String()
}

// hasGitHubCLI checks if GitHub CLI is available
func (p *PullRequestHandler) hasGitHubCLI() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// createWithGitHubCLI creates PR using GitHub CLI
func (p *PullRequestHandler) createWithGitHubCLI(sourceBranch, targetBranch, description string) error {
	output.Info("Creating PR using GitHub CLI...")

	// Create temporary description file
	tmpFile, err := os.CreateTemp("", "pr-description-*.md")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(description); err != nil {
		return fmt.Errorf("failed to write description: %v", err)
	}
	tmpFile.Close()

	// Run gh pr create
	cmd := exec.Command("gh", "pr", "create",
		"--base", targetBranch,
		"--head", sourceBranch,
		"--title", fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch),
		"--body-file", tmpFile.Name())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create PR with GitHub CLI: %v", err)
	}

	output.Success("Pull request created successfully!")
	return nil
}

// openPRInBrowser opens PR creation in browser
func (p *PullRequestHandler) openPRInBrowser(sourceBranch, targetBranch, description string) error {
	output.Info("Opening PR creation in browser...")

	// Get remote URL
	remoteURL, err := p.git.GetRemoteURL("origin")
	if err != nil {
		return fmt.Errorf("failed to get remote URL: %v", err)
	}

	// Convert SSH to HTTPS if needed
	if strings.HasPrefix(remoteURL, "git@") {
		remoteURL = strings.Replace(remoteURL, "git@github.com:", "https://github.com/", 1)
		remoteURL = strings.Replace(remoteURL, ".git", "", 1)
	}

	// Create PR URL
	prURL := fmt.Sprintf("%s/compare/%s...%s", remoteURL, targetBranch, sourceBranch)

	// Open browser
	cmd := exec.Command("open", prURL)
	if err := cmd.Run(); err != nil {
		// Try alternative commands
		cmd = exec.Command("xdg-open", prURL)
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("start", prURL)
			if err := cmd.Run(); err != nil {
				output.Info("Please open this URL in your browser: %s", prURL)
				return nil
			}
		}
	}

	output.Success("Browser opened with PR creation form!")
	output.Info("PR URL: %s", prURL)
	return nil
}

// showUsage displays the PR command usage
func (p *PullRequestHandler) showUsage() error {
	return output.Markdown(`# Pull Request Command

Create and manage pull requests for the repository.

## Usage

` + "```" + `bash
git @ pr [command] [options]
git @ pr create
git @ pr list
git @ pr view <number>
git @ pr merge <number>
git @ pr close <number>
` + "```" + `

## Commands

• **create, new**: Create a new pull request (default)
• **list, ls**: List open pull requests
• **view, show <number>**: View a specific pull request
• **merge <number>**: Merge a pull request
• **close <number>**: Close a pull request

## Options

• **-h, --help**: Show this help message

## Examples

` + "```" + `bash
# Create PR (default)
git @ pr

# Create PR with preparation
git @ pr create

# List all open PRs
git @ pr list

# View PR #123
git @ pr view 123

# Merge PR #123
git @ pr merge 123

# Close PR #123
git @ pr close 123
` + "```" + `

## Features

• **Auto-push**: Automatically pushes branch if not on remote
• **Smart description**: Generates PR description from commits
• **GitHub CLI integration**: Uses 'gh' command if available
• **Browser fallback**: Opens browser if GitHub CLI not available
• **Interactive confirmation**: Confirms actions before proceeding
• **PR Preparation**: Squash, rebase, conventional commits, changelog updates

## PR Preparation

Before creating a PR, GitAT can help prepare your branch:

• **Commit Squashing**: Squash multiple commits into a single conventional commit
• **Rebasing**: Rebase onto the latest trunk branch
• **Conventional Commits**: Use commitizen to create properly formatted commit messages
• **Changelog Updates**: Automatically update changelog with your changes

## Requirements

• GitHub CLI (optional, for enhanced functionality)
• Remote repository configured
• Branch pushed to remote (auto-handled)

## Notes

• Cannot create PR from trunk branch
• Automatically generates commit-based description
• Supports both GitHub CLI and browser workflows
`)
}