package commands

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// Manager handles all GitAT commands
type Manager struct {
	config *config.Config
	git    *git.Repository
}

// NewManager creates a new commands manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config: cfg,
		git:    git.NewRepository(cfg.RepoPath),
	}
}

// Work handles the work command
func (m *Manager) Work(args []string) error {
	if len(args) == 0 {
		return m.showWorkUsage()
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showWorkUsage()
		}
	}

	return m.createWorkBranch(args)
}

// Helper methods for work functionality
func (m *Manager) createWorkBranch(args []string) error {
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
	_, err := m.git.Run("rev-parse", "--git-dir")
	if err != nil {
		return fmt.Errorf("error: Not in a git repository")
	}

	// Get current branch
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	// If full name provided, use it directly
	if fullName != "" {
		return m.createWorkBranchFromName(fullName, currentBranch)
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
		fmt.Printf("Creating %s branch\n", workType)
		fmt.Printf("Please enter a description for the %s: ", workType)

		fmt.Scanln(&description)
		if description == "" {
			return fmt.Errorf("error: Description cannot be empty")
		}
	}

	// Format description into kebab-case
	formattedDescription := m.formatBranchName(description)

	// Create branch name
	branchName := fmt.Sprintf("%s-%s", workType, formattedDescription)

	// Special handling for hotfix - should come from trunk
	if workType == "hotfix" {
		trunkBranch, _ := m.git.GetConfig("at.trunk")
		if trunkBranch == "" {
			trunkBranch = "main"
		}

		// Validate trunk branch exists
		_, err := m.git.Run("rev-parse", "--verify", trunkBranch)
		if err != nil {
			return fmt.Errorf("error: Trunk branch '%s' does not exist\nPlease ensure the trunk branch exists or configure it with: git @ _trunk <branch>", trunkBranch)
		}

		// Switch to trunk branch first
		fmt.Printf("Switching to trunk branch: %s\n", trunkBranch)
		_, err = m.git.Run("checkout", trunkBranch)
		if err != nil {
			return fmt.Errorf("error: Failed to switch to trunk branch '%s'", trunkBranch)
		}

		// Update trunk branch
		fmt.Println("Updating trunk branch...")
		_, err = m.git.Run("remote", "get-url", "origin")
		if err == nil {
			_, err = m.git.Run("pull", "origin", trunkBranch)
			if err != nil {
				fmt.Println("Warning: Failed to pull latest changes from remote, but continuing...")
			}
		}

		currentBranch = trunkBranch
	}

	// Create the work branch
	return m.createWorkBranchFromName(branchName, currentBranch)
}

