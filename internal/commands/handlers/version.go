package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// VersionHandler handles version-related commands
type VersionHandler struct {
	BaseHandler
}

// NewVersionHandler creates a new version handler
func NewVersionHandler(cfg *config.Config, gitRepo *git.Repository) *VersionHandler {
	return &VersionHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the version command
func (v *VersionHandler) Execute(args []string) error {
	if len(args) == 0 {
		// Show current version
		version, err := v.getVersion()
		if err != nil {
			return fmt.Errorf("failed to get version: %w", err)
		}

		output.Title("📦 GitAT Version")
		output.Info("Current version: %s", version)

		// Show version components
		major, _ := v.git.GetConfig("at.major")
		minor, _ := v.git.GetConfig("at.minor")
		fix, _ := v.git.GetConfig("at.fix")

		output.Table(
			[]string{"Component", "Value"},
			[][]string{
				{"Major", major},
				{"Minor", minor},
				{"Fix", fix},
			},
		)

		return nil
	}

	if len(args) >= 1 {
		// Handle help first
		if args[0] == "-h" || args[0] == "--help" || args[0] == "help" || args[0] == "h" {
			return v.showUsage()
		}

		// Handle tag and reset (these are single operations)
		if args[0] == "-t" || args[0] == "--tag" {
			version, err := v.getVersion()
			if err != nil {
				return fmt.Errorf("failed to get version: %w", err)
			}
			output.Code("v" + version)
			return nil
		}

		if args[0] == "--reset" {
			return v.resetVersion()
		}

		if args[0] == "--set" {
			return v.setVersion()
		}

		// Handle multiple increment flags
		var operations []string

		// Parse all arguments for increment operations
		for _, arg := range args {
			// Handle combined flags like -Mmb
			if strings.HasPrefix(arg, "-") && len(arg) > 1 {
				// Parse each character in the flag
				for i := 1; i < len(arg); i++ {
					switch arg[i] {
					case 'M':
						operations = append(operations, "major")
					case 'm':
						operations = append(operations, "minor")
					case 'b':
						operations = append(operations, "fix")
					}
				}
			} else {
				// Handle individual flags
				switch arg {
				case "-M", "--major":
					operations = append(operations, "major")
				case "-m", "--minor":
					operations = append(operations, "minor")
				case "-b", "--bump":
					operations = append(operations, "fix")
				}
			}
		}

		// Execute operations in order: major, minor, fix
		if len(operations) > 0 {
			// Sort operations to ensure correct order
			orderedOps := make([]string, 0)

			// Add major first if present
			for _, op := range operations {
				if op == "major" {
					orderedOps = append(orderedOps, op)
				}
			}

			// Add minor second if present
			for _, op := range operations {
				if op == "minor" {
					orderedOps = append(orderedOps, op)
				}
			}

			// Add fix last if present
			for _, op := range operations {
				if op == "fix" {
					orderedOps = append(orderedOps, op)
				}
			}

			// Execute operations
			for _, op := range orderedOps {
				switch op {
				case "major":
					if err := v.incrementMajor(); err != nil {
						return err
					}
				case "minor":
					if err := v.incrementMinor(); err != nil {
						return err
					}
				case "fix":
					if err := v.incrementFix(); err != nil {
						return err
					}
				}
			}

			// Show final version
			version, err := v.getVersion()
			if err != nil {
				return fmt.Errorf("failed to get final version: %w", err)
			}
			output.Info("Final version: %s", version)
			return nil
		}
	}

	return v.showUsage()
}

// getVersion gets the current version
func (v *VersionHandler) getVersion() (string, error) {
	major, err := v.git.GetConfig("at.major")
	if err != nil {
		return "", fmt.Errorf("failed to get major version: %w", err)
	}

	minor, err := v.git.GetConfig("at.minor")
	if err != nil {
		return "", fmt.Errorf("failed to get minor version: %w", err)
	}

	fix, err := v.git.GetConfig("at.fix")
	if err != nil {
		return "", fmt.Errorf("failed to get fix version: %w", err)
	}

	return fmt.Sprintf("%s.%s.%s", major, minor, fix), nil
}

// resetVersion resets the version to 0.0.0
func (v *VersionHandler) resetVersion() error {
	output.Warning("This will reset the version to 0.0.0")

	var confirmed bool
	err := huh.NewConfirm().
		Title("Reset Version").
		Description("Are you sure you want to reset the version to 0.0.0?").
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("Version reset cancelled")
		return nil
	}

	// Reset version components
	if err := v.git.SetConfig("at.major", "0"); err != nil {
		return fmt.Errorf("failed to reset major version: %w", err)
	}
	if err := v.git.SetConfig("at.minor", "0"); err != nil {
		return fmt.Errorf("failed to reset minor version: %w", err)
	}
	if err := v.git.SetConfig("at.fix", "0"); err != nil {
		return fmt.Errorf("failed to reset fix version: %w", err)
	}

	// Log the change
	if err := v.writeVersionLog("Version reset to 0.0.0"); err != nil {
		output.Warning("Failed to log version change: %v", err)
	}

	output.Success("Version reset to 0.0.0")
	return nil
}

