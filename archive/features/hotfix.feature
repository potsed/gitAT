Feature: Hotfix Branch Creation
  As a developer using GitAT
  I want to create hotfix branches for urgent fixes
  So that I can quickly address critical issues that need immediate deployment

  Background:
    Given I am in a git repository
    And I have a trunk branch configured
    And I am on a feature branch

  Scenario: Create hotfix with interactive name prompt
    When I run "git @ hotfix"
    Then it should prompt for a hotfix branch name
    And it should show recommended naming format
    And when I enter "hotfix/fix-login-bug"
    Then it should create the hotfix branch successfully
    And it should switch to the new hotfix branch
    And it should set the working branch to the hotfix branch

  Scenario: Create hotfix with specific name
    When I run "git @ hotfix fix-login-bug"
    Then it should create a hotfix branch named "fix-login-bug"
    And it should switch to the new hotfix branch
    And it should set the working branch to the hotfix branch

  Scenario: Create hotfix with name option
    When I run "git @ hotfix -n security-patch"
    Then it should create a hotfix branch named "security-patch"
    And it should switch to the new hotfix branch

  Scenario: Create hotfix with --name option
    When I run "git @ hotfix --name urgent-database-fix"
    Then it should create a hotfix branch named "urgent-database-fix"
    And it should switch to the new hotfix branch

  Scenario: Create hotfix with recommended naming format
    When I run "git @ hotfix fix-login-bug"
    Then it should create a hotfix branch named "hotfix-fix-login-bug"
    And it should switch to the new hotfix branch

  Scenario: Show help for hotfix command
    When I run "git @ hotfix -h"
    Then it should display the usage information
    And it should show all available options
    And it should show examples

  Scenario: Show help with --help
    When I run "git @ hotfix --help"
    Then it should display the same usage information as -h

  Scenario: Show help with help
    When I run "git @ hotfix help"
    Then it should display the same usage information as -h

  Scenario: Error when not in git repository
    Given I am not in a git repository
    When I run "git @ hotfix fix-login-bug"
    Then it should show "Error: Not in a git repository"

  Scenario: Error when in detached HEAD state
    Given I am in a detached HEAD state
    When I run "git @ hotfix fix-login-bug"
    Then it should show "Error: Not on a branch (detached HEAD state)"

  Scenario: Error when trunk branch does not exist
    Given the trunk branch does not exist
    When I run "git @ hotfix fix-login-bug"
    Then it should show "Error: Trunk branch 'main' does not exist"
    And it should suggest configuring the trunk branch

  Scenario: Error when hotfix branch already exists
    Given a hotfix branch "fix-login-bug" already exists
    When I run "git @ hotfix fix-login-bug"
    Then it should show "Error: Branch 'fix-login-bug' already exists"

  Scenario: Error for invalid branch name
    When I run "git @ hotfix invalid@name"
    Then it should show "Error: Invalid branch name 'invalid@name'"
    And it should explain valid characters

  Scenario: Error for empty branch name
    When I run "git @ hotfix"
    And I enter an empty name
    Then it should show "Error: Hotfix name cannot be empty"

  Scenario: Error for reserved branch names
    When I run "git @ hotfix HEAD"
    Then it should show "Error: Invalid branch name 'HEAD'"

  Scenario: Error for reserved branch names - master
    When I run "git @ hotfix master"
    Then it should show "Error: Invalid branch name 'master'"

  Scenario: Error for reserved branch names - main
    When I run "git @ hotfix main"
    Then it should show "Error: Invalid branch name 'main'"

  Scenario: Error for reserved branch names - develop
    When I run "git @ hotfix develop"
    Then it should show "Error: Invalid branch name 'develop'"

  Scenario: Warning for uncommitted changes
    Given I have uncommitted changes
    When I run "git @ hotfix fix-login-bug"
    Then it should show a warning about uncommitted changes
    And it should ask for confirmation to continue

  Scenario: Continue with uncommitted changes after confirmation
    Given I have uncommitted changes
    When I run "git @ hotfix fix-login-bug" and confirm with "y"
    Then it should proceed to create the hotfix branch
    And it should save the uncommitted changes to WIP

  Scenario: Abort with uncommitted changes after rejection
    Given I have uncommitted changes
    When I run "git @ hotfix fix-login-bug" and reject with "n"
    Then it should abort the operation

  Scenario: Save WIP state before creating hotfix
    Given I have uncommitted changes
    When I run "git @ hotfix fix-login-bug"
    Then it should save the current work state using git @ wip -s
    And it should show "Saving current work state..."

  Scenario: Switch to trunk branch
    When I run "git @ hotfix fix-login-bug"
    Then it should switch to the trunk branch
    And it should show "Switching to trunk branch: main"

  Scenario: Update trunk branch from remote
    Given I have a remote origin configured
    When I run "git @ hotfix fix-login-bug"
    Then it should pull latest changes from remote trunk branch
    And it should show "Updating trunk branch..."

  Scenario: Create hotfix branch from trunk
    When I run "git @ hotfix fix-login-bug"
    Then it should create a new branch from the trunk branch
    And it should switch to the new hotfix branch
    And it should show "Creating hotfix branch: fix-login-bug"

  Scenario: Set working branch to hotfix branch
    When I run "git @ hotfix fix-login-bug"
    Then it should set the working branch to the hotfix branch
    And it should show "Setting working branch to hotfix branch..."

  Scenario: Success message with status information
    When I run "git @ hotfix fix-login-bug"
    Then it should show "✅ Hotfix branch 'fix-login-bug' created successfully!"
    And it should show current status information
    And it should show next steps

  Scenario: Next steps guidance
    When I run "git @ hotfix fix-login-bug"
    Then it should show next steps including:
    And it should mention "git @ save 'Fix description'"
    And it should mention "git @ pr 'Hotfix: description'"
    And it should mention "git @ release -p (patch release)"

  Scenario: Use configured trunk branch
    Given the trunk branch is configured as "develop"
    When I run "git @ hotfix fix-login-bug"
    Then it should use "develop" as the trunk branch
    And it should create the hotfix branch from develop

  Scenario: Fallback to main when no trunk configured
    Given no trunk branch is configured
    When I run "git @ hotfix fix-login-bug"
    Then it should use "main" as the trunk branch
    And it should create the hotfix branch from main

  Scenario: Error for missing name option value
    When I run "git @ hotfix -n"
    Then it should show "Error: --name requires a value"

  Scenario: Error for missing --name option value
    When I run "git @ hotfix --name"
    Then it should show "Error: --name requires a value"

  Scenario: Error for unknown option
    When I run "git @ hotfix --unknown"
    Then it should show "Error: Unknown option '--unknown'"

  Scenario: Integration with existing GitAT commands
    When I run "git @ hotfix fix-login-bug"
    Then it should use git @ wip -s to save work state
    And it should use git @ branch to set working branch
    And it should use git config at.trunk to get trunk branch

  Scenario: Hotfix branch naming recommendations
    When I run "git @ hotfix"
    Then it should show "Format: description (e.g., fix-login-bug)"
    And it should show "Branch will be created as: hotfix-description"

  Scenario: Hotfix workflow integration
    Given I have created a hotfix branch
    When I make changes to fix the issue
    And I run "git @ save 'Fix critical login bug'"
    And I run "git @ pr 'Hotfix: Fix critical login bug'"
    Then the hotfix should be ready for review and deployment

  Scenario: Hotfix with patch release suggestion
    Given I have created a hotfix branch
    When I complete the hotfix
    Then it should suggest "git @ release -p (patch release)"
    And this should create a patch version bump 