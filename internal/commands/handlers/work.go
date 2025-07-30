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

// WorkHandler handles work-related commands
type WorkHandler struct {
	BaseHandler
}

// NewWorkHandler creates a new work handler
func NewWorkHandler(cfg *config.Config, gitRepo *git.Repository) *WorkHandler {
	return &WorkHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the work command
func (w *WorkHandler) Execute(args []string) error {
	if len(args) == 0 {
		return w.showUsage()
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return w.showUsage()
		}
	}

	return w.createWorkBranch(args)
}

// createWorkBranch creates a new work branch
func (w *WorkHandler) createWorkBranch(args []string) error {
	var workType, description, fullName string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-n", "--name":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				fullName = args[i+1]
				i++ // Skip next argument
			} else {
				return fmt.Errorf("error: --name requires a value")
			}
		default:
			if workType == "" {
				workType = arg
			} else if description == "" {
				description = arg
			} else {
				return fmt.Errorf("error: Too many arguments")
			}
		}
	}

	// Validate we're in a git repository
	_, err := w.git.Run("rev-parse", "--git-dir")
	if err != nil {
		return fmt.Errorf("error: Not in a git repository")
	}

	// Get current branch
	currentBranch, err := w.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	// If full name provided, use it directly
	if fullName != "" {
		return w.createWorkBranchFromName(fullName, currentBranch)
	}

	// Validate work type
	if workType == "" {
		return fmt.Errorf("error: Work type is required\nAvailable types: hotfix, feature, bugfix, release, chore, docs, style, refactor, perf, test, ci, build, revert")
	}

	// Validate work type against allowed types
	allowedTypes := []string{"hotfix", "feature", "bugfix", "release", "chore", "docs", "style", "refactor", "perf", "test", "ci", "build", "revert"}
	validType := false
	for _, t := range allowedTypes {
		if workType == t {
			validType = true
			break
		}
	}

	if !validType {
		return fmt.Errorf("error: Invalid work type '%s'\nAvailable types: %s", workType, strings.Join(allowedTypes, ", "))
	}

	// Prompt for description if not provided
	if description == "" {
		output.Title("🚀 Creating " + workType + " Branch")
		output.Info("Please enter a description for the %s:", workType)

		err := huh.NewInput().
			Title("Description").
			Description("Enter a description for your work").
			Value(&description).
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("description cannot be empty")
				}
				return nil
			}).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get description: %w", err)
		}
	}

	// Format description into kebab-case
	formattedDescription := w.formatBranchName(description)

	// Create branch name
	branchName := fmt.Sprintf("%s-%s", workType, formattedDescription)

	// Special handling for hotfix - should come from trunk
	if workType == "hotfix" {
		trunkBranch, _ := w.git.GetConfig("at.trunk")
		if trunkBranch == "" {
			trunkBranch = "main"
		}

		// Validate trunk branch exists
		_, err := w.git.Run("rev-parse", "--verify", trunkBranch)
		if err != nil {
			return fmt.Errorf("error: Trunk branch '%s' does not exist\nPlease ensure the trunk branch exists or configure it with: git @ _trunk <branch>", trunkBranch)
		}

		// Switch to trunk branch first
		output.Info("Switching to trunk branch: %s", trunkBranch)
		_, err = w.git.Run("checkout", trunkBranch)
		if err != nil {
			return fmt.Errorf("error: Failed to switch to trunk branch '%s'", trunkBranch)
		}

		// Update trunk branch
		output.Info("Updating trunk branch...")
		_, err = w.git.Run("remote", "get-url", "origin")
		if err == nil {
			_, err = w.git.Run("pull", "origin", trunkBranch)
			if err != nil {
				output.Warning("Failed to pull latest changes from remote, but continuing...")
			}
		}

		currentBranch = trunkBranch
	}

	// Create the work branch
	return w.createWorkBranchFromName(branchName, currentBranch)
}