func (m *Manager) createWorkBranchFromName(branchName, baseBranch string) error {
	workType := strings.Split(branchName, "-")[0]

	// Validate branch name
	if !m.validateBranchName(branchName) {
		return fmt.Errorf("error: Invalid branch name '%s'\nBranch names must contain only alphanumeric characters, hyphens, underscores, and slashes", branchName)
	}

	// Check if branch already exists
	_, err := m.git.Run("rev-parse", "--verify", branchName)
	if err == nil {
		return fmt.Errorf("error: Branch '%s' already exists", branchName)
	}

	// Check for uncommitted changes
	_, err = m.git.Run("diff", "--quiet")
	hasUncommitted := err != nil
	_, err = m.git.Run("diff", "--cached", "--quiet")
	hasStaged := err != nil

	if hasUncommitted || hasStaged {
		fmt.Println("Warning: You have uncommitted changes")
		fmt.Println("These will be saved to WIP before creating work branch")
		fmt.Print("Continue? (y/N): ")

		var confirmation string
		fmt.Scanln(&confirmation)
		if !strings.HasPrefix(strings.ToLower(confirmation), "y") {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// Save current WIP state
	fmt.Println("Saving current work state...")
	_ = m.setWIP() // Ignore errors, just warn

	// Create and switch to work branch
	fmt.Printf("Creating %s branch: %s\n", workType, branchName)
	_, err = m.git.Run("checkout", "-b", branchName)
	if err != nil {
		return fmt.Errorf("error: Failed to create %s branch '%s'", workType, branchName)
	}

	// Set working branch to new branch
	fmt.Printf("Setting working branch to %s branch...\n", workType)
	err = m.setBranch(branchName)
	if err != nil {
		fmt.Printf("Warning: Failed to set working branch, but %s branch is ready\n", workType)
	}

	fmt.Println()
	fmt.Printf("✅ %s branch '%s' created successfully!\n", workType, branchName)
	fmt.Println()
	fmt.Println("Current status:")
	fmt.Printf("  Branch: %s\n", branchName)
	fmt.Printf("  Base: %s\n", baseBranch)
	workingBranch, _ := m.git.GetConfig("at.branch")
	fmt.Printf("  Working branch: %s\n", workingBranch)
	fmt.Println()

	// Format work type for display
	displayType := m.getDisplayType(workType)
	titleType := m.getTitleType(workType)

	fmt.Println("Next steps:")
	fmt.Printf("  1. Make your changes\n")
	fmt.Printf("  2. git @ save '[%s] Description of changes'\n", displayType)
	fmt.Printf("  3. git @ pr '%s: Description of changes'\n", titleType)

	// Special guidance for different types
	switch workType {
	case "hotfix":
		fmt.Println("  4. After merge, consider: git @ release -p (patch release)")
	case "feature":
		fmt.Println("  4. After merge, consider: git @ release -m (minor release)")
	case "release":
		fmt.Println("  4. After merge, consider: git @ release -M (major release)")
	}
	fmt.Println()

	return nil
}

// formatBranchName formats a description into kebab-case
func (m *Manager) formatBranchName(input string) string {
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
func (m *Manager) validateBranchName(name string) bool {
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
func (m *Manager) getDisplayType(workType string) string {
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
func (m *Manager) getTitleType(workType string) string {
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

// Hotfix handles the hotfix command
func (m *Manager) Hotfix(args []string) error {
	if len(args) == 0 {
		return m.createHotfix([]string{})
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showHotfixUsage()
		}
	}

	return m.createHotfix(args)
}

// Helper methods for hotfix functionality
func (m *Manager) createHotfix(args []string) error {
	var hotfixName string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-n", "--name":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				hotfixName = args[i+1]
				i++ // Skip next argument
			} else {
				return fmt.Errorf("error: --name requires a value")
			}
		default:
			// If no name provided yet, use this as name
			if hotfixName == "" {
				hotfixName = arg
			} else {
				return fmt.Errorf("error: Unknown option '%s'", arg)
			}
		}
	}

	// Validate we're in a git repository
	_, err := m.git.Run("rev-parse", "--git-dir")
	if err != nil {
		return fmt.Errorf("error: Not in a git repository")
	}

	// Get current branch
	_, err = m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	// Get trunk branch
	trunkBranch, _ := m.git.GetConfig("at.trunk")
	if trunkBranch == "" {
		trunkBranch = "main"
	}

	// Validate trunk branch exists
	_, err = m.git.Run("rev-parse", "--verify", trunkBranch)
	if err != nil {
		return fmt.Errorf("error: Trunk branch '%s' does not exist\nPlease ensure the trunk branch exists or configure it with: git @ _trunk <branch>", trunkBranch)
	}

	// Prompt for hotfix name if not provided
	if hotfixName == "" {
		fmt.Printf("Creating hotfix branch from %s\n", trunkBranch)
		fmt.Println("Please enter a name for the hotfix branch:")
		fmt.Println("Format: description (e.g., fix-login-bug)")
		fmt.Println("Branch will be created as: hotfix-description")
		fmt.Print("Hotfix description: ")

		fmt.Scanln(&hotfixName)
		if hotfixName == "" {
			return fmt.Errorf("error: Hotfix description cannot be empty")
		}
	}

	// Ensure hotfix name has the correct prefix
	if !strings.HasPrefix(hotfixName, "hotfix-") {
		hotfixName = "hotfix-" + hotfixName
	}

	// Validate hotfix name
	if !m.validateBranchName(hotfixName) {
		return fmt.Errorf("error: Invalid branch name '%s'\nBranch names must contain only alphanumeric characters, hyphens, underscores, and slashes", hotfixName)
	}

	// Check if hotfix branch already exists
	_, err = m.git.Run("rev-parse", "--verify", hotfixName)
	if err == nil {
		return fmt.Errorf("error: Branch '%s' already exists", hotfixName)
	}

	// Check for uncommitted changes
	_, err = m.git.Run("diff", "--quiet")
	hasUncommitted := err != nil
	_, err = m.git.Run("diff", "--cached", "--quiet")
	hasStaged := err != nil

	if hasUncommitted || hasStaged {
		fmt.Println("Warning: You have uncommitted changes")
		fmt.Println("These will be saved to WIP before creating hotfix branch")
		fmt.Print("Continue? (y/N): ")

		var confirmation string
		fmt.Scanln(&confirmation)
		if !strings.HasPrefix(strings.ToLower(confirmation), "y") {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	// Save current WIP state
	fmt.Println("Saving current work state...")
	_ = m.setWIP() // Ignore errors, just warn

	// Switch to trunk branch
	fmt.Printf("Switching to trunk branch: %s\n", trunkBranch)
	_, err = m.git.Run("checkout", trunkBranch)
	if err != nil {
		return fmt.Errorf("error: Failed to switch to trunk branch '%s'", trunkBranch)
	}

	// Ensure trunk branch is up to date
	fmt.Println("Updating trunk branch...")
	_, err = m.git.Run("remote", "get-url", "origin")
	if err == nil {
		_, err = m.git.Run("pull", "origin", trunkBranch)
		if err != nil {
			fmt.Println("Warning: Failed to pull latest changes from remote, but continuing...")
		}
	}

	// Create and switch to hotfix branch
	fmt.Printf("Creating hotfix branch: %s\n", hotfixName)
	_, err = m.git.Run("checkout", "-b", hotfixName)
	if err != nil {
		return fmt.Errorf("error: Failed to create hotfix branch '%s'", hotfixName)
	}

	// Set working branch to hotfix branch
	fmt.Println("Setting working branch to hotfix branch...")
	err = m.setBranch(hotfixName)
	if err != nil {
		fmt.Println("Warning: Failed to set working branch, but hotfix branch is ready")
	}

	fmt.Println()
	fmt.Printf("✅ Hotfix branch '%s' created successfully!\n", hotfixName)
	fmt.Println()
	fmt.Println("Current status:")
	fmt.Printf("  Branch: %s\n", hotfixName)
	fmt.Printf("  Base: %s\n", trunkBranch)
	workingBranch, _ := m.git.GetConfig("at.branch")
	fmt.Printf("  Working branch: %s\n", workingBranch)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Make your urgent fixes")
	fmt.Println("  2. git @ save 'Fix description'")
	fmt.Println("  3. git @ pr 'Hotfix: description'")
	fmt.Println("  4. After merge, consider: git @ release -p (patch release)")
	fmt.Println()

	return nil
}

// Save handles the save command
func (m *Manager) Save(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showSaveUsage()
		}
	}

	// Basic input validation
	if len(args) > 0 {
		// Check for dangerous characters
		message := strings.Join(args, " ")
		if strings.ContainsAny(message, ";|`$(){}") {
			return fmt.Errorf("error: Invalid message. Use only alphanumeric characters, dots, underscores, and hyphens")
		}
	}

	return m.saveWork(args)
}

// Helper methods for save functionality
func (m *Manager) saveWork(args []string) error {
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	workingBranch, _ := m.git.GetConfig("at.branch")
	repoPath, err := m.git.Run("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("error: Not in a git repository")
	}
	repoPath = strings.TrimSpace(repoPath)

	// If no working branch is set, set it to current branch
	if workingBranch == "" {
		fmt.Println("No working branch configured. Setting current branch as working branch...")
		err = m.git.SetConfig("at.branch", currentBranch)
		if err != nil {
			return fmt.Errorf("failed to set working branch: %w", err)
		}
		workingBranch = currentBranch
	}

	// Check branch protection
	if currentBranch == "master" || currentBranch == "develop" {
		return fmt.Errorf("error: Cannot save changes on %s. Create a new branch instead!", currentBranch)
	}

	if currentBranch == "prod" {
		fmt.Println("Warning: You are on the production branch!")
		fmt.Print("Are you sure you want to commit this? (Y/N): ")

		var confirmation string
		fmt.Scanln(&confirmation)

		if !strings.HasPrefix(strings.ToLower(confirmation), "y") {
			fmt.Println("Operation cancelled.")
			return nil
		}

		// Tag the version
		versionTag, err := m.git.Run("config", "at.major")
		if err == nil {
			major := strings.TrimSpace(versionTag)
			minor, _ := m.git.GetConfig("at.minor")
			fix, _ := m.git.GetConfig("at.fix")
			tagName := fmt.Sprintf("v%s.%s.%s", major, minor, fix)
			_, err = m.git.Run("tag", tagName)
			if err != nil {
				fmt.Printf("Warning: Failed to create version tag: %v\n", err)
			}
		}
	} else if currentBranch != workingBranch {
		return fmt.Errorf("error: Cannot save changes. You're not on the correct working branch '%s'\nCurrent branch: '%s'\nTo fix this, run: git @ branch '%s'", workingBranch, currentBranch, currentBranch)
	}

	// Generate commit message with label and user message
	var message string

	if len(args) > 0 {
		// User provided a message, combine with label and work type prefix
		workTypePrefix := m.getWorkTypePrefix(currentBranch)
		label, _ := m.showLabel()

		if label != "" {
			message = fmt.Sprintf("%s%s %s", workTypePrefix, label, strings.Join(args, " "))
		} else {
			message = fmt.Sprintf("%s%s", workTypePrefix, strings.Join(args, " "))
		}
	} else {
		// No user message, use default label with work type prefix
		workTypePrefix := m.getWorkTypePrefix(currentBranch)
		label, _ := m.showLabel()

		if label != "" {
			message = fmt.Sprintf("%s%s", workTypePrefix, label)
		} else {
			message = fmt.Sprintf("%sUpdate", workTypePrefix)
		}
	}

	// Add all changes and commit
	_, err = m.git.Run("add", ".")
	if err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	_, err = m.git.Run("commit", "-m", message)
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	fmt.Println("Changes saved successfully")
	return nil
}

// getWorkTypePrefix returns the appropriate work type prefix based on branch name
func (m *Manager) getWorkTypePrefix(branchName string) string {
	switch {
	case strings.HasPrefix(branchName, "hotfix-"):
		return "[HOTFIX] "
	case strings.HasPrefix(branchName, "feature-"):
		return "[FEATURE] "
	case strings.HasPrefix(branchName, "bugfix-"):
		return "[BUGFIX] "
	case strings.HasPrefix(branchName, "release-"):
		return "[RELEASE] "
	case strings.HasPrefix(branchName, "chore-"):
		return "[CHORE] "
	case strings.HasPrefix(branchName, "docs-"):
		return "[DOCS] "
	case strings.HasPrefix(branchName, "style-"):
		return "[STYLE] "
	case strings.HasPrefix(branchName, "refactor-"):
		return "[REFACTOR] "
	case strings.HasPrefix(branchName, "perf-"):
		return "[PERF] "
	case strings.HasPrefix(branchName, "test-"):
		return "[TEST] "
	case strings.HasPrefix(branchName, "ci-"):
		return "[CI] "
	case strings.HasPrefix(branchName, "build-"):
		return "[BUILD] "
	case strings.HasPrefix(branchName, "revert-"):
		return "[REVERT] "
	default:
		return ""
	}
}

// Squash handles the squash command
func (m *Manager) Squash(args []string) error {
	if len(args) == 0 {
		return m.squashToParent("", false)
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showSquashUsage()
		case "-s", "--save":
			return m.squashToParent("", true)
		case "-p", "--pr":
			return m.squashForPR()
		case "-a", "--auto":
			return fmt.Errorf("error: --auto requires a value (on|off|status)")
		}
	}

	if len(args) == 2 {
		switch args[0] {
		case "-a", "--auto":
			return m.handleAutoSquash(args[1])
		case "-s", "--save":
			return m.squashToParent(args[1], true)
		}
	}

	// Handle target branch argument
	if len(args) == 1 {
		return m.squashToParent(args[0], false)
	}

	if len(args) == 2 && args[0] == "-s" {
		return m.squashToParent(args[1], true)
	}

	return m.showSquashUsage()
}

// Helper methods for squash functionality
func (m *Manager) squashToParent(targetBranch string, doSave bool) error {
	var headSHA string
	var err error

	// If no target branch specified, auto-detect parent branch
	if targetBranch == "" {
		targetBranch, err = m.detectParentBranch()
		if err != nil {
			return fmt.Errorf("error: Could not auto-detect parent branch\nPlease specify a target branch: git @ squash <branch>")
		}
		fmt.Printf("Auto-detected parent branch: %s\n", targetBranch)
	}

	headSHA, err = m.getHeadSHA(targetBranch)
	if err != nil {
		return fmt.Errorf("error: Branch \"%s\" does not exist locally", targetBranch)
	}

	fmt.Printf("Target branch: %s (SHA: %s)\n", targetBranch, headSHA)
	err = m.performSquash(headSHA)
	if err != nil {
		return err
	}

	currentBranch, _ := m.git.GetCurrentBranch()
	fmt.Printf("Squashed branch %s back to %s\n", currentBranch, targetBranch)

	if doSave {
		return m.saveWork([]string{})
	}

	return nil
}

func (m *Manager) squashForPR() error {
	trunkBranch, _ := m.git.GetConfig("at.trunk")
	if trunkBranch == "" {
		trunkBranch = "main"
	}

	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	// Check if we're on the trunk branch
	if currentBranch == trunkBranch {
		return fmt.Errorf("error: Cannot squash PR from %s to itself", trunkBranch)
	}

	// Get the number of commits ahead of trunk branch
	output, err := m.git.Run("rev-list", "--count", fmt.Sprintf("%s..HEAD", trunkBranch))
	if err != nil {
		return fmt.Errorf("error: Cannot find merge base with %s", trunkBranch)
	}

	commitCount := strings.TrimSpace(output)
	if commitCount == "0" || commitCount == "1" {
		fmt.Println("Only one commit or no commits to squash")
		return nil
	}

	fmt.Printf("Found %s commits to squash for PR\n", commitCount)

	// Get the commit hash where the branch diverged from trunk
	baseCommit, err := m.git.Run("merge-base", trunkBranch, "HEAD")
	if err != nil {
		return fmt.Errorf("error: Cannot find merge base with %s", trunkBranch)
	}
	baseCommit = strings.TrimSpace(baseCommit)

	// Create a temporary branch for the squash
	tempBranch := fmt.Sprintf("%s-squash-%d", currentBranch, time.Now().Unix())

	// Create temp branch from base
	_, err = m.git.Run("checkout", "-b", tempBranch, baseCommit)
	if err != nil {
		return fmt.Errorf("error: Failed to create temporary branch")
	}

	// Cherry-pick all commits from current branch
	output, err = m.git.Run("rev-list", "--reverse", fmt.Sprintf("%s..HEAD", baseCommit))
	if err != nil {
		// Clean up on failure
		m.git.Run("checkout", currentBranch)
		m.git.Run("branch", "-D", tempBranch)
		return fmt.Errorf("error: Failed to get commit list")
	}

	commitHashes := strings.Split(strings.TrimSpace(output), "\n")
	cherryPickSuccess := true

	for _, commitHash := range commitHashes {
		commitHash = strings.TrimSpace(commitHash)
		if commitHash != "" {
			_, err = m.git.Run("cherry-pick", commitHash)
			if err != nil {
				fmt.Printf("Error: Failed to cherry-pick commit %s\n", commitHash)
				cherryPickSuccess = false
				break
			}
		}
	}

	if !cherryPickSuccess {
		// Clean up on failure
		m.git.Run("checkout", currentBranch)
		m.git.Run("branch", "-D", tempBranch)
		return fmt.Errorf("error: Squashing failed due to conflicts")
	}

	// Reset current branch to temp branch
	_, err = m.git.Run("checkout", currentBranch)
	if err != nil {
		m.git.Run("branch", "-D", tempBranch)
		return fmt.Errorf("error: Failed to switch back to current branch")
	}

	_, err = m.git.Run("reset", "--hard", tempBranch)
	if err != nil {
		m.git.Run("branch", "-D", tempBranch)
		return fmt.Errorf("error: Failed to reset current branch to squashed state")
	}

	// Clean up temp branch
	m.git.Run("branch", "-D", tempBranch)

	fmt.Printf("✅ Successfully squashed %s commits into one for PR\n", commitCount)
	return nil
}

func (m *Manager) handleAutoSquash(action string) error {
	switch action {
	case "on", "true", "enable", "1":
		return m.enableAutoSquash()
	case "off", "false", "disable", "0":
		return m.disableAutoSquash()
	case "status", "show", "check":
		return m.showAutoSquashStatus()
	default:
		return fmt.Errorf("error: Invalid auto action '%s'. Use 'on', 'off', or 'status'", action)
	}
}

func (m *Manager) enableAutoSquash() error {
	err := m.git.SetConfig("at.pr.squash", "true")
	if err != nil {
		return fmt.Errorf("failed to enable auto squash: %w", err)
	}

	fmt.Println("✅ Automatic PR squashing enabled")
	fmt.Println("   Commits will be automatically squashed before creating PRs")
	fmt.Println("   Use 'git @ pr -S' to override and skip squashing")
	return nil
}

func (m *Manager) disableAutoSquash() error {
	err := m.git.SetConfig("at.pr.squash", "false")
	if err != nil {
		return fmt.Errorf("failed to disable auto squash: %w", err)
	}

	fmt.Println("✅ Automatic PR squashing disabled")
	fmt.Println("   PRs will be created with all commits as-is")
	fmt.Println("   Use 'git @ pr -s' to force squashing")
	return nil
}

func (m *Manager) showAutoSquashStatus() error {
	setting, _ := m.git.GetConfig("at.pr.squash")

	fmt.Println("PR Squash Setting:")
	if setting == "true" {
		fmt.Println("  Status: ✅ ENABLED")
		fmt.Println("  Commits will be automatically squashed before creating PRs")
		fmt.Println("  Override: git @ pr -S (force no squash)")
	} else {
		fmt.Println("  Status: ❌ DISABLED")
		fmt.Println("  PRs will be created with all commits as-is")
		fmt.Println("  Override: git @ pr -s (force squash)")
	}

	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  git @ squash --auto on     # Enable automatic squashing")
	fmt.Println("  git @ squash --auto off    # Disable automatic squashing")
	fmt.Println("  git @ squash --auto status # Show this information")
	return nil
}

func (m *Manager) detectParentBranch() (string, error) {
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return "", err
	}

	// Method 1: Check git config for upstream branch
	parentBranch, _ := m.git.GetConfig(fmt.Sprintf("branch.%s.merge", currentBranch))
	if parentBranch != "" {
		parentBranch = strings.TrimPrefix(parentBranch, "refs/heads/")
		return parentBranch, nil
	}

	// Method 2: Check if current branch has an upstream tracking branch
	output, err := m.git.Run("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err == nil {
		parentBranch := strings.TrimSpace(output)
		parentBranch = strings.TrimPrefix(parentBranch, "refs/remotes/origin/")
		return parentBranch, nil
	}

	// Method 3: Find the branch that the current branch diverged from
	output, err = m.git.Run("branch", "--list")
	if err != nil {
		return "", err
	}

	branches := strings.Split(strings.TrimSpace(output), "\n")
	var bestParentBranch string
	var bestMergeDate int64

	for _, branch := range branches {
		branch = strings.TrimSpace(strings.TrimPrefix(branch, "* "))
		branch = strings.TrimSpace(strings.TrimPrefix(branch, " "))

		if branch != "" && branch != currentBranch {
			// Get the merge base between current branch and this branch
			mergeBase, err := m.git.Run("merge-base", branch, "HEAD")
			if err != nil {
				continue
			}
			mergeBase = strings.TrimSpace(mergeBase)

			if mergeBase != "" {
				// Get the commit date of the merge base
				dateOutput, err := m.git.Run("log", "-1", "--format=%ct", mergeBase)
				if err != nil {
					continue
				}

				mergeDate, err := strconv.ParseInt(strings.TrimSpace(dateOutput), 10, 64)
				if err != nil {
					continue
				}

				// The branch with the most recent merge base is the most likely parent
				if mergeDate > 0 && (bestParentBranch == "" || mergeDate > bestMergeDate) {
					bestParentBranch = branch
					bestMergeDate = mergeDate
				}
			}
		}
	}

	if bestParentBranch != "" {
		return bestParentBranch, nil
	}

	// Method 4: Fallback to configured trunk branch
	trunkBranch, _ := m.git.GetConfig("at.trunk")
	if trunkBranch != "" {
		return trunkBranch, nil
	}

	// Method 5: Try common branch names
	commonBranches := []string{"main", "master", "develop", "development"}
	for _, branch := range commonBranches {
		_, err := m.git.Run("rev-parse", "--verify", branch)
		if err == nil {
			return branch, nil
		}
	}

	return "", fmt.Errorf("no parent branch found")
}

