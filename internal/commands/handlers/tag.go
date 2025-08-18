package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// TagHandler handles tag-related commands
type TagHandler struct {
	BaseHandler
}

// NewTagHandler creates a new tag handler
func NewTagHandler(cfg *config.Config, gitRepo *git.Repository) *TagHandler {
	return &TagHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the tag command
func (t *TagHandler) Execute(args []string) error {
	// Check for help flags
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return t.showUsage()
	}

	// If no arguments, show interactive form
	if len(args) == 0 {
		return t.showInteractiveForm()
	}

	switch args[0] {
	case "create":
		return t.showInteractiveForm()
	case "list":
		return t.listTags()
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("tag name required for delete command")
		}
		return t.deleteTag(args[1])
	case "push":
		if len(args) >= 2 {
			return t.pushTag(args[1])
		}
		return t.pushAllTags()
	default:
		// If no subcommand, treat as tag name with message
		if len(args) >= 2 {
			tagName := args[0]
			message := strings.Join(args[1:], " ")
			return t.createTag(tagName, message)
		}
		return fmt.Errorf("tag message required when specifying tag name")
	}
}

// createTag creates a new git tag
func (t *TagHandler) createTag(tagName, message string) error {
	if tagName == "" {
		return fmt.Errorf("tag name cannot be empty")
	}

	if message == "" {
		return fmt.Errorf("tag message cannot be empty")
	}

	// Check if tag already exists
	exists, err := t.tagExists(tagName)
	if err != nil {
		return fmt.Errorf("failed to check if tag exists: %v", err)
	}

	if exists {
		var overwrite bool
		err := huh.NewConfirm().
			Title("Tag Already Exists").
			Description(fmt.Sprintf("Tag '%s' already exists. Overwrite?", tagName)).
			Value(&overwrite).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get user input: %v", err)
		}

		if !overwrite {
			output.Info("Tag creation cancelled")
			return nil
		}

		// Delete existing tag
		if _, err := t.git.Run("tag", "-d", tagName); err != nil {
			return fmt.Errorf("failed to delete existing tag: %v", err)
		}
	}

	// Create annotated tag
	_, err = t.git.Run("tag", "-a", tagName, "-m", message)
	if err != nil {
		return fmt.Errorf("failed to create tag: %v", err)
	}

	output.Success("Created tag: %s", tagName)
	output.Info("Message: %s", message)

	// Ask if user wants to push the tag
	var pushTag bool
	err = huh.NewConfirm().
		Title("Push Tag").
		Description("Push tag to remote repository?").
		Value(&pushTag).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if pushTag {
		return t.pushTag(tagName)
	}

	return nil
}

// createTagFromVersion creates a tag using the current version
func (t *TagHandler) createTagFromVersion(message string) error {
	version, err := t.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %v", err)
	}

	if version == "0.0.0" {
		return fmt.Errorf("no version set. Use 'git @ version' to set a version first")
	}

	return t.createTag(version, message)
}

// getCurrentVersion gets the current version from git config
func (t *TagHandler) getCurrentVersion() (string, error) {
	major, err := t.git.GetConfig("at.major")
	if err != nil || major == "" {
		major = "0"
	}

	minor, err := t.git.GetConfig("at.minor")
	if err != nil || minor == "" {
		minor = "0"
	}

	fix, err := t.git.GetConfig("at.fix")
	if err != nil || fix == "" {
		fix = "0"
	}

	return fmt.Sprintf("%s.%s.%s", major, minor, fix), nil
}

// tagExists checks if a tag already exists
func (t *TagHandler) tagExists(tagName string) (bool, error) {
	output, err := t.git.Run("tag", "-l", tagName)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(output) == tagName, nil
}

// listTags lists all tags
func (t *TagHandler) listTags() error {
	output.Title("📋 Git Tags")

	tags, err := t.git.Run("tag", "-l", "--sort=-version:refname")
	if err != nil {
		return fmt.Errorf("failed to list tags: %v", err)
	}

	if strings.TrimSpace(tags) == "" {
		output.Info("No tags found")
		return nil
	}

	tagList := strings.Split(strings.TrimSpace(tags), "\n")
	for _, tag := range tagList {
		if tag != "" {
			// Get tag message if it's an annotated tag
			message, err := t.git.Run("tag", "-l", "--format=%(contents:subject)", tag)
			if err == nil && strings.TrimSpace(message) != "" {
				output.Info("🏷️  %s - %s", tag, strings.TrimSpace(message))
			} else {
				output.Info("🏷️  %s", tag)
			}
		}
	}

	return nil
}

