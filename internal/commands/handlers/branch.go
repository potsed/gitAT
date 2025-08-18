package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// BranchHandler handles branch-related commands
type BranchHandler struct {
	BaseHandler
}

// NewBranchHandler creates a new branch handler
func NewBranchHandler(cfg *config.Config, gitRepo *git.Repository) *BranchHandler {
	return &BranchHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the branch command
func (s *BranchHandler) Execute(args []string) error {
	// Check for help flags
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return s.showUsage()
	}

	// If no arguments, show interactive form
	if len(args) == 0 {
		return s.showInteractiveForm()
	}

	// Parse options
	action := ""
	branchName := ""
	force := false
	remote := false

	for i, arg := range args {
		switch arg {
		case "create", "new":
			action = "create"
			if i+1 < len(args) {
				branchName = args[i+1]
			}
		case "delete", "del", "remove":
			action = "delete"
			if i+1 < len(args) {
				branchName = args[i+1]
			}
		case "switch", "checkout":
			action = "switch"
			if i+1 < len(args) {
				branchName = args[i+1]
			}
		case "list", "ls":
			action = "list"
		case "--force", "-f":
			force = true
		case "--remote", "-r":
			remote = true
		}
	}

	// If no action specified, show interactive form
	if action == "" {
		return s.showInteractiveForm()
	}

	switch action {
	case "list":
		return s.listBranches(remote)
	case "create":
		return s.createBranch(branchName, force)
	case "delete":
		return s.deleteBranch(branchName, force)
	case "switch":
		return s.switchBranch(branchName)
	default:
		return s.listBranches(remote)
	}
}

// listBranches lists all branches
func (s *BranchHandler) listBranches(remote bool) error {
	output.Title("🌿 Branch List")

	if remote {
		// List remote branches
		branches, err := s.git.Run("branch", "-r")
		if err != nil {
			return fmt.Errorf("failed to get remote branches: %w", err)
		}

		branchList := strings.Split(strings.TrimSpace(branches), "\n")
		output.Info("Remote branches:")
		for _, branch := range branchList {
			branch = strings.TrimSpace(branch)
			if branch != "" {
				output.Info("  %s", branch)
			}
		}
	} else {
		// List local branches
		branches, err := s.git.Run("branch")
		if err != nil {
			return fmt.Errorf("failed to get local branches: %w", err)
		}

		branchList := strings.Split(strings.TrimSpace(branches), "\n")
		output.Info("Local branches:")
		for _, branch := range branchList {
			branch = strings.TrimSpace(branch)
			if branch != "" {
				// Highlight current branch
				if strings.HasPrefix(branch, "* ") {
					output.Success("  %s (current)", strings.TrimPrefix(branch, "* "))
				} else {
					output.Info("  %s", branch)
				}
			}
		}
	}

	return nil
}

// createBranch creates a new branch
func (s *BranchHandler) createBranch(branchName string, force bool) error {
	if branchName == "" {
		var input string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Branch Name").
					Description("Enter the name for the new branch").
					Value(&input).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("branch name cannot be empty")
						}
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get branch name: %w", err)
		}
		branchName = input
	}

	// Check if branch already exists
	_, err := s.git.Run("show-ref", "--verify", "--quiet", "refs/heads/"+branchName)
	if err == nil && !force {
		return fmt.Errorf("branch '%s' already exists. Use --force to overwrite", branchName)
	}

	// Create branch
	args := []string{"checkout", "-b", branchName}
	if force {
		args = append(args, "--force")
	}

	_, err = s.git.Run(args...)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	output.Success("Created and switched to branch: %s", branchName)
	return nil
}

// deleteBranch deletes a branch
func (s *BranchHandler) deleteBranch(branchName string, force bool) error {
	if branchName == "" {
		return fmt.Errorf("branch name required for deletion")
	}

	// Get current branch
	currentBranch, err := s.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Don't delete current branch
	if branchName == currentBranch {
		return fmt.Errorf("cannot delete current branch '%s'. Switch to another branch first", branchName)
	}

	// Confirm deletion
	if !force {
		var confirmed bool
		err := huh.NewConfirm().
			Title("Confirm Deletion").
			Description(fmt.Sprintf("Delete branch '%s'?", branchName)).
			Value(&confirmed).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !confirmed {
			output.Info("Branch deletion cancelled")
			return nil
		}
	}

	// Delete branch
	args := []string{"branch"}
	if force {
		args = append(args, "-D")
	} else {
		args = append(args, "-d")
	}
	args = append(args, branchName)

	_, err = s.git.Run(args...)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	output.Success("Deleted branch: %s", branchName)
	return nil
}

// switchBranch switches to a different branch
func (s *BranchHandler) switchBranch(branchName string) error {
	if branchName == "" {
		return fmt.Errorf("branch name required for switching")
	}

	// Check if branch exists
	_, err := s.git.Run("show-ref", "--verify", "--quiet", "refs/heads/"+branchName)
	if err != nil {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}

	// Switch to branch
	_, err = s.git.Run("checkout", branchName)
	if err != nil {
		return fmt.Errorf("failed to switch to branch: %w", err)
	}

	output.Success("Switched to branch: %s", branchName)
	return nil
}