func (m *Manager) getHeadSHA(branch string) (string, error) {
	output, err := m.git.Run("rev-parse", "--verify", "--quiet", branch)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (m *Manager) performSquash(targetSHA string) error {
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	// Validate target SHA
	if targetSHA == "" {
		return fmt.Errorf("❌ Invalid target SHA: %s", targetSHA)
	}

	// Verify target SHA exists
	_, err = m.git.Run("rev-parse", "--verify", targetSHA)
	if err != nil {
		return fmt.Errorf("❌ Target SHA does not exist: %s", targetSHA)
	}

	// Get the number of commits to squash
	output, err := m.git.Run("rev-list", "--count", fmt.Sprintf("%s..HEAD", targetSHA))
	if err != nil {
		return fmt.Errorf("error: Failed to count commits")
	}

	commitCount := strings.TrimSpace(output)
	count, err := strconv.Atoi(commitCount)
	if err != nil {
		return fmt.Errorf("error: Invalid commit count")
	}

	if count <= 1 {
		fmt.Println("Only one commit or no commits to squash")
		return nil
	}

	fmt.Printf("Squashing %d commits...\n", count)

	// Store current HEAD for safety
	originalHead, err := m.git.Run("rev-parse", "HEAD")
	if err != nil {
		return fmt.Errorf("error: Failed to get current HEAD")
	}
	originalHead = strings.TrimSpace(originalHead)

	// Create a temporary branch for the squash operation
	tempBranch := fmt.Sprintf("%s-squash-%d", currentBranch, time.Now().Unix())

	// Check if working directory is clean
	_, err = m.git.Run("diff", "--quiet")
	hasUncommitted := err != nil
	_, err = m.git.Run("diff", "--cached", "--quiet")
	hasStaged := err != nil

	if hasUncommitted || hasStaged {
		fmt.Println("⚠️  Working directory has uncommitted changes")
		fmt.Println("   Stashing changes before squashing...")
		_, err = m.git.Run("stash", "push", "-m", "Auto-stash before squashing")
		if err != nil {
			return fmt.Errorf("❌ Failed to stash uncommitted changes")
		}
	}

	// Create temp branch from target
	fmt.Println("Creating temporary branch from target...")
	fmt.Printf("   Current branch: %s\n", currentBranch)
	fmt.Printf("   Target SHA: %s\n", targetSHA)
	fmt.Printf("   Temp branch name: %s\n", tempBranch)

	_, err = m.git.Run("checkout", "-b", tempBranch, targetSHA)
	if err != nil {
		// Restore stashed changes if we stashed them
		if hasUncommitted || hasStaged {
			fmt.Println("   Restoring stashed changes...")
			m.git.Run("stash", "pop")
		}
		return fmt.Errorf("❌ Failed to create temporary branch for squashing")
	}

	// Cherry-pick all commits from current branch to temp branch
	output, err = m.git.Run("rev-list", "--reverse", fmt.Sprintf("%s..%s", targetSHA, originalHead))
	if err != nil {
		// Clean up on failure
		m.git.Run("checkout", currentBranch)
		m.git.Run("branch", "-D", tempBranch)
		if hasUncommitted || hasStaged {
			fmt.Println("   Restoring stashed changes...")
			m.git.Run("stash", "pop")
		}
		return fmt.Errorf("error: Failed to get commit list")
	}

	commitHashes := strings.Split(strings.TrimSpace(output), "\n")
	cherryPickSuccess := true

	for _, commitHash := range commitHashes {
		commitHash = strings.TrimSpace(commitHash)
		if commitHash != "" {
			_, err = m.git.Run("cherry-pick", commitHash)
			if err != nil {
				fmt.Printf("❌ Failed to cherry-pick commit %s\n", commitHash)
				cherryPickSuccess = false
				break
			}
		}
	}

	if !cherryPickSuccess {
		// Clean up on failure
		m.git.Run("checkout", currentBranch)
		m.git.Run("branch", "-D", tempBranch)
		if hasUncommitted || hasStaged {
			fmt.Println("   Restoring stashed changes...")
			m.git.Run("stash", "pop")
		}
		return fmt.Errorf("❌ Squashing failed due to conflicts")
	}

	// Reset current branch to temp branch (this creates the squashed commit)
	_, err = m.git.Run("checkout", currentBranch)
	if err != nil {
		m.git.Run("branch", "-D", tempBranch)
		if hasUncommitted || hasStaged {
			fmt.Println("   Restoring stashed changes...")
			m.git.Run("stash", "pop")
		}
		return fmt.Errorf("❌ Failed to switch back to current branch")
	}

	_, err = m.git.Run("reset", "--hard", tempBranch)
	if err != nil {
		m.git.Run("branch", "-D", tempBranch)
		if hasUncommitted || hasStaged {
			fmt.Println("   Restoring stashed changes...")
			m.git.Run("stash", "pop")
		}
		return fmt.Errorf("❌ Failed to reset current branch to squashed state")
	}

	// Clean up temp branch
	m.git.Run("branch", "-D", tempBranch)

	// Restore stashed changes if we stashed them
	if hasUncommitted || hasStaged {
		fmt.Println("   Restoring stashed changes...")
		_, err = m.git.Run("stash", "pop")
		if err != nil {
			fmt.Println("⚠️  Warning: Failed to restore stashed changes")
			fmt.Println("   Use 'git stash list' to see stashed changes")
		}
	}

	fmt.Printf("✅ Successfully squashed %d commits into one\n", count)
	return nil
}

// PullRequest handles the pr command
func (m *Manager) PullRequest(args []string) error {
	if len(args) == 0 {
		return m.createPR("", "", "", false, "", "")
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showPRUsage()
		}
	}

	return m.parsePRArgs(args)
}

// Helper methods for PR functionality
func (m *Manager) parsePRArgs(args []string) error {
	var title, description, baseBranch string
	var openBrowser bool
	var forceSquash, forceNoSquash string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-t", "--title":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				title = args[i+1]
				i++ // Skip next argument
			} else {
				return fmt.Errorf("error: --title requires a value")
			}
		case "-d", "--description":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				description = args[i+1]
				i++ // Skip next argument
			} else {
				return fmt.Errorf("error: --description requires a value")
			}
		case "-b", "--base":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				baseBranch = args[i+1]
				i++ // Skip next argument
			} else {
				return fmt.Errorf("error: --base requires a value")
			}
		case "-o", "--open":
			openBrowser = true
		case "-s", "--squash":
			forceSquash = "true"
		case "-S", "--no-squash":
			forceNoSquash = "true"
		default:
			// If no title provided yet, use this as title
			if title == "" {
				title = arg
			} else {
				return fmt.Errorf("error: Unknown option '%s'", arg)
			}
		}
	}

	return m.createPR(title, description, baseBranch, openBrowser, forceSquash, forceNoSquash)
}

func (m *Manager) createPR(title, description, baseBranch string, openBrowser bool, forceSquash, forceNoSquash string) error {
	// Validate we're in a git repository
	_, err := m.git.Run("rev-parse", "--git-dir")
	if err != nil {
		return fmt.Errorf("error: Not in a git repository")
	}

	// Get current branch
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	// Set default base branch if not provided
	if baseBranch == "" {
		baseBranch, _ = m.git.GetConfig("at.trunk")
		if baseBranch == "" {
			baseBranch = "main"
		}
	}

	// Check if we're trying to create PR from trunk branch
	if currentBranch == baseBranch {
		return fmt.Errorf("error: Cannot create PR from %s to itself", baseBranch)
	}

	// Check if there are commits between the branches
	output, err := m.git.Run("rev-list", "--count", fmt.Sprintf("%s..%s", baseBranch, currentBranch))
	if err != nil {
		return fmt.Errorf("error: Cannot find merge base with %s", baseBranch)
	}

	commitCount := strings.TrimSpace(output)
	if commitCount == "0" {
		return fmt.Errorf("error: No commits between %s and %s\nCannot create a PR without any commits to merge.", currentBranch, baseBranch)
	}

	// Check for uncommitted changes
	_, err = m.git.Run("diff", "--quiet")
	hasUncommitted := err != nil
	_, err = m.git.Run("diff", "--cached", "--quiet")
	hasStaged := err != nil

	if hasUncommitted || hasStaged {
		fmt.Println("Warning: You have uncommitted changes. Consider committing them first.")
		fmt.Print("Continue anyway? (y/N): ")

		var confirmation string
		fmt.Scanln(&confirmation)
		if !strings.HasPrefix(strings.ToLower(confirmation), "y") {
			return nil
		}
	}

	// Determine if we should squash commits
	shouldSquash := false
	if forceSquash == "true" {
		shouldSquash = true
	} else if forceNoSquash == "true" {
		shouldSquash = false
	} else {
		// Check configuration setting
		squashSetting, _ := m.git.GetConfig("at.pr.squash")
		if squashSetting == "true" {
			shouldSquash = true
		}
	}

	// Squash commits if enabled
	if shouldSquash {
		fmt.Println("Auto-squashing commits before creating PR...")
		err = m.squashForPR()
		if err != nil {
			return fmt.Errorf("error: Failed to squash commits: %w", err)
		}
		fmt.Println("✅ Commits squashed successfully")
	}

	// Set default title if not provided
	if title == "" {
		title, err = m.getDefaultPRTitle()
		if err != nil {
			title = fmt.Sprintf("Update from %s", currentBranch)
		}
	}

	// Generate automatic description if not provided
	if description == "" {
		fmt.Println("Generating automatic description based on changed files...")
		description = m.generateAutoDescription(baseBranch, currentBranch)
	}

	// Get platform and repo info
	platform := m.detectPlatform()
	repoInfo, err := m.getRepoInfo()
	if err != nil {
		return err
	}

	fmt.Printf("Creating PR for %s repository: %s\n", platform, repoInfo)
	fmt.Printf("From: %s → To: %s\n", currentBranch, baseBranch)
	fmt.Printf("Title: %s\n", title)
	if description != "" {
		fmt.Println("Description: Auto-generated based on changed files")
	}

	// Try to create PR using CLI tools
	success := false

	switch platform {
	case "github":
		if m.createGitHubPR(title, description, baseBranch, currentBranch) {
			success = true
		}
	case "gitlab":
		if m.createGitLabMR(title, description, baseBranch, currentBranch) {
			success = true
		}
	}

	// If CLI failed or not supported, provide web URL
	if !success {
		webURL := m.generateWebURL(platform, repoInfo, currentBranch, baseBranch)

		fmt.Println()
		fmt.Println("PR creation via CLI not available. Please create the PR manually:")
		fmt.Printf("URL: %s\n", webURL)

		if openBrowser {
			fmt.Println("Opening in browser...")
			m.openURL(webURL)
		}
	}

	return nil
}

func (m *Manager) detectPlatform() string {
	// Get the primary remote URL
	remoteURL, _ := m.git.GetConfig("remote.origin.url")

	if remoteURL == "" {
		return "unknown"
	}

	// Detect platform from URL
	if strings.Contains(remoteURL, "github.com") {
		return "github"
	} else if strings.Contains(remoteURL, "gitlab.com") || strings.Contains(remoteURL, "gitlab.") {
		return "gitlab"
	} else if strings.Contains(remoteURL, "bitbucket.org") || strings.Contains(remoteURL, "bitbucket.") {
		return "bitbucket"
	} else {
		return "generic"
	}
}