// deleteTag deletes a tag
func (t *TagHandler) deleteTag(tagName string) error {
	exists, err := t.tagExists(tagName)
	if err != nil {
		return fmt.Errorf("failed to check if tag exists: %v", err)
	}

	if !exists {
		return fmt.Errorf("tag '%s' does not exist", tagName)
	}

	var confirm bool
	err = huh.NewConfirm().
		Title("Delete Tag").
		Description(fmt.Sprintf("Delete tag '%s'?", tagName)).
		Value(&confirm).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if !confirm {
		output.Info("Tag deletion cancelled")
		return nil
	}

	// Delete local tag
	_, err = t.git.Run("tag", "-d", tagName)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %v", err)
	}

	output.Success("Deleted tag: %s", tagName)

	// Ask if user wants to delete from remote
	var deleteRemote bool
	err = huh.NewConfirm().
		Title("Delete Remote Tag").
		Description("Delete tag from remote repository?").
		Value(&deleteRemote).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if deleteRemote {
		_, err = t.git.Run("push", "origin", ":refs/tags/"+tagName)
		if err != nil {
			output.Warning("Failed to delete remote tag: %v", err)
		} else {
			output.Success("Deleted remote tag: %s", tagName)
		}
	}

	return nil
}

// pushTag pushes a specific tag to remote
func (t *TagHandler) pushTag(tagName string) error {
	exists, err := t.tagExists(tagName)
	if err != nil {
		return fmt.Errorf("failed to check if tag exists: %v", err)
	}

	if !exists {
		return fmt.Errorf("tag '%s' does not exist", tagName)
	}

	_, err = t.git.Run("push", "origin", tagName)
	if err != nil {
		return fmt.Errorf("failed to push tag: %v", err)
	}

	output.Success("Pushed tag: %s", tagName)
	return nil
}

// pushAllTags pushes all tags to remote
func (t *TagHandler) pushAllTags() error {
	_, err := t.git.Run("push", "origin", "--tags")
	if err != nil {
		return fmt.Errorf("failed to push tags: %v", err)
	}

	output.Success("Pushed all tags to remote")
	return nil
}

// showUsage displays the usage information
func (t *TagHandler) showUsage() error {
	usage := `# Tag Command

Manages Git tags with version integration and interactive forms.

## Usage

  git @ tag
  git @ tag create
  git @ tag list
  git @ tag delete <name>
  git @ tag push [<name>]
  git @ tag <name> <message>

## Commands

• **create**: Create tag interactively (default)
• **list**: List all tags with messages
• **delete <name>**: Delete a tag (local and optionally remote)
• **push [<name>]**: Push specific tag or all tags to remote
• **<name> <message>**: Create tag directly with name and message

## Options

• **-h, --help**: Show this help message

## Examples

  # Create tag interactively (shows current version)
  git @ tag
  git @ tag create

  # Create tag from current version with custom message
  git @ tag "Release v2.0.5 with new features"

  # Create specific tag with message
  git @ tag v1.0.0 "Initial release"

  # List all tags
  git @ tag list

  # Delete a tag
  git @ tag delete v1.0.0

  # Push specific tag
  git @ tag push v1.0.0

  # Push all tags
  git @ tag push

## Features

• **Version Integration**: Automatically uses current GitAT version
• **Interactive Forms**: Guided tag creation with validation
• **Annotated Tags**: Creates annotated tags with messages
• **Remote Management**: Push/delete tags from remote repositories
• **Confirmation Prompts**: Safety prompts for destructive operations
• **Tag Validation**: Checks for existing tags before creation

## Integration

• Works with **git @ version** command
• Uses current version as default tag name
• Supports conventional version tagging
• Integrates with GitAT workflow

## Use Cases

• **Release Tagging**: Tag releases with version numbers
• **Milestone Marking**: Mark important project milestones
• **Version Management**: Sync tags with semantic versions
• **Release Notes**: Annotate releases with descriptions

## Notes

• Tags are created as annotated tags with messages
• Supports both local and remote tag management
• Integrates with GitAT version system
• Follows semantic versioning conventions
`

	return output.Markdown(usage)
}

// showInteractiveForm shows an interactive form for tag creation
func (t *TagHandler) showInteractiveForm() error {
	// Get current version
	currentVersion, err := t.getCurrentVersion()
	if err != nil {
		currentVersion = "0.0.0"
	}

	var tagName string
	var message string
	var useCurrentVersion bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Use Current Version").
				Description(fmt.Sprintf("Use current version '%s' as tag name?", currentVersion)).
				Value(&useCurrentVersion),
		).WithHideFunc(func() bool {
			return currentVersion == "0.0.0"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Tag Name").
				Description("Enter the tag name").
				Placeholder("e.g., v1.0.0").
				Value(&tagName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("tag name cannot be empty")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return useCurrentVersion
		}),

		huh.NewGroup(
			huh.NewText().
				Title("Tag Message").
				Description("Enter the tag message/description").
				Placeholder("e.g., Release with new features and bug fixes").
				Value(&message).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("tag message cannot be empty")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to show form: %w", err)
	}

	// Use current version if selected
	if useCurrentVersion {
		tagName = currentVersion
	}

	return t.createTag(tagName, message)
}
