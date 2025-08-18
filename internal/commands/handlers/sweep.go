package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// SweepHandler handles sweep-related commands
type SweepHandler struct {
	BaseHandler
}

// NewSweepHandler creates a new sweep handler
func NewSweepHandler(cfg *config.Config, gitRepo *git.Repository) *SweepHandler {
	return &SweepHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the sweep command
func (s *SweepHandler) Execute(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return s.showUsage()
	}

	// Parse options
	dryRun := false
	force := false
	mergedOnly := false
	remoteDeletedOnly := false

	for _, arg := range args {
		switch arg {
		case "--dry-run":
			dryRun = true
		case "--force":
			force = true
		case "--merged":
			mergedOnly = true
		case "--remote-deleted":
			remoteDeletedOnly = true
		}
	}

	// Get branches to sweep
	branchesToDelete, err := s.getBranchesToSweep(mergedOnly, remoteDeletedOnly)
	if err != nil {
		return fmt.Errorf("failed to analyze branches: %w", err)
	}

	if len(branchesToDelete) == 0 {
		output.Info("No branches to sweep")
		return nil
	}

	// Show preview
	output.Title("🧹 Branch Sweep Preview")
	output.Info("Found %d branches to delete:", len(branchesToDelete))

	for _, branch := range branchesToDelete {
		output.Info("  - %s", branch)
	}

	if dryRun {
		output.Success("Dry run completed - no branches were deleted")
		return nil
	}

	// Confirm deletion
	if !force {
		var confirmed bool
		err := huh.NewConfirm().
			Title("Confirm Deletion").
			Description(fmt.Sprintf("Delete %d branches?", len(branchesToDelete))).
			Value(&confirmed).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !confirmed {
			output.Info("Sweep cancelled")
			return nil
		}
	}

	// Delete branches
	deletedCount := 0
	for _, branch := range branchesToDelete {
		_, err := s.git.Run("branch", "-D", branch)
		if err != nil {
			output.Warning("Failed to delete branch %s: %v", branch, err)
		} else {
			deletedCount++
			output.Info("Deleted branch: %s", branch)
		}
	}

	output.Success("Sweep completed - deleted %d branches", deletedCount)
	return nil
}

// getBranchesToSweep analyzes and returns branches that should be deleted
func (s *SweepHandler) getBranchesToSweep(mergedOnly, remoteDeletedOnly bool) ([]string, error) {
	var branchesToDelete []string

	// Get current branch
	currentBranch, err := s.git.GetCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	// Get all local branches
	branches, err := s.git.Run("branch", "--format=%(refname:short)")
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	branchList := strings.Split(strings.TrimSpace(branches), "\n")

	for _, branch := range branchList {
		branch = strings.TrimSpace(branch)
		if branch == "" || branch == currentBranch {
			continue
		}

		// Skip protected branches
		if s.isProtectedBranch(branch) {
			continue
		}

		shouldDelete := false

		// Check if merged
		if !remoteDeletedOnly {
			_, err := s.git.Run("branch", "--merged", branch)
			if err == nil {
				shouldDelete = true
			}
		}

		// Check if remote deleted
		if !mergedOnly && !shouldDelete {
			_, err := s.git.Run("ls-remote", "--heads", "origin", branch)
			if err != nil {
				shouldDelete = true
			}
		}

		if shouldDelete {
			branchesToDelete = append(branchesToDelete, branch)
		}
	}

	return branchesToDelete, nil
}

// isProtectedBranch checks if a branch should be protected from deletion
func (s *SweepHandler) isProtectedBranch(branch string) bool {
	protectedBranches := []string{"master", "main", "develop", "trunk"}
	for _, protected := range protectedBranches {
		if branch == protected {
			return true
		}
	}
	return false
}

// showUsage displays the sweep command usage
func (s *SweepHandler) showUsage() error {
	usage := "# Sweep Command\n\n" +
		"Cleans up local branches by removing merged branches and branches that have been deleted on the remote.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ sweep [options]\n" +
		"git @ sweep --dry-run\n" +
		"git @ sweep --force\n" +
		"```\n\n" +
		"## Options\n\n" +
		"- **--dry-run**: Preview what would be deleted without actually deleting\n" +
		"- **--force**: Force deletion without confirmation prompts\n" +
		"- **--merged**: Only delete branches that have been merged\n" +
		"- **--remote-deleted**: Only delete branches that no longer exist on remote\n" +
		"- **-h, --help**: Show this help message\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# Preview what would be swept\n" +
		"git @ sweep --dry-run\n\n" +
		"# Sweep merged branches only\n" +
		"git @ sweep --merged\n\n" +
		"# Force sweep without confirmation\n" +
		"git @ sweep --force\n\n" +
		"# Sweep all (merged + remote-deleted)\n" +
		"git @ sweep\n" +
		"```\n\n" +
		"## What Gets Swept\n\n" +
		"### Merged Branches\n" +
		"- Local branches that have been merged into the current branch\n" +
		"- Feature branches that are no longer needed\n" +
		"- Hotfix branches that have been deployed\n\n" +
		"### Remote-Deleted Branches\n" +
		"- Local branches whose remote counterparts have been deleted\n" +
		"- Branches that were merged and deleted on remote\n" +
		"- Stale feature branches\n\n" +
		"## Safety Features\n\n" +
		"- **Dry Run**: Always preview changes before applying\n" +
		"- **Confirmation**: Prompts for confirmation before deletion\n" +
		"- **Protected Branches**: Never deletes trunk branches (master, main, develop)\n" +
		"- **Current Branch**: Never deletes the currently checked out branch\n\n" +
		"## Workflow\n\n" +
		"1. **Analysis**: Scans for branches to clean up\n" +
		"2. **Preview**: Shows what would be deleted (with --dry-run)\n" +
		"3. **Confirmation**: Prompts for user confirmation\n" +
		"4. **Cleanup**: Removes the selected branches\n" +
		"5. **Report**: Shows summary of actions taken\n\n" +
		"## Best Practices\n\n" +
		"- Always run with `--dry-run` first to preview changes\n" +
		"- Review the list of branches before confirming deletion\n" +
		"- Keep your local repository clean by running sweep regularly\n" +
		"- Use `--merged` if you only want to clean up completed features\n"
	return output.Markdown(usage)
}
