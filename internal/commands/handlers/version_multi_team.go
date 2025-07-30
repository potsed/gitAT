package handlers

import (
	"fmt"
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

// MultiTeamVersionHandler handles version management for multi-team scenarios
type MultiTeamVersionHandler struct {
	BaseHandler
}

// NewMultiTeamVersionHandler creates a new multi-team version handler
func NewMultiTeamVersionHandler(cfg *config.Config, gitRepo *git.Repository) *MultiTeamVersionHandler {
	return &MultiTeamVersionHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles multi-team version commands
func (m *MultiTeamVersionHandler) Execute(args []string) error {
	if len(args) == 0 {
		return m.showMultiTeamStatus()
	}

	switch args[0] {
	case "-h", "--help", "help", "h":
		return m.showUsage()
	case "sync":
		return m.syncVersions()
	case "lock":
		return m.lockVersion()
	case "unlock":
		return m.unlockVersion()
	case "propose":
		return m.proposeVersion(args[1:])
	case "approve":
		return m.approveVersion(args[1:])
	case "reject":
		return m.rejectVersion(args[1:])
	case "teams":
		return m.showTeamStatus()
	case "history":
		return m.showVersionHistory()
	default:
		return m.showUsage()
	}
}

// showMultiTeamStatus shows the current multi-team version status
func (m *MultiTeamVersionHandler) showMultiTeamStatus() error {
	output.Title("🌐 Multi-Team Version Status")

	// Get current version
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Get version lock status
	isLocked, lockInfo := m.getVersionLockStatus()

	// Get team approvals
	approvals := m.getTeamApprovals()

	// Display status
	statusData := [][]string{
		{"Current Version", currentVersion},
		{"Lock Status", m.formatLockStatus(isLocked, lockInfo)},
		{"Team Approvals", m.formatApprovals(approvals)},
	}
	output.Table([]string{"Property", "Value"}, statusData)

	// Show pending proposals
	m.showPendingProposals()

	return nil
}

// syncVersions synchronizes versions across all teams
func (m *MultiTeamVersionHandler) syncVersions() error {
	output.Title("🔄 Synchronizing Versions Across Teams")

	// Check if version is locked
	isLocked, lockInfo := m.getVersionLockStatus()
	if isLocked {
		output.Warning("Version is currently locked by %s", lockInfo.team)
		output.Info("Lock reason: %s", lockInfo.reason)
		output.Info("Lock expires: %s", lockInfo.expiresAt)

		var forceSync bool
		err := huh.NewConfirm().
			Title("Force Sync?").
			Description("Do you want to force sync despite the lock?").
			Value(&forceSync).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !forceSync {
			output.Info("Sync cancelled")
			return nil
		}
	}

	// Get all team versions
	teamVersions := m.getAllTeamVersions()

	// Find the highest version
	highestVersion := m.findHighestVersion(teamVersions)

	// Propose sync to all teams
	output.Info("Proposing sync to version: %s", highestVersion)

	var confirmed bool
	err := huh.NewConfirm().
		Title("Confirm Sync").
		Description(fmt.Sprintf("Sync all teams to version %s?", highestVersion)).
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("Sync cancelled")
		return nil
	}

	// Perform the sync
	if err := m.performVersionSync(highestVersion); err != nil {
		return fmt.Errorf("failed to sync versions: %w", err)
	}

	output.Success("Successfully synchronized all teams to version %s", highestVersion)
	return nil
}

// lockVersion locks the version to prevent changes
func (m *MultiTeamVersionHandler) lockVersion() error {
	output.Title("🔒 Lock Version")

	var reason, duration string
	var team string

	// Get current team from config or prompt
	currentTeam, _ := m.git.GetConfig("at.team")
	if currentTeam == "" {
		output.Warning("No team configured. Please set your team first.")
		output.Info("Use: git @ version teams set <team-name>")
		return nil
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Team").
				Description("Your team name").
				Value(&team).
				Placeholder(currentTeam),
			huh.NewInput().
				Title("Lock Reason").
				Description("Why are you locking the version?").
				Value(&reason).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("lock reason cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Lock Duration").
				Description("How long should the lock last?").
				Options(
					huh.NewOption("1 hour", "1h"),
					huh.NewOption("4 hours", "4h"),
					huh.NewOption("1 day", "24h"),
					huh.NewOption("1 week", "168h"),
					huh.NewOption("Until manually unlocked", "0"),
				).
				Value(&duration),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get lock details: %w", err)
	}

	// Calculate expiration time
	var expiresAt time.Time
	if duration != "0" {
		durationHours, _ := strconv.Atoi(strings.TrimSuffix(duration, "h"))
		expiresAt = time.Now().Add(time.Duration(durationHours) * time.Hour)
	} else {
		expiresAt = time.Now().AddDate(0, 0, 30) // 30 days for manual unlock
	}

	// Set lock information
	lockInfo := fmt.Sprintf("%s|%s|%s|%s", team, reason, expiresAt.Format(time.RFC3339), duration)

	if err := m.git.SetConfig("at.version.lock", lockInfo); err != nil {
		return fmt.Errorf("failed to set version lock: %w", err)
	}

	// Log the lock
	if err := m.writeVersionLog(fmt.Sprintf("Version locked by %s: %s (expires: %s)", team, reason, expiresAt.Format("2006-01-02 15:04:05"))); err != nil {
		output.Warning("Failed to log version lock: %v", err)
	}

	output.Success("Version locked by %s", team)
	output.Info("Reason: %s", reason)
	output.Info("Expires: %s", expiresAt.Format("2006-01-02 15:04:05"))

	return nil
}

// unlockVersion unlocks the version
func (m *MultiTeamVersionHandler) unlockVersion() error {
	output.Title("🔓 Unlock Version")

	// Check if version is locked
	isLocked, lockInfo := m.getVersionLockStatus()
	if !isLocked {
		output.Info("Version is not currently locked")
		return nil
	}

	// Check if current user can unlock
	currentTeam, _ := m.git.GetConfig("at.team")
	if currentTeam != lockInfo.team {
		output.Warning("Version is locked by %s, but you are from team %s", lockInfo.team, currentTeam)

		var forceUnlock bool
		err := huh.NewConfirm().
			Title("Force Unlock?").
			Description("Do you want to force unlock the version?").
			Value(&forceUnlock).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !forceUnlock {
			output.Info("Unlock cancelled")
			return nil
		}
	}

	// Remove lock
	if err := m.git.SetConfig("at.version.lock", ""); err != nil {
		return fmt.Errorf("failed to remove version lock: %w", err)
	}

	// Log the unlock
	if err := m.writeVersionLog(fmt.Sprintf("Version unlocked by %s", currentTeam)); err != nil {
		output.Warning("Failed to log version unlock: %v", err)
	}

	output.Success("Version unlocked successfully")
	return nil
}

// proposeVersion proposes a version change
func (m *MultiTeamVersionHandler) proposeVersion(args []string) error {
	output.Title("📋 Propose Version Change")

	var newVersion, reason string
	var teams []string

	// Parse arguments
	if len(args) >= 1 {
		newVersion = args[0]
	}
	if len(args) >= 2 {
		reason = strings.Join(args[1:], " ")
	}

	// If not provided via args, prompt for them
	if newVersion == "" {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("New Version").
					Description("Enter the proposed version (e.g., 2.1.0)").
					Value(&newVersion).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("version cannot be empty")
						}
						// Validate semantic version format
						parts := strings.Split(s, ".")
						if len(parts) != 3 {
							return fmt.Errorf("version must be in format MAJOR.MINOR.FIX")
						}
						for _, part := range parts {
							if _, err := strconv.Atoi(part); err != nil {
								return fmt.Errorf("version components must be numbers")
							}
						}
						return nil
					}),
				huh.NewText().
					Title("Reason").
					Description("Explain why this version change is needed").
					Value(&reason).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("reason cannot be empty")
						}
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed to get proposal details: %w", err)
		}
	}

	// Get current team
	currentTeam, _ := m.git.GetConfig("at.team")
	if currentTeam == "" {
		output.Warning("No team configured. Please set your team first.")
		return nil
	}

	// Get all teams that need to approve
	teams = m.getTeamsForApproval()

	// Create proposal
	proposalID := m.generateProposalID()
	proposal := fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		proposalID, newVersion, reason, currentTeam, time.Now().Format(time.RFC3339), strings.Join(teams, ","))

	// Store proposal
	if err := m.git.SetConfig("at.version.proposal."+proposalID, proposal); err != nil {
		return fmt.Errorf("failed to store proposal: %w", err)
	}

	// Log the proposal
	if err := m.writeVersionLog(fmt.Sprintf("Version proposal %s: %s -> %s by %s", proposalID, m.getCurrentVersionString(), newVersion, currentTeam)); err != nil {
		output.Warning("Failed to log proposal: %v", err)
	}

	output.Success("Version proposal created: %s", proposalID)
	output.Info("Proposed version: %s", newVersion)
	output.Info("Reason: %s", reason)
	output.Info("Teams to approve: %s", strings.Join(teams, ", "))

	return nil
}