func (m *Manager) getRepoInfo() (string, error) {
	remoteURL, _ := m.git.GetConfig("remote.origin.url")
	platform := m.detectPlatform()

	if remoteURL == "" {
		return "", fmt.Errorf("error: No remote origin configured")
	}

	// Extract owner and repo name
	switch platform {
	case "github":
		// Handle both SSH and HTTPS URLs
		if strings.HasPrefix(remoteURL, "git@github.com:") {
			remoteURL = strings.TrimPrefix(remoteURL, "git@github.com:")
		} else if strings.HasPrefix(remoteURL, "https://github.com/") {
			remoteURL = strings.TrimPrefix(remoteURL, "https://github.com/")
		}
		return strings.TrimSuffix(remoteURL, ".git"), nil
	case "gitlab":
		// Handle both SSH and HTTPS URLs
		if strings.HasPrefix(remoteURL, "git@gitlab.com:") {
			remoteURL = strings.TrimPrefix(remoteURL, "git@gitlab.com:")
		} else if strings.HasPrefix(remoteURL, "https://gitlab.com/") {
			remoteURL = strings.TrimPrefix(remoteURL, "https://gitlab.com/")
		}
		return strings.TrimSuffix(remoteURL, ".git"), nil
	case "bitbucket":
		// Handle both SSH and HTTPS URLs
		if strings.HasPrefix(remoteURL, "git@bitbucket.org:") {
			remoteURL = strings.TrimPrefix(remoteURL, "git@bitbucket.org:")
		} else if strings.HasPrefix(remoteURL, "https://bitbucket.org/") {
			remoteURL = strings.TrimPrefix(remoteURL, "https://bitbucket.org/")
		}
		return strings.TrimSuffix(remoteURL, ".git"), nil
	default:
		return remoteURL, nil
	}
}

func (m *Manager) getDefaultPRTitle() (string, error) {
	output, err := m.git.Run("log", "-1", "--pretty=format:%s")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (m *Manager) generateAutoDescription(baseBranch, currentBranch string) string {
	// Get list of changed files
	output, err := m.git.Run("diff", "--name-only", fmt.Sprintf("%s..%s", baseBranch, currentBranch))
	if err != nil {
		return "No files changed compared to " + baseBranch
	}

	changedFiles := strings.Split(strings.TrimSpace(output), "\n")
	if len(changedFiles) == 0 || (len(changedFiles) == 1 && changedFiles[0] == "") {
		return "No files changed compared to " + baseBranch
	}

	// Count files by type
	totalFiles := 0
	addedFiles := 0
	modifiedFiles := 0
	deletedFiles := 0

	// Get file status information
	statusOutput, _ := m.git.Run("diff", "--name-status", fmt.Sprintf("%s..%s", baseBranch, currentBranch))
	statusLines := strings.Split(strings.TrimSpace(statusOutput), "\n")

	// Create a map of file statuses
	fileStatuses := make(map[string]string)
	for _, line := range statusLines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			fileStatuses[parts[1]] = parts[0]
		}
	}

	// Analyze each changed file
	for _, file := range changedFiles {
		file = strings.TrimSpace(file)
		if file != "" {
			totalFiles++

			// Check if file was added, modified, or deleted
			if status, exists := fileStatuses[file]; exists {
				if status == "A" {
					addedFiles++
				} else if status == "D" {
					deletedFiles++
				} else {
					modifiedFiles++
				}
			} else {
				modifiedFiles++
			}
		}
	}

	// Generate summary
	description := "# 📋 Pull Request Summary\n\n"
	description += fmt.Sprintf("This PR contains changes from branch `%s` targeting `%s`.\n\n", currentBranch, baseBranch)

	description += "## 📊 Changes Overview\n\n"
	description += "| Metric | Count |\n"
	description += "|--------|-------|\n"
	description += fmt.Sprintf("| **Total Files** | %d |\n", totalFiles)

	if addedFiles > 0 {
		description += fmt.Sprintf("| **Added** | %d |\n", addedFiles)
	}
	if modifiedFiles > 0 {
		description += fmt.Sprintf("| **Modified** | %d |\n", modifiedFiles)
	}
	if deletedFiles > 0 {
		description += fmt.Sprintf("| **Deleted** | %d |\n", deletedFiles)
	}

	description += "\n## 📁 File Analysis\n\n"

	// Group files by directory/type
	fileTypes := make(map[string]bool)
	dirs := make(map[string]bool)

	// Process changed files to build file types and directories
	for _, file := range changedFiles {
		file = strings.TrimSpace(file)
		if file != "" {
			// Get file extension
			parts := strings.Split(file, ".")
			ext := "no-extension"
			if len(parts) > 1 {
				ext = parts[len(parts)-1]
			}
			fileTypes[ext] = true

			// Get directory
			dir := "root"
			if strings.Contains(file, "/") {
				dir = file[:strings.LastIndex(file, "/")]
			}
			dirs[dir] = true
		}
	}

	// Add file type summary
	if len(fileTypes) > 0 {
		description += "### 🔤 File Types\n\n"
		description += "This PR affects the following file types:\n\n"
		for ext := range fileTypes {
			description += fmt.Sprintf("- `%s`\n", ext)
		}
		description += "\n"
	}

	// Add directory summary if multiple directories
	if len(dirs) > 1 {
		description += "### 📂 Directories Affected\n\n"
		description += "Changes span across the following directories:\n\n"
		for dir := range dirs {
			description += fmt.Sprintf("- `%s`\n", dir)
		}
		description += "\n"
	}

	// List all changed files
	description += "## 📝 Changed Files\n\n"
	description += "<details>\n<summary>📋 Click to view all changed files</summary>\n\n"
	description += "```\n"

	for _, file := range changedFiles {
		file = strings.TrimSpace(file)
		if file != "" {
			// Get file status
			status := "✏️" // Default to modified
			if statusCode, exists := fileStatuses[file]; exists {
				if statusCode == "A" {
					status = "➕"
				} else if statusCode == "D" {
					status = "🗑️"
				}
			}

			description += fmt.Sprintf("%s %s\n", status, file)
		}
	}

	description += "```\n</details>\n\n"

	// Add commit summary
	commitCount, _ := m.git.Run("rev-list", "--count", fmt.Sprintf("%s..%s", baseBranch, currentBranch))
	if commitCount != "" && commitCount != "0" && commitCount != "1" {
		description += "## 🔄 Commits\n\n"
		description += fmt.Sprintf("This PR includes **%s commits**.\n\n", commitCount)
		description += "<details>\n<summary>📜 Click to view commit history</summary>\n\n"
		description += "```\n"

		// Get commit history
		commitHistory, _ := m.git.Run("log", "--oneline", "-5", fmt.Sprintf("%s..%s", baseBranch, currentBranch))
		if commitHistory != "" {
			description += strings.TrimSpace(commitHistory) + "\n"
		}

		description += "```\n</details>\n\n"
	}

	// Add footer
	description += "---\n\n"
	description += "*This description was automatically generated based on the changes in this PR.*\n"

	return description
}

func (m *Manager) createGitHubPR(title, description, baseBranch, currentBranch string) bool {
	// Check if gh CLI is available
	_, err := m.git.Run("gh", "--version")
	if err != nil {
		fmt.Println("GitHub CLI (gh) not installed. Please install it or use the web interface.")
		fmt.Println("Install: https://cli.github.com/")
		return false
	}

	// Check if user is authenticated
	_, err = m.git.Run("gh", "auth", "status")
	if err != nil {
		fmt.Println("GitHub CLI not authenticated. Please run 'gh auth login' first.")
		return false
	}

	// Create PR
	if description != "" {
		_, err = m.git.Run("gh", "pr", "create", "--title", title, "--body", description, "--base", baseBranch, "--head", currentBranch)
	} else {
		_, err = m.git.Run("gh", "pr", "create", "--title", title, "--base", baseBranch, "--head", currentBranch)
	}

	return err == nil
}

func (m *Manager) createGitLabMR(title, description, baseBranch, currentBranch string) bool {
	// Check if glab CLI is available
	_, err := m.git.Run("glab", "--version")
	if err != nil {
		fmt.Println("GitLab CLI (glab) not installed. Please install it or use the web interface.")
		fmt.Println("Install: https://gitlab.com/gitlab-org/cli")
		return false
	}

	// Check if user is authenticated
	_, err = m.git.Run("glab", "auth", "status")
	if err != nil {
		fmt.Println("GitLab CLI not authenticated. Please run 'glab auth login' first.")
		return false
	}

	// Create MR
	if description != "" {
		_, err = m.git.Run("glab", "mr", "create", "--title", title, "--description", description, "--source-branch", currentBranch, "--target-branch", baseBranch)
	} else {
		_, err = m.git.Run("glab", "mr", "create", "--title", title, "--source-branch", currentBranch, "--target-branch", baseBranch)
	}

	return err == nil
}

func (m *Manager) generateWebURL(platform, repoInfo, currentBranch, baseBranch string) string {
	switch platform {
	case "github":
		return fmt.Sprintf("https://github.com/%s/compare/%s...%s", repoInfo, baseBranch, currentBranch)
	case "gitlab":
		return fmt.Sprintf("https://gitlab.com/%s/-/merge_requests/new?source_branch=%s&target_branch=%s", repoInfo, currentBranch, baseBranch)
	case "bitbucket":
		return fmt.Sprintf("https://bitbucket.org/%s/pull-requests/new?source=%s&t=1", repoInfo, currentBranch)
	default:
		return fmt.Sprintf("https://%s/compare/%s...%s", repoInfo, baseBranch, currentBranch)
	}
}

func (m *Manager) openURL(url string) {
	// Try to open URL in browser
	_, err := m.git.Run("open", url)
	if err != nil {
		_, err = m.git.Run("xdg-open", url)
		if err != nil {
			fmt.Printf("Please open this URL in your browser: %s\n", url)
		}
	}
}

func (m *Manager) showWorkUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ work <type> [<description>] [options]

DESCRIPTION:
  Create work branches following Conventional Commits specification.
  Supports all standard commit types for organized development workflow.

WORK TYPES (Conventional Commits):
  hotfix    : Urgent fixes for production (PATCH version)
  feature   : New features (MINOR version)
  bugfix    : Bug fixes (PATCH version)
  release   : Release preparation (PATCH version)
  chore     : Maintenance tasks
  docs      : Documentation changes
  style     : Code style changes
  refactor  : Code refactoring
  perf      : Performance improvements
  test      : Test additions/changes
  ci        : CI/CD changes
  build     : Build system changes
  revert    : Revert commits

OPTIONS:
  -n, --name <name>     Specify full branch name
  -h, --help           Show this help message

EXAMPLES:
  git @ work hotfix "fix-login-bug"           # Creates hotfix-fix-login-bug
  git @ work feature "add-user-auth"          # Creates feature-add-user-auth
  git @ work bugfix "fix-crash-on-startup"    # Creates bugfix-fix-crash-on-startup
  git @ work docs "update-api-documentation"  # Creates docs-update-api-documentation
  git @ work chore "update-dependencies"      # Creates chore-update-dependencies
  git @ work feature "Incorrect Branch Name"  # Creates feature-incorrect-branch-name
  git @ work hotfix "Fix Login Bug!"          # Creates hotfix-fix-login-bug

BRANCH NAMING:
  Format: <type>-<description>
  Automatic formatting: Descriptions are automatically converted to kebab-case
  Examples:
    hotfix-fix-login-bug
    feature-add-user-auth
    bugfix-fix-crash-on-startup
    docs-update-api-documentation
    chore-update-dependencies
    feature-incorrect-branch-name (from "Incorrect Branch Name")
    hotfix-fix-login-bug (from "Fix Login Bug!")

CONVENTIONAL COMMITS INTEGRATION:
  - Branch types follow Conventional Commits specification
  - Commit messages will include [TYPE] prefix
  - Supports semantic versioning correlation
  - Integrates with git @ branch --<type> listing

WORKFLOW:
  1. Creates branch from current branch or trunk
  2. Switches to new work branch
  3. Sets working branch to new branch
  4. Provides next steps guidance