// setVersion sets the version interactively
func (v *VersionHandler) setVersion() error {
	var major, minor, fix string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Major Version").
				Description("Enter the major version number").
				Value(&major).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("major version cannot be empty")
					}
					if _, err := strconv.Atoi(s); err != nil {
						return fmt.Errorf("major version must be a number")
					}
					return nil
				}),
			huh.NewInput().
				Title("Minor Version").
				Description("Enter the minor version number").
				Value(&minor).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("minor version cannot be empty")
					}
					if _, err := strconv.Atoi(s); err != nil {
						return fmt.Errorf("minor version must be a number")
					}
					return nil
				}),
			huh.NewInput().
				Title("Fix Version").
				Description("Enter the fix version number").
				Value(&fix).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("fix version cannot be empty")
					}
					if _, err := strconv.Atoi(s); err != nil {
						return fmt.Errorf("fix version must be a number")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get version input: %w", err)
	}

	// Set version components
	if err := v.git.SetConfig("at.major", major); err != nil {
		return fmt.Errorf("failed to set major version: %w", err)
	}
	if err := v.git.SetConfig("at.minor", minor); err != nil {
		return fmt.Errorf("failed to set minor version: %w", err)
	}
	if err := v.git.SetConfig("at.fix", fix); err != nil {
		return fmt.Errorf("failed to set fix version: %w", err)
	}

	// Log the change
	if err := v.writeVersionLog(fmt.Sprintf("Version set to %s.%s.%s", major, minor, fix)); err != nil {
		output.Warning("Failed to log version change: %v", err)
	}

	output.Success("Version set to %s.%s.%s", major, minor, fix)
	return nil
}

// incrementMajor increments the major version
func (v *VersionHandler) incrementMajor() error {
	major, err := v.git.GetConfig("at.major")
	if err != nil {
		return fmt.Errorf("failed to get major version: %w", err)
	}

	majorNum, err := strconv.Atoi(major)
	if err != nil {
		return fmt.Errorf("failed to parse major version: %w", err)
	}

	newMajor := strconv.Itoa(majorNum + 1)
	if err := v.git.SetConfig("at.major", newMajor); err != nil {
		return fmt.Errorf("failed to set major version: %w", err)
	}

	// Reset minor and fix to 0
	if err := v.git.SetConfig("at.minor", "0"); err != nil {
		return fmt.Errorf("failed to reset minor version: %w", err)
	}
	if err := v.git.SetConfig("at.fix", "0"); err != nil {
		return fmt.Errorf("failed to reset fix version: %w", err)
	}

	// Log the change
	if err := v.writeVersionLog(fmt.Sprintf("Major version incremented to %s", newMajor)); err != nil {
		output.Warning("Failed to log version change: %v", err)
	}

	output.Success("Major version incremented to %s", newMajor)
	return nil
}

