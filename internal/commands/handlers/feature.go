package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// FeatureHandler handles feature-related commands
type FeatureHandler struct {
	BaseHandler
}

// NewFeatureHandler creates a new feature handler
func NewFeatureHandler(cfg *config.Config, gitRepo *git.Repository) *FeatureHandler {
	return &FeatureHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the feature command
func (f *FeatureHandler) Execute(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return f.showUsage()
	}

	// Parse options
	action := "create" // default action
	featureName := ""
	description := ""

	for i, arg := range args {
		switch arg {
		case "create", "new":
			action = "create"
			if i+1 < len(args) {
				featureName = args[i+1]
			}
		case "finish", "complete":
			action = "finish"
			if i+1 < len(args) {
				featureName = args[i+1]
			}
		case "list", "show":
			action = "list"
		case "start", "begin":
			action = "start"
			if i+1 < len(args) {
				featureName = args[i+1]
			}
		}
	}

	// Get description if creating
	if action == "create" && featureName != "" && len(args) > 2 {
		description = strings.Join(args[2:], " ")
	}

	switch action {
	case "create":
		return f.createFeature(featureName, description)
	case "finish":
		return f.finishFeature(featureName)
	case "list":
		return f.listFeatures()
	case "start":
		return f.startFeature(featureName)
	default:
		return f.createFeature(featureName, description)
	}
}

// createFeature creates a new feature branch
func (f *FeatureHandler) createFeature(featureName, description string) error {
	// Get develop branch
	developBranch, err := f.git.GetConfig("at.develop")
	if err != nil || developBranch == "" {
		developBranch = "develop" // default fallback
	}

	// Generate feature name if not provided
	if featureName == "" {
		var input string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Feature Name").
					Description("Enter a name for the feature (e.g., user-authentication)").
					Value(&input).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("feature name cannot be empty")
						}
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get feature name: %w", err)
		}
		featureName = input
	}

	// Get description if not provided
	if description == "" {
		var input string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Feature Description").
					Description("Describe what this feature implements").
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

	// Check if we're on develop branch
	currentBranch, err := f.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if currentBranch != developBranch {
		output.Warning("You're not on the develop branch (%s). Switching to %s...", developBranch, developBranch)
		_, err = f.git.Run("checkout", developBranch)
		if err != nil {
			return fmt.Errorf("failed to switch to develop branch: %w", err)
		}
	}

	// Pull latest changes
	output.Info("Pulling latest changes from %s...", developBranch)
	_, err = f.git.Run("pull", "origin", developBranch)
	if err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	// Create feature branch
	featureBranchName := fmt.Sprintf("feature/%s", featureName)
	_, err = f.git.Run("checkout", "-b", featureBranchName)
	if err != nil {
		return fmt.Errorf("failed to create feature branch: %w", err)
	}

	// Set branch description
	err = f.git.SetConfig(fmt.Sprintf("branch.%s.description", featureBranchName), description)
	if err != nil {
		output.Warning("Failed to set branch description: %v", err)
	}

	output.Success("Created feature branch: %s", featureBranchName)
	output.Info("Description: %s", description)
	output.Info("Base branch: %s", developBranch)

	return nil
}

// startFeature starts working on an existing feature
func (f *FeatureHandler) startFeature(featureName string) error {
	if featureName == "" {
		return fmt.Errorf("feature name required")
	}

	featureBranchName := fmt.Sprintf("feature/%s", featureName)

	// Check if feature branch exists
	_, err := f.git.Run("show-ref", "--verify", "--quiet", "refs/heads/"+featureBranchName)
	if err != nil {
		return fmt.Errorf("feature branch '%s' does not exist", featureBranchName)
	}

	// Switch to feature branch
	_, err = f.git.Run("checkout", featureBranchName)
	if err != nil {
		return fmt.Errorf("failed to switch to feature branch: %w", err)
	}

	// Pull latest changes
	_, err = f.git.Run("pull", "origin", featureBranchName)
	if err != nil {
		output.Warning("Failed to pull latest changes: %v", err)
	}

	output.Success("Started working on feature: %s", featureBranchName)
	return nil
}

// finishFeature finishes a feature branch
func (f *FeatureHandler) finishFeature(featureName string) error {
	// Get current branch if not specified
	if featureName == "" {
		currentBranch, err := f.git.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		featureName = currentBranch
	}

	// Validate it's a feature branch
	if !strings.HasPrefix(featureName, "feature/") {
		return fmt.Errorf("branch '%s' is not a feature branch", featureName)
	}

	// Get develop branch
	developBranch, err := f.git.GetConfig("at.develop")
	if err != nil || developBranch == "" {
		developBranch = "develop" // default fallback
	}

	// Check if we have uncommitted changes
	_, err = f.git.Run("diff", "--quiet")
	if err != nil {
		output.Warning("You have uncommitted changes. Please commit or stash them first.")
		return nil
	}

	// Confirm finishing
	var confirmed bool
	err = huh.NewConfirm().
		Title("Finish Feature").
		Description(fmt.Sprintf("Finish feature branch '%s' and merge to %s?", featureName, developBranch)).
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("Feature finishing cancelled")
		return nil
	}

	// Switch to develop branch
	_, err = f.git.Run("checkout", developBranch)
	if err != nil {
		return fmt.Errorf("failed to switch to develop branch: %w", err)
	}

	// Pull latest changes
	_, err = f.git.Run("pull", "origin", developBranch)
	if err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	// Merge feature branch
	_, err = f.git.Run("merge", "--no-ff", featureName)
	if err != nil {
		return fmt.Errorf("failed to merge feature branch: %w", err)
	}

	// Delete feature branch
	_, err = f.git.Run("branch", "-d", featureName)
	if err != nil {
		output.Warning("Failed to delete feature branch: %v", err)
	}

	output.Success("Feature '%s' finished successfully", featureName)
	output.Info("Merged to: %s", developBranch)

	return nil
}