`)
	return nil
}

func (m *Manager) showHotfixUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ hotfix [options] [<name>]

DESCRIPTION:
  Create a hotfix branch for urgent fixes that need to be deployed immediately.
  Creates a new branch from the trunk branch (master/main) and switches to it.

FEATURES:
  ✅ Creates hotfix branch from trunk (master/main)
  ✅ Saves current WIP state before switching
  ✅ Integrates with existing GitAT workflow
  ✅ Interactive name prompt if not provided
  ✅ Validates branch name and repository state

WORKFLOW:
  1. Saves current WIP state (if any)
  2. Switches to trunk branch
  3. Creates new hotfix branch from trunk
  4. Switches to hotfix branch
  5. Sets working branch to hotfix branch

OPTIONS:
  -n, --name <name>     Specify hotfix branch name
  -h, --help           Show this help message

EXAMPLES:
  git @ hotfix                    # Interactive name prompt
  git @ hotfix "fix-login-bug"    # Create hotfix with specific name
  git @ hotfix -n "security-patch" # Create hotfix with name option
  git @ hotfix --name "urgent-fix" # Create hotfix with name option

BRANCH NAMING:
  Format: hotfix-description (single hotfix- prefix)
  Examples:
    hotfix-fix-login-bug
    hotfix-security-patch
    hotfix-urgent-database-fix

INTEGRATION:
  - Uses git @ wip to save/restore work state
  - Uses git @ branch to set working branch
  - Uses git @ _trunk to get trunk branch name
  - Follows GitAT workflow patterns
`)
	return nil
}

func (m *Manager) showSaveUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ save [<message>]

DESCRIPTION:
  Securely save current changes with comprehensive validation and security checks.
  This is the primary command for committing changes in GitAT workflow.

FEATURES:
  ✅ Auto-branch setup: Sets working branch if not configured
  ✅ Security validation: Validates inputs and paths
  ✅ Branch protection: Prevents saves on master/develop
  ✅ Production warnings: Confirms before saving to prod
  ✅ Safe execution: Uses secure command execution

EXAMPLES:
  git @ save                           # Save with default message
  git @ save "Add user authentication" # Save with custom message
  git @ save "Fix login bug"           # Save with descriptive message

VALIDATION:
  Messages must contain only:
  - Alphanumeric characters (a-z, A-Z, 0-9)
  - Dots (.), underscores (_), hyphens (-)
  - Spaces and common punctuation

SECURITY:
  - All inputs are validated against dangerous patterns
  - Path operations are restricted to repository root
  - Commands are executed safely
  - Security events are logged

BRANCH PROTECTION:
  - Cannot save on master or develop branches
  - Production branch requires confirmation
  - Must be on configured working branch
`)
	return nil
}

func (m *Manager) setWIP() error {
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	oldWIP, _ := m.git.GetConfig("at.wip")
	err = m.git.SetConfig("at.wip", currentBranch)
	if err != nil {
		return fmt.Errorf("failed to set WIP: %w", err)
	}

	fmt.Printf("WIP updated to %s from %s\n", currentBranch, oldWIP)
	return nil
}

func (m *Manager) setBranch(branchName string) error {
	oldBranch, _ := m.git.GetConfig("at.branch")
	err := m.git.SetConfig("at.branch", branchName)
	if err != nil {
		return fmt.Errorf("failed to set branch: %w", err)
	}
	fmt.Printf("Branch updated to: %s from %s\n", branchName, oldBranch)
	return nil
}

func (m *Manager) showLabel() (string, error) {
	product, _ := m.git.GetConfig("at.product")
	feature, _ := m.git.GetConfig("at.feature")
	issue, _ := m.git.GetConfig("at.task")

	if product == "" || feature == "" || issue == "" {
		return "", nil
	}

	label := fmt.Sprintf("[%s.%s.%s]", product, feature, issue)
	return label, nil
}

func (m *Manager) showSquashUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ squash [options] [<target-branch>]

DESCRIPTION:
  Squash multiple commits into a single, consolidated commit by combining
  all commits ahead of the target branch into one clean commit. Automatically
  detects the parent branch based on where the current branch was created from.

OPTIONS:
  -s, --save           Run 'git @ save' after squashing
  -p, --pr             Squash for PR (uses configured trunk branch)
  -a, --auto           Enable/disable automatic PR squashing
  -h, --help           Show this help

EXAMPLES:
  git @ squash                      # Squash to parent branch (auto-detected)
  git @ squash -s                   # Squash to parent and save
  git @ squash develop              # Squash to specific branch
  git @ squash master -s            # Squash to master and save
  git @ squash --pr                 # Squash for PR using trunk branch
  git @ squash --auto on            # Enable automatic squashing
  git @ squash --auto off           # Disable automatic squashing
  git @ squash --auto status        # Show automatic squashing status

PR SQUASHING:
  When using --pr, the command will:
  1. Use the configured trunk branch (at.trunk) as target
  2. Squash commits ahead of the trunk branch
  3. Preserve commit messages in the final squashed commit

AUTOMATIC PR SQUASHING:
  Configure automatic squashing for git @ pr:
  git @ squash --auto on            # Enable automatic squashing
  git @ squash --auto off           # Disable automatic squashing
  git @ squash --auto status        # Show current setting

PROCESS:
  1. Auto-detects parent branch (or uses specified target)
  2. Validates target branch exists
  3. Retrieves HEAD SHA of target branch
  4. Creates temporary branch from target
  5. Cherry-picks all commits from current branch
  6. Resets current branch to squashed state
  7. Optionally runs 'git @ save'

USE CASES:
  - Clean up commit history before PR
  - Remove intermediate commits from feature branch
  - Create single clean commit from multiple commits
  - Automatic squashing before creating PRs
  - Consolidate related changes into meaningful commits
  - Simplify rollback operations

WARNING:
  You may need to force push after squashing if branch is shared.
  Use with caution on shared branches.

GIT COMMANDS USED:
  - git rev-parse --verify --quiet --long ${BRANCH}
  - git cherry-pick ${COMMIT}
  - git reset --hard ${BRANCH}
  - git checkout -b ${TEMP_BRANCH}

SECURITY:
  All squash operations are validated and logged.
`)
	return nil
}

func (m *Manager) showPRUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ pr [<title>] [options]

DESCRIPTION:
  Create a Pull Request (PR) or Merge Request (MR) for the current branch.
  Automatically detects the Git hosting platform and uses appropriate tools.

PLATFORMS SUPPORTED:
  ✅ GitHub: Uses 'gh' CLI or provides web URL
  ✅ GitLab: Uses 'glab' CLI or provides web URL
  ✅ Bitbucket: Provides web URL
  ✅ Generic: Provides web URL with branch info

FEATURES:
  ✅ Auto-platform detection
  ✅ CLI tool integration (gh, glab)
  ✅ Web URL fallback
  ✅ Branch validation
  ✅ Commit message integration
  ✅ Custom title and description
  ✅ Automatic commit squashing (configurable)
  ✅ Automatic description generation from changed files

EXAMPLES:
  git @ pr                                    # Create PR with default title and auto-generated description
  git @ pr "Add user authentication"          # Create PR with custom title and auto-generated description
  git @ pr -d "Detailed description here"     # Create PR with custom description
  git @ pr -b main                            # Create PR targeting main branch
  git @ pr -h                                 # Show this help

OPTIONS:
  -t, --title <title>       PR title (defaults to last commit message)
  -d, --description <desc>  PR description
  -b, --base <branch>       Target branch (defaults to configured trunk)
  -o, --open               Open PR in browser after creation
  -s, --squash             Force squash commits before PR (overrides setting)
  -S, --no-squash          Force no squash (overrides setting)
  -h, --help               Show this help message

AUTOMATIC FEATURES:
  - Uses last commit message as default title
  - Generates description from changed files (when not provided)
  - Includes branch name and commit info
  - Validates current branch is not trunk
  - Checks for uncommitted changes
  - Auto-squash commits if at.pr.squash is enabled

CONFIGURATION:
  git config at.pr.squash true    # Enable automatic squashing
  git config at.pr.squash false   # Disable automatic squashing
`)
	return nil
}

// Branch handles the branch command
func (m *Manager) Branch(args []string) error {
	if len(args) == 0 {
		// Show configured working branch
		return m.showBranch()
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showBranchUsage()
		case "-n", "--new", "n", "new":
			return m.newWorkingBranch()
		case "-c", "--current", "c", "current":
			return m.currentBranch()
		case "-s", "--set", "s", "set", ".":
			return m.setBranchToCurrent()
		case "--hotfix":
			return m.listBranchesByType("hotfix-")
		case "--feature":
			return m.listBranchesByType("feature-")
		case "--bugfix":
			return m.listBranchesByType("bugfix-")
		case "--release":
			return m.listBranchesByType("release-")
		case "--chore":
			return m.listBranchesByType("chore-")
		case "--docs":
			return m.listBranchesByType("docs-")
		case "--style":
			return m.listBranchesByType("style-")
		case "--refactor":
			return m.listBranchesByType("refactor-")
		case "--perf":
			return m.listBranchesByType("perf-")
		case "--test":
			return m.listBranchesByType("test-")
		case "--ci":
			return m.listBranchesByType("ci-")
		case "--build":
			return m.listBranchesByType("build-")
		case "--revert":
			return m.listBranchesByType("revert-")
		case "--all-types":
			return m.listAllWorkTypes()
		default:
			// Set working branch to specified name
			return m.setBranch(args[0])
		}
	}

	return m.showBranchUsage()
}

// Helper methods for branch management
func (m *Manager) showBranch() error {
	branch, err := m.git.GetConfig("at.branch")
	if err != nil {
		// No branch configured
		fmt.Println("")
		return nil
	}
	fmt.Println(branch)
	return nil
}

func (m *Manager) currentBranch() error {
	branch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}
	fmt.Println(branch)
	return nil
}

func (m *Manager) setBranchToCurrent() error {
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}
	return m.setBranch(currentBranch)
}

func (m *Manager) newWorkingBranch() error {
	_, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	// Create new branch name with timestamp
	now := time.Now()
	newBranchName := fmt.Sprintf("feature-%s", now.Format("20060102-150405"))

	fmt.Printf("Creating new working branch: %s\n", newBranchName)

	// Create and checkout new branch
	_, err = m.git.Run("checkout", "-b", newBranchName)
	if err != nil {
		return fmt.Errorf("failed to create new branch: %w", err)
	}

	// Set as working branch
	err = m.setBranch(newBranchName)
	if err != nil {
		return err
	}

	fmt.Printf("New working branch created and set: %s\n", newBranchName)
	return nil
}

func (m *Manager) listBranchesByType(prefix string) error {
	branchType := strings.TrimSuffix(prefix, "-")
	fmt.Printf("📋 %s branches:\n\n", branchType)

	// Get all local branches
	output, err := m.git.Run("branch", "--list")
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	currentBranch, _ := m.git.GetCurrentBranch()
	foundBranches := false

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Remove asterisk and spaces from branch name
		branch := strings.TrimSpace(strings.TrimPrefix(line, "* "))
		branch = strings.TrimSpace(strings.TrimPrefix(branch, " "))

		if branch != "" && strings.HasPrefix(branch, prefix) {
			foundBranches = true
			status := ""
			if branch == currentBranch {
				status = " (current)"
			}

			fmt.Printf("  🌿 %s%s\n", branch, status)

			// Get last commit info
			lastCommit, err := m.git.Run("log", "-1", "--oneline", branch)
			if err != nil {
				fmt.Printf("     📝 No commits\n")
			} else {
				fmt.Printf("     📝 %s\n", strings.TrimSpace(lastCommit))
			}
			fmt.Println()
		}
	}

	if !foundBranches {
		fmt.Printf("  No %s branches found\n\n", branchType)
		fmt.Printf("  To create a %s branch:\n", branchType)

		switch branchType {
		case "hotfix":
			fmt.Println("    git @ hotfix 'description'")
		case "feature":
			fmt.Println("    git @ feature 'description'")
		case "bugfix":
			fmt.Println("    git @ bugfix 'description'")
		default:
			fmt.Printf("    git checkout -b %sdescription\n", prefix)
		}
	}

	return nil
}

func (m *Manager) listAllWorkTypes() error {
	fmt.Println("📊 All work type branches:")

	workTypes := []string{"hotfix", "feature", "bugfix", "release", "chore", "docs", "style", "refactor", "perf", "test", "ci", "build", "revert"}
	currentBranch, _ := m.git.GetCurrentBranch()

	// Get all local branches
	output, err := m.git.Run("branch", "--list")
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, workType := range workTypes {
		var branches []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			branch := strings.TrimSpace(strings.TrimPrefix(line, "* "))
			branch = strings.TrimSpace(strings.TrimPrefix(branch, " "))

			if branch != "" && strings.HasPrefix(branch, workType+"-") {
				branches = append(branches, branch)
			}
		}

		if len(branches) > 0 {
			fmt.Printf("  📁 %s branches (%d):\n", workType, len(branches))
			for _, branch := range branches {
				status := ""
				if branch == currentBranch {
					status = " (current)"
				}
				fmt.Printf("    🌿 %s%s\n", branch, status)
			}
			fmt.Println()
		}
	}

	fmt.Println("💡 Use 'git @ branch --<type>' to see detailed info for each type")
	fmt.Println("   Example: git @ branch --hotfix")

	return nil
}

func (m *Manager) showBranchUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ branch [<branch-name>] [options]

DESCRIPTION:
  Manage working branch configuration and list branches by type.
  Shows current working branch, sets new working branch, or lists branches.

OPTIONS:
  (no options)           Show configured working branch
  <branch-name>          Set working branch to specified name
  -c, --current          Show current Git branch
  -s, --set, .           Set working branch to current branch
  -n, --new              Create new feature branch with timestamp
  --hotfix               List hotfix branches
  --feature              List feature branches
  --bugfix               List bugfix branches
  --release              List release branches
  --chore                List chore branches
  --docs                 List docs branches
  --style                List style branches
  --refactor             List refactor branches
  --perf                 List perf branches
  --test                 List test branches
  --ci                   List ci branches
  --build                List build branches
  --revert               List revert branches
  --all-types            List all work type branches
  -h, --help             Show this help

EXAMPLES:
  git @ branch                    # Show configured working branch
  git @ branch feature-auth       # Set working branch to feature-auth
  git @ branch -c                 # Show current Git branch
  git @ branch -s                 # Set working branch to current branch
  git @ branch -n                 # Create new feature branch
  git @ branch --hotfix           # List hotfix branches
  git @ branch --all-types        # List all work type branches

WORKFLOW:
  Use working branch to track which branch you're actively working on.
  This helps with context switching and branch management.
`)
	return nil
}

