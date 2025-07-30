package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// HotfixHandler handles hotfix-related commands
type HotfixHandler struct {
	BaseHandler
}

// NewHotfixHandler creates a new hotfix handler
func NewHotfixHandler(cfg *config.Config, gitRepo *git.Repository) *HotfixHandler {
	return &HotfixHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the hotfix command
func (h *HotfixHandler) Execute(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return h.showUsage()
	}

	// Parse options
	action := "create" // default action
	branchName := ""
	description := ""

	for i, arg := range args {
		switch arg {
		case "create", "new":
			action = "create"
			if i+1 < len(args) {
				branchName = args[i+1]
			}
		case "finish", "complete":
			action = "finish"
			if i+1 < len(args) {
				branchName = args[i+1]
			}
		case "list", "show":
			action = "list"
		}
	}

	// Get description if creating
	if action == "create" && branchName != "" && len(args) > 2 {
		description = strings.Join(args[2:], " ")
	}

	switch action {
	case "create":
		return h.createHotfix(branchName, description)
	case "finish":
		return h.finishHotfix(branchName)
	case "list":
		return h.listHotfixes()
	default:
		return h.createHotfix(branchName, description)
	}
}

// createHotfix creates a new hotfix branch
func (h *HotfixHandler) createHotfix(branchName, description string) error {
	// Get trunk branch
	trunkBranch, err := h.git.GetConfig("at.trunk")
	if err != nil || trunkBranch == "" {
		trunkBranch = "main" // default fallback
	}

	// Generate branch name if not provided
	if branchName == "" {
		var input string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Hotfix Name").
					Description("Enter a name for the hotfix (e.g., fix-login-bug)").
					Value(&input).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("hotfix name cannot be empty")
						}
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get hotfix name: %w", err)
		}
		branchName = input
	}

	// Get description if not provided
	if description == "" {
		var input string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Hotfix Description").
					Description("Describe what this hotfix addresses").
					Value(&input).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("description cannot be empty")
						}
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get description: %w", err)
		}
		description = input
	}

	// Check if we're on trunk branch
	currentBranch, err := h.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if currentBranch != trunkBranch {
		output.Warning("You're not on the trunk branch (%s). Switching to %s...", trunkBranch, trunkBranch)
		_, err = h.git.Run("checkout", trunkBranch)
		if err != nil {
			return fmt.Errorf("failed to switch to trunk branch: %w", err)
		}
	}

	// Pull latest changes
	output.Info("Pulling latest changes from %s...", trunkBranch)
	_, err = h.git.Run("pull", "origin", trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	// Create hotfix branch
	hotfixBranchName := fmt.Sprintf("hotfix/%s", branchName)
	_, err = h.git.Run("checkout", "-b", hotfixBranchName)
	if err != nil {
		return fmt.Errorf("failed to create hotfix branch: %w", err)
	}

	// Set branch description
	err = h.git.SetConfig(fmt.Sprintf("branch.%s.description", hotfixBranchName), description)
	if err != nil {
		output.Warning("Failed to set branch description: %v", err)
	}

	output.Success("Created hotfix branch: %s", hotfixBranchName)
	output.Info("Description: %s", description)
	output.Info("Base branch: %s", trunkBranch)

	return nil
}

// finishHotfix finishes a hotfix branch
func (h *HotfixHandler) finishHotfix(branchName string) error {
	// Get current branch if not specified
	if branchName == "" {
		currentBranch, err := h.git.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		branchName = currentBranch
	}

	// Validate it's a hotfix branch
	if !strings.HasPrefix(branchName, "hotfix/") {
		return fmt.Errorf("branch '%s' is not a hotfix branch", branchName)
	}

	// Get trunk branch
	trunkBranch, err := h.git.GetConfig("at.trunk")
	if err != nil || trunkBranch == "" {
		trunkBranch = "main" // default fallback
	}

	// Check if we have uncommitted changes
	_, err = h.git.Run("diff", "--quiet")
	if err != nil {
		output.Warning("You have uncommitted changes. Please commit or stash them first.")
		return nil
	}

	// Confirm finishing
	var confirmed bool
	err = huh.NewConfirm().
		Title("Finish Hotfix").
		Description(fmt.Sprintf("Finish hotfix branch '%s' and merge to %s?", branchName, trunkBranch)).
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("Hotfix finishing cancelled")
		return nil
	}

	// Switch to trunk branch
	_, err = h.git.Run("checkout", trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to switch to trunk branch: %w", err)
	}

	// Pull latest changes
	_, err = h.git.Run("pull", "origin", trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	// Merge hotfix branch
	_, err = h.git.Run("merge", "--no-ff", branchName)
	if err != nil {
		return fmt.Errorf("failed to merge hotfix branch: %w", err)
	}

	// Delete hotfix branch
	_, err = h.git.Run("branch", "-d", branchName)
	if err != nil {
		output.Warning("Failed to delete hotfix branch: %v", err)
	}

	output.Success("Hotfix '%s' finished successfully", branchName)
	output.Info("Merged to: %s", trunkBranch)

	return nil
}