// createWorkBranchFromName creates a work branch with the given name
func (w *WorkHandler) createWorkBranchFromName(branchName, baseBranch string) error {
	workType := strings.Split(branchName, "-")[0]

	// Validate branch name
	if !w.validateBranchName(branchName) {
		return fmt.Errorf("error: Invalid branch name '%s'\nBranch names must contain only alphanumeric characters, hyphens, underscores, and slashes", branchName)
	}

	// Check if branch already exists
	_, err := w.git.Run("rev-parse", "--verify", branchName)
	if err == nil {
		return fmt.Errorf("error: Branch '%s' already exists", branchName)
	}

	// Check for uncommitted changes
	_, err = w.git.Run("diff", "--quiet")
	hasUncommitted := err != nil
	_, err = w.git.Run("diff", "--cached", "--quiet")
	hasStaged := err != nil

	if hasUncommitted || hasStaged {
		output.Warning("You have uncommitted changes")
		output.Info("These will be saved to WIP before creating work branch")

		var confirmed bool
		err := huh.NewConfirm().
			Title("Continue?").
			Description("Do you want to continue and save changes to WIP?").
			Value(&confirmed).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !confirmed {
			output.Info("Operation cancelled")
			return nil
		}
	}

	// Save current WIP state
	output.Info("Saving current work state...")
	_ = w.setWIP() // Ignore errors, just warn

	// Create and switch to work branch
	output.Info("Creating %s branch: %s", workType, branchName)
	_, err = w.git.Run("checkout", "-b", branchName)
	if err != nil {
		return fmt.Errorf("error: Failed to create %s branch '%s'", workType, branchName)
	}

	// Set working branch to new branch
	output.Info("Setting working branch to %s branch...", workType)
	err = w.setBranch(branchName)
	if err != nil {
		output.Warning("Failed to set working branch, but %s branch is ready", workType)
	}

	output.Success("%s branch '%s' created successfully!", workType, branchName)

	// Show current status
	output.Title("📊 Current Status")
	workingBranch, _ := w.git.GetConfig("at.branch")

	statusData := [][]string{
		{"Branch", branchName},
		{"Base", baseBranch},
		{"Working Branch", workingBranch},
	}
	output.Table([]string{"Property", "Value"}, statusData)

	// Format work type for display
	displayType := w.getDisplayType(workType)
	titleType := w.getTitleType(workType)

	output.Title("🎯 Next Steps")
	steps := []string{
		"Make your changes",
		fmt.Sprintf("git @ save '[%s] Description of changes'", displayType),
		fmt.Sprintf("git @ pr '%s: Description of changes'", titleType),
	}

	// Special guidance for different types
	switch workType {
	case "hotfix":
		steps = append(steps, "After merge, consider: git @ release -p (patch release)")
	case "feature":
		steps = append(steps, "After merge, consider: git @ release -m (minor release)")
	case "release":
		steps = append(steps, "After merge, consider: git @ release -M (major release)")
	}

	for i, step := range steps {
		output.Info("%d. %s", i+1, step)
	}

	return nil
}

// formatBranchName formats a description into kebab-case
func (w *WorkHandler) formatBranchName(input string) string {
	// Convert to lowercase
	input = strings.ToLower(input)

	// Replace spaces, underscores, and other separators with hyphens
	input = strings.ReplaceAll(input, " ", "-")
	input = strings.ReplaceAll(input, "_", "-")
	input = strings.ReplaceAll(input, ".", "-")
	input = strings.ReplaceAll(input, "\\", "-")

	// Remove any non-alphanumeric characters except hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	input = reg.ReplaceAllString(input, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	input = reg.ReplaceAllString(input, "-")

	// Remove leading and trailing hyphens
	input = strings.Trim(input, "-")

	// If empty after formatting, use "update"
	if input == "" {
		input = "update"
	}

	return input
}