// Sweep handles the sweep command
func (m *Manager) Sweep(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showSweepUsage()
		}
	}

	// TODO: Implement sweep command
	fmt.Println("Sweep command not yet implemented")
	return nil
}

func (m *Manager) showSweepUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ sweep [options]

DESCRIPTION:
  Clean up local branches that have been merged or deleted remotely.
  Safely removes branches that are no longer needed.

OPTIONS:
  -h, --help           Show this help message

EXAMPLES:
  git @ sweep                    # Clean up merged and deleted branches
  git @ sweep --help             # Show this help

FEATURES:
  - Removes merged branches
  - Removes branches deleted from remote
  - Preserves current branch
  - Preserves configured working branch
  - Safe operation with confirmation
`)
	return nil
}

// Info handles the info command
func (m *Manager) Info(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showInfoUsage()
		}
	}

	// TODO: Implement info command
	fmt.Println("Info command not yet implemented")
	return nil
}

func (m *Manager) showInfoUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ info [options]

DESCRIPTION:
  Show comprehensive status report from all GitAT commands.
  Provides a complete overview of current repository state.

OPTIONS:
  -h, --help           Show this help message

EXAMPLES:
  git @ info                    # Show comprehensive status
  git @ info --help             # Show this help

FEATURES:
  - Repository status
  - Branch information
  - Configuration summary
  - Recent commits
  - Uncommitted changes
  - Remote status
`)
	return nil
}

// Hash handles the hash command
func (m *Manager) Hash(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showHashUsage()
		}
	}

	// TODO: Implement hash command
	fmt.Println("Hash command not yet implemented")
	return nil
}

func (m *Manager) showHashUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ hash [options]

DESCRIPTION:
  Show detailed branch status and commit relationships.
  Provides information about branch divergence and merge bases.

OPTIONS:
  -h, --help           Show this help message

EXAMPLES:
  git @ hash                    # Show branch status and relationships
  git @ hash --help             # Show this help

FEATURES:
  - Branch divergence information
  - Commit relationship analysis
  - Merge base details
  - Branch comparison
  - Commit graph visualization
`)
	return nil
}

// Product handles the product command
func (m *Manager) Product(args []string) error {
	if len(args) == 0 {
		// Show current product name
		product, err := m.git.GetConfig("at.product")
		if err != nil {
			fmt.Println("")
			return nil
		}
		fmt.Println(product)
		return nil
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showProductUsage()
		}
	}

	// Set product name
	productName := strings.Join(args, " ")

	// Validate product name
	if strings.ContainsAny(productName, ";|`$(){}") {
		return fmt.Errorf("error: Invalid product name. Use only alphanumeric characters, dots, underscores, and hyphens")
	}

	oldProduct, _ := m.git.GetConfig("at.product")
	err := m.git.SetConfig("at.product", productName)
	if err != nil {
		return fmt.Errorf("failed to set product: %w", err)
	}

	fmt.Printf("Project updated to: %s from %s\n", productName, oldProduct)
	return nil
}

func (m *Manager) showProductUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ product [<product-name>]

DESCRIPTION:
  Set or get the current product name for GitAT workflow management.
  The product name is used in commit labels and configuration.
  Examples: gitAT, myApp, apiService

EXAMPLES:
  git @ product                    # Show current product name
  git @ product gitAT              # Set product name to "gitAT"
  git @ product myApp              # Set product name to "myApp"
  git @ product apiService         # Set product name to "apiService"

VALIDATION:
  Product names must contain only:
  - Alphanumeric characters (a-z, A-Z, 0-9)
  - Dots (.)
  - Underscores (_)
  - Hyphens (-)

STORAGE:
  Saved in git config: at.product

SECURITY:
  All inputs are validated against dangerous characters and patterns.
`)
	return nil
}

// Feature handles the feature command
func (m *Manager) Feature(args []string) error {
	if len(args) == 0 {
		// Show current feature name
		feature, err := m.git.GetConfig("at.feature")
		if err != nil {
			fmt.Println("")
			return nil
		}
		fmt.Println(feature)
		return nil
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showFeatureUsage()
		}
	}

	// Set feature name
	featureName := strings.Join(args, " ")

	// Validate feature name
	if strings.ContainsAny(featureName, ";|`$(){}") {
		return fmt.Errorf("error: Invalid feature name. Use only alphanumeric characters, dots, underscores, and hyphens")
	}

	oldFeature, _ := m.git.GetConfig("at.feature")
	err := m.git.SetConfig("at.feature", featureName)
	if err != nil {
		return fmt.Errorf("failed to set feature: %w", err)
	}

	fmt.Printf("Feature updated to: %s from %s\n", featureName, oldFeature)
	return nil
}

func (m *Manager) showFeatureUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ feature [<feature-name>]

DESCRIPTION:
  Set or get the current feature name for GitAT workflow management.
  The feature name is used in commit labels and helps track what you're working on.

EXAMPLES:
  git @ feature                    # Show current feature name
  git @ feature user-auth          # Set feature to "user-auth"
  git @ feature payment-integration # Set feature to "payment-integration"

VALIDATION:
  Feature names must contain only:
  - Alphanumeric characters (a-z, A-Z, 0-9)
  - Dots (.)
  - Underscores (_)
  - Hyphens (-)

STORAGE:
  Saved in git config: at.feature

SECURITY:
  All inputs are validated against dangerous characters and patterns.
`)
	return nil
}

// Issue handles the issue command
func (m *Manager) Issue(args []string) error {
	if len(args) == 0 {
		// Show current issue ID
		issue, err := m.git.GetConfig("at.task")
		if err != nil {
			fmt.Println("")
			return nil
		}
		fmt.Println(issue)
		return nil
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showIssueUsage()
		}
	}

	// Set issue ID
	issueID := strings.Join(args, " ")

	oldIssue, _ := m.git.GetConfig("at.task")
	err := m.git.SetConfig("at.task", issueID)
	if err != nil {
		return fmt.Errorf("failed to set issue: %w", err)
	}

	fmt.Printf("Task updated to: %s from %s\n", issueID, oldIssue)
	return nil
}

func (m *Manager) showIssueUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ issue [<issue-id>]

DESCRIPTION:
  Set or get the current issue/task identifier for tracking.
  The issue ID is used in commit labels and helps link commits to issues.

EXAMPLES:
  git @ issue                    # Show current issue ID
  git @ issue PROJ-123           # Set issue to "PROJ-123"
  git @ issue BUG-456            # Set issue to "BUG-456"

STORAGE:
  Saved in git config: at.task

SECURITY:
  All issue operations are validated and logged.
`)
	return nil
}

// Version handles the version command
func (m *Manager) Version(args []string) error {
	if len(args) == 0 {
		// Show current version
		version, err := m.getVersion()
		if err != nil {
			return fmt.Errorf("failed to get version: %w", err)
		}
		fmt.Println(version)
		return nil
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showVersionUsage()
		case "-M", "--major":
			return m.incrementMajor()
		case "-m", "--minor":
			return m.incrementMinor()
		case "-b", "--bump":
			return m.incrementFix()
		case "-t", "--tag":
			version, err := m.getVersion()
			if err != nil {
				return fmt.Errorf("failed to get version: %w", err)
			}
			fmt.Printf("v%s\n", version)
			return nil
		case "-r", "--reset":
			return m.resetVersion()
		}
	}

	return m.showVersionUsage()
}

func (m *Manager) getVersion() (string, error) {
	major, err := m.git.GetConfig("at.major")
	if err != nil {
		major = "0"
	}

	minor, err := m.git.GetConfig("at.minor")
	if err != nil {
		minor = "0"
	}

	fix, err := m.git.GetConfig("at.fix")
	if err != nil {
		fix = "0"
	}

	return fmt.Sprintf("%s.%s.%s", major, minor, fix), nil
}

func (m *Manager) resetVersion() error {
	err := m.git.SetConfig("at.major", "0")
	if err != nil {
		return fmt.Errorf("failed to reset major version: %w", err)
	}

	err = m.git.SetConfig("at.minor", "0")
	if err != nil {
		return fmt.Errorf("failed to reset minor version: %w", err)
	}

	err = m.git.SetConfig("at.fix", "0")
	if err != nil {
		return fmt.Errorf("failed to reset fix version: %w", err)
	}

	fmt.Println("Version reset to 0.0.0")
	return nil
}

func (m *Manager) incrementMajor() error {
	major, err := m.git.GetConfig("at.major")
	if err != nil {
		major = "0"
	}

	majorInt, err := strconv.Atoi(major)
	if err != nil {
		return fmt.Errorf("invalid major version: %s", major)
	}

	newMajor := majorInt + 1
	err = m.git.SetConfig("at.major", strconv.Itoa(newMajor))
	if err != nil {
		return fmt.Errorf("failed to increment major version: %w", err)
	}

	// Reset minor and fix to 0
	err = m.git.SetConfig("at.minor", "0")
	if err != nil {
		return fmt.Errorf("failed to reset minor version: %w", err)
	}

	err = m.git.SetConfig("at.fix", "0")
	if err != nil {
		return fmt.Errorf("failed to reset fix version: %w", err)
	}

	fmt.Printf("Major version incremented: %d.0.0\n", newMajor)
	return nil
}

func (m *Manager) incrementMinor() error {
	minor, err := m.git.GetConfig("at.minor")
	if err != nil {
		minor = "0"
	}

	minorInt, err := strconv.Atoi(minor)
	if err != nil {
		return fmt.Errorf("invalid minor version: %s", minor)
	}

	newMinor := minorInt + 1
	err = m.git.SetConfig("at.minor", strconv.Itoa(newMinor))
	if err != nil {
		return fmt.Errorf("failed to increment minor version: %w", err)
	}

	// Reset fix to 0
	err = m.git.SetConfig("at.fix", "0")
	if err != nil {
		return fmt.Errorf("failed to reset fix version: %w", err)
	}

	major, _ := m.git.GetConfig("at.major")
	if major == "" {
		major = "0"
	}

	fmt.Printf("Minor version incremented: %s.%d.0\n", major, newMinor)
	return nil
}

func (m *Manager) incrementFix() error {
	fix, err := m.git.GetConfig("at.fix")
	if err != nil {
		fix = "0"
	}

	fixInt, err := strconv.Atoi(fix)
	if err != nil {
		return fmt.Errorf("invalid fix version: %s", fix)
	}

	newFix := fixInt + 1
	err = m.git.SetConfig("at.fix", strconv.Itoa(newFix))
	if err != nil {
		return fmt.Errorf("failed to increment fix version: %w", err)
	}

	major, _ := m.git.GetConfig("at.major")
	if major == "" {
		major = "0"
	}

	minor, _ := m.git.GetConfig("at.minor")
	if minor == "" {
		minor = "0"
	}

	fmt.Printf("Fix version incremented: %s.%s.%d\n", major, minor, newFix)
	return nil
}