// showUsage displays the branch command usage
func (s *BranchHandler) showUsage() error {
	usage := "# Branch Command\n\n" +
		"Manages Git branches with interactive creation, deletion, and switching.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ branch [action] [branch-name] [options]\n" +
		"git @ branch create [name]\n" +
		"git @ branch delete [name]\n" +
		"git @ branch switch [name]\n" +
		"```\n\n" +
		"## Actions\n\n" +
		"- **list** (default): List all branches\n" +
		"- **create**: Create a new branch\n" +
		"- **delete**: Delete a branch\n" +
		"- **switch**: Switch to a different branch\n\n" +
		"## Arguments\n\n" +
		"- **branch-name**: Name of the branch to operate on\n\n" +
		"## Options\n\n" +
		"- **--force, -f**: Force the operation (overwrite existing, force delete)\n" +
		"- **--remote, -r**: Show remote branches when listing\n" +
		"- **-h, --help**: Show this help message\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# List all local branches\n" +
		"git @ branch\n\n" +
		"# List remote branches\n" +
		"git @ branch --remote\n\n" +
		"# Create a new branch\n" +
		"git @ branch create feature/new-feature\n\n" +
		"# Create with interactive name\n" +
		"git @ branch create\n\n" +
		"# Switch to a branch\n" +
		"git @ branch switch develop\n\n" +
		"# Delete a branch\n" +
		"git @ branch delete old-feature\n\n" +
		"# Force delete a branch\n" +
		"git @ branch delete old-feature --force\n" +
		"```\n\n" +
		"## Workflow\n\n" +
		"### Creating Branches\n" +
		"1. **Name**: Provide branch name or use interactive prompt\n" +
		"2. **Validation**: Check if branch already exists\n" +
		"3. **Creation**: Create and switch to new branch\n" +
		"4. **Confirmation**: Show success message\n\n" +
		"### Deleting Branches\n" +
		"1. **Validation**: Ensure not deleting current branch\n" +
		"2. **Confirmation**: Prompt for user confirmation\n" +
		"3. **Deletion**: Remove the branch\n" +
		"4. **Result**: Show success message\n\n" +
		"### Switching Branches\n" +
		"1. **Validation**: Check if target branch exists\n" +
		"2. **Switch**: Change to target branch\n" +
		"3. **Confirmation**: Show success message\n\n" +
		"## Safety Features\n\n" +
		"- **Current Branch Protection**: Cannot delete currently checked out branch\n" +
		"- **Confirmation**: Prompts before destructive operations\n" +
		"- **Validation**: Checks branch existence before operations\n" +
		"- **Force Option**: Explicit flag required for dangerous operations\n\n" +
		"## Best Practices\n\n" +
		"- Use descriptive branch names (feature/, bugfix/, hotfix/)\n" +
		"- Delete branches after merging to keep repository clean\n" +
		"- Use force delete only when absolutely necessary\n" +
		"- Switch to a safe branch before deleting others\n"
	return output.Markdown(usage)
}

// showInteractiveForm shows an interactive form for branch operations
func (s *BranchHandler) showInteractiveForm() error {
	var action string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Branch Action").
				Description("What would you like to do?").
				Options(
					huh.NewOption("List branches", "list"),
					huh.NewOption("Create new branch", "create"),
					huh.NewOption("Switch to branch", "switch"),
					huh.NewOption("Delete branch", "delete"),
				).
				Value(&action),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to show form: %w", err)
	}

	switch action {
	case "list":
		return s.listBranches(false)
	case "create":
		return s.showCreateBranchForm()
	case "switch":
		return s.showSwitchBranchForm()
	case "delete":
		return s.showDeleteBranchForm()
	default:
		return s.listBranches(false)
	}
}

// showCreateBranchForm shows a form to create a new branch
func (s *BranchHandler) showCreateBranchForm() error {
	var branchName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Branch Name").
				Description("Enter the name for the new branch").
				Placeholder("e.g., feature/new-feature").
				Value(&branchName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("branch name cannot be empty")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to show form: %w", err)
	}

	return s.createBranch(branchName, false)
}

// showSwitchBranchForm shows a form to switch to a branch
func (s *BranchHandler) showSwitchBranchForm() error {
	// Get available branches
	branches, err := s.git.Run("branch", "--format=%(refname:short)")
	if err != nil {
		return fmt.Errorf("failed to get branches: %w", err)
	}

	branchList := strings.Split(strings.TrimSpace(branches), "\n")
	var branchOptions []huh.Option[string]

	for _, branch := range branchList {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			branchOptions = append(branchOptions, huh.NewOption(branch, branch))
		}
	}

	if len(branchOptions) == 0 {
		return fmt.Errorf("no branches available")
	}

	var selectedBranch string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Branch").
				Description("Choose a branch to switch to").
				Options(branchOptions...).
				Value(&selectedBranch),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to show form: %w", err)
	}

	return s.switchBranch(selectedBranch)
}

// showDeleteBranchForm shows a form to delete a branch
func (s *BranchHandler) showDeleteBranchForm() error {
	// Get available branches (excluding current)
	currentBranch, err := s.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	branches, err := s.git.Run("branch", "--format=%(refname:short)")
	if err != nil {
		return fmt.Errorf("failed to get branches: %w", err)
	}

	branchList := strings.Split(strings.TrimSpace(branches), "\n")
	var branchOptions []huh.Option[string]

	for _, branch := range branchList {
		branch = strings.TrimSpace(branch)
		if branch != "" && branch != currentBranch {
			branchOptions = append(branchOptions, huh.NewOption(branch, branch))
		}
	}

	if len(branchOptions) == 0 {
		return fmt.Errorf("no branches available to delete")
	}

	var selectedBranch string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Branch to Delete").
				Description("Choose a branch to delete (current branch excluded)").
				Options(branchOptions...).
				Value(&selectedBranch),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to show form: %w", err)
	}

	return s.deleteBranch(selectedBranch, false)
}
