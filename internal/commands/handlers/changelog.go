// changelog.go - Changelog command handler
package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// ChangelogHandler handles changelog-related commands
type ChangelogHandler struct {
	BaseHandler
}

// NewChangelogHandler creates a new changelog handler
func NewChangelogHandler(cfg *config.Config, gitRepo *git.Repository) *ChangelogHandler {
	return &ChangelogHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the changelog command
func (c *ChangelogHandler) Execute(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "generate", "gen":
			return c.generateChangelog()
		case "update":
			return c.updateChangelog()
		case "preview":
			return c.previewChangelog()
		case "release":
			if len(args) < 2 {
				return fmt.Errorf("release version required")
			}
			return c.releaseChangelog(args[1])
		case "-h", "--help", "help":
			return c.showUsage()
		default:
			return c.generateChangelog()
		}
	}
	
	return c.generateChangelog()
}

// generateChangelog generates a new changelog
func (c *ChangelogHandler) generateChangelog() error {
	// Get trunk branch
	trunkBranch := c.getTrunkBranch()

	// Get commits
	commits, err := c.getCommits(trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	// Parse into categories
	categories := c.parseCommits(commits)

	// Generate changelog
	changelog := c.formatChangelog(categories)

	// Show changelog
	output.Title("Generated Changelog")
	fmt.Println(changelog)

	// Try to save automatically in non-interactive environments
	// Check if we're in a non-interactive environment
	if isInteractiveEnvironment() {
		// Ask to save
		var save bool
		err = huh.NewConfirm().
			Title("Save Changelog?").
			Description("Save this changelog to CHANGELOG.md?").
			Value(&save).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get user input: %w", err)
		}

		if save {
			return c.saveChangelog(changelog)
		}
	} else {
		// In non-interactive environments, save automatically
		output.Info("Non-interactive environment detected. Saving changelog automatically.")
		return c.saveChangelog(changelog)
	}

	return nil
}

// updateChangelog updates existing changelog
func (c *ChangelogHandler) updateChangelog() error {
	output.Info("Updating existing changelog...")
	// Implementation would read existing and merge
	return c.generateChangelog()
}