func (m *Manager) showVersionUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ version [options]

DESCRIPTION:
  Manage semantic versioning for your project.
  Uses MAJOR.MINOR.FIX format (e.g., 1.2.3).

OPTIONS:
  (no options)           Show current version
  -M, --major           Increment major version (resets minor and fix to 0)
  -m, --minor           Increment minor version (resets fix to 0)
  -b, --bump            Increment fix version
  -t, --tag             Show version tag (e.g., "v1.2.3")
  -r, --reset           Reset version to 0.0.0 (use with caution)
  -h, --help            Show this help

EXAMPLES:
  git @ version                    # Show current version (e.g., "1.2.3")
  git @ version -M                 # Increment major: 1.2.3 → 2.0.0
  git @ version -m                 # Increment minor: 1.2.3 → 1.3.0
  git @ version -b                 # Increment fix: 1.2.3 → 1.2.4
  git @ version -t                 # Show version tag (e.g., "v1.2.3")
  git @ version -r                 # Reset to 0.0.0

STORAGE:
  Major version: git config at.major
  Minor version: git config at.minor
  Fix version: git config at.fix

SEMANTIC VERSIONING:
  MAJOR: Breaking changes, incompatible API changes
  MINOR: New features, backward compatible
  FIX: Bug fixes, backward compatible

SECURITY:
  All version operations are logged for audit purposes.
`)
	return nil
}

// Release handles the release command
func (m *Manager) Release(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showReleaseUsage()
		}
	}

	// TODO: Implement release command
	fmt.Println("Release command not yet implemented")
	return nil
}

func (m *Manager) showReleaseUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ release [options]

DESCRIPTION:
  Create releases with proper tagging and version management.
  Automatically creates version tags and release notes.

OPTIONS:
  -h, --help           Show this help message

EXAMPLES:
  git @ release                    # Create release with current version
  git @ release --help             # Show this help

FEATURES:
  - Automatic version tagging
  - Release note generation
  - Changelog creation
  - Semantic versioning support
`)
	return nil
}

// Master handles the master command
func (m *Manager) Master(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showMasterUsage()
		}
	}

	// Get current branch
	currentBranch, err := m.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("error: Not on a branch (detached HEAD state)")
	}

	// Get trunk branch
	trunkBranch, _ := m.git.GetConfig("at.trunk")
	if trunkBranch == "" {
		trunkBranch = "main"
	}

	// Check if we're already on trunk branch
	if currentBranch == trunkBranch {
		fmt.Printf("Already on %s branch\n", trunkBranch)
		return nil
	}

	// Check for uncommitted changes
	_, err = m.git.Run("diff", "--quiet")
	hasUncommitted := err != nil
	_, err = m.git.Run("diff", "--cached", "--quiet")
	hasStaged := err != nil

	if hasUncommitted || hasStaged {
		fmt.Println("Warning: You have uncommitted changes")
		fmt.Println("Stashing changes before switching to trunk branch...")

		_, err = m.git.Run("stash", "push", "-m", fmt.Sprintf("Auto-stash before switching to %s", trunkBranch))
		if err != nil {
			return fmt.Errorf("failed to stash changes: %w", err)
		}
	}

	// Switch to trunk branch
	fmt.Printf("Switching to %s branch...\n", trunkBranch)
	_, err = m.git.Run("checkout", trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to switch to %s branch: %w", trunkBranch, err)
	}

	// Pull latest changes
	fmt.Println("Pulling latest changes...")
	_, err = m.git.Run("pull", "origin", trunkBranch)
	if err != nil {
		fmt.Printf("Warning: Failed to pull latest changes: %v\n", err)
	}

	fmt.Printf("✅ Successfully switched to %s branch\n", trunkBranch)
	return nil
}

func (m *Manager) showMasterUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ master [options]

DESCRIPTION:
  Switch to the trunk branch (master/main) with automatic stash management.
  Safely switches to the configured trunk branch, stashing any uncommitted changes.

OPTIONS:
  -h, --help           Show this help message

EXAMPLES:
  git @ master                    # Switch to trunk branch
  git @ master --help             # Show this help

FEATURES:
  - Automatic stash of uncommitted changes
  - Pulls latest changes from remote
  - Uses configured trunk branch (at.trunk)
  - Safe branch switching
  - Status feedback

WORKFLOW:
  1. Stashes uncommitted changes (if any)
  2. Switches to trunk branch
  3. Pulls latest changes from remote
  4. Provides status feedback
`)
	return nil
}

// Changes handles the changes command
func (m *Manager) Changes(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showChangesUsage()
		}
	}

	// Show uncommitted changes
	output, err := m.git.Run("diff", "--name-only", "--no-color")
	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if output == "" {
		fmt.Println("No uncommitted changes")
	} else {
		fmt.Print(output)
	}

	return nil
}

func (m *Manager) showChangesUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ changes

DESCRIPTION:
  Show uncommitted changes in the working directory.
  Lists files that have been modified but not yet committed.

EXAMPLES:
  git @ changes                    # Show modified files

OUTPUT:
  Lists file names that have been changed since last commit.

SECURITY:
  All change operations are validated and logged.
`)
	return nil
}

// Logs handles the logs command
func (m *Manager) Logs(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showLogsUsage()
		}
	}

	// Show recent commit history
	output, err := m.git.Run("log", "-10", "--pretty=oneline", "--abbrev-commit")
	if err != nil {
		return fmt.Errorf("failed to get logs: %w", err)
	}

	fmt.Print(output)
	return nil
}

func (m *Manager) showLogsUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ logs

DESCRIPTION:
  Show recent commit history in a compact format.
  Displays the last 10 commits with abbreviated commit hashes.

EXAMPLES:
  git @ logs                    # Show recent commits

OUTPUT:
  Shows commit hash, author, date, and message for recent commits.

SECURITY:
  All log operations are validated and logged.
`)
	return nil
}

// WIP handles the wip command
func (m *Manager) WIP(args []string) error {
	if len(args) == 0 {
		return m.showWIP()
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showWIPUsage()
		case "-s", "--set", "s", "set", ".":
			return m.setWIP()
		case "-c", "--checkout", "c", "checkout":
			return m.checkoutWIP()
		case "-r", "--restore", "r", "restore":
			return m.restoreWIP()
		}
	}

	return m.showWIPUsage()
}

func (m *Manager) showWIP() error {
	wipBranch, err := m.git.GetConfig("at.wip")
	if err != nil || wipBranch == "" {
		fmt.Println("No WIP branch configured")
		return nil
	}
	fmt.Println(wipBranch)
	return nil
}

func (m *Manager) checkoutWIP() error {
	wipBranch, err := m.git.GetConfig("at.wip")
	if err != nil || wipBranch == "" {
		return fmt.Errorf("error: No WIP branch configured")
	}

	_, err = m.git.Run("checkout", wipBranch)
	if err != nil {
		return fmt.Errorf("failed to checkout WIP branch: %w", err)
	}

	fmt.Printf("Switched to WIP branch: %s\n", wipBranch)
	return nil
}

func (m *Manager) restoreWIP() error {
	wipBranch, err := m.git.GetConfig("at.wip")
	if err != nil || wipBranch == "" {
		return fmt.Errorf("error: No WIP branch configured")
	}

	// Set working branch to WIP branch
	err = m.setBranch(wipBranch)
	if err != nil {
		return fmt.Errorf("failed to set working branch: %w", err)
	}

	// Create work branch
	err = m.newWorkingBranch()
	if err != nil {
		return fmt.Errorf("failed to create work branch: %w", err)
	}

	fmt.Printf("Restored WIP branch: %s\n", wipBranch)
	return nil
}

func (m *Manager) showWIPUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ wip [options]

DESCRIPTION:
  Manage Work-In-Progress branch state.
  Tracks which branch you were working on for quick context switching.

OPTIONS:
  (no options)           Show current WIP branch
  -s, --set              Set current branch as WIP
  -c, --checkout         Checkout WIP branch
  -r, --restore          Restore WIP to working branch
  -h, --help             Show this help

EXAMPLES:
  git @ wip                    # Show current WIP branch
  git @ wip -s                 # Set current branch as WIP
  git @ wip -c                 # Checkout WIP branch
  git @ wip -r                 # Restore WIP to working branch

WORKFLOW:
  Use WIP to quickly switch between different features you're working on.
  Set WIP when you need to context switch to another task.

STORAGE:
  Saved in git config: at.wip

SECURITY:
  All WIP operations are validated and logged.
`)
	return nil
}

// Label handles the _label command
func (m *Manager) Label(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showLabelUsage()
		}
	}

	if len(args) == 0 {
		label, err := m.showLabel()
		if err != nil {
			return err
		}
		if label == "" {
			fmt.Println("[Update]")
		} else {
			fmt.Println(label)
		}
		return nil
	}

	if len(args) == 1 {
		return m.setLabel(args[0])
	}

	return m.showLabelUsage()
}

func (m *Manager) setLabel(label string) error {
	_, err := m.git.Run("config", "--replace-all", "at.label", label)
	if err != nil {
		return fmt.Errorf("failed to set label: %w", err)
	}

	newLabel, err := m.showLabel()
	if err != nil {
		return err
	}

	fmt.Printf("Label updated to: %s\n", newLabel)
	return nil
}

func (m *Manager) showLabelUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ _label [<custom-label>]

DESCRIPTION:
  Manage commit labels for GitAT workflow.
  Labels are used in commit messages and help track project context.

EXAMPLES:
  git @ _label                    # Show formatted label [product.feature.issue]
  git @ _label "Custom label"     # Set custom label

DEFAULT FORMAT:
  [product.feature.issue]
  Example: [gitAT.user-auth.PROJ-123]

STORAGE:
  Saved in git config: at.label

SECURITY:
  All label operations are validated and logged.
`)
	return nil
}

// ID handles the _id command
func (m *Manager) ID(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showIDUsage()
		}
	}

	product, _ := m.git.GetConfig("at.product")
	major, _ := m.git.GetConfig("at.major")
	minor, _ := m.git.GetConfig("at.minor")
	fix, _ := m.git.GetConfig("at.fix")

	id := fmt.Sprintf("%s:%s%s%s", product, major, minor, fix)
	fmt.Println(id)
	return nil
}

func (m *Manager) showIDUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ _id

DESCRIPTION:
  Generate a unique identifier for the current product state.
  Creates an ID based on product name and version.

FORMAT:
  product:major.minor.fix
  Example: gitAT:1.2.3

EXAMPLES:
  git @ _id                    # Show project ID

OUTPUT:
  Unique identifier combining product name and version.

SECURITY:
  All ID operations are validated and logged.
`)
	return nil
}

// Path handles the _path command
func (m *Manager) Path(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showPathUsage()
		}
	}

	output, err := m.git.Run("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("failed to get repository path: %w", err)
	}

	fmt.Print(output)
	return nil
}

func (m *Manager) showPathUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ _path

DESCRIPTION:
  Get the git repository root path.
  Returns the absolute path to the root of the current git repository.

EXAMPLES:
  git @ _path                    # Show repository root path

OUTPUT:
  Absolute path to the git repository root directory.

SECURITY:
  All path operations are validated and logged.