// approveVersion approves a version proposal
func (m *MultiTeamVersionHandler) approveVersion(args []string) error {
	output.Title("✅ Approve Version Proposal")

	if len(args) == 0 {
		return fmt.Errorf("proposal ID required")
	}

	proposalID := args[0]

	// Get proposal
	proposal, err := m.git.GetConfig("at.version.proposal." + proposalID)
	if err != nil || proposal == "" {
		return fmt.Errorf("proposal %s not found", proposalID)
	}

	// Parse proposal
	parts := strings.Split(proposal, "|")
	if len(parts) < 6 {
		return fmt.Errorf("invalid proposal format")
	}

	newVersion := parts[1]
	reason := parts[2]
	proposingTeam := parts[3]
	teams := strings.Split(parts[5], ",")

	// Get current team
	currentTeam, _ := m.git.GetConfig("at.team")
	if currentTeam == "" {
		output.Warning("No team configured. Please set your team first.")
		return nil
	}

	// Check if current team can approve
	canApprove := false
	for _, team := range teams {
		if team == currentTeam {
			canApprove = true
			break
		}
	}

	if !canApprove {
		return fmt.Errorf("team %s cannot approve this proposal", currentTeam)
	}

	// Show proposal details
	output.Info("Proposal ID: %s", proposalID)
	output.Info("Proposed by: %s", proposingTeam)
	output.Info("New version: %s", newVersion)
	output.Info("Reason: %s", reason)
	output.Info("Teams to approve: %s", strings.Join(teams, ", "))

	var confirmed bool
	err = huh.NewConfirm().
		Title("Approve Proposal").
		Description(fmt.Sprintf("Approve version change to %s?", newVersion)).
		Value(&confirmed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		output.Info("Approval cancelled")
		return nil
	}

	// Record approval
	approvalKey := fmt.Sprintf("at.version.approval.%s.%s", proposalID, currentTeam)
	if err := m.git.SetConfig(approvalKey, time.Now().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("failed to record approval: %w", err)
	}

	// Check if all teams have approved
	allApproved := m.checkAllTeamsApproved(proposalID, teams)

	if allApproved {
		// Apply the version change
		if err := m.applyVersionChange(newVersion); err != nil {
			return fmt.Errorf("failed to apply version change: %w", err)
		}

		// Clean up proposal
		if err := m.git.SetConfig("at.version.proposal."+proposalID, ""); err != nil {
			output.Warning("Failed to clean up proposal: %v", err)
		}

		output.Success("Version change approved and applied: %s", newVersion)
	} else {
		output.Success("Approval recorded. Waiting for other teams...")
	}

	return nil
}

