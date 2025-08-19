// commitizen.go - Commitizen command handler for conventional commits
package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// CommitizenHandler handles commitizen-related commands for conventional commits
type CommitizenHandler struct {
	BaseHandler
}

// NewCommitizenHandler creates a new commitizen handler
func NewCommitizenHandler(cfg *config.Config, gitRepo *git.Repository) *CommitizenHandler {
	return &CommitizenHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the commitizen command
func (c *CommitizenHandler) Execute(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "commit", "c":
			return c.commit(args[1:])
		case "init":
			return c.initCommitizen()
		case "config":
			return c.configure()
		case "prompt":
			return c.promptCommit()
		case "-h", "--help", "help":
			return c.showUsage()
		default:
			return c.commit(args)
		}
	}
	
	return c.promptCommit()
}

// commit creates a commit with conventional commit format
func (c *CommitizenHandler) commit(args []string) error {
	// If message provided as args, use it directly
	if len(args) > 0 {
		message := strings.Join(args, " ")
		return c.createCommit(message)
	}
	
	// Otherwise prompt for conventional commit
	return c.promptCommit()
}

// promptCommit prompts user for conventional commit details
func (c *CommitizenHandler) promptCommit() error {
	var commitType, scope, subject, body, footer string
	var isBreaking bool

	// Define commit types
	commitTypes := []struct {
		Value string
		Desc  string
	}{
		{"feat", "A new feature"},
		{"fix", "A bug fix"},
		{"docs", "Documentation only changes"},
		{"style", "Changes that do not affect the meaning of the code"},
		{"refactor", "A code change that neither fixes a bug nor adds a feature"},
		{"perf", "A code change that improves performance"},
		{"test", "Adding missing tests or correcting existing tests"},
		{"build", "Changes that affect the build system or external dependencies"},
		{"ci", "Changes to our CI configuration files and scripts"},
		{"chore", "Other changes that don't modify src or test files"},
		{"revert", "Reverts a previous commit"},
	}

	// Create form for conventional commit
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the type of change that you're committing:").
				OptionsFunc(func() []huh.Option[string] {
					var options []huh.Option[string]
					for _, ct := range commitTypes {
						options = append(options, huh.NewOption(fmt.Sprintf("%s: %s", ct.Value, ct.Desc), ct.Value))
					}
					return options
				}, &commitTypes).
				Value(&commitType),
			
			huh.NewInput().
				Title("What is the scope of this change (e.g. component or file name): (press enter to skip)").
				Value(&scope),
			
			huh.NewInput().
				Title("Write a short, imperative tense description of the change:").
				Placeholder("feat: add 'graphiteWidth' option").
				Value(&subject).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("subject cannot be empty")
					}
					if len(s) > 72 {
						return fmt.Errorf("subject must be less than 72 characters")
					}
					return nil
				}),
			
			huh.NewConfirm().
				Title("Are there any breaking changes?").
				Value(&isBreaking),
			
			huh.NewText().
				Title("Provide a longer description of the change: (press enter to skip)").
				Placeholder("Provide a longer description of the change.").
				Value(&body),
			
			huh.NewText().
				Title("List any breaking changes or issues closed by this change: (press enter to skip)").
				Placeholder("BREAKING CHANGE: Close #123").
				Value(&footer),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get commit details: %w", err)
	}

	// Build conventional commit message
	message := c.buildConventionalCommit(commitType, scope, subject, body, footer, isBreaking)

	// Show preview
	output.Title("Commit Message Preview")
	output.Info("%s", message)

	// Confirm commit
	var proceed bool
	err := huh.NewConfirm().
		Title("Commit Changes?").
		Description("Do you want to commit with this message?").
		Value(&proceed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %w", err)
	}

	if !proceed {
		output.Info("Commit cancelled")
		return nil
	}

	// Create the commit
	return c.createCommit(message)
}