// listHotfixes lists all hotfix branches
func (h *HotfixHandler) listHotfixes() error {
	output.Title("🔥 Hotfix Branches")

	// Get all branches
	branches, err := h.git.Run("branch", "--list", "hotfix/*")
	if err != nil {
		return fmt.Errorf("failed to get hotfix branches: %w", err)
	}

	if strings.TrimSpace(branches) == "" {
		output.Info("No hotfix branches found")
		return nil
	}

	branchList := strings.Split(strings.TrimSpace(branches), "\n")
	output.Info("Found %d hotfix branches:", len(branchList))

	for i, branch := range branchList {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			// Get branch description
			description, _ := h.git.GetConfig(fmt.Sprintf("branch.%s.description", branch))
			if description != "" {
				output.Info("  %d. %s - %s", i+1, branch, description)
			} else {
				output.Info("  %d. %s", i+1, branch)
			}
		}
	}

	return nil
}

// showUsage displays the hotfix command usage
func (h *HotfixHandler) showUsage() error {
	usage := "# Hotfix Command\n\n" +
		"Manages hotfix branches for critical bug fixes that need immediate deployment.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ hotfix [action] [name] [description]\n" +
		"git @ hotfix create [name] [description]\n" +
		"git @ hotfix finish [name]\n" +
		"git @ hotfix list\n" +
		"```\n\n" +
		"## Actions\n\n" +
		"- **create** (default): Create a new hotfix branch\n" +
		"- **finish**: Finish and merge a hotfix branch\n" +
		"- **list**: List all hotfix branches\n\n" +
		"## Arguments\n\n" +
		"- **name**: Name of the hotfix (e.g., fix-login-bug)\n" +
		"- **description**: Description of what the hotfix addresses\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# Create a hotfix with interactive prompts\n" +
		"git @ hotfix\n\n" +
		"# Create with name and description\n" +
		"git @ hotfix create fix-login-bug \"Fix authentication failure\"\n\n" +
		"# Finish current hotfix branch\n" +
		"git @ hotfix finish\n\n" +
		"# Finish specific hotfix branch\n" +
		"git @ hotfix finish hotfix/fix-login-bug\n\n" +
		"# List all hotfix branches\n" +
		"git @ hotfix list\n" +
		"```\n\n" +
		"## Workflow\n\n" +
		"### Creating Hotfix\n" +
		"1. **Validation**: Ensure on trunk branch\n" +
		"2. **Sync**: Pull latest changes from trunk\n" +
		"3. **Create**: Create hotfix branch from trunk\n" +
		"4. **Describe**: Set branch description\n" +
		"5. **Confirm**: Show success message\n\n" +
		"### Finishing Hotfix\n" +
		"1. **Validation**: Check for uncommitted changes\n" +
		"2. **Confirm**: Prompt for user confirmation\n" +
		"3. **Switch**: Switch to trunk branch\n" +
		"4. **Sync**: Pull latest changes\n" +
		"5. **Merge**: Merge hotfix branch (no-ff)\n" +
		"6. **Cleanup**: Delete hotfix branch\n" +
		"7. **Result**: Show success message\n\n" +
		"### Listing Hotfixes\n" +
		"1. **Search**: Find all hotfix branches\n" +
		"2. **Display**: Show formatted list with descriptions\n" +
		"3. **Count**: Display total number of hotfixes\n\n" +
		"## Safety Features\n\n" +
		"- **Trunk Validation**: Ensures hotfixes are created from trunk\n" +
		"- **Change Check**: Warns about uncommitted changes\n" +
		"- **Confirmation**: Prompts before finishing hotfixes\n" +
		"- **No-FF Merge**: Preserves hotfix history\n\n" +
		"## Best Practices\n\n" +
		"- Use hotfixes only for critical production issues\n" +
		"- Keep hotfixes small and focused\n" +
		"- Test thoroughly before finishing\n" +
		"- Deploy immediately after finishing\n" +
		"- Use descriptive names and descriptions\n"
	return output.Markdown(usage)
}