// validateBranchName validates a branch name
func (w *WorkHandler) validateBranchName(name string) bool {
	// Check if name is empty
	if name == "" {
		return false
	}

	// Check for dangerous characters
	if strings.ContainsAny(name, ";|`$(){}") {
		return false
	}

	// Check for path traversal attempts
	if strings.Contains(name, "../") {
		return false
	}

	// Check for valid characters (alphanumeric, hyphens, underscores, slashes)
	reg := regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
	if !reg.MatchString(name) {
		return false
	}

	// Check for reserved names
	reservedNames := []string{"HEAD", "head", "master", "main", "develop", "development"}
	for _, reserved := range reservedNames {
		if name == reserved {
			return false
		}
	}

	return true
}

// getDisplayType returns the display type for a work type
func (w *WorkHandler) getDisplayType(workType string) string {
	switch workType {
	case "hotfix":
		return "HOTFIX"
	case "feature":
		return "FEATURE"
	case "bugfix":
		return "BUGFIX"
	case "release":
		return "RELEASE"
	case "chore":
		return "CHORE"
	case "docs":
		return "DOCS"
	case "style":
		return "STYLE"
	case "refactor":
		return "REFACTOR"
	case "perf":
		return "PERF"
	case "test":
		return "TEST"
	case "ci":
		return "CI"
	case "build":
		return "BUILD"
	case "revert":
		return "REVERT"
	default:
		return strings.ToUpper(workType)
	}
}

// getTitleType returns the title type for a work type
func (w *WorkHandler) getTitleType(workType string) string {
	switch workType {
	case "hotfix":
		return "Hotfix"
	case "feature":
		return "Feature"
	case "bugfix":
		return "Bugfix"
	case "release":
		return "Release"
	case "chore":
		return "Chore"
	case "docs":
		return "Docs"
	case "style":
		return "Style"
	case "refactor":
		return "Refactor"
	case "perf":
		return "Perf"
	case "test":
		return "Test"
	case "ci":
		return "CI"
	case "build":
		return "Build"
	case "revert":
		return "Revert"
	default:
		return strings.Title(workType)
	}
}

// setWIP saves current work to WIP
func (w *WorkHandler) setWIP() error {
	// Check if there are changes to save
	_, err := w.git.Run("diff", "--quiet")
	if err == nil {
		// No changes to save
		return nil
	}

	// Get current branch
	currentBranch, err := w.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Create WIP commit
	_, err = w.git.Run("add", "-A")
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	_, err = w.git.Run("commit", "-m", fmt.Sprintf("[WIP] %s", currentBranch))
	if err != nil {
		return fmt.Errorf("failed to create WIP commit: %w", err)
	}

	return nil
}

// setBranch sets the working branch
func (w *WorkHandler) setBranch(branchName string) error {
	return w.git.SetConfig("at.branch", branchName)
}

// showUsage displays the work command usage
func (w *WorkHandler) showUsage() error {
	return output.Markdown(`# Work Command

Creates a new work branch for development.

## Usage

` + "```" + `bash
git @ work <type> [description]
git @ work --name <full-branch-name>
` + "```" + `

## Arguments

- **type**: The type of work (required unless using --name)
- **description**: Description of the work (optional, will prompt if not provided)

## Options

- **-n, --name**: Specify the full branch name directly
- **-h, --help**: Show this help message

## Work Types

- **feature**: New features
- **bugfix**: Bug fixes
- **hotfix**: Critical bug fixes (creates from trunk)
- **release**: Release preparation
- **chore**: Maintenance tasks
- **docs**: Documentation updates
- **style**: Code style changes
- **refactor**: Code refactoring
- **perf**: Performance improvements
- **test**: Test additions/changes
- **ci**: CI/CD changes
- **build**: Build system changes
- **revert**: Revert commits

## Examples

` + "```" + `bash
# Create a feature branch
git @ work feature add-user-authentication

# Create with custom name
git @ work --name feature/custom-branch-name

# Will prompt for description
git @ work bugfix
` + "```" + `

## Workflow

1. Creates a new branch from current branch (or trunk for hotfix)
2. Switches to the new branch
3. Sets it as the working branch
4. Provides next steps for your workflow

## Special Cases

- **Hotfix branches**: Automatically switch to trunk branch first
- **Uncommitted changes**: Will prompt to save to WIP before creating branch
`)
}