// rejectVersion rejects a version proposal
func (m *MultiTeamVersionHandler) rejectVersion(args []string) error {
	output.Title("❌ Reject Version Proposal")

	if len(args) == 0 {
		return fmt.Errorf("proposal ID required")
	}

	proposalID := args[0]

	// Get proposal
	proposal, err := m.git.GetConfig("at.version.proposal." + proposalID)
	if err != nil || proposal == "" {
		return fmt.Errorf("proposal %s not found", proposalID)
	}

	// Parse proposal
	parts := strings.Split(proposal, "|")
	if len(parts) < 6 {
		return fmt.Errorf("invalid proposal format")
	}

	newVersion := parts[1]
	reason := parts[2]
	proposingTeam := parts[3]
	teams := strings.Split(parts[5], ",")

	// Get current team
	currentTeam, _ := m.git.GetConfig("at.team")
	if currentTeam == "" {
		output.Warning("No team configured. Please set your team first.")
		return nil
	}

	// Check if current team can reject
	canReject := false
	for _, team := range teams {
		if team == currentTeam {
			canReject = true
			break
		}
	}

	if !canReject {
		return fmt.Errorf("team %s cannot reject this proposal", currentTeam)
	}

	// Show proposal details
	output.Info("Proposal ID: %s", proposalID)
	output.Info("Proposed by: %s", proposingTeam)
	output.Info("New version: %s", newVersion)
	output.Info("Reason: %s", reason)
	output.Info("Teams to approve: %s", strings.Join(teams, ", "))

	var rejectionReason string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Rejection Reason").
				Description("Why are you rejecting this proposal?").
				Value(&rejectionReason).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("rejection reason cannot be empty")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed to get rejection reason: %w", err)
	}

	// Record rejection
	rejectionKey := fmt.Sprintf("at.version.rejection.%s.%s", proposalID, currentTeam)
	rejectionData := fmt.Sprintf("%s|%s", time.Now().Format(time.RFC3339), rejectionReason)

	if err := m.git.SetConfig(rejectionKey, rejectionData); err != nil {
		return fmt.Errorf("failed to record rejection: %w", err)
	}

	// Clean up proposal
	if err := m.git.SetConfig("at.version.proposal."+proposalID, ""); err != nil {
		output.Warning("Failed to clean up proposal: %v", err)
	}

	// Log the rejection
	if err := m.writeVersionLog(fmt.Sprintf("Version proposal %s rejected by %s: %s", proposalID, currentTeam, rejectionReason)); err != nil {
		output.Warning("Failed to log rejection: %v", err)
	}

	output.Success("Version proposal rejected")
	output.Info("Reason: %s", rejectionReason)

	return nil
}

