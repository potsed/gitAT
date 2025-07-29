Feature: Enhanced Squash Command with Auto-Detection
  As a developer using GitAT
  I want to squash commits for clean history and PR preparation
  So that I can maintain clean commit history and prepare branches for review

  Background:
    Given I am in a git repository
    And I have a feature branch with multiple commits

  Scenario: Auto-detect parent branch and squash
    Given I have a feature branch created from "develop"
    When I run "git @ squash"
    Then it should auto-detect the parent branch as "develop"
    And it should show "Auto-detected parent branch: develop"
    And it should reset the current branch to develop HEAD
    And it should keep all changes staged for commit
    And it should show "Squashed branch [current] back to develop"

  Scenario: Auto-detect parent branch with upstream tracking
    Given I have a feature branch with upstream tracking to "main"
    When I run "git @ squash"
    Then it should auto-detect the parent branch as "main"
    And it should show "Auto-detected parent branch: main"
    And it should reset the current branch to main HEAD

  Scenario: Auto-detect parent branch with git config merge
    Given I have a feature branch with git config branch.<name>.merge set to "develop"
    When I run "git @ squash"
    Then it should auto-detect the parent branch as "develop"
    And it should show "Auto-detected parent branch: develop"

  Scenario: Auto-detect parent branch with branch divergence analysis
    Given I have a feature branch that diverged from "main"
    When I run "git @ squash"
    Then it should auto-detect the parent branch as "main"
    And it should show "Auto-detected parent branch: main"

  Scenario: Auto-detect parent branch fallback to configured trunk
    Given no upstream tracking is configured
    And the trunk branch is configured as "main"
    When I run "git @ squash"
    Then it should auto-detect the parent branch as "main"
    And it should show "Auto-detected parent branch: main"

  Scenario: Auto-detect parent branch fallback to common names
    Given no upstream tracking is configured
    And no trunk branch is configured
    And a "main" branch exists
    When I run "git @ squash"
    Then it should auto-detect the parent branch as "main"
    And it should show "Auto-detected parent branch: main"

  Scenario: Error when no parent can be detected
    Given I have an isolated branch with no parent relationship
    When I run "git @ squash"
    Then it should show "Error: Could not auto-detect parent branch"
    And it should suggest specifying a target branch
    And it should exit with error

  Scenario: Basic squash to explicit target branch
    When I run "git @ squash develop"
    Then it should reset the current branch to develop HEAD
    And it should keep all changes staged for commit
    And it should show "Squashed branch [current] back to develop"

  Scenario: Squash to master branch
    When I run "git @ squash master"
    Then it should reset the current branch to master HEAD
    And it should keep all changes staged for commit
    And it should show "Squashed branch [current] back to master"

  Scenario: Squash and save with auto-detection
    When I run "git @ squash -s"
    Then it should auto-detect the parent branch
    And it should reset the current branch to parent HEAD
    And it should automatically run "git @ save"
    And it should create a new commit with all changes

  Scenario: Squash and save with explicit branch
    When I run "git @ squash develop -s"
    Then it should reset the current branch to develop HEAD
    And it should automatically run "git @ save"
    And it should create a new commit with all changes

  Scenario: Squash and save with --save flag
    When I run "git @ squash develop --save"
    Then it should reset the current branch to develop HEAD
    And it should automatically run "git @ save"
    And it should create a new commit with all changes

  Scenario: Error for non-existent target branch
    When I run "git @ squash nonexistent-branch"
    Then it should show "ERROR: Branch \"nonexistent-branch\" does not exist locally"
    And it should exit with error

  Scenario: Show help for squash command
    When I run "git @ squash -h"
    Then it should display the usage information
    And it should show all available options
    And it should show examples including auto-detection

  Scenario: Show help with --help
    When I run "git @ squash --help"
    Then it should display the same usage information as -h

  Scenario: Show help with help
    When I run "git @ squash help"
    Then it should display the same usage information as -h

  Scenario: PR squash mode
    Given I have multiple commits on my feature branch
    And the trunk branch is configured as "main"
    When I run "git @ squash --pr"
    Then it should use the configured trunk branch as target
    And it should squash commits ahead of the trunk branch
    And it should show "Successfully squashed [X] commits into one for PR"

  Scenario: PR squash mode with --pr flag
    Given I have multiple commits on my feature branch
    And the trunk branch is configured as "main"
    When I run "git @ squash -p"
    Then it should use the configured trunk branch as target
    And it should squash commits ahead of the trunk branch
    And it should show "Successfully squashed [X] commits into one for PR"

  Scenario: PR squash with no commits to squash
    Given I have only one commit on my feature branch
    When I run "git @ squash --pr"
    Then it should show "Only one commit or no commits to squash"
    And it should not perform any squashing

  Scenario: PR squash error when on trunk branch
    Given I am on the trunk branch
    When I run "git @ squash --pr"
    Then it should show "Error: Cannot squash PR from [trunk] to itself"
    And it should exit with error

  Scenario: PR squash error in detached HEAD
    Given I am in a detached HEAD state
    When I run "git @ squash --pr"
    Then it should show "Error: Not on a branch (detached HEAD state)"
    And it should exit with error

  Scenario: Enable automatic PR squashing
    When I run "git @ squash --auto on"
    Then it should enable the squash setting
    And it should show "Automatic PR squashing enabled"
    And it should store the setting as "at.pr.squash true"

  Scenario: Disable automatic PR squashing
    When I run "git @ squash --auto off"
    Then it should disable the squash setting
    And it should show "Automatic PR squashing disabled"
    And it should store the setting as "at.pr.squash false"

  Scenario: Show automatic PR squashing status
    When I run "git @ squash --auto status"
    Then it should show the current squashing status
    And it should show whether automatic squashing is enabled or disabled
    And it should show override commands

  Scenario: Error for invalid auto action
    When I run "git @ squash --auto invalid"
    Then it should show "Error: Invalid auto action 'invalid'"
    And it should suggest valid options: 'on', 'off', or 'status'
    And it should exit with error

  Scenario: Error for missing auto action value
    When I run "git @ squash --auto"
    Then it should show "Error: --auto requires a value (on|off|status)"
    And it should exit with error

  Scenario: Auto-detection with orphan branch
    Given I have an orphan branch with no parent
    When I run "git @ squash"
    Then it should show "Error: Could not auto-detect parent branch"
    And it should suggest specifying a target branch

  Scenario: Auto-detection priority order
    Given I have multiple detection methods available
    When I run "git @ squash"
    Then it should prioritize git config branch.<name>.merge
    And if not available, it should use upstream tracking (@{u})
    And if not available, it should use branch divergence analysis
    And if not available, it should use configured trunk branch
    And if not available, it should use common branch names

  Scenario: Integration with existing GitAT workflow
    Given I have created a feature branch using "git @ work feature"
    And I have made multiple commits
    When I run "git @ squash"
    Then it should auto-detect the parent branch
    And it should squash all commits into one
    And it should prepare the branch for PR creation

  Scenario: Squash workflow with Conventional Commits
    Given I am on a feature branch with Conventional Commits prefix
    When I run "git @ squash -s"
    Then it should auto-detect the parent branch
    And it should create a commit with the appropriate type prefix
    And it should maintain the Conventional Commits format 