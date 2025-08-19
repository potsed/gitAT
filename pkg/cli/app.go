package cli

import (
	"fmt"
	"os"

	"github.com/potsed/gitAT/internal/commands"
	"github.com/potsed/gitAT/internal/config"
)

// Version information set at build time via ldflags
var (
	Version   = "unknown"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

// App represents the CLI application
type App struct {
	config *config.Config
	cmds   *commands.Manager
	version string
	commitHash string
	buildDate string
}

// NewApp creates a new CLI application
func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		cmds:   commands.NewManager(cfg),
		version: Version,
		commitHash: CommitHash,
		buildDate: BuildDate,
	}
}

// Run executes the CLI application with the given arguments
func (a *App) Run(args []string) error {
	if len(args) == 0 {
		// Requirement 4.1: When user executes `git @` with no arguments, show Git help
		return a.cmds.Execute("", []string{})
	}

	command := args[0]
	commandArgs := args[1:]

	// Handle help and version commands - these should show Gitat info, not fall through
	// Requirements 4.2 and 4.3: --version and --help should show Gitat information
	switch command {
	case "help", "-h", "--help":
		return a.showUsage()
	case "-v", "--version":
		return a.showVersion()
	}

	// Execute the command using the new Manager structure
	return a.cmds.Execute(command, commandArgs)
}

// showUsage displays the usage information
func (a *App) showUsage() error {
	fmt.Fprintf(os.Stdout, `GitAT - Git Workflow Management Tool

Usage: git @ <command> [options]

  work <type> <description>    Create work branches following Conventional Commits
  hotfix <description>         Create hotfix branches for urgent fixes
  save "message"               Securely save changes with validation
  squash [branch]              Squash commits with auto-detection of parent branch
  pr [options]                 Create Pull Requests with auto-description generation
  sprout                       Manage working branch configuration
  sweep                        Clean up local branches (merged + remote-deleted)
  info                         Comprehensive status report from all commands
  shasum                       Detailed branch status and commit relationships
  product [<name>]             Product name configuration
  feature [<name>]             Feature name configuration
  issue [<id>]                 Issue/task identifier configuration
  semver                       Semantic versioning management
  dub                          Enhanced tag creation with version integration
  release                      Create releases with proper tagging
  master, root                 Switch to trunk branches (main/master)
  wip                          Work in progress management
  changes                      View uncommitted changes
  logz                         View commit history
  ignore                       Add patterns to .gitignore
  setup-local                  Initialize local repository with branch structure
  setup-remote                 Initialize remote repository with basic structure
  changelog                    Generate and manage changelogs
  rebase                       Rebase current branch onto another branch
  commitizen, cz               Create conventional commits with interactive prompt

Options:
  -h, --help                   Show this help message
  -v, --version                Show version information (GitAT {{VERSION}})

Examples:
  git @ work feature add-user-authentication
  git @ save "Add user authentication system"
  git @ pr
  git @ sweep

For more information, visit: https://github.com/potsed/gitAT

## Man Pages

GitAT includes comprehensive man pages:
  man git-@        # Main command
  man git-@-work   # Work branch creation
  man git-@-save   # Commit changes
  man git-@-sprout # Branch management
  man git-@-sweep  # Branch cleanup
  man git-@-logz   # Commit history
  man git-@-shasum # Branch status
  man git-@-main   # Switch to trunk branch
  man git-@-changelog # Changelog management
  man git-@-rebase    # Branch rebasing
  man git-@-commitizen # Conventional commits
`)
	return nil
}

// showVersion displays the version information
func (a *App) showVersion() error {
	fmt.Fprintf(os.Stdout, "GitAT %s\n", a.version)
	if a.commitHash != "unknown" {
		fmt.Fprintf(os.Stdout, "Commit: %s\n", a.commitHash)
	}
	if a.buildDate != "unknown" {
		fmt.Fprintf(os.Stdout, "Built: %s\n", a.buildDate)
	}
	return nil
}
