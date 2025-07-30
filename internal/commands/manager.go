package commands

import (
	"fmt"

	"github.com/potsed/gitAT/internal/commands/handlers"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// Manager handles all GitAT commands
type Manager struct {
	config *config.Config
	git    *git.Repository
	// Command handlers
	work     *handlers.WorkHandler
	version  *handlers.VersionHandler
	save     *handlers.SaveHandler
	squash   *handlers.SquashHandler
	pr       *handlers.PullRequestHandler
	branch   *handlers.BranchHandler
	sweep    *handlers.SweepHandler
	wip      *handlers.WIPHandler
	hotfix   *handlers.HotfixHandler
	utility  *handlers.UtilityHandler
	info     *handlers.InfoHandler
	release  *handlers.ReleaseHandler
	feature  *handlers.FeatureHandler
	product  *handlers.ProductHandler
	issue    *handlers.IssueHandler
	label    *handlers.LabelHandler
	trunk    *handlers.TrunkHandler
	ignore   *handlers.IgnoreHandler
	init     *handlers.InitHandler
	security *handlers.SecurityHandler
	goCmd    *handlers.GoHandler
}

// NewManager creates a new commands manager
func NewManager(cfg *config.Config) *Manager {
	gitRepo := git.NewRepository(cfg.RepoPath)

	return &Manager{
		config:   cfg,
		git:      gitRepo,
		work:     handlers.NewWorkHandler(cfg, gitRepo),
		version:  handlers.NewVersionHandler(cfg, gitRepo),
		save:     handlers.NewSaveHandler(cfg, gitRepo),
		squash:   handlers.NewSquashHandler(cfg, gitRepo),
		pr:       handlers.NewPullRequestHandler(cfg, gitRepo),
		branch:   handlers.NewBranchHandler(cfg, gitRepo),
		sweep:    handlers.NewSweepHandler(cfg, gitRepo),
		wip:      handlers.NewWIPHandler(cfg, gitRepo),
		hotfix:   handlers.NewHotfixHandler(cfg, gitRepo),
		utility:  handlers.NewUtilityHandler(cfg, gitRepo),
		info:     handlers.NewInfoHandler(cfg, gitRepo),
		release:  handlers.NewReleaseHandler(cfg, gitRepo),
		feature:  handlers.NewFeatureHandler(cfg, gitRepo),
		product:  handlers.NewProductHandler(cfg, gitRepo),
		issue:    handlers.NewIssueHandler(cfg, gitRepo),
		label:    handlers.NewLabelHandler(cfg, gitRepo),
		trunk:    handlers.NewTrunkHandler(cfg, gitRepo),
		ignore:   handlers.NewIgnoreHandler(cfg, gitRepo),
		init:     handlers.NewInitHandler(cfg, gitRepo),
		security: handlers.NewSecurityHandler(cfg, gitRepo),
		goCmd:    handlers.NewGoHandler(cfg, gitRepo),
	}
}

// Execute dispatches commands to appropriate handlers
func (m *Manager) Execute(command string, args []string) error {
	switch command {
	case "work":
		return m.work.Execute(args)
	case "version":
		return m.version.Execute(args)
	case "save":
		return m.save.Execute(args)
	case "squash":
		return m.squash.Execute(args)
	case "pr":
		return m.pr.Execute(args)
	case "branch":
		return m.branch.Execute(args)
	case "sweep":
		return m.sweep.Execute(args)
	case "wip":
		return m.wip.Execute(args)
	case "hotfix":
		return m.hotfix.Execute(args)
	case "info":
		return m.info.Execute(args)
	case "release":
		return m.release.Execute(args)
	case "feature":
		return m.feature.Execute(args)
	case "product":
		return m.product.Execute(args)
	case "issue":
		return m.issue.Execute(args)
	case "label":
		return m.label.Execute(args)
	case "trunk":
		return m.trunk.Execute(args)
	case "ignore":
		return m.ignore.Execute(args)
	case "initlocal":
		return m.init.ExecuteLocal(args)
	case "initremote":
		return m.init.ExecuteRemote(args)
	case "security":
		return m.security.Execute(args)
	case "_go":
		return m.goCmd.Execute(args)
	// Utility commands
	case "changes":
		return m.utility.ExecuteChanges(args)
	case "logs":
		return m.utility.ExecuteLogs(args)
	case "hash":
		return m.utility.ExecuteHash(args)
	case "id":
		return m.utility.ExecuteID(args)
	case "path":
		return m.utility.ExecutePath(args)
	case "master":
		return m.utility.ExecuteMaster(args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// GetGit returns the git repository instance
func (m *Manager) GetGit() *git.Repository {
	return m.git
}

// GetConfig returns the config instance
func (m *Manager) GetConfig() *config.Config {
	return m.config
}
