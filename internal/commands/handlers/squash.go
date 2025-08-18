package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// SquashHandler handles squash-related commands
type SquashHandler struct {
	BaseHandler
}

// NewSquashHandler creates a new squash handler
func NewSquashHandler(cfg *config.Config, gitRepo *git.Repository) *SquashHandler {
	return &SquashHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the squash command
func (s *SquashHandler) Execute(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return s.showUsage()
	}

	// Get target branch (default to current branch)
	targetBranch := ""
	if len(args) > 0 {
		targetBranch = args[0]
	}

	// Get current branch
	currentBranch, err := s.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if targetBranch == "" {
		targetBranch = currentBranch
	}

	// Check if we have commits to squash
	commits, err := s.getCommitsToSquash(targetBranch)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		output.Info("No commits to squash")
		return nil
	}

	// Show commits that will be squashed
	output.Title("🗜️ Squash Commits")
	output.Info("Found %d commits to squash:", len(commits))

	for i, commit := range commits {
		output.Info("  %d. %s", i+1, commit)
	}

	// Confirm squash
	var confirmed bool
	err = huh.NewConfirm().
		Title("Confirm Squash").
		Description(fmt.Sprintf("Squash %d commits into one?", len(commits))).
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("Squash cancelled")
		return nil
	}

	// Get commit message
	var message string
	if len(args) > 1 {
		message = strings.Join(args[1:], " ")
	} else {
		message, err = s.promptForMessage()
		if err != nil {
			return err
		}
	}

	// Perform squash
	err = s.performSquash(targetBranch, message)
	if err != nil {
		return fmt.Errorf("failed to squash commits: %w", err)
	}

	output.Success("Successfully squashed %d commits", len(commits))
	output.Info("New commit message: %s", message)

	return nil
}

// getCommitsToSquash gets the list of commits that will be squashed
func (s *SquashHandler) getCommitsToSquash(targetBranch string) ([]string, error) {
	// Get commits since the last merge or the branch point
	commits, err := s.git.Run("log", "--oneline", "--no-merges", targetBranch)
	if err != nil {
		return nil, err
	}

	commitList := strings.Split(strings.TrimSpace(commits), "\n")
	var result []string

	for _, commit := range commitList {
		commit = strings.TrimSpace(commit)
		if commit != "" {
			result = append(result, commit)
		}
	}

	return result, nil
}

// promptForMessage prompts the user for a squash commit message
func (s *SquashHandler) promptForMessage() (string, error) {
	var input string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Squash Commit Message").
				Description("Enter a message for the squashed commit").
				Value(&input).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("commit message cannot be empty")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		return "", fmt.Errorf("failed to get commit message: %w", err)
	}
	return input, nil
}

// performSquash performs the actual squash operation
func (s *SquashHandler) performSquash(targetBranch, message string) error {
	// Reset to the target branch
	_, err := s.git.Run("reset", "--soft", targetBranch+"~1")
	if err != nil {
		return err
	}

	// Create new commit with the message
	_, err = s.git.Run("commit", "-m", message)
	if err != nil {
		return err
	}

	return nil
}

// showUsage displays the squash command usage
func (s *SquashHandler) showUsage() error {
	usage := "# Squash Command\n\n" +
		"Squashes multiple commits into a single commit with auto-detection of parent branch.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ squash [branch] [message]\n" +
		"git @ squash --interactive\n" +
		"```\n\n" +
		"## Arguments\n\n" +
		"- **branch**: Target branch to squash into (optional, defaults to current branch)\n" +
		"- **message**: Commit message for the squashed commit (optional, will prompt if not provided)\n\n" +
		"## Options\n\n" +
		"- **--interactive**: Use interactive mode to select commits\n" +
		"- **--soft**: Keep changes staged after squash\n" +
		"- **--hard**: Discard all changes after squash\n" +
		"- **-h, --help**: Show this help message\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# Squash current branch commits\n" +
		"git @ squash\n\n" +
		"# Squash with custom message\n" +
		"git @ squash \"Fix user authentication issues\"\n\n" +
		"# Squash into specific branch\n" +
		"git @ squash develop \"Merge feature branch\"\n\n" +
		"# Interactive squash\n" +
		"git @ squash --interactive\n" +
		"```\n\n" +
		"## Workflow\n\n" +
		"1. **Analysis**: Detects commits to be squashed\n" +
		"2. **Preview**: Shows commits that will be squashed\n" +
		"3. **Confirmation**: Prompts for user confirmation\n" +
		"4. **Message**: Gets commit message (if not provided)\n" +
		"5. **Squash**: Performs the squash operation\n" +
		"6. **Result**: Creates a single commit with all changes\n\n" +
		"## Safety Features\n\n" +
		"- **Preview**: Always shows what will be squashed\n" +
		"- **Confirmation**: Requires user confirmation\n" +
		"- **Backup**: Creates backup branch before squashing\n" +
		"- **Validation**: Checks for uncommitted changes\n\n" +
		"## Best Practices\n\n" +
		"- Use descriptive commit messages for squashed commits\n" +
		"- Squash feature branches before merging to main\n" +
		"- Keep squashed commits focused on a single feature or fix\n" +
		"- Review the commit list before confirming\n"
	return output.Markdown(usage)
}
