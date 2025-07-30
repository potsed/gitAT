package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// WIPHandler handles WIP-related commands
type WIPHandler struct {
	BaseHandler
}

// NewWIPHandler creates a new WIP handler
func NewWIPHandler(cfg *config.Config, gitRepo *git.Repository) *WIPHandler {
	return &WIPHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the WIP command
func (w *WIPHandler) Execute(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return w.showUsage()
	}

	// Parse options
	action := "save" // default action
	message := ""

	for i, arg := range args {
		switch arg {
		case "save", "commit":
			action = "save"
			if i+1 < len(args) {
				message = strings.Join(args[i+1:], " ")
			}
		case "list", "show":
			action = "list"
		case "restore", "apply":
			action = "restore"
			if i+1 < len(args) {
				message = args[i+1]
			}
		case "clear", "clean":
			action = "clear"
		}
	}

	switch action {
	case "save":
		return w.saveWIP(message)
	case "list":
		return w.listWIP()
	case "restore":
		return w.restoreWIP(message)
	case "clear":
		return w.clearWIP()
	default:
		return w.saveWIP(message)
	}
}

// saveWIP saves current work in progress
func (w *WIPHandler) saveWIP(message string) error {
	// Check if there are changes to save
	_, err := w.git.Run("diff", "--quiet")
	if err == nil {
		output.Warning("No changes to save")
		return nil
	}

	// Generate WIP message if not provided
	if message == "" {
		message = fmt.Sprintf("WIP: %s", time.Now().Format("2006-01-02 15:04:05"))
	} else {
		message = fmt.Sprintf("WIP: %s", message)
	}

	// Stage all changes
	_, err = w.git.Run("add", "-A")
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Create WIP commit
	_, err = w.git.Run("commit", "-m", message)
	if err != nil {
		return fmt.Errorf("failed to create WIP commit: %w", err)
	}

	// Get current branch
	currentBranch, err := w.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	output.Success("Work in progress saved")
	output.Info("Branch: %s", currentBranch)
	output.Info("Message: %s", message)

	return nil
}

// listWIP lists all WIP commits
func (w *WIPHandler) listWIP() error {
	output.Title("📝 WIP Commits")

	// Get WIP commits
	commits, err := w.git.Run("log", "--oneline", "--grep=^WIP:")
	if err != nil {
		return fmt.Errorf("failed to get WIP commits: %w", err)
	}

	if strings.TrimSpace(commits) == "" {
		output.Info("No WIP commits found")
		return nil
	}

	commitList := strings.Split(strings.TrimSpace(commits), "\n")
	output.Info("Found %d WIP commits:", len(commitList))

	for i, commit := range commitList {
		commit = strings.TrimSpace(commit)
		if commit != "" {
			output.Info("  %d. %s", i+1, commit)
		}
	}

	return nil
}

// restoreWIP restores a WIP commit
func (w *WIPHandler) restoreWIP(commitRef string) error {
	if commitRef == "" {
		// Show available WIP commits and let user choose
		commits, err := w.git.Run("log", "--oneline", "--grep=^WIP:")
		if err != nil {
			return fmt.Errorf("failed to get WIP commits: %w", err)
		}

		if strings.TrimSpace(commits) == "" {
			output.Info("No WIP commits found")
			return nil
		}

		commitList := strings.Split(strings.TrimSpace(commits), "\n")
		var options []string
		for _, commit := range commitList {
			commit = strings.TrimSpace(commit)
			if commit != "" {
				options = append(options, commit)
			}
		}

		var selected string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select WIP Commit").
					Description("Choose which WIP commit to restore").
					Options(huh.NewOptions(options...)...).
					Value(&selected),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to select WIP commit: %w", err)
		}

		// Extract commit hash from selection
		parts := strings.Fields(selected)
		if len(parts) > 0 {
			commitRef = parts[0]
		}
	}

	// Confirm restoration
	var confirmed bool
	err := huh.NewConfirm().
		Title("Confirm Restoration").
		Description(fmt.Sprintf("Restore WIP commit %s?", commitRef)).
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("WIP restoration cancelled")
		return nil
	}

	// Restore the WIP commit
	_, err = w.git.Run("cherry-pick", "--no-commit", commitRef)
	if err != nil {
		return fmt.Errorf("failed to restore WIP commit: %w", err)
	}

	output.Success("WIP commit restored successfully")
	output.Info("Commit: %s", commitRef)

	return nil
}

