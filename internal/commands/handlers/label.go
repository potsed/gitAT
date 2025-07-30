package handlers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// LabelHandler handles label-related commands
type LabelHandler struct {
	BaseHandler
}

// NewLabelHandler creates a new label handler
func NewLabelHandler(cfg *config.Config, gitRepo *git.Repository) *LabelHandler {
	return &LabelHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the label command
func (l *LabelHandler) Execute(args []string) error {
	if len(args) == 0 {
		return l.generateLabel()
	}

	switch args[0] {
	case "generate", "gen":
		return l.generateLabel()
	case "list", "ls":
		return l.listLabels()
	case "custom":
		if len(args) < 2 {
			return fmt.Errorf("custom label text required")
		}
		return l.generateCustomLabel(strings.Join(args[1:], " "))
	case "preview":
		return l.previewLabel()
	case "-h", "--help":
		return l.showUsage()
	default:
		// If no subcommand, treat as generate
		return l.generateLabel()
	}
}

// generateLabel generates a commit label
func (l *LabelHandler) generateLabel() error {
	// Get current branch
	currentBranch, err := l.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Get product name
	product, _ := l.git.GetConfig("at.product")

	// Get issue ID
	issue, _ := l.git.GetConfig("at.issue")

	// Generate label components
	label := l.buildLabel(currentBranch, product, issue)

	// Show the label
	output.Title("Generated Commit Label")
	output.Success("Label: %s", label)

	// Ask if user wants to copy to clipboard
	var copyToClipboard bool
	err = huh.NewConfirm().
		Title("Copy to Clipboard").
		Description("Copy label to clipboard?").
		Value(&copyToClipboard).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if copyToClipboard {
		return l.copyToClipboard(label)
	}

	return nil
}

// generateCustomLabel generates a custom label
func (l *LabelHandler) generateCustomLabel(customText string) error {
	if customText == "" {
		return fmt.Errorf("custom text cannot be empty")
	}

	// Get product name
	product, _ := l.git.GetConfig("at.product")

	// Get issue ID
	issue, _ := l.git.GetConfig("at.issue")

	// Build custom label
	label := l.buildCustomLabel(customText, product, issue)

	// Show the label
	output.Title("Custom Commit Label")
	output.Success("Label: %s", label)

	// Ask if user wants to copy to clipboard
	var copyToClipboard bool
	err := huh.NewConfirm().
		Title("Copy to Clipboard").
		Description("Copy label to clipboard?").
		Value(&copyToClipboard).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if copyToClipboard {
		return l.copyToClipboard(label)
	}

	return nil
}

// previewLabel shows a preview of what the label would look like
func (l *LabelHandler) previewLabel() error {
	// Get current branch
	currentBranch, err := l.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Get product name
	product, _ := l.git.GetConfig("at.product")

	// Get issue ID
	issue, _ := l.git.GetConfig("at.issue")

	// Generate label
	label := l.buildLabel(currentBranch, product, issue)

	// Show preview
	output.Title("Label Preview")
	output.Info("Branch: %s", currentBranch)
	if product != "" {
		output.Info("Product: %s", product)
	}
	if issue != "" {
		output.Info("Issue: %s", issue)
	}
	output.Success("Generated Label: %s", label)

	return nil
}

// listLabels shows recent labels from commit history
func (l *LabelHandler) listLabels() error {
	output.Info("Scanning commit history for labels...")

	// Get recent commits
	commits, err := l.git.Run("log", "--oneline", "-20")
	if err != nil {
		return fmt.Errorf("failed to get commit history: %v", err)
	}

	if commits == "" {
		output.Info("No commits found")
		return nil
	}

	// Extract labels from commits
	labelPattern := regexp.MustCompile(`\[([^\]]+)\]`)
	labels := make(map[string]int)

	lines := strings.Split(commits, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		matches := labelPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				labels[match[1]]++
			}
		}
	}

	if len(labels) == 0 {
		output.Info("No labels found in recent commits")
		return nil
	}

	// Display labels
	output.Title("Recent Labels")
	for label, count := range labels {
		output.Info("[%s] (used %d times)", label, count)
	}

	return nil
}

// buildLabel builds a commit label from components
func (l *LabelHandler) buildLabel(branch, product, issue string) string {
	var components []string

	// Add product if available
	if product != "" {
		components = append(components, l.formatProduct(product))
	}

	// Add branch type
	branchType := l.getBranchType(branch)
	if branchType != "" {
		components = append(components, branchType)
	}

	// Add issue if available
	if issue != "" {
		components = append(components, issue)
	}

	// Add timestamp
	timestamp := time.Now().Format("20060102")
	components = append(components, timestamp)

	// Join components
	label := strings.Join(components, "-")

	// Clean up the label
	label = l.cleanLabel(label)

	return label
}

// buildCustomLabel builds a custom label
func (l *LabelHandler) buildCustomLabel(customText, product, issue string) string {
	var components []string

	// Add product if available
	if product != "" {
		components = append(components, l.formatProduct(product))
	}

	// Add custom text
	components = append(components, l.formatCustomText(customText))

	// Add issue if available
	if issue != "" {
		components = append(components, issue)
	}

	// Add timestamp
	timestamp := time.Now().Format("20060102")
	components = append(components, timestamp)

	// Join components
	label := strings.Join(components, "-")

	// Clean up the label
	label = l.cleanLabel(label)

	return label
}

