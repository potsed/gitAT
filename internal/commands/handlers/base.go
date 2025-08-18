package handlers

import (
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// Handler defines the interface for all command handlers
type Handler interface {
	Execute(args []string) error
}

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	config *config.Config
	git    *git.Repository
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(cfg *config.Config, gitRepo *git.Repository) BaseHandler {
	return BaseHandler{
		config: cfg,
		git:    gitRepo,
	}
}

// GetConfig returns the config instance
func (b *BaseHandler) GetConfig() *config.Config {
	return b.config
}

// GetGit returns the git repository instance
func (b *BaseHandler) GetGit() *git.Repository {
	return b.git
}
