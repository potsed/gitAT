package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// SaveHandler handles save-related commands
type SaveHandler struct {
	BaseHandler
}

// NewSaveHandler creates a new save handler
func NewSaveHandler(cfg *config.Config, gitRepo *git.Repository) *SaveHandler {
	return &SaveHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the save command
func (s *SaveHandler) Execute(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return s.showUsage()
	}

	var message string
	if len(args) > 0 {
		message = strings.Join(args, " ")
	}

	// Check if there are changes to save
	_, err := s.git.Run("diff", "--quiet")
	if err == nil {
		output.Warning("No changes to save")
		return nil
	}

	// Stage all changes
	_, err = s.git.Run("add", "-A")
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// If no message provided, prompt for one
	if message == "" {
		message, err = s.promptForMessage()
		if err != nil {
			return err
		}
	}

	// Create commit
	_, err = s.git.Run("commit", "-m", message)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Get current branch and show success
	return s.showSuccess(message)
}

// promptForMessage prompts the user for a commit message
func (s *SaveHandler) promptForMessage() (string, error) {
	var input string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Commit Message").
				Description("Enter a commit message for your changes").
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

// showSuccess displays success information
func (s *SaveHandler) showSuccess(message string) error {
	currentBranch, err := s.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	output.Success("Changes saved successfully")
	output.Info("Branch: %s", currentBranch)
	output.Info("Commit: %s", message)

	return nil
}

// showUsage displays the save command usage
func (s *SaveHandler) showUsage() error {
	usage := "# Save Command\n\n" +
		"Saves current work with a commit message.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ save [message]\n" +
		"git @ save -m \"commit message\"\n" +
		"```\n\n" +
		"## Arguments\n\n" +
		"- **message**: Commit message (optional, will prompt if not provided)\n\n" +
		"## Options\n\n" +
		"- **-m, --message**: Specify commit message directly\n" +
		"- **-h, --help**: Show this help message\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# Save with interactive message prompt\n" +
		"git @ save\n\n" +
		"# Save with direct message\n" +
		"git @ save \"Add user authentication feature\"\n\n" +
		"# Save with flag\n" +
		"git @ save -m \"Fix login bug\"\n" +
		"```\n\n" +
		"## Workflow\n\n" +
		"1. Stages all changes in the working directory\n" +
		"2. Prompts for commit message if not provided\n" +
		"3. Creates a commit with the message\n" +
		"4. Updates the working branch reference\n\n" +
		"## Commit Message Format\n\n" +
		"Follows conventional commit format:\n" +
		"- **feat**: New features\n" +
		"- **fix**: Bug fixes\n" +
		"- **docs**: Documentation changes\n" +
		"- **style**: Code style changes\n" +
		"- **refactor**: Code refactoring\n" +
		"- **test**: Test changes\n" +
		"- **chore**: Maintenance tasks\n"
	return output.Markdown(usage)
}