// clearWIP clears all WIP commits
func (w *WIPHandler) clearWIP() error {
	// Get WIP commits
	commits, err := w.git.Run("log", "--oneline", "--grep=^WIP:")
	if err != nil {
		return fmt.Errorf("failed to get WIP commits: %w", err)
	}

	if strings.TrimSpace(commits) == "" {
		output.Info("No WIP commits to clear")
		return nil
	}

	commitList := strings.Split(strings.TrimSpace(commits), "\n")
	wipCount := len(commitList)

	// Confirm clearing
	var confirmed bool
	err = huh.NewConfirm().
		Title("Confirm Clear").
		Description(fmt.Sprintf("Clear %d WIP commits?", wipCount)).
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("WIP clear cancelled")
		return nil
	}

	// Clear WIP commits by resetting to the last non-WIP commit
	_, err = w.git.Run("reset", "--hard", "HEAD~"+fmt.Sprintf("%d", wipCount))
	if err != nil {
		return fmt.Errorf("failed to clear WIP commits: %w", err)
	}

	output.Success("Cleared %d WIP commits", wipCount)
	return nil
}

// showUsage displays the WIP command usage
func (w *WIPHandler) showUsage() error {
	usage := "# WIP Command\n\n" +
		"Manages Work In Progress commits for temporary saving and restoration.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ wip [action] [options]\n" +
		"git @ wip save [message]\n" +
		"git @ wip list\n" +
		"git @ wip restore [commit]\n" +
		"git @ wip clear\n" +
		"```\n\n" +
		"## Actions\n\n" +
		"- **save** (default): Save current work as WIP commit\n" +
		"- **list**: List all WIP commits\n" +
		"- **restore**: Restore a WIP commit\n" +
		"- **clear**: Clear all WIP commits\n\n" +
		"## Arguments\n\n" +
		"- **message**: Custom message for WIP commit\n" +
		"- **commit**: Commit hash or reference to restore\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# Save current work with auto-generated message\n" +
		"git @ wip\n\n" +
		"# Save with custom message\n" +
		"git @ wip save \"Working on user auth\"\n\n" +
		"# List all WIP commits\n" +
		"git @ wip list\n\n" +
		"# Restore specific WIP commit\n" +
		"git @ wip restore abc1234\n\n" +
		"# Restore with interactive selection\n" +
		"git @ wip restore\n\n" +
		"# Clear all WIP commits\n" +
		"git @ wip clear\n" +
		"```\n\n" +
		"## Workflow\n\n" +
		"### Saving WIP\n" +
		"1. **Check**: Verify there are changes to save\n" +
		"2. **Stage**: Add all changes to staging\n" +
		"3. **Commit**: Create WIP commit with timestamp or custom message\n" +
		"4. **Confirm**: Show success message with details\n\n" +
		"### Listing WIP\n" +
		"1. **Search**: Find all commits with WIP prefix\n" +
		"2. **Display**: Show formatted list of WIP commits\n" +
		"3. **Count**: Display total number of WIP commits\n\n" +
		"### Restoring WIP\n" +
		"1. **Select**: Choose WIP commit to restore (interactive if not specified)\n" +
		"2. **Confirm**: Prompt for user confirmation\n" +
		"3. **Apply**: Cherry-pick the WIP commit\n" +
		"4. **Result**: Show success message\n\n" +
		"### Clearing WIP\n" +
		"1. **Count**: Determine number of WIP commits\n" +
		"2. **Confirm**: Prompt for user confirmation\n" +
		"3. **Reset**: Hard reset to remove WIP commits\n" +
		"4. **Result**: Show success message\n\n" +
		"## Safety Features\n\n" +
		"- **Confirmation**: Prompts before destructive operations\n" +
		"- **Validation**: Checks for changes before saving\n" +
		"- **Interactive**: User-friendly selection for restoration\n" +
		"- **Backup**: WIP commits are regular Git commits\n\n" +
		"## Best Practices\n\n" +
		"- Use WIP commits for temporary work preservation\n" +
		"- Clear WIP commits when work is complete\n" +
		"- Use descriptive messages for easier identification\n" +
		"- Don't rely on WIP for long-term storage\n"
	return output.Markdown(usage)
}