// buildConventionalCommit builds a conventional commit message
func (c *CommitizenHandler) buildConventionalCommit(commitType, scope, subject, body, footer string, isBreaking bool) string {
	var builder strings.Builder

	// Build header
	builder.WriteString(commitType)
	if scope != "" {
		builder.WriteString(fmt.Sprintf("(%s)", scope))
	}
	if isBreaking {
		builder.WriteString("!")
	}
	builder.WriteString(": ")
	builder.WriteString(subject)

	// Add body if present
	if body != "" {
		builder.WriteString("\n\n")
		builder.WriteString(body)
	}

	// Add footer if present
	if footer != "" {
		builder.WriteString("\n\n")
		builder.WriteString(footer)
	}

	return builder.String()
}

// createCommit creates a git commit with the message
func (c *CommitizenHandler) createCommit(message string) error {
	// Stage all changes
	output.Info("Staging all changes...")
	_, err := c.git.Run("add", "-A")
	if err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Create commit
	output.Info("Creating commit...")
	_, err = c.git.Run("commit", "-m", message)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Show success
	currentBranch, err := c.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	output.Success("Changes committed successfully!")
	output.Info("Branch: %s", currentBranch)
	output.Info("Message: %s", strings.Split(message, "\n")[0]) // First line only

	return nil
}

// initCommitizen initializes commitizen configuration
func (c *CommitizenHandler) initCommitizen() error {
	output.Title("Initializing Commitizen")
	
	// This would set up commitizen configuration
	// For now, just show what would be done
	output.Info("Would initialize commitizen configuration")
	output.Info("This would typically:")
	output.Info("  • Install commitizen dependencies")
	output.Info("  • Create .czrc configuration file")
	output.Info("  • Set up commit message hooks")
	
	return nil
}

// configure allows configuration of commitizen settings
func (c *CommitizenHandler) configure() error {
	output.Title("Configuring Commitizen")
	
	// This would allow configuration of commitizen settings
	// For now, just show what would be done
	output.Info("Would allow configuration of commitizen settings")
	output.Info("This would typically:")
	output.Info("  • Configure commit types")
	output.Info("  • Set default scopes")
	output.Info("  • Configure message templates")
	
	return nil
}

// showUsage displays the commitizen command usage
func (c *CommitizenHandler) showUsage() error {
	return output.Markdown(`# Commitizen Command

Create conventional commits with an interactive prompt.

## Usage

` + "```" + `bash
git @ commitizen [command] [options]
git @ commitizen commit [message]
git @ commitizen prompt
git @ commitizen init
git @ commitizen config
` + "```" + `

## Commands

• **commit, c [message]**: Create a commit (with prompt if no message)
• **prompt**: Interactive prompt for conventional commit (default)
• **init**: Initialize commitizen configuration
• **config**: Configure commitizen settings

## Options

• **-h, --help**: Show this help message

## Examples

` + "```" + `bash
# Interactive commit prompt
git @ commitizen

# Commit with direct message
git @ commitizen commit "feat: add new feature"

# Initialize commitizen
git @ commitizen init

# Configure settings
git @ commitizen config
` + "```" + `

## Conventional Commit Format

The interactive prompt guides you through creating commits that follow the conventional commit specification:

` + "```" + `
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
` + "```" + `

### Types

• **feat**: A new feature (MINOR version bump)
• **fix**: A bug fix (PATCH version bump)
• **docs**: Documentation only changes
• **style**: Changes that do not affect the meaning of the code
• **refactor**: A code change that neither fixes a bug nor adds a feature
• **perf**: A code change that improves performance
• **test**: Adding missing tests or correcting existing tests
• **build**: Changes that affect the build system or external dependencies
• **ci**: Changes to our CI configuration files and scripts
• **chore**: Other changes that don't modify src or test files
• **revert**: Reverts a previous commit

### Features

• **Interactive Prompt**: Guided process for creating conventional commits
• **Validation**: Ensures commit messages follow the specification
• **Breaking Changes**: Special handling for breaking changes
• **Scope Support**: Optional scope for more detailed categorization
• **Body and Footer**: Support for detailed descriptions and references

### Best Practices

• Use present tense ("add feature" not "added feature")
• Use imperative mood ("move cursor to..." not "moves cursor to...")
• Limit the subject line to 72 characters
• Reference issues and pull requests in the footer
• Mark breaking changes with "!" after the type/scope
`)
}