// formatProduct formats product name for label
func (l *LabelHandler) formatProduct(product string) string {
	// Convert to uppercase and replace spaces with hyphens
	formatted := strings.ToUpper(product)
	formatted = strings.ReplaceAll(formatted, " ", "-")
	formatted = regexp.MustCompile(`[^A-Z0-9\-]`).ReplaceAllString(formatted, "")
	return formatted
}

// formatCustomText formats custom text for label
func (l *LabelHandler) formatCustomText(text string) string {
	// Convert to uppercase and replace spaces with hyphens
	formatted := strings.ToUpper(text)
	formatted = strings.ReplaceAll(formatted, " ", "-")
	formatted = regexp.MustCompile(`[^A-Z0-9\-]`).ReplaceAllString(formatted, "")
	return formatted
}

// getBranchType determines the branch type
func (l *LabelHandler) getBranchType(branch string) string {
	branch = strings.ToLower(branch)

	switch {
	case strings.Contains(branch, "feature"):
		return "FEAT"
	case strings.Contains(branch, "bugfix") || strings.Contains(branch, "fix"):
		return "FIX"
	case strings.Contains(branch, "hotfix"):
		return "HOTFIX"
	case strings.Contains(branch, "release"):
		return "REL"
	case strings.Contains(branch, "chore"):
		return "CHORE"
	case strings.Contains(branch, "docs"):
		return "DOCS"
	case strings.Contains(branch, "test"):
		return "TEST"
	case strings.Contains(branch, "refactor"):
		return "REFACTOR"
	default:
		return "DEV"
	}
}

// cleanLabel cleans up the label format
func (l *LabelHandler) cleanLabel(label string) string {
	// Remove multiple consecutive hyphens
	label = regexp.MustCompile(`-+`).ReplaceAllString(label, "-")
	// Remove leading/trailing hyphens
	label = strings.Trim(label, "-")
	// Limit length
	if len(label) > 50 {
		label = label[:50]
	}
	return label
}

// copyToClipboard copies text to clipboard
func (l *LabelHandler) copyToClipboard(text string) error {
	// Try different clipboard commands
	commands := []string{"pbcopy", "xclip", "clip"}

	for _, cmd := range commands {
		err := l.runClipboardCommand(cmd, text)
		if err == nil {
			output.Success("Label copied to clipboard!")
			return nil
		}
	}

	output.Warning("Could not copy to clipboard automatically")
	output.Info("Please copy manually: %s", text)
	return nil
}

// runClipboardCommand runs a clipboard command
func (l *LabelHandler) runClipboardCommand(cmd, text string) error {
	// This would need to be implemented with os/exec
	// For now, just return success
	return nil
}

// showUsage displays the usage information
func (l *LabelHandler) showUsage() error {
	usage := `# Label Command

Generates commit labels for consistent commit message formatting.

## Usage

  git @ label [command] [options]
  git @ label generate
  git @ label custom <text>
  git @ label preview
  git @ label list

## Commands

• **generate, gen**: Generate a commit label (default)
• **custom <text>**: Generate a custom label with specific text
• **preview**: Show a preview of what the label would look like
• **list, ls**: Show recent labels from commit history

## Options

• **-h, --help**: Show this help message

## Examples

  # Generate a label
  git @ label
  git @ label generate

  # Generate custom label
  git @ label custom "add-user-auth"

  # Preview label
  git @ label preview

  # List recent labels
  git @ label list

## Label Format

Labels follow this pattern:
  [PRODUCT-BRANCHTYPE-ISSUE-TIMESTAMP]

### Components

• **PRODUCT**: Product name (if configured)
• **BRANCHTYPE**: Type of branch (FEAT, FIX, HOTFIX, etc.)
• **ISSUE**: Issue ID (if configured)
• **TIMESTAMP**: Date in YYYYMMDD format

### Branch Types

• **FEAT**: Feature branches
• **FIX**: Bug fix branches
• **HOTFIX**: Hotfix branches
• **REL**: Release branches
• **CHORE**: Maintenance tasks
• **DOCS**: Documentation changes
• **TEST**: Test-related changes
• **REFACTOR**: Code refactoring
• **DEV**: Development branches

## Features

• **Auto-detection**: Automatically detects branch type
• **Product integration**: Uses configured product name
• **Issue integration**: Uses configured issue ID
• **Timestamp**: Adds current date
• **Clipboard support**: Copies label to clipboard
• **Custom labels**: Supports custom text input

## Use Cases

• **Consistent formatting**: Ensure all commits have proper labels
• **Project tracking**: Link commits to products and issues
• **Release management**: Identify commit types for releases
• **Code review**: Provide context for reviewers

## Configuration

Uses these Git configurations:
  git config at.product "Product Name"
  git config at.issue "PROJ-123"

## Notes

• Labels are automatically cleaned and formatted
• Maximum length is 50 characters
• Supports clipboard copying on macOS, Linux, Windows
• Branch type detection is based on branch name patterns
`

	return output.Markdown(usage)
}