// listFeatures lists all feature branches
func (f *FeatureHandler) listFeatures() error {
	output.Title("✨ Feature Branches")

	// Get all branches
	branches, err := f.git.Run("branch", "--list", "feature/*")
	if err != nil {
		return fmt.Errorf("failed to get feature branches: %w", err)
	}

	if strings.TrimSpace(branches) == "" {
		output.Info("No feature branches found")
		return nil
	}

	branchList := strings.Split(strings.TrimSpace(branches), "\n")
	output.Info("Found %d feature branches:", len(branchList))

	for i, branch := range branchList {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			// Get branch description
			description, _ := f.git.GetConfig(fmt.Sprintf("branch.%s.description", branch))
			if description != "" {
				output.Info("  %d. %s - %s", i+1, branch, description)
			} else {
				output.Info("  %d. %s", i+1, branch)
			}
		}
	}

	return nil
}

// showUsage displays the feature command usage
func (f *FeatureHandler) showUsage() error {
	usage := "# Feature Command\n\n" +
		"Manages feature branches for new functionality development.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ feature [action] [name] [description]\n" +
		"git @ feature create [name] [description]\n" +
		"git @ feature start [name]\n" +
		"git @ feature finish [name]\n" +
		"git @ feature list\n" +
		"```\n\n" +
		"## Actions\n\n" +
		"- **create** (default): Create a new feature branch\n" +
		"- **start**: Start working on an existing feature\n" +
		"- **finish**: Finish and merge a feature branch\n" +
		"- **list**: List all feature branches\n\n" +
		"## Arguments\n\n" +
		"- **name**: Name of the feature (e.g., user-authentication)\n" +
		"- **description**: Description of what the feature implements\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# Create a feature with interactive prompts\n" +
		"git @ feature\n\n" +
		"# Create with name and description\n" +
		"git @ feature create user-auth \"Add user authentication system\"\n\n" +
		"# Start working on existing feature\n" +
		"git @ feature start user-auth\n\n" +
		"# Finish current feature branch\n" +
		"git @ feature finish\n\n" +
		"# Finish specific feature branch\n" +
		"git @ feature finish feature/user-auth\n\n" +
		"# List all feature branches\n" +
		"git @ feature list\n" +
		"```\n\n" +
		"## Workflow\n\n" +
		"### Creating Features\n" +
		"1. **Validation**: Ensure on develop branch\n" +
		"2. **Sync**: Pull latest changes from develop\n" +
		"3. **Create**: Create feature branch from develop\n" +
		"4. **Describe**: Set branch description\n" +
		"5. **Confirm**: Show success message\n\n" +
		"### Starting Features\n" +
		"1. **Validation**: Check if feature branch exists\n" +
		"2. **Switch**: Switch to feature branch\n" +
		"3. **Sync**: Pull latest changes\n" +
		"4. **Confirm**: Show success message\n\n" +
		"### Finishing Features\n" +
		"1. **Validation**: Check for uncommitted changes\n" +
		"2. **Confirm**: Prompt for user confirmation\n" +
		"3. **Switch**: Switch to develop branch\n" +
		"4. **Sync**: Pull latest changes\n" +
		"5. **Merge**: Merge feature branch (no-ff)\n" +
		"6. **Cleanup**: Delete feature branch\n" +
		"7. **Result**: Show success message\n\n" +
		"### Listing Features\n" +
		"1. **Search**: Find all feature branches\n" +
		"2. **Display**: Show formatted list with descriptions\n" +
		"3. **Count**: Display total number of features\n\n" +
		"## Safety Features\n\n" +
		"- **Develop Validation**: Ensures features are created from develop\n" +
		"- **Change Check**: Warns about uncommitted changes\n" +
		"- **Confirmation**: Prompts before finishing features\n" +
		"- **No-FF Merge**: Preserves feature history\n\n" +
		"## Best Practices\n\n" +
		"- Use descriptive feature names\n" +
		"- Keep features focused on single functionality\n" +
		"- Test thoroughly before finishing\n" +
		"- Use descriptive branch descriptions\n" +
		"- Finish features promptly when complete\n"
	return output.Markdown(usage)
}
