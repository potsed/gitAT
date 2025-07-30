package cli

import (
	"fmt"
	"os"

	"github.com/potsed/gitAT/internal/commands"
	"github.com/potsed/gitAT/internal/config"
)

// App represents the CLI application
type App struct {
	config *config.Config
	cmds   *commands.Manager
}

// NewApp creates a new CLI application
func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		cmds:   commands.NewManager(cfg),
	}
}

// Run executes the CLI application with the given arguments
func (a *App) Run(args []string) error {
	if len(args) == 0 {
		return a.showUsage()
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "work":
		return a.cmds.Work(commandArgs)
	case "hotfix":
		return a.cmds.Hotfix(commandArgs)
	case "save":
		return a.cmds.Save(commandArgs)
	case "squash":
		return a.cmds.Squash(commandArgs)
	case "pr":
		return a.cmds.PullRequest(commandArgs)
	case "branch":
		return a.cmds.Branch(commandArgs)
	case "sweep":
		return a.cmds.Sweep(commandArgs)
	case "info":
		return a.cmds.Info(commandArgs)
	case "hash":
		return a.cmds.Hash(commandArgs)
	case "product":
		return a.cmds.Product(commandArgs)
	case "feature":
		return a.cmds.Feature(commandArgs)
	case "issue":
		return a.cmds.Issue(commandArgs)
	case "version":
		return a.cmds.Version(commandArgs)
	case "release":
		return a.cmds.Release(commandArgs)
	case "master", "root":
		return a.cmds.Master(commandArgs)
	case "wip":
		return a.cmds.WIP(commandArgs)
	case "changes":
		return a.cmds.Changes(commandArgs)
	case "logs":
		return a.cmds.Logs(commandArgs)
	case "_label":
		return a.cmds.Label(commandArgs)
	case "_id":
		return a.cmds.ID(commandArgs)
	case "_path":
		return a.cmds.Path(commandArgs)
	case "_trunk":
		return a.cmds.Trunk(commandArgs)
	case "ignore":
		return a.cmds.Ignore(commandArgs)
	case "initlocal":
		return a.cmds.InitLocal(commandArgs)
	case "initremote":
		return a.cmds.InitRemote(commandArgs)
	case "_security":
		return a.cmds.Security(commandArgs)
	case "_go":
		return a.cmds.Go(commandArgs)
	case "help", "-h", "--help":
		return a.showUsage()
	case "-v", "--version":
		return a.showVersion()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// showUsage displays the usage information
func (a *App) showUsage() error {
	fmt.Fprintf(os.Stdout, `GitAT - Git Workflow Management Tool

Usage: git @ <command> [options]

Commands:
  work <type> <description>    Create work branches following Conventional Commits
  hotfix <description>         Create hotfix branches for urgent fixes
  save "message"               Securely save changes with validation
  squash [branch]              Squash commits with auto-detection of parent branch
  pr [options]                 Create Pull Requests with auto-description generation
  branch                       Manage working branch configuration
  sweep                        Clean up local branches (merged + remote-deleted)
  info                         Comprehensive status report from all commands
  hash                         Detailed branch status and commit relationships
  product [<name>]             Product name configuration
  feature [<name>]             Feature name configuration
  issue [<id>]                 Issue/task identifier configuration
  version                      Semantic versioning management
  release                      Create releases with proper tagging
  master, root                 Switch to trunk branches
  wip                          Work in progress management
  changes                      View uncommitted changes
  logs                         View commit history
  _label                       Generate commit labels
  _id                          Generate unique project identifiers
  _path                        Get repository path
  _trunk                       Manage trunk branch configuration
  ignore                       Add patterns to .gitignore
  initlocal                    Initialize local repository with branch structure
  initremote                   Initialize remote repository with basic structure
  _security                    Security utilities and status
  _go                          Initialize GitAT for current repository

Options:
  -h, --help                   Show this help message
  -v, --version                Show version information

Examples:
  git @ work feature add-user-authentication
  git @ save "Add user authentication system"
  git @ pr
  git @ sweep

For more information, visit: https://github.com/potsed/gitAT
`)
	return nil
}

// showVersion displays the version information
func (a *App) showVersion() error {
	fmt.Fprintf(os.Stdout, "GitAT v1.1.0\n")
	return nil
}
