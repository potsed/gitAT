package handlers

import (
	"fmt"
	"strings"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// ShasumHandler handles shasum-related commands
type ShasumHandler struct {
	BaseHandler
}

// NewShasumHandler creates a new shasum handler
func NewShasumHandler(cfg *config.Config, gitRepo *git.Repository) *ShasumHandler {
	return &ShasumHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the shasum command
func (s *ShasumHandler) Execute(args []string) error {
	// Check for help flags
	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help", "help", "h":
			return s.showUsage()
		}
	}

	// Default to showing current branch hash info
	currentBranch, err := s.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Get the current commit hash
	currentHash, err := s.git.GetCommitHash("HEAD")
	if err != nil {
		return fmt.Errorf("failed to get current commit hash: %w", err)
	}

	// Get merge base with trunk if we can determine trunk
	trunkBranch := s.config.Trunk
	if trunkBranch == "" {
		// Try to get from git config
		trunkBranch, _ = s.git.GetConfig("at.trunk")
	}
	if trunkBranch == "" {
		// Default fallbacks
		trunkBranch = "main"
		// Check if main exists
		_, err := s.git.Run("rev-parse", "--verify", "main")
		if err != nil {
			trunkBranch = "master"
			// Check if master exists
			_, err := s.git.Run("rev-parse", "--verify", "master")
			if err != nil {
				trunkBranch = ""
			}
		}
	}

	var mergeBase string
	if trunkBranch != "" {
		// Verify trunk branch exists
		_, err := s.git.Run("rev-parse", "--verify", trunkBranch)
		if err == nil {
			mergeBase, _ = s.git.GetMergeBase("HEAD", trunkBranch)
		}
	}

	// Display hash information
	output.Title("Branch Hash Information")
	
	data := [][]string{
		{"Current Branch", currentBranch},
		{"Current Hash", currentHash},
	}
	
	if mergeBase != "" && trunkBranch != "" {
		data = append(data, []string{"Merge Base with " + trunkBranch, mergeBase})
		// Show how many commits ahead/behind
		aheadCount, _ := s.git.Run("rev-list", "--count", trunkBranch+"..HEAD")
		behindCount, _ := s.git.Run("rev-list", "--count", "HEAD.."+trunkBranch)
		if aheadCount != "" && behindCount != "" {
			data = append(data, []string{"Commits Ahead", aheadCount})
			data = append(data, []string{"Commits Behind", behindCount})
		}
	}
	
	output.Table([]string{"Property", "Value"}, data)

	// Show recent commits for context
	output.Title("Recent Commits")
	logOutput, err := s.git.Run("log", "--oneline", "-5")
	if err == nil && logOutput != "" {
		lines := strings.Split(logOutput, "\n")
		for _, line := range lines {
			if line != "" {
				output.Info("%s", line)
			}
		}
	}

	return nil
}

// showUsage displays the shasum command usage
func (s *ShasumHandler) showUsage() error {
	return output.Markdown(`# Shasum Command

Display detailed branch status and commit relationships.

## Usage

` + "```" + `bash
git @ shasum
` + "```" + `

## Description

This command provides detailed information about the current branch's
commit hash, its relationship to the trunk branch, and recent commit history.

## Information Displayed

- Current branch name
- Current commit hash (SHA)
- Merge base with trunk branch (if trunk can be determined)
- Number of commits ahead/behind trunk
- Recent commit history (last 5 commits)

## Examples

` + "```" + `bash
# Show hash information for current branch
git @ shasum
` + "```" + `
`)
}