// showTeamStatus shows the status of all teams
func (m *MultiTeamVersionHandler) showTeamStatus() error {
	output.Title("👥 Team Status")

	teams := m.getAllTeams()

	var teamData [][]string
	for _, team := range teams {
		version, _ := m.getTeamVersion(team)
		lastSync, _ := m.getTeamLastSync(team)
		status := m.getTeamStatus(team)

		teamData = append(teamData, []string{
			team,
			version,
			lastSync,
			status,
		})
	}

	output.Table([]string{"Team", "Version", "Last Sync", "Status"}, teamData)
	return nil
}

// showVersionHistory shows the version change history
func (m *MultiTeamVersionHandler) showVersionHistory() error {
	output.Title("📜 Version History")

	// Read version log file
	logFile := m.getVersionLogPath()
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		output.Info("No version history found")
		return nil
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		return fmt.Errorf("failed to read version history: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Show last 20 entries
	start := len(lines) - 20
	if start < 0 {
		start = 0
	}

	output.Info("Last 20 version changes:")
	for i := start; i < len(lines); i++ {
		if lines[i] != "" {
			output.Dim("%s", lines[i])
		}
	}

	return nil
}

// Helper methods

func (m *MultiTeamVersionHandler) getCurrentVersion() (string, error) {
	major, err := m.git.GetConfig("at.major")
	if err != nil {
		return "", err
	}
	minor, err := m.git.GetConfig("at.minor")
	if err != nil {
		return "", err
	}
	fix, err := m.git.GetConfig("at.fix")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s.%s", major, minor, fix), nil
}

func (m *MultiTeamVersionHandler) getCurrentVersionString() string {
	version, _ := m.getCurrentVersion()
	return version
}

func (m *MultiTeamVersionHandler) getVersionLockStatus() (bool, lockInfo) {
	lockData, err := m.git.GetConfig("at.version.lock")
	if err != nil || lockData == "" {
		return false, lockInfo{}
	}

	parts := strings.Split(lockData, "|")
	if len(parts) < 3 {
		return false, lockInfo{}
	}

	expiresAt, _ := time.Parse(time.RFC3339, parts[2])

	// Check if lock has expired
	if time.Now().After(expiresAt) {
		// Remove expired lock
		m.git.SetConfig("at.version.lock", "")
		return false, lockInfo{}
	}

	return true, lockInfo{
		team:      parts[0],
		reason:    parts[1],
		expiresAt: expiresAt,
	}
}

func (m *MultiTeamVersionHandler) getTeamApprovals() map[string]string {
	approvals := make(map[string]string)

	// Get all approval configs
	configs, err := m.git.Run("config", "--get-regexp", "at.version.approval")
	if err != nil {
		return approvals
	}

	lines := strings.Split(strings.TrimSpace(configs), "\n")
	for _, line := range lines {
		if line != "" {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				approvalKey := parts[0]
				approvalTime := parts[1]

				// Extract team name from key
				keyParts := strings.Split(approvalKey, ".")
				if len(keyParts) >= 4 {
					team := keyParts[len(keyParts)-1]
					approvals[team] = approvalTime
				}
			}
		}
	}

	return approvals
}

