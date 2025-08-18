package handlers

import (
	"fmt"
	"strings"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// LogzHandler handles logz-related commands
type LogzHandler struct {
	BaseHandler
}

// NewLogzHandler creates a new logz handler
func NewLogzHandler(cfg *config.Config, gitRepo *git.Repository) *LogzHandler {
	return &LogzHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the logz command
func (l *LogzHandler) Execute(args []string) error {
	// Check for help flags
	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return l.showUsage()
		}
	}

	// Default to showing last 10 commits with a nice format
	gitArgs := []string{"log", "--oneline", "-10"}
	
	// Parse additional arguments
	if len(args) > 0 {
		gitArgs = append(gitArgs, args...)
	}

	// Execute git log command
	outputStr, err := l.git.Run(gitArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute git log: %w", err)
	}

	// Display the output
	if outputStr != "" {
		lines := strings.Split(outputStr, "\n")
		output.Title("Commit History")
		for _, line := range lines {
			if line != "" {
				output.Info("%s", line)
			}
		}
	} else {
		output.Info("%s", "No commits found")
	}

	return nil
}

// showUsage displays the logz command usage
func (l *LogzHandler) showUsage() error {
	return output.Markdown(`# Logz Command

View commit history with enhanced formatting.

## Usage

` + "```" + `bash
git @ logz [options]
` + "```" + `

## Description

This command provides a nicely formatted view of the commit history.
It's a GitAT-enhanced version of 'git log' with better output formatting.

## Options

All standard 'git log' options are supported. Some common ones:

- **-n, --max-count=<number>**: Limit the number of commits to output
- **--oneline**: Show each commit as a single line
- **--graph**: Draw a text-based graphical representation of the commit history
- **--all**: Show all branches
- **--author=<pattern>**: Show only commits by authors matching the pattern

## Examples

` + "```" + `bash
# Show last 10 commits (default)
git @ logz

# Show last 5 commits
git @ logz -5

# Show commits with graph
git @ logz --graph --oneline

# Show all commits by a specific author
git @ logz --author="John Doe"
` + "```" + `
`)
}