Feature: Enhanced Squash Command
  As a developer using GitAT
  I want to squash commits for clean history and PR preparation
  So that I can maintain clean commit history and prepare branches for review

  Background:
    Given I am in a git repository
    And I have a feature branch with multiple commits

  Scenario: Basic squash to target branch
    When I run "git @ squash develop"
    Then it should reset the current branch to develop HEAD
    And it should keep all changes staged for commit
    And it should show "Squashed branch [current] back to develop"

  Scenario: Squash to master branch
    When I run "git @ squash master"
    Then it should reset the current branch to master HEAD
    And it should keep all changes staged for commit
    And it should show "Squashed branch [current] back to master"

  Scenario: Squash and save
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
    And it should show examples

  Scenario: Show help with --help
    When I run "git @ squash --help"
    Then it should display the same usage information as -h

  Scenario: Show help with help
    When I run "git @ squash help"
    Then it should display the same usage information as -h

  Scenario: Error for missing target branch
    When I run "git @ squash"
    Then it should show usage information
    And it should exit with error

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

  Scenario: Show automatic squashing status when enabled
    Given automatic PR squashing is enabled
    When I run "git @ squash --auto status"
    Then it should show "Status: ✅ ENABLED"
    And it should show "Commits will be automatically squashed before creating PRs"
    And it should show override instructions

  Scenario: Show automatic squashing status when disabled
    Given automatic PR squashing is disabled
    When I run "git @ squash --auto status"
    Then it should show "Status: ❌ DISABLED"
    And it should show "PRs will be created with all commits as-is"
    And it should show override instructions

  Scenario: Enable with alternative commands
    When I run "git @ squash --auto true"
    Then it should enable the squash setting

  Scenario: Enable with enable command
    When I run "git @ squash --auto enable"
    Then it should enable the squash setting

  Scenario: Enable with 1 command
    When I run "git @ squash --auto 1"
    Then it should enable the squash setting

  Scenario: Disable with alternative commands
    When I run "git @ squash --auto false"
    Then it should disable the squash setting

  Scenario: Disable with disable command
    When I run "git @ squash --auto disable"
    Then it should disable the squash setting

  Scenario: Disable with 0 command
    When I run "git @ squash --auto 0"
    Then it should disable the squash setting

  Scenario: Show status with show command
    When I run "git @ squash --auto show"
    Then it should show the current status

  Scenario: Show status with check command
    When I run "git @ squash --auto check"
    Then it should show the current status

  Scenario: Error for invalid auto action
    When I run "git @ squash --auto invalid"
    Then it should show "Error: Invalid auto action 'invalid'"
    And it should suggest valid options

  Scenario: Error for missing auto value
    When I run "git @ squash --auto"
    Then it should show "Error: --auto requires a value (on|off|status)"
    And it should exit with error

  Scenario: Default status when no setting configured
    Given no PR squash setting is configured
    When I run "git @ squash --auto status"
    Then it should show "Status: ❌ DISABLED"

  Scenario: Configuration persistence
    When I run "git @ squash --auto on"
    And I run "git @ squash --auto status"
    Then it should show "Status: ✅ ENABLED"
    And the setting should persist across git sessions

  Scenario: Integration with git @ pr when enabled
    Given automatic PR squashing is enabled
    And I have multiple commits on my feature branch
    When I run "git @ pr 'Test PR'"
    Then it should show "Auto-squashing commits before creating PR"
    And it should squash the commits
    And it should show "Commits squashed successfully"
    And it should create the PR with squashed commits

  Scenario: Integration with git @ pr when disabled
    Given automatic PR squashing is disabled
    And I have multiple commits on my feature branch
    When I run "git @ pr 'Test PR'"
    Then it should not show squashing messages
    And it should create the PR with all commits as-is

  Scenario: Force squash override when disabled
    Given automatic PR squashing is disabled
    And I have multiple commits on my feature branch
    When I run "git @ pr 'Test PR' -s"
    Then it should show "Auto-squashing commits before creating PR"
    And it should squash the commits despite the setting being disabled

  Scenario: Force no squash override when enabled
    Given automatic PR squashing is enabled
    And I have multiple commits on my feature branch
    When I run "git @ pr 'Test PR' -S"
    Then it should not show squashing messages
    And it should create the PR with all commits as-is despite the setting being enabled

  Scenario: PR squash with cherry-pick conflicts
    Given I have conflicting commits that cannot be cherry-picked
    When I run "git @ squash --pr"
    Then it should show "Error: Failed to cherry-pick commit"
    And it should clean up temporary branches
    And it should restore the original branch state

  Scenario: PR squash cleanup on failure
    Given I have commits that will cause squash to fail
    When I run "git @ squash --pr"
    Then it should clean up temporary branches
    And it should restore the original branch state

  Scenario: PR squash preserves commit messages
    Given I have multiple commits with meaningful messages
    When I run "git @ squash --pr"
    Then it should squash the commits
    And the final commit should have a meaningful message

  Scenario: PR squash with custom trunk branch
    Given the trunk branch is configured as "develop"
    And I have multiple commits on my feature branch
    When I run "git @ squash --pr"
    Then it should use "develop" as the target branch
    And it should squash commits ahead of develop 