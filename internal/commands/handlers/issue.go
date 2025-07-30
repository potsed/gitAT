package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// IssueHandler handles issue-related commands
type IssueHandler struct {
	BaseHandler
}

// NewIssueHandler creates a new issue handler
func NewIssueHandler(cfg *config.Config, gitRepo *git.Repository) *IssueHandler {
	return &IssueHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the issue command
func (i *IssueHandler) Execute(args []string) error {
	if len(args) == 0 {
		return i.showIssue()
	}

	switch args[0] {
	case "set":
		if len(args) < 2 {
			return fmt.Errorf("issue ID required")
		}
		return i.setIssue(strings.Join(args[1:], " "))
	case "get":
		return i.showIssue()
	case "clear", "unset":
		return i.clearIssue()
	case "list":
		return i.listIssues()
	case "link":
		if len(args) < 2 {
			return fmt.Errorf("issue ID required for link command")
		}
		return i.linkIssue(args[1])
	case "-h", "--help":
		return i.showUsage()
	default:
		// If no subcommand, treat as set
		return i.setIssue(strings.Join(args, " "))
	}
}

// setIssue sets the issue ID
func (i *IssueHandler) setIssue(issueID string) error {
	if issueID == "" {
		return fmt.Errorf("issue ID cannot be empty")
	}

	// Validate issue ID format
	if !i.isValidIssueID(issueID) {
		return fmt.Errorf("invalid issue ID format: %s (use format like PROJ-123, #123, etc.)", issueID)
	}

	// Get current issue
	currentIssue, _ := i.git.GetConfig("at.issue")

	// If same issue, just confirm
	if currentIssue == issueID {
		output.Success("Issue already set to: %s", issueID)
		return nil
	}

	// Show confirmation if changing
	if currentIssue != "" {
		var proceed bool
		err := huh.NewConfirm().
			Title("Change Issue").
			Description(fmt.Sprintf("Change issue from '%s' to '%s'?", currentIssue, issueID)).
			Value(&proceed).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get user input: %v", err)
		}

		if !proceed {
			output.Info("Issue change cancelled")
			return nil
		}
	}

	// Set the issue
	if err := i.git.SetConfig("at.issue", issueID); err != nil {
		return fmt.Errorf("failed to set issue: %v", err)
	}

	output.Success("Issue set to: %s", issueID)
	return nil
}

// showIssue shows the current issue
func (i *IssueHandler) showIssue() error {
	issue, err := i.git.GetConfig("at.issue")
	if err != nil {
		output.Info("No issue configured")
		return nil
	}

	if issue == "" {
		output.Info("No issue configured")
		return nil
	}

	output.Success("Current issue: %s", issue)
	return nil
}

// clearIssue clears the issue setting
func (i *IssueHandler) clearIssue() error {
	currentIssue, _ := i.git.GetConfig("at.issue")
	if currentIssue == "" {
		output.Info("No issue configured")
		return nil
	}

	var proceed bool
	err := huh.NewConfirm().
		Title("Clear Issue").
		Description(fmt.Sprintf("Clear issue setting '%s'?", currentIssue)).
		Value(&proceed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if !proceed {
		output.Info("Issue clear cancelled")
		return nil
	}

	if err := i.git.SetConfig("at.issue", ""); err != nil {
		return fmt.Errorf("failed to clear issue: %v", err)
	}

	output.Success("Issue cleared")
	return nil
}

// listIssues shows recent issues from commit history
func (i *IssueHandler) listIssues() error {
	output.Info("Scanning commit history for issue references...")

	// Get recent commits
	commits, err := i.git.Run("log", "--oneline", "-20")
	if err != nil {
		return fmt.Errorf("failed to get commit history: %v", err)
	}

	if commits == "" {
		output.Info("No commits found")
		return nil
	}

	// Extract issue IDs from commits
	issuePattern := regexp.MustCompile(`(?i)(?:issue|bug|task|proj|jira|gh)-?\d+`)
	issues := make(map[string]int)

	lines := strings.Split(commits, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		matches := issuePattern.FindAllString(line, -1)
		for _, match := range matches {
			issues[strings.ToUpper(match)]++
		}
	}

	if len(issues) == 0 {
		output.Info("No issue references found in recent commits")
		return nil
	}

	// Display issues
	output.Title("Recent Issue References")
	for issue, count := range issues {
		output.Info("%s (mentioned %d times)", issue, count)
	}

	return nil
}

// linkIssue creates a link to the issue
func (i *IssueHandler) linkIssue(issueID string) error {
	if !i.isValidIssueID(issueID) {
		return fmt.Errorf("invalid issue ID format: %s", issueID)
	}

	// Try to determine issue tracker from remote URL
	remoteURL, err := i.git.GetRemoteURL("origin")
	if err != nil {
		output.Warning("Could not determine issue tracker from remote URL")
		output.Info("Issue ID: %s", issueID)
		return nil
	}

	// Generate issue URL based on remote
	issueURL := i.generateIssueURL(remoteURL, issueID)
	if issueURL != "" {
		output.Success("Issue URL: %s", issueURL)
	} else {
		output.Info("Issue ID: %s", issueID)
	}

	return nil
}

// isValidIssueID validates issue ID format
func (i *IssueHandler) isValidIssueID(issueID string) bool {
	// Common issue ID patterns
	patterns := []string{
		`^[A-Z]+-\d+$`,       // PROJ-123
		`^#\d+$`,             // #123
		`^[A-Z]+\d+$`,        // PROJ123
		`^[A-Z]+-\d+-\d+$`,   // PROJ-123-456
		`^[A-Z]+-\d+[A-Z]+$`, // PROJ-123A
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, strings.ToUpper(issueID)); matched {
			return true
		}
	}

	return false
}

