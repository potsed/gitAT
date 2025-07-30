package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Repository represents a Git repository
type Repository struct {
	Path string
}

// NewRepository creates a new Git repository instance
func NewRepository(path string) *Repository {
	return &Repository{
		Path: path,
	}
}

// Run executes a Git command
func (r *Repository) Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git command failed: %w", err)
	}
	
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch returns the current branch name
func (r *Repository) GetCurrentBranch() (string, error) {
	return r.Run("rev-parse", "--abbrev-ref", "HEAD")
}

// GetConfig gets a Git configuration value
func (r *Repository) GetConfig(key string) (string, error) {
	return r.Run("config", "--get", key)
}

// SetConfig sets a Git configuration value
func (r *Repository) SetConfig(key, value string) error {
	_, err := r.Run("config", key, value)
	return err
}

// CreateBranch creates a new branch
func (r *Repository) CreateBranch(name string) error {
	_, err := r.Run("checkout", "-b", name)
	return err
}

// SwitchBranch switches to an existing branch
func (r *Repository) SwitchBranch(name string) error {
	_, err := r.Run("checkout", name)
	return err
}

// GetBranches returns all local branches
func (r *Repository) GetBranches() ([]string, error) {
	output, err := r.Run("branch", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	
	if output == "" {
		return []string{}, nil
	}
	
	return strings.Split(output, "\n"), nil
}

// GetMergedBranches returns branches merged into the given branch
func (r *Repository) GetMergedBranches(branch string) ([]string, error) {
	output, err := r.Run("branch", "--merged", branch, "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	
	if output == "" {
		return []string{}, nil
	}
	
	return strings.Split(output, "\n"), nil
}

// DeleteBranch deletes a branch
func (r *Repository) DeleteBranch(name string, force bool) error {
	args := []string{"branch"}
	if force {
		args = append(args, "-D")
	} else {
		args = append(args, "-d")
	}
	args = append(args, name)
	
	_, err := r.Run(args...)
	return err
}

// Add adds files to staging
func (r *Repository) Add(files ...string) error {
	args := append([]string{"add"}, files...)
	_, err := r.Run(args...)
	return err
}

// Commit creates a commit with the given message
func (r *Repository) Commit(message string) error {
	_, err := r.Run("commit", "-m", message)
	return err
}

// GetStatus returns the current Git status
func (r *Repository) GetStatus() (string, error) {
	return r.Run("status", "--porcelain")
}

// GetLog returns the Git log
func (r *Repository) GetLog(format string, limit int) (string, error) {
	args := []string{"log"}
	if format != "" {
		args = append(args, "--format="+format)
	}
	if limit > 0 {
		args = append(args, fmt.Sprintf("-%d", limit))
	}
	
	return r.Run(args...)
}

// GetDiff returns the diff between branches or commits
func (r *Repository) GetDiff(from, to string) (string, error) {
	return r.Run("diff", "--name-only", from, to)
}

// GetRemoteURL returns the remote URL
func (r *Repository) GetRemoteURL(remote string) (string, error) {
	return r.Run("remote", "get-url", remote)
}

// Push pushes to remote
func (r *Repository) Push(remote, branch string) error {
	_, err := r.Run("push", remote, branch)
	return err
}

// Pull pulls from remote
func (r *Repository) Pull(remote, branch string) error {
	_, err := r.Run("pull", remote, branch)
	return err
}

// Fetch fetches from remote
func (r *Repository) Fetch(remote string) error {
	_, err := r.Run("fetch", remote)
	return err
}

// PruneRemote removes remote tracking branches that no longer exist on remote
func (r *Repository) PruneRemote(remote string) error {
	_, err := r.Run("remote", "prune", remote)
	return err
}

// GetMergeBase returns the merge base of two commits
func (r *Repository) GetMergeBase(commit1, commit2 string) (string, error) {
	return r.Run("merge-base", commit1, commit2)
}

// GetCommitHash returns the hash of a commit
func (r *Repository) GetCommitHash(ref string) (string, error) {
	return r.Run("rev-parse", ref)
}

// GetCommitDate returns the commit date
func (r *Repository) GetCommitDate(ref string) (string, error) {
	return r.Run("log", "-1", "--format=%ct", ref)
}

// GetCommitAuthor returns the commit author
func (r *Repository) GetCommitAuthor(ref string) (string, error) {
	return r.Run("log", "-1", "--format=%an", ref)
}

// GetCommitMessage returns the commit message
func (r *Repository) GetCommitMessage(ref string) (string, error) {
	return r.Run("log", "-1", "--format=%s", ref)
} 