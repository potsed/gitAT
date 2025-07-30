package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// ReleaseHandler handles release-related commands
type ReleaseHandler struct {
	BaseHandler
}

// NewReleaseHandler creates a new release handler
func NewReleaseHandler(cfg *config.Config, gitRepo *git.Repository) *ReleaseHandler {
	return &ReleaseHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the release command
func (r *ReleaseHandler) Execute(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return r.showUsage()
	}

	// Parse options
	action := "create" // default action
	releaseName := ""
	version := ""

	for i, arg := range args {
		switch arg {
		case "create", "new":
			action = "create"
			if i+1 < len(args) {
				releaseName = args[i+1]
			}
		case "finish", "complete":
			action = "finish"
			if i+1 < len(args) {
				releaseName = args[i+1]
			}
		case "list", "show":
			action = "list"
		case "version":
			action = "version"
			if i+1 < len(args) {
				version = args[i+1]
			}
		}
	}

	switch action {
	case "create":
		return r.createRelease(releaseName, version)
	case "finish":
		return r.finishRelease(releaseName)
	case "list":
		return r.listReleases()
	case "version":
		return r.setVersion(version)
	default:
		return r.createRelease(releaseName, version)
	}
}

// createRelease creates a new release branch
func (r *ReleaseHandler) createRelease(releaseName, version string) error {
	// Get develop branch
	developBranch, err := r.git.GetConfig("at.develop")
	if err != nil || developBranch == "" {
		developBranch = "develop" // default fallback
	}

	// Generate release name if not provided
	if releaseName == "" {
		var input string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Release Name").
					Description("Enter a name for the release (e.g., v1.2.0)").
					Value(&input).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("release name cannot be empty")
						}
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get release name: %w", err)
		}
		releaseName = input
	}

	// Get version if not provided
	if version == "" {
		var input string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Version").
					Description("Enter the version number (e.g., 1.2.0)").
					Value(&input).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("version cannot be empty")
						}
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get version: %w", err)
		}
		version = input
	}

	// Check if we're on develop branch
	currentBranch, err := r.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if currentBranch != developBranch {
		output.Warning("You're not on the develop branch (%s). Switching to %s...", developBranch, developBranch)
		_, err = r.git.Run("checkout", developBranch)
		if err != nil {
			return fmt.Errorf("failed to switch to develop branch: %w", err)
		}
	}

	// Pull latest changes
	output.Info("Pulling latest changes from %s...", developBranch)
	_, err = r.git.Run("pull", "origin", developBranch)
	if err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	// Create release branch
	releaseBranchName := fmt.Sprintf("release/%s", releaseName)
	_, err = r.git.Run("checkout", "-b", releaseBranchName)
	if err != nil {
		return fmt.Errorf("failed to create release branch: %w", err)
	}

	// Set version in configuration
	err = r.git.SetConfig("at.version", version)
	if err != nil {
		output.Warning("Failed to set version: %v", err)
	}

	output.Success("Created release branch: %s", releaseBranchName)
	output.Info("Version: %s", version)
	output.Info("Base branch: %s", developBranch)

	return nil
}

// finishRelease finishes a release branch
func (r *ReleaseHandler) finishRelease(releaseName string) error {
	// Get current branch if not specified
	if releaseName == "" {
		currentBranch, err := r.git.GetCurrentBranch()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		releaseName = currentBranch
	}

	// Validate it's a release branch
	if !strings.HasPrefix(releaseName, "release/") {
		return fmt.Errorf("branch '%s' is not a release branch", releaseName)
	}

	// Get trunk and develop branches
	trunkBranch, err := r.git.GetConfig("at.trunk")
	if err != nil || trunkBranch == "" {
		trunkBranch = "main" // default fallback
	}

	developBranch, err := r.git.GetConfig("at.develop")
	if err != nil || developBranch == "" {
		developBranch = "develop" // default fallback
	}

	// Check if we have uncommitted changes
	_, err = r.git.Run("diff", "--quiet")
	if err != nil {
		output.Warning("You have uncommitted changes. Please commit or stash them first.")
		return nil
	}

	// Get version
	version, _ := r.git.GetConfig("at.version")
	if version == "" {
		version = "unknown"
	}

	// Confirm finishing
	var confirmed bool
	err = huh.NewConfirm().
		Title("Finish Release").
		Description(fmt.Sprintf("Finish release branch '%s' (v%s) and merge to %s and %s?", releaseName, version, trunkBranch, developBranch)).
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("Release finishing cancelled")
		return nil
	}

	// Merge to trunk
	_, err = r.git.Run("checkout", trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to switch to trunk branch: %w", err)
	}

	_, err = r.git.Run("pull", "origin", trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	_, err = r.git.Run("merge", "--no-ff", releaseName)
	if err != nil {
		return fmt.Errorf("failed to merge to trunk: %w", err)
	}

	// Tag the release
	tagName := fmt.Sprintf("v%s", version)
	_, err = r.git.Run("tag", "-a", tagName, "-m", fmt.Sprintf("Release %s", tagName))
	if err != nil {
		output.Warning("Failed to create tag: %v", err)
	}

	// Merge to develop
	_, err = r.git.Run("checkout", developBranch)
	if err != nil {
		return fmt.Errorf("failed to switch to develop branch: %w", err)
	}

	_, err = r.git.Run("pull", "origin", developBranch)
	if err != nil {
		return fmt.Errorf("failed to pull latest changes: %w", err)
	}

	_, err = r.git.Run("merge", "--no-ff", releaseName)
	if err != nil {
		return fmt.Errorf("failed to merge to develop: %w", err)
	}

	// Delete release branch
	_, err = r.git.Run("branch", "-d", releaseName)
	if err != nil {
		output.Warning("Failed to delete release branch: %v", err)
	}

	output.Success("Release '%s' finished successfully", releaseName)
	output.Info("Version: %s", version)
	output.Info("Tagged as: %s", tagName)
	output.Info("Merged to: %s and %s", trunkBranch, developBranch)

	return nil
}

