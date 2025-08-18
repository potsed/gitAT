package handlers

import (
	"fmt"
	"strings"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// InfoHandler handles info-related commands
type InfoHandler struct {
	BaseHandler
}

// NewInfoHandler creates a new info handler
func NewInfoHandler(cfg *config.Config, gitRepo *git.Repository) *InfoHandler {
	return &InfoHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the info command
func (i *InfoHandler) Execute(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		return i.showUsage()
	}

	// Parse options
	showStatus := true
	showBranches := false
	showConfig := false
	showRemote := false

	for _, arg := range args {
		switch arg {
		case "--status", "-s":
			showStatus = true
		case "--branches", "-b":
			showBranches = true
		case "--config", "-c":
			showConfig = true
		case "--remote", "-r":
			showRemote = true
		case "--all", "-a":
			showStatus = true
			showBranches = true
			showConfig = true
			showRemote = true
		}
	}

	// If no specific options, show status by default
	if !showBranches && !showConfig && !showRemote {
		showStatus = true
	}

	if showStatus {
		if err := i.showStatus(); err != nil {
			return err
		}
	}

	if showBranches {
		if err := i.showBranches(); err != nil {
			return err
		}
	}

	if showConfig {
		if err := i.showConfig(); err != nil {
			return err
		}
	}

	if showRemote {
		if err := i.showRemote(); err != nil {
			return err
		}
	}

	return nil
}

// showStatus shows repository status
func (i *InfoHandler) showStatus() error {
	output.Title("📊 Repository Status")

	// Get current branch
	currentBranch, err := i.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Get repository status
	status, err := i.git.Run("status", "--porcelain")
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	// Get last commit info
	lastCommit, err := i.git.Run("log", "-1", "--oneline")
	if err != nil {
		return fmt.Errorf("failed to get last commit: %w", err)
	}

	output.Info("Current Branch: %s", currentBranch)
	output.Info("Last Commit: %s", strings.TrimSpace(lastCommit))

	// Parse status
	statusLines := strings.Split(strings.TrimSpace(status), "\n")
	if len(statusLines) == 1 && statusLines[0] == "" {
		output.Success("Working directory clean")
	} else {
		output.Warning("Working directory has changes:")
		for _, line := range statusLines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					status := parts[0]
					file := parts[1]
					switch status[0] {
					case 'M':
						output.Info("  Modified: %s", file)
					case 'A':
						output.Info("  Added: %s", file)
					case 'D':
						output.Info("  Deleted: %s", file)
					case 'R':
						output.Info("  Renamed: %s", file)
					case 'C':
						output.Info("  Copied: %s", file)
					case 'U':
						output.Info("  Unmerged: %s", file)
					case '?':
						output.Info("  Untracked: %s", file)
					}
				}
			}
		}
	}

	return nil
}

// showBranches shows branch information
func (i *InfoHandler) showBranches() error {
	output.Title("🌿 Branch Information")

	// Get current branch
	currentBranch, err := i.git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Get all local branches
	localBranches, err := i.git.Run("branch")
	if err != nil {
		return fmt.Errorf("failed to get local branches: %w", err)
	}

	// Get all remote branches
	remoteBranches, err := i.git.Run("branch", "-r")
	if err != nil {
		return fmt.Errorf("failed to get remote branches: %w", err)
	}

	output.Info("Current Branch: %s", currentBranch)

	// Parse local branches
	localBranchList := strings.Split(strings.TrimSpace(localBranches), "\n")
	output.Info("Local Branches (%d):", len(localBranchList))
	for _, branch := range localBranchList {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			if strings.HasPrefix(branch, "* ") {
				output.Success("  %s (current)", strings.TrimPrefix(branch, "* "))
			} else {
				output.Info("  %s", branch)
			}
		}
	}

	// Parse remote branches
	remoteBranchList := strings.Split(strings.TrimSpace(remoteBranches), "\n")
	output.Info("Remote Branches (%d):", len(remoteBranchList))
	for _, branch := range remoteBranchList {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			output.Info("  %s", branch)
		}
	}

	return nil
}

// showConfig shows Git configuration
func (i *InfoHandler) showConfig() error {
	output.Title("⚙️ Git Configuration")

	// Get user info
	userName, _ := i.git.GetConfig("user.name")
	userEmail, _ := i.git.GetConfig("user.email")

	// Get repository info
	remoteOrigin, _ := i.git.GetConfig("remote.origin.url")

	// Get GitAT specific config
	trunkBranch, _ := i.git.GetConfig("at.trunk")
	developBranch, _ := i.git.GetConfig("at.develop")

	output.Info("User Configuration:")
	if userName != "" {
		output.Info("  Name: %s", userName)
	}
	if userEmail != "" {
		output.Info("  Email: %s", userEmail)
	}

	output.Info("Repository Configuration:")
	if remoteOrigin != "" {
		output.Info("  Remote Origin: %s", remoteOrigin)
	}

	output.Info("GitAT Configuration:")
	if trunkBranch != "" {
		output.Info("  Trunk Branch: %s", trunkBranch)
	} else {
		output.Info("  Trunk Branch: main (default)")
	}
	if developBranch != "" {
		output.Info("  Develop Branch: %s", developBranch)
	} else {
		output.Info("  Develop Branch: develop (default)")
	}

	return nil
}

// showRemote shows remote repository information
func (i *InfoHandler) showRemote() error {
	output.Title("🌐 Remote Information")

	// Get remote URLs
	remotes, err := i.git.Run("remote", "-v")
	if err != nil {
		return fmt.Errorf("failed to get remotes: %w", err)
	}

	if strings.TrimSpace(remotes) == "" {
		output.Info("No remote repositories configured")
		return nil
	}

	remoteLines := strings.Split(strings.TrimSpace(remotes), "\n")
	output.Info("Remote Repositories:")
	for _, remote := range remoteLines {
		remote = strings.TrimSpace(remote)
		if remote != "" {
			parts := strings.Fields(remote)
			if len(parts) >= 2 {
				name := parts[0]
				url := parts[1]
				purpose := "fetch"
				if len(parts) >= 3 && parts[2] == "(push)" {
					purpose = "push"
				}
				output.Info("  %s (%s): %s", name, purpose, url)
			}
		}
	}

	return nil
}

// showUsage displays the info command usage
func (i *InfoHandler) showUsage() error {
	usage := "# Info Command\n\n" +
		"Displays comprehensive information about the Git repository and its state.\n\n" +
		"## Usage\n\n" +
		"```bash\n" +
		"git @ info [options]\n" +
		"git @ info --status\n" +
		"git @ info --branches\n" +
		"git @ info --config\n" +
		"git @ info --remote\n" +
		"```\n\n" +
		"## Options\n\n" +
		"- **--status, -s**: Show repository status (default)\n" +
		"- **--branches, -b**: Show branch information\n" +
		"- **--config, -c**: Show Git configuration\n" +
		"- **--remote, -r**: Show remote repository information\n" +
		"- **--all, -a**: Show all information\n" +
		"- **-h, --help**: Show this help message\n\n" +
		"## Examples\n\n" +
		"```bash\n" +
		"# Show repository status (default)\n" +
		"git @ info\n\n" +
		"# Show all information\n" +
		"git @ info --all\n\n" +
		"# Show only branch information\n" +
		"git @ info --branches\n\n" +
		"# Show configuration\n" +
		"git @ info --config\n\n" +
		"# Show remote information\n" +
		"git @ info --remote\n" +
		"```\n\n" +
		"## Information Displayed\n\n" +
		"### Status Information\n" +
		"- Current branch\n" +
		"- Last commit\n" +
		"- Working directory status\n" +
		"- Modified, added, deleted files\n" +
		"- Untracked files\n\n" +
		"### Branch Information\n" +
		"- Current branch (highlighted)\n" +
		"- All local branches\n" +
		"- All remote branches\n" +
		"- Branch counts\n\n" +
		"### Configuration Information\n" +
		"- User name and email\n" +
		"- Remote origin URL\n" +
		"- GitAT specific configuration\n" +
		"- Trunk and develop branch settings\n\n" +
		"### Remote Information\n" +
		"- Remote repository URLs\n" +
		"- Fetch and push URLs\n" +
		"- Remote names and purposes\n\n" +
		"## Use Cases\n\n" +
		"- **Quick Status Check**: See what's changed and current state\n" +
		"- **Branch Overview**: Understand branch structure\n" +
		"- **Configuration Review**: Verify Git and GitAT settings\n" +
		"- **Remote Verification**: Check remote repository setup\n" +
		"- **Troubleshooting**: Diagnose repository issues\n"
	return output.Markdown(usage)
}