func (m *MultiTeamVersionHandler) getAllTeamVersions() map[string]string {
	// This would typically query a central registry or database
	// For now, we'll simulate with local config
	teams := m.getAllTeams()
	versions := make(map[string]string)

	for _, team := range teams {
		version, _ := m.getTeamVersion(team)
		versions[team] = version
	}

	return versions
}

func (m *MultiTeamVersionHandler) getAllTeams() []string {
	// This would typically come from a central configuration
	// For now, return a default set
	return []string{"frontend", "backend", "mobile", "devops", "qa"}
}

func (m *MultiTeamVersionHandler) getTeamVersion(team string) (string, error) {
	// This would typically query team-specific repositories or configs
	// For now, return current version
	return m.getCurrentVersion()
}

func (m *MultiTeamVersionHandler) getTeamLastSync(team string) (string, error) {
	// This would typically query team sync history
	// For now, return current time
	return time.Now().Format("2006-01-02 15:04:05"), nil
}

func (m *MultiTeamVersionHandler) getTeamStatus(team string) string {
	// This would typically check team repository status
	// For now, return "active"
	return "active"
}

func (m *MultiTeamVersionHandler) getTeamsForApproval() []string {
	// This would typically come from project configuration
	// For now, return all teams except current team
	currentTeam, _ := m.git.GetConfig("at.team")
	allTeams := m.getAllTeams()

	var teams []string
	for _, team := range allTeams {
		if team != currentTeam {
			teams = append(teams, team)
		}
	}

	return teams
}

func (m *MultiTeamVersionHandler) generateProposalID() string {
	return fmt.Sprintf("prop_%d", time.Now().Unix())
}

func (m *MultiTeamVersionHandler) checkAllTeamsApproved(proposalID string, teams []string) bool {
	for _, team := range teams {
		approvalKey := fmt.Sprintf("at.version.approval.%s.%s", proposalID, team)
		approval, _ := m.git.GetConfig(approvalKey)
		if approval == "" {
			return false
		}
	}
	return true
}

func (m *MultiTeamVersionHandler) applyVersionChange(newVersion string) error {
	parts := strings.Split(newVersion, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid version format")
	}

	if err := m.git.SetConfig("at.major", parts[0]); err != nil {
		return err
	}
	if err := m.git.SetConfig("at.minor", parts[1]); err != nil {
		return err
	}
	if err := m.git.SetConfig("at.fix", parts[2]); err != nil {
		return err
	}

	return nil
}

func (m *MultiTeamVersionHandler) findHighestVersion(versions map[string]string) string {
	highest := "0.0.0"

	for _, version := range versions {
		if m.compareVersions(version, highest) > 0 {
			highest = version
		}
	}

	return highest
}

func (m *MultiTeamVersionHandler) compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < 3; i++ {
		num1 := 0
		num2 := 0

		if i < len(parts1) {
			num1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			num2, _ = strconv.Atoi(parts2[i])
		}

		if num1 > num2 {
			return 1
		}
		if num1 < num2 {
			return -1
		}
	}

	return 0
}

func (m *MultiTeamVersionHandler) performVersionSync(targetVersion string) error {
	// This would typically sync versions across all team repositories
	// For now, just update local version
	return m.applyVersionChange(targetVersion)
}

func (m *MultiTeamVersionHandler) showPendingProposals() {
	// Get all pending proposals
	configs, err := m.git.Run("config", "--get-regexp", "at.version.proposal")
	if err != nil {
		return
	}

	lines := strings.Split(strings.TrimSpace(configs), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return
	}

	output.Title("📋 Pending Proposals")

	for _, line := range lines {
		if line != "" {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				proposalKey := parts[0]
				proposalData := parts[1]

				// Extract proposal ID
				keyParts := strings.Split(proposalKey, ".")
				if len(keyParts) >= 3 {
					proposalID := keyParts[len(keyParts)-1]

					// Parse proposal data
					dataParts := strings.Split(proposalData, "|")
					if len(dataParts) >= 5 {
						newVersion := dataParts[1]
						reason := dataParts[2]
						proposingTeam := dataParts[3]

						output.Info("Proposal %s: %s -> %s by %s", proposalID, m.getCurrentVersionString(), newVersion, proposingTeam)
						output.Dim("Reason: %s", reason)
					}
				}
			}
		}
	}
}