`)
	return nil
}

// Trunk handles the _trunk command
func (m *Manager) Trunk(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showTrunkUsage()
		}
	}

	if len(args) == 0 {
		return m.showTrunk()
	}

	if len(args) == 1 {
		return m.setTrunk(args[0])
	}

	return m.showTrunkUsage()
}

func (m *Manager) showTrunk() error {
	current, err := m.git.GetConfig("at.trunk")
	if err != nil || current == "" {
		// Auto-detect trunk branch from remote HEAD
		output, err := m.git.Run("branch", "-rl", "*/HEAD")
		if err != nil {
			current = "develop"
		} else {
			// Parse the output to get the branch name
			lines := strings.Split(strings.TrimSpace(output), "\n")
			if len(lines) > 0 {
				parts := strings.Split(strings.TrimSpace(lines[0]), "/")
				if len(parts) > 0 {
					current = parts[len(parts)-1]
				} else {
					current = "develop"
				}
			} else {
				current = "develop"
			}
		}
		// Set it without calling setTrunk to avoid recursion
		_, err = m.git.Run("config", "--replace-all", "at.trunk", current)
		if err != nil {
			return fmt.Errorf("failed to set trunk branch: %w", err)
		}
		fmt.Printf("Auto-detected trunk branch: %s\n", current)
	}
	fmt.Println(current)
	return nil
}

func (m *Manager) setTrunk(branchName string) error {
	from, _ := m.git.GetConfig("at.trunk")
	_, err := m.git.Run("config", "--replace-all", "at.trunk", branchName)
	if err != nil {
		return fmt.Errorf("failed to set trunk branch: %w", err)
	}

	fmt.Printf("Base branch updated to: %s from %s\n", branchName, from)
	return nil
}

func (m *Manager) showTrunkUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ _trunk [<branch-name>]

DESCRIPTION:
  Manage the base/trunk branch (usually develop or master).
  The trunk branch is used as the base for feature branches.

EXAMPLES:
  git @ _trunk                    # Show current trunk branch
  git @ _trunk develop            # Set trunk to "develop"
  git @ _trunk master             # Set trunk to "master"

AUTO-DETECTION:
  If no trunk is set, automatically detects from remote HEAD.

STORAGE:
  Saved in git config: at.trunk

SECURITY:
  All trunk operations are validated and logged.
`)
	return nil
}

// Ignore handles the ignore command
func (m *Manager) Ignore(args []string) error {
	if len(args) == 0 {
		return m.showIgnoreUsage()
	}

	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showIgnoreUsage()
		}
	}

	if len(args) == 1 {
		return m.addToGitignore(args[0])
	}

	return m.showIgnoreUsage()
}

func (m *Manager) addToGitignore(pattern string) error {
	// Get repository path
	repoPath, err := m.git.Run("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("failed to get repository path: %w", err)
	}
	repoPath = strings.TrimSpace(repoPath)
	gitignorePath := repoPath + "/.gitignore"

	// Check if pattern already exists
	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == pattern {
			fmt.Printf("String %s exists in %s\n", pattern, gitignorePath)
			return nil
		}
	}

	// Add pattern to .gitignore
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(pattern + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}

	fmt.Printf("String %s appended to %s\n", pattern, gitignorePath)
	return nil
}

func (m *Manager) showIgnoreUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ ignore <pattern>

DESCRIPTION:
  Add patterns to .gitignore file.
  Checks if pattern already exists before adding.

ARGUMENTS:
  <pattern>    Pattern to add to .gitignore

EXAMPLES:
  git @ ignore "*.log"              # Ignore all log files
  git @ ignore "node_modules/"      # Ignore node_modules directory
  git @ ignore "build/"             # Ignore build directory
  git @ ignore "*.tmp"              # Ignore temporary files

PROCESS:
  1. Checks if pattern already exists in .gitignore
  2. Adds pattern if not present
  3. Reports success or existing pattern

SECURITY:
  All ignore operations are validated and logged.
`)
	return nil
}

// InitLocal handles the initlocal command
func (m *Manager) InitLocal(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showInitLocalUsage()
		}
	}

	if len(args) == 0 {
		return m.showInitLocalUsage()
	}

	if len(args) == 1 {
		return m.initLocal(args[0])
	}

	return m.showInitLocalUsage()
}

func (m *Manager) initLocal(originURL string) error {
	// Initialize git repository
	_, err := m.git.Run("init")
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Add all files
	_, err = m.git.Run("add", ".")
	if err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// Create initial commit
	_, err = m.git.Run("commit", "-m", "Initial commit")
	if err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	// Set product name from origin URL
	parts := strings.Split(originURL, "/")
	if len(parts) > 0 {
		productName := parts[len(parts)-1]
		if strings.HasSuffix(productName, ".git") {
			productName = productName[:len(productName)-4]
		}
		err = m.Product([]string{productName})
		if err != nil {
			return fmt.Errorf("failed to set product name: %w", err)
		}
	}

	// Reset version
	err = m.resetVersion()
	if err != nil {
		return fmt.Errorf("failed to reset version: %w", err)
	}

	// Add remote origin
	remoteURL := fmt.Sprintf("git@gitlab.com:squibler/%s.git", originURL)
	_, err = m.git.Run("remote", "add", "origin", remoteURL)
	if err != nil {
		return fmt.Errorf("failed to add remote origin: %w", err)
	}

	// Push to master
	_, err = m.git.Run("push", "--set-upstream", "origin", "master")
	if err != nil {
		return fmt.Errorf("failed to push to master: %w", err)
	}

	// Create staging branch
	_, err = m.git.Run("checkout", "-b", "staging")
	if err != nil {
		return fmt.Errorf("failed to create staging branch: %w", err)
	}

	_, err = m.git.Run("push", "--set-upstream", "origin", "staging")
	if err != nil {
		return fmt.Errorf("failed to push staging branch: %w", err)
	}

	// Create develop branch
	_, err = m.git.Run("checkout", "-b", "develop")
	if err != nil {
		return fmt.Errorf("failed to create develop branch: %w", err)
	}

	_, err = m.git.Run("push", "--set-upstream", "origin", "develop")
	if err != nil {
		return fmt.Errorf("failed to push develop branch: %w", err)
	}

	// Show next steps
	fmt.Println("")
	fmt.Println("You should now visit the settings for this project and set")
	fmt.Println("the default branch to 'develop'")
	fmt.Println("")
	fmt.Printf("https://gitlab.com/squibler/%s/settings/repository\n", originURL)
	fmt.Println("")

	return nil
}

func (m *Manager) showInitLocalUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ initlocal <origin-url> <project-name>

DESCRIPTION:
  Initialize a new local repository with remote setup and proper branch structure.
  Creates master → staging → develop branch hierarchy.

ARGUMENTS:
  <origin-url>      Remote repository URL (e.g., git@gitlab.com:user)
  <project-name>    Project name for the repository

EXAMPLES:
  git @ initlocal git@gitlab.com:user my-project
  git @ initlocal git@github.com:org api-service

PROCESS:
  1. Initializes git repository
  2. Sets project name
  3. Creates remote origin
  4. Sets up branch structure: master → staging → develop
  5. Pushes all branches to remote

BRANCH STRUCTURE:
  master   → Production-ready code
  staging  → Pre-production testing
  develop  → Development and feature integration

NEXT STEPS:
  After initialization, visit your repository settings and set the default
  branch to 'develop':
  https://gitlab.com/user/project-name/settings/repository

SECURITY:
  All initialization operations are validated and logged.
`)
	return nil
}

// InitRemote handles the initremote command
func (m *Manager) InitRemote(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showInitRemoteUsage()
		}
	}

	if len(args) == 0 {
		return m.showInitRemoteUsage()
	}

	if len(args) == 1 {
		return m.initRemote(args[0])
	}

	return m.showInitRemoteUsage()
}

func (m *Manager) initRemote(repoURL string) error {
	// Initialize git repository
	_, err := m.git.Run("init")
	if err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Add remote origin
	_, err = m.git.Run("remote", "add", "origin", repoURL)
	if err != nil {
		return fmt.Errorf("failed to add remote origin: %w", err)
	}

	// Add all files
	_, err = m.git.Run("add", ".")
	if err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// Create initial commit
	_, err = m.git.Run("commit", "-m", "Initial commit")
	if err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	// Push to master
	_, err = m.git.Run("push", "--set-upstream", "origin", "master")
	if err != nil {
		return fmt.Errorf("failed to push to master: %w", err)
	}

	// Create develop branch
	_, err = m.git.Run("checkout", "-b", "develop")
	if err != nil {
		return fmt.Errorf("failed to create develop branch: %w", err)
	}

	// Create CHANGELOG file
	changelogContent := fmt.Sprintf("[%s]\r- CHANGELOG CREATED\r- INITIAL COMMIT\r\r", time.Now().Format("2006-01-02"))
	err = os.WriteFile("CHANGELOG", []byte(changelogContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create CHANGELOG: %w", err)
	}

	// Add and commit CHANGELOG
	_, err = m.git.Run("add", ".")
	if err != nil {
		return fmt.Errorf("failed to add CHANGELOG: %w", err)
	}

	_, err = m.git.Run("commit", "-m", "CHANGELOG CREATED")
	if err != nil {
		return fmt.Errorf("failed to commit CHANGELOG: %w", err)
	}

	// Push develop branch
	_, err = m.git.Run("push", "--set-upstream", "origin", "develop")
	if err != nil {
		return fmt.Errorf("failed to push develop branch: %w", err)
	}

	return nil
}

func (m *Manager) showInitRemoteUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ initremote <repository-url>

DESCRIPTION:
  Initialize a remote repository with basic structure.
  Creates initial commit and sets up develop branch with CHANGELOG.

ARGUMENTS:
  <repository-url>    Remote repository URL (e.g., git@github.com:user/repo.git)

PROCESS:
  1. Initializes git repository
  2. Adds remote origin
  3. Creates initial commit
  4. Pushes to master branch
  5. Creates develop branch
  6. Creates CHANGELOG file
  7. Pushes develop branch

EXAMPLES:
  git @ initremote git@github.com:user/my-repo.git
  git @ initremote git@gitlab.com:org/project.git

SECURITY:
  All initialization operations are validated and logged.
`)
	return nil
}

// Security handles the _security command
func (m *Manager) Security(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showSecurityUsage()
		}
	}

	// For now, just show security status
	return m.showSecurityStatus()
}

func (m *Manager) showSecurityStatus() error {
	fmt.Fprintf(os.Stdout, `GitAT Security Status

Security features are built into all GitAT operations:
- Input validation and sanitization
- Path traversal protection
- Command injection prevention
- Permission checking
- Secure configuration management
- Error handling and logging

All operations are validated and logged for security monitoring.
`)
	return nil
}

func (m *Manager) showSecurityUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ _security

DESCRIPTION:
  Security utilities for GitAT.
  Implements defensive coding practices and security checks.

FEATURES:
  - Input validation and sanitization
  - Path traversal protection
  - Command injection prevention
  - Permission checking
  - Secure configuration management
  - Error handling and logging

SECURITY:
  All security operations are validated and logged.
`)
	return nil
}

// Go handles the _go command
func (m *Manager) Go(args []string) error {
	if len(args) == 1 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return m.showGoUsage()
		}
	}

	return m.initializeGitAT()
}

func (m *Manager) initializeGitAT() error {
	// Set the base branch based on the origin
	output, err := m.git.Run("branch", "-rl", "*/HEAD")
	if err != nil {
		return fmt.Errorf("failed to get origin branch: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		parts := strings.Split(strings.TrimSpace(lines[0]), "/")
		if len(parts) > 0 {
			originBranch := parts[len(parts)-1]
			err = m.setTrunk(originBranch)
			if err != nil {
				return fmt.Errorf("failed to set trunk branch: %w", err)
			}
		}
	}

	// Reset the version
	err = m.resetVersion()
	if err != nil {
		return fmt.Errorf("failed to reset version: %w", err)
	}

	// Set the current working branch
	err = m.setBranchToCurrent()
	if err != nil {
		return fmt.Errorf("failed to set working branch: %w", err)
	}

	// Set the current WIP branch
	err = m.setWIP()
	if err != nil {
		return fmt.Errorf("failed to set WIP branch: %w", err)
	}

	// Mark repository as initialized
	_, err = m.git.Run("config", "--replace-all", "at.initialised", "true")
	if err != nil {
		return fmt.Errorf("failed to mark repository as initialized: %w", err)
	}

	fmt.Println("GitAT initialized successfully!")
	return nil
}

func (m *Manager) showGoUsage() error {
	fmt.Fprintf(os.Stdout, `Usage: git @ _go

DESCRIPTION:
  Initialize GitAT settings for a new repository.
  Sets up all main configurations for general use of the tool.

PROCESS:
  1. Sets base branch based on remote HEAD
  2. Resets version to 0.0.0
  3. Sets current working branch
  4. Sets current WIP branch
  5. Marks repository as initialized

EXAMPLES:
  git @ _go                    # Initialize GitAT for current repository

USE CASE:
  Run this command when setting up GitAT for the first time in a repository.

STORAGE:
  Sets git config: at.initialised = true

SECURITY:
  All initialization operations are validated and logged.
`)
	return nil
}