// listReleases lists all release branches
func (r *ReleaseHandler) listReleases() error {
	output.Title("🚀 Release Branches")

	// Get all branches
	branches, err := r.git.Run("branch", "--list", "release/*")
	if err != nil {
		return fmt.Errorf("failed to get release branches: %w", err)
	}

	if strings.TrimSpace(branches) == "" {
		output.Info("No release branches found")
		return nil
	}

	branchList := strings.Split(strings.TrimSpace(branches), "\n")
	output.Info("Found %d release branches:", len(branchList))

	for i, branch := range branchList {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			// Get version for this branch
			version, _ := r.git.GetConfig("at.version")
			if version != "" {
				output.Info("  %d. %s (v%s)", i+1, branch, version)
			} else {
				output.Info("  %d. %s", i+1, branch)
			}
		}
	}

	return nil
}

// setVersion sets the version for the current release
func (r *ReleaseHandler) setVersion(version string) error {
	if version == "" {
		return fmt.Errorf("version required")
	}

	err := r.git.SetConfig("at.version", version)
	if err != nil {
		return fmt.Errorf("failed to set version: %w", err)
	}

	output.Success("Version set to: %s", version)
	return nil
}

// showUsage displays the release command usage
func (r *ReleaseHandler) showUsage() error {
	usage := "# Release Command\n\n" +
		"Manages release branches for preparing new versions of the software.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ release [action] [name] [version]\n" +
		"git @ release create [name] [version]\n" +
		"git @ release finish [name]\n" +
		"git @ release list\n" +
		"git @ release version [version]\n" +
		"```\n\n" +
		"## Actions\n\n" +
		"- **create** (default): Create a new release branch\n" +
		"- **finish**: Finish and merge a release branch\n" +
		"- **list**: List all release branches\n" +
		"- **version**: Set version for current release\n\n" +
		"## Arguments\n\n" +
		"- **name**: Name of the release (e.g., v1.2.0)\n" +
		"- **version**: Version number (e.g., 1.2.0)\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# Create a release with interactive prompts\n" +
		"git @ release\n\n" +
		"# Create with name and version\n" +
		"git @ release create v1.2.0 1.2.0\n\n" +
		"# Finish current release branch\n" +
		"git @ release finish\n\n" +
		"# Finish specific release branch\n" +
		"git @ release finish release/v1.2.0\n\n" +
		"# List all release branches\n" +
		"git @ release list\n\n" +
		"# Set version for current release\n" +
		"git @ release version 1.2.0\n" +
		"```\n\n" +
		"## Workflow\n\n" +
		"### Creating Releases\n" +
		"1. **Validation**: Ensure on develop branch\n" +
		"2. **Sync**: Pull latest changes from develop\n" +
		"3. **Create**: Create release branch from develop\n" +
		"4. **Version**: Set version number\n" +
		"5. **Confirm**: Show success message\n\n" +
		"### Finishing Releases\n" +
		"1. **Validation**: Check for uncommitted changes\n" +
		"2. **Confirm**: Prompt for user confirmation\n" +
		"3. **Merge to Trunk**: Merge to main/master branch\n" +
		"4. **Tag**: Create version tag\n" +
		"5. **Merge to Develop**: Merge back to develop branch\n" +
		"6. **Cleanup**: Delete release branch\n" +
		"7. **Result**: Show success message\n\n" +
		"### Listing Releases\n" +
		"1. **Search**: Find all release branches\n" +
		"2. **Display**: Show formatted list with versions\n" +
		"3. **Count**: Display total number of releases\n\n" +
		"## Safety Features\n\n" +
		"- **Develop Validation**: Ensures releases are created from develop\n" +
		"- **Change Check**: Warns about uncommitted changes\n" +
		"- **Confirmation**: Prompts before finishing releases\n" +
		"- **No-FF Merge**: Preserves release history\n" +
		"- **Version Tagging**: Creates semantic version tags\n\n" +
		"## Best Practices\n\n" +
		"- Use semantic versioning (MAJOR.MINOR.PATCH)\n" +
		"- Test thoroughly before finishing releases\n" +
		"- Update version numbers consistently\n" +
		"- Create meaningful release tags\n" +
		"- Document release notes\n"
	return output.Markdown(usage)
}