// previewChangelog shows a preview without saving
func (c *ChangelogHandler) previewChangelog() error {
	// Get trunk branch
	trunkBranch := c.getTrunkBranch()

	// Get commits
	commits, err := c.getCommits(trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	// Parse into categories
	categories := c.parseCommits(commits)

	// Show preview
	output.Title("Changelog Preview")
	output.Info("Categories found: %d", len(categories))
	
	for category, entries := range categories {
		if len(entries) > 0 {
			output.Subtitle(fmt.Sprintf("%s (%d)", category, len(entries)))
			for _, entry := range entries {
				output.Info("  • %s", entry)
			}
		}
	}

	return nil
}

// releaseChangelog generates a release entry
func (c *ChangelogHandler) releaseChangelog(version string) error {
	// Get trunk branch
	trunkBranch := c.getTrunkBranch()

	// Get commits
	commits, err := c.getCommits(trunkBranch)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	// Parse into categories
	categories := c.parseCommits(commits)

	// Generate release entry
	date := time.Now().Format("2006-01-02")
	releaseEntry := c.formatReleaseEntry(version, date, categories)

	// Show release entry
	output.Title("Release Changelog Entry")
	fmt.Println(releaseEntry)

	return nil
}

// Helper methods
func (c *ChangelogHandler) getTrunkBranch() string {
	trunkBranch := c.config.Trunk
	if trunkBranch == "" {
		trunkBranch, _ = c.git.GetConfig("at.trunk")
	}
	if trunkBranch == "" {
		trunkBranch = "main"
		_, err := c.git.Run("rev-parse", "--verify", "main")
		if err != nil {
			trunkBranch = "master"
		}
	}
	return trunkBranch
}

func (c *ChangelogHandler) getCommits(trunkBranch string) ([]string, error) {
	// Try to get latest tag
	latestTag, err := c.git.Run("describe", "--tags", "--abbrev=0")
	if err != nil {
		// No tags, get all commits from trunk
		commits, err := c.git.Run("log", trunkBranch, "--oneline")
		if err != nil {
			return nil, err
		}
		if commits == "" {
			return []string{}, nil
		}
		return strings.Split(commits, "\n"), nil
	}

	// Get commits since last tag
	commits, err := c.git.Run("log", fmt.Sprintf("%s..HEAD", latestTag), "--oneline")
	if err != nil {
		return nil, err
	}
	
	if commits == "" {
		return []string{}, nil
	}
	
	return strings.Split(commits, "\n"), nil
}

func (c *ChangelogHandler) parseCommits(commits []string) map[string][]string {
	categories := make(map[string][]string)
	
	// Initialize categories
	categories["Features"] = []string{}
	categories["Bug Fixes"] = []string{}
	categories["Documentation"] = []string{}
	categories["Refactoring"] = []string{}
	categories["Performance"] = []string{}
	categories["Tests"] = []string{}
	categories["Build System"] = []string{}
	categories["CI"] = []string{}
	categories["Chores"] = []string{}
	categories["Deprecations"] = []string{}
	categories["Security"] = []string{}
	categories["Removals"] = []string{}
	categories["Other Changes"] = []string{}

	for _, commit := range commits {
		if commit == "" {
			continue
		}
		
		// Parse conventional commit
		category, message := c.parseCommit(commit)
		categories[category] = append(categories[category], message)
	}

	return categories
}

func (c *ChangelogHandler) parseCommit(commit string) (string, string) {
	// Extract message (skip hash)
	parts := strings.SplitN(commit, " ", 2)
	if len(parts) < 2 {
		return "Other Changes", commit
	}
	
	message := parts[1]
	
	// Categorize based on conventional commits
	switch {
	case strings.HasPrefix(message, "feat:"):
		return "Features", strings.TrimPrefix(message, "feat:")
	case strings.HasPrefix(message, "fix:"):
		return "Bug Fixes", strings.TrimPrefix(message, "fix:")
	case strings.HasPrefix(message, "docs:"):
		return "Documentation", strings.TrimPrefix(message, "docs:")
	case strings.HasPrefix(message, "refactor:"):
		return "Refactoring", strings.TrimPrefix(message, "refactor:")
	case strings.HasPrefix(message, "perf:"):
		return "Performance", strings.TrimPrefix(message, "perf:")
	case strings.HasPrefix(message, "test:"):
		return "Tests", strings.TrimPrefix(message, "test:")
	case strings.HasPrefix(message, "build:"):
		return "Build System", strings.TrimPrefix(message, "build:")
	case strings.HasPrefix(message, "ci:"):
		return "CI", strings.TrimPrefix(message, "ci:")
	case strings.HasPrefix(message, "chore:"):
		return "Chores", strings.TrimPrefix(message, "chore:")
	case strings.HasPrefix(message, "deprecate:"):
		return "Deprecations", strings.TrimPrefix(message, "deprecate:")
	case strings.HasPrefix(message, "security:"):
		return "Security", strings.TrimPrefix(message, "security:")
	case strings.HasPrefix(message, "remove:"):
		return "Removals", strings.TrimPrefix(message, "remove:")
	default:
		return "Other Changes", message
	}
}

func (c *ChangelogHandler) formatChangelog(categories map[string][]string) string {
	var builder strings.Builder
	
	builder.WriteString("# Changelog\n\n")
	builder.WriteString("All notable changes to this project will be documented in this file.\n\n")
	builder.WriteString("The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),\n")
	builder.WriteString("and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).\n\n")
	
	date := time.Now().Format("2006-01-02")
	builder.WriteString(fmt.Sprintf("## [Unreleased] - %s\n\n", date))
	
	// Categories mapped to Keep a Changelog 1.1.0 standard
	categoryMapping := map[string]string{
		"Features":       "Added",
		"Bug Fixes":      "Fixed",
		"Documentation":  "Changed", // Documentation changes are typically considered changes to existing functionality
		"Refactoring":    "Changed", // Refactoring is a change to existing functionality
		"Performance":    "Changed", // Performance improvements are changes to existing functionality
		"Tests":          "Changed", // Test improvements are changes to existing functionality
		"Build System":   "Changed", // Build system changes are changes to existing functionality
		"CI":             "Changed", // CI changes are changes to existing functionality
		"Chores":         "Changed", // Chores are typically changes to existing functionality
		"Other Changes":  "Changed", // Other changes are typically changes to existing functionality
	}
	
	// Standard Keep a Changelog categories in order
	standardCategories := []string{"Added", "Changed", "Deprecated", "Removed", "Fixed", "Security"}
	
	// Create a map to group entries by standard categories
	standardEntries := make(map[string][]string)
	for _, cat := range standardCategories {
		standardEntries[cat] = []string{}
	}
	
	// Map our categories to standard categories
	for ourCategory, entries := range categories {
		if standardCategory, exists := categoryMapping[ourCategory]; exists {
			standardEntries[standardCategory] = append(standardEntries[standardCategory], entries...)
		} else {
			// Default to "Changed" for unmapped categories
			standardEntries["Changed"] = append(standardEntries["Changed"], entries...)
		}
	}
	
	// Write entries in standard category order
	for _, category := range standardCategories {
		entries := standardEntries[category]
		if len(entries) > 0 {
			builder.WriteString(fmt.Sprintf("### %s\n\n", category))
			for _, entry := range entries {
				builder.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(entry)))
			}
			builder.WriteString("\n")
		}
	}
	
	return builder.String()
}