// generateIssueURL generates issue URL based on remote
func (i *IssueHandler) generateIssueURL(remoteURL, issueID string) string {
	// GitHub
	if strings.Contains(remoteURL, "github.com") {
		// Extract owner/repo from URL
		parts := strings.Split(remoteURL, "/")
		if len(parts) >= 2 {
			owner := parts[len(parts)-2]
			repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
			return fmt.Sprintf("https://github.com/%s/%s/issues/%s", owner, repo, strings.TrimPrefix(issueID, "#"))
		}
	}

	// GitLab
	if strings.Contains(remoteURL, "gitlab.com") {
		parts := strings.Split(remoteURL, "/")
		if len(parts) >= 2 {
			owner := parts[len(parts)-2]
			repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
			return fmt.Sprintf("https://gitlab.com/%s/%s/-/issues/%s", owner, repo, strings.TrimPrefix(issueID, "#"))
		}
	}

	// Jira (common patterns)
	if strings.Contains(remoteURL, "jira") || strings.Contains(issueID, "JIRA") {
		// This would need to be configured per organization
		return fmt.Sprintf("https://jira.company.com/browse/%s", issueID)
	}

	return ""
}

// showUsage displays the usage information
func (i *IssueHandler) showUsage() error {
	usage := `# Issue Command

Manages issue/task identifiers for the repository.

## Usage

  git @ issue [<id>]
  git @ issue set <id>
  git @ issue get
  git @ issue clear
  git @ issue list
  git @ issue link <id>

## Commands

• **set <id>**: Set the issue ID
• **get**: Show current issue ID (default)
• **clear, unset**: Clear the issue setting
• **list**: Show recent issue references from commits
• **link <id>**: Generate issue URL for given ID

## Options

• **-h, --help**: Show this help message

## Examples

  # Set issue ID
  git @ issue PROJ-123
  git @ issue set PROJ-123

  # Show current issue
  git @ issue
  git @ issue get

  # Clear issue setting
  git @ issue clear

  # List recent issues
  git @ issue list

  # Link to issue
  git @ issue link PROJ-123

## Features

• **Validation**: Ensures issue IDs follow common formats
• **Confirmation**: Confirms changes when updating existing issue
• **Persistence**: Stores in Git configuration
• **Integration**: Works with other GitAT commands
• **URL Generation**: Creates links to issue trackers

## Issue ID Formats

• **PROJ-123**: Project prefix with dash
• **#123**: Simple number with hash
• **PROJ123**: Project prefix without dash
• **PROJ-123-456**: Multi-level issue
• **PROJ-123A**: Issue with suffix

## Use Cases

• **Task Tracking**: Link commits to specific issues/tasks
• **Code Review**: Reference issues in commit messages
• **Release Notes**: Generate release notes from issues
• **Project Management**: Track work across multiple issues

## Configuration

Issue IDs are stored in Git configuration:
  git config at.issue "PROJ-123"

## Notes

• Issue IDs are case-insensitive (stored as uppercase)
• Changes are stored in repository configuration
• Can be used by other GitAT commands for context
• Supports GitHub, GitLab, and Jira URL generation
`

	return output.Markdown(usage)
}