// incrementMinor increments the minor version
func (v *VersionHandler) incrementMinor() error {
	minor, err := v.git.GetConfig("at.minor")
	if err != nil {
		return fmt.Errorf("failed to get minor version: %w", err)
	}

	minorNum, err := strconv.Atoi(minor)
	if err != nil {
		return fmt.Errorf("failed to parse minor version: %w", err)
	}

	newMinor := strconv.Itoa(minorNum + 1)
	if err := v.git.SetConfig("at.minor", newMinor); err != nil {
		return fmt.Errorf("failed to set minor version: %w", err)
	}

	// Reset fix to 0
	if err := v.git.SetConfig("at.fix", "0"); err != nil {
		return fmt.Errorf("failed to reset fix version: %w", err)
	}

	// Log the change
	if err := v.writeVersionLog(fmt.Sprintf("Minor version incremented to %s", newMinor)); err != nil {
		output.Warning("Failed to log version change: %v", err)
	}

	output.Success("Minor version incremented to %s", newMinor)
	return nil
}

// incrementFix increments the fix version
func (v *VersionHandler) incrementFix() error {
	fix, err := v.git.GetConfig("at.fix")
	if err != nil {
		return fmt.Errorf("failed to get fix version: %w", err)
	}

	fixNum, err := strconv.Atoi(fix)
	if err != nil {
		return fmt.Errorf("failed to parse fix version: %w", err)
	}

	newFix := strconv.Itoa(fixNum + 1)
	if err := v.git.SetConfig("at.fix", newFix); err != nil {
		return fmt.Errorf("failed to set fix version: %w", err)
	}

	// Log the change
	if err := v.writeVersionLog(fmt.Sprintf("Fix version incremented to %s", newFix)); err != nil {
		output.Warning("Failed to log version change: %v", err)
	}

	output.Success("Fix version incremented to %s", newFix)
	return nil
}

// writeVersionLog writes version change logs to a file
func (v *VersionHandler) writeVersionLog(message string) error {
	// Get git root directory
	gitRoot, err := v.git.Run("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("failed to get git root: %w", err)
	}
	gitRoot = strings.TrimSpace(gitRoot)

	// Create logs directory if it doesn't exist
	logsDir := filepath.Join(gitRoot, ".git", "gitat-logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file path
	logFile := filepath.Join(logsDir, "version-changes.log")

	// Open file in append mode, create if doesn't exist
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Write log entry with timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	_, err = io.WriteString(file, logEntry)
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

// showUsage displays the version command usage
func (v *VersionHandler) showUsage() error {
	usage := `# Version Command

Manages semantic versioning for the project.

## Usage

git @ version                    # Show current version
git @ version -t                 # Show version tag (v1.2.3)
git @ version --reset            # Reset version to 0.0.0
git @ version --set              # Set version interactively
git @ version -M                 # Increment major version
git @ version -m                 # Increment minor version
git @ version -b                 # Increment fix version
git @ version -Mmb               # Increment all (major, minor, fix)

## Options

- **-t, --tag**: Show version as tag (e.g., v1.2.3)
- **--reset**: Reset version to 0.0.0 (requires confirmation)
- **--set**: Set version interactively using a form
- **-M, --major**: Increment major version (resets minor and fix to 0)
- **-m, --minor**: Increment minor version (resets fix to 0)
- **-b, --bump**: Increment fix version
- **-h, --help**: Show this help message

## Examples

# Show current version
git @ version

# Increment major version
git @ version -M

# Increment multiple components
git @ version -Mmb

# Reset version
git @ version --reset

# Set version interactively
git @ version --set

## Version Format

Uses semantic versioning: **MAJOR.MINOR.FIX**

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **FIX**: Bug fixes (backward compatible)

## Logging

All version changes are logged to .git/gitat-logs/version-changes.log for audit purposes.`

	return output.Markdown(usage)
}