func (c *ChangelogHandler) formatReleaseEntry(version, date string, categories map[string][]string) string {
	var builder strings.Builder
	
	builder.WriteString(fmt.Sprintf("## [%s] - %s\n\n", version, date))
	
	// Categories mapped to Keep a Changelog 1.1.0 standard
	categoryMapping := map[string]string{
		"Features":       "Added",
		"Bug Fixes":      "Fixed",
		"Documentation":  "Changed", // Documentation changes are typically considered changes to existing functionality
		"Refactoring":    "Changed", // Refactoring is a change to existing functionality
		"Performance":    "Changed", // Performance improvements are changes to existing functionality
		"Tests":          "Changed", // Test improvements are changes to existing functionality
		"Build System":   "Changed", // Build system changes are changes to existing functionality
		"CI":             "Changed", // CI changes are changes to existing functionality
		"Chores":         "Changed", // Chores are typically changes to existing functionality
		"Other Changes":  "Changed", // Other changes are typically changes to existing functionality
	}
	
	// Standard Keep a Changelog categories in order
	standardCategories := []string{"Added", "Changed", "Deprecated", "Removed", "Fixed", "Security"}
	
	// Create a map to group entries by standard categories
	standardEntries := make(map[string][]string)
	for _, cat := range standardCategories {
		standardEntries[cat] = []string{}
	}
	
	// Map our categories to standard categories
	for ourCategory, entries := range categories {
		if standardCategory, exists := categoryMapping[ourCategory]; exists {
			standardEntries[standardCategory] = append(standardEntries[standardCategory], entries...)
		} else {
			// Default to "Changed" for unmapped categories
			standardEntries["Changed"] = append(standardEntries["Changed"], entries...)
		}
	}
	
	// Write entries in standard category order
	for _, category := range standardCategories {
		entries := standardEntries[category]
		if len(entries) > 0 {
			builder.WriteString(fmt.Sprintf("### %s\n\n", category))
			for _, entry := range entries {
				builder.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(entry)))
			}
			builder.WriteString("\n")
		}
	}
	
	return builder.String()
}

func (c *ChangelogHandler) saveChangelog(content string) error {
	// Actually save to file using standard Go file operations
	filename := "CHANGELOG.md"
	filePath := filepath.Join(c.git.Path, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, read existing content
		existingContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read existing changelog: %w", err)
		}
		
		// If file starts with # Changelog, we need to merge content
		if strings.HasPrefix(string(existingContent), "# Changelog") {
			// For simplicity, we'll prepend our new content to the existing content
			// but skip the header lines from our new content
			lines := strings.Split(content, "\n")
			if len(lines) > 3 {
				// Skip first 3 lines (# Changelog, blank, and description)
				newContent := strings.Join(lines[3:], "\n") + "\n" + string(existingContent)
				content = newContent
			}
		} else {
			// Prepend our content to existing content
			content = content + "\n" + string(existingContent)
		}
	}
	// If file doesn't exist, we'll create it with just our content
	
	// Write to file
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write changelog: %w", err)
	}
	
	output.Success("Changelog saved to %s", filename)
	return nil
}

// showUsage displays help
func (c *ChangelogHandler) showUsage() error {
	return output.Markdown(`# Changelog Command

Generate and manage changelogs based on conventional commits.

## Usage

` + "```" + `bash
git @ changelog [command]
git @ changelog generate
git @ changelog update
git @ changelog preview
git @ changelog release <version>
` + "```" + `

## Commands

• **generate, gen**: Generate new changelog (default)
• **update**: Update existing changelog
• **preview**: Show preview without saving
• **release <version>**: Generate release entry

## Examples

` + "```" + `bash
# Generate changelog
git @ changelog

# Update existing changelog
git @ changelog update

# Preview changelog
git @ changelog preview

# Generate release entry
git @ changelog release v1.2.0
` + "```" + `

## Features

• **Conventional Commits**: Auto-categorizes based on commit types
• **Release Management**: Generates formatted release entries
• **File Integration**: Reads from and writes to CHANGELOG.md
• **Interactive Prompts**: Asks for confirmation before saving
`)
}

// isInteractiveEnvironment checks if we're running in an interactive environment
func isInteractiveEnvironment() bool {
	// Check if stdin is a terminal
	_, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	
	// Check if we can access /dev/tty
	_, err = os.Open("/dev/tty")
	return err == nil
}