func (m *MultiTeamVersionHandler) formatLockStatus(isLocked bool, lockInfo lockInfo) string {
	if !isLocked {
		return "Unlocked"
	}
	return fmt.Sprintf("Locked by %s (expires: %s)", lockInfo.team, lockInfo.expiresAt.Format("2006-01-02 15:04:05"))
}

func (m *MultiTeamVersionHandler) formatApprovals(approvals map[string]string) string {
	if len(approvals) == 0 {
		return "None"
	}

	var teams []string
	for team := range approvals {
		teams = append(teams, team)
	}
	return strings.Join(teams, ", ")
}

func (m *MultiTeamVersionHandler) getVersionLogPath() string {
	gitRoot, err := m.git.Run("rev-parse", "--show-toplevel")
	if err != nil {
		return ""
	}
	gitRoot = strings.TrimSpace(gitRoot)
	return filepath.Join(gitRoot, ".git", "gitat-logs", "version-changes.log")
}

func (m *MultiTeamVersionHandler) writeVersionLog(message string) error {
	logFile := m.getVersionLogPath()
	if logFile == "" {
		return fmt.Errorf("failed to get log file path")
	}

	// Create logs directory if it doesn't exist
	logsDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Open file in append mode
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Write log entry with timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	_, err = file.WriteString(logEntry)
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

func (m *MultiTeamVersionHandler) showUsage() error {
	return output.Markdown(`# Multi-Team Version Management

Manages semantic versioning across multiple teams with approval workflows.

## Usage

git @ version multi-team                    # Show multi-team status
git @ version multi-team sync               # Synchronize versions across teams
git @ version multi-team lock               # Lock version to prevent changes
git @ version multi-team unlock             # Unlock version
git @ version multi-team propose <version>  # Propose version change
git @ version multi-team approve <id>       # Approve version proposal
git @ version multi-team reject <id>        # Reject version proposal
git @ version multi-team teams              # Show team status
git @ version multi-team history            # Show version history

## Commands

### Status & Sync
- **multi-team**: Show current multi-team version status
- **sync**: Synchronize versions across all teams
- **teams**: Show status of all teams

### Version Locking
- **lock**: Lock version to prevent changes (with expiration)
- **unlock**: Unlock version (only by locking team or force)

### Approval Workflow
- **propose**: Create a version change proposal
- **approve**: Approve a version proposal
- **reject**: Reject a version proposal with reason

### History
- **history**: Show version change history

## Examples

# Show current status
git @ version multi-team

# Propose a version change
git @ version multi-team propose 2.1.0 "Add new API endpoints"

# Approve a proposal
git @ version multi-team approve prop_1234567890

# Lock version for release
git @ version multi-team lock

# Sync all teams to latest version
git @ version multi-team sync

## Team Configuration

Set your team with:
git config at.team <team-name>

Available teams: frontend, backend, mobile, devops, qa

## Approval Process

1. **Propose**: Create a version change proposal
2. **Review**: Teams review the proposal
3. **Approve/Reject**: Teams approve or reject with reasons
4. **Apply**: If all teams approve, version is automatically applied
5. **Sync**: All teams are synchronized to the new version

## Locking Mechanism

- **Purpose**: Prevent version conflicts during critical periods
- **Duration**: Configurable (1h, 4h, 1d, 1w, or manual)
- **Override**: Force unlock available for emergencies
- **Logging**: All lock/unlock actions are logged

## Best Practices

1. **Coordinate**: Communicate version changes with all teams
2. **Lock During Releases**: Lock versions during release windows
3. **Review Proposals**: Carefully review all version proposals
4. **Document Changes**: Provide clear reasons for version changes
5. **Regular Syncs**: Regularly sync versions across teams
`)
}

// lockInfo represents version lock information
type lockInfo struct {
	team      string
	reason    string
	expiresAt time.Time
}
