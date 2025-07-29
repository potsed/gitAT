Feature: Work Branch Creation (Conventional Commits)
  As a developer using GitAT
  I want to create work branches following Conventional Commits specification
  So that I can maintain organized development workflow with proper commit types

  Background:
    Given I am in a git repository
    And I have a trunk branch configured
    And I am on a feature branch

  Scenario: Create feature branch with description
    When I run "git @ work feature add-user-auth"
    Then it should create a feature branch named "feature-add-user-auth"
    And it should switch to the new feature branch
    And it should set the working branch to the feature branch

  Scenario: Create hotfix branch with description
    When I run "git @ work hotfix fix-login-bug"
    Then it should create a hotfix branch named "hotfix-fix-login-bug"
    And it should switch to trunk branch first
    And it should create the hotfix branch from trunk
    And it should switch to the new hotfix branch

  Scenario: Create bugfix branch with description
    When I run "git @ work bugfix fix-crash-on-startup"
    Then it should create a bugfix branch named "bugfix-fix-crash-on-startup"
    And it should switch to the new bugfix branch

  Scenario: Create docs branch with description
    When I run "git @ work docs update-api-documentation"
    Then it should create a docs branch named "docs-update-api-documentation"
    And it should switch to the new docs branch

  Scenario: Create chore branch with description
    When I run "git @ work chore update-dependencies"
    Then it should create a chore branch named "chore-update-dependencies"
    And it should switch to the new chore branch

  Scenario: Create interactive work branch
    When I run "git @ work feature"
    Then it should prompt for a description
    And when I enter "add-payment-gateway"
    Then it should create a feature branch named "feature-add-payment-gateway"

  Scenario: Create work branch with full name option
    When I run "git @ work -n custom-feature-branch"
    Then it should create a branch named "custom-feature-branch"
    And it should switch to the new branch

  Scenario: Show help for work command
    When I run "git @ work -h"
    Then it should display the usage information
    And it should show all available work types
    And it should show examples

  Scenario: Show help with --help
    When I run "git @ work --help"
    Then it should display the same usage information as -h

  Scenario: Show help with help
    When I run "git @ work help"
    Then it should display the same usage information as -h

  Scenario: Error when not in git repository
    Given I am not in a git repository
    When I run "git @ work feature add-auth"
    Then it should show "Error: Not in a git repository"

  Scenario: Error when in detached HEAD state
    Given I am in a detached HEAD state
    When I run "git @ work feature add-auth"
    Then it should show "Error: Not on a branch (detached HEAD state)"

  Scenario: Error when work type is missing
    When I run "git @ work"
    Then it should show "Error: Work type is required"
    And it should list available work types

  Scenario: Error for invalid work type
    When I run "git @ work invalid-type add-auth"
    Then it should show "Error: Invalid work type 'invalid-type'"
    And it should list available work types

  Scenario: Error for missing name option value
    When I run "git @ work -n"
    Then it should show "Error: --name requires a value"

  Scenario: Error for missing --name option value
    When I run "git @ work --name"
    Then it should show "Error: --name requires a value"

  Scenario: Error for too many arguments
    When I run "git @ work feature add-auth extra-arg"
    Then it should show "Error: Too many arguments"

  Scenario: Error when branch already exists
    Given a feature branch "feature-add-auth" already exists
    When I run "git @ work feature add-auth"
    Then it should show "Error: Branch 'feature-add-auth' already exists"

  Scenario: Error for invalid branch name
    When I run "git @ work feature invalid@name"
    Then it should show "Error: Invalid branch name 'feature-invalid@name'"
    And it should explain valid characters

  Scenario: Warning for uncommitted changes
    Given I have uncommitted changes
    When I run "git @ work feature add-auth"
    Then it should show a warning about uncommitted changes
    And it should ask for confirmation to continue

  Scenario: Continue with uncommitted changes after confirmation
    Given I have uncommitted changes
    When I run "git @ work feature add-auth" and confirm with "y"
    Then it should proceed to create the feature branch
    And it should save the uncommitted changes to WIP

  Scenario: Abort with uncommitted changes after rejection
    Given I have uncommitted changes
    When I run "git @ work feature add-auth" and reject with "n"
    Then it should abort the operation

  Scenario: Save WIP state before creating work branch
    Given I have uncommitted changes
    When I run "git @ work feature add-auth"
    Then it should save the current work state using git @ wip -s
    And it should show "Saving current work state..."

  Scenario: Set working branch to new work branch
    When I run "git @ work feature add-auth"
    Then it should set the working branch to the new feature branch
    And it should show "Setting working branch to feature branch..."

  Scenario: Success message with status information
    When I run "git @ work feature add-auth"
    Then it should show "✅ feature branch 'feature-add-auth' created successfully!"
    And it should show current status information
    And it should show next steps

  Scenario: Next steps guidance for feature
    When I run "git @ work feature add-auth"
    Then it should show next steps including:
    And it should mention "git @ save '[FEATURE] Description of changes'"
    And it should mention "git @ pr 'Feature: Description of changes'"
    And it should mention "git @ release -m (minor release)"

  Scenario: Next steps guidance for hotfix
    When I run "git @ work hotfix fix-bug"
    Then it should show next steps including:
    And it should mention "git @ save '[HOTFIX] Description of changes'"
    And it should mention "git @ pr 'Hotfix: Description of changes'"
    And it should mention "git @ release -p (patch release)"

  Scenario: Next steps guidance for release
    When I run "git @ work release v2.0.0"
    Then it should show next steps including:
    And it should mention "git @ save '[RELEASE] Description of changes'"
    And it should mention "git @ pr 'Release: Description of changes'"
    And it should mention "git @ release -M (major release)"

  Scenario: Hotfix branch creation from trunk
    When I run "git @ work hotfix fix-critical-bug"
    Then it should switch to trunk branch first
    And it should update trunk branch from remote
    And it should create hotfix branch from trunk
    And it should switch to the new hotfix branch

  Scenario: Error when trunk branch does not exist for hotfix
    Given the trunk branch does not exist
    When I run "git @ work hotfix fix-bug"
    Then it should show "Error: Trunk branch 'main' does not exist"
    And it should suggest configuring the trunk branch

  Scenario: Allowed work types
    When I run "git @ work hotfix test"
    Then it should create a hotfix branch successfully

    When I run "git @ work feature test"
    Then it should create a feature branch successfully

    When I run "git @ work bugfix test"
    Then it should create a bugfix branch successfully

    When I run "git @ work release test"
    Then it should create a release branch successfully

    When I run "git @ work chore test"
    Then it should create a chore branch successfully

    When I run "git @ work docs test"
    Then it should create a docs branch successfully

    When I run "git @ work style test"
    Then it should create a style branch successfully

    When I run "git @ work refactor test"
    Then it should create a refactor branch successfully

    When I run "git @ work perf test"
    Then it should create a perf branch successfully

    When I run "git @ work test test"
    Then it should create a test branch successfully

    When I run "git @ work ci test"
    Then it should create a ci branch successfully

    When I run "git @ work build test"
    Then it should create a build branch successfully

    When I run "git @ work revert test"
    Then it should create a revert branch successfully

  Scenario: Integration with existing GitAT commands
    When I run "git @ work feature add-auth"
    Then it should use git @ wip -s to save work state
    And it should use git @ branch to set working branch
    And it should use git config at.trunk to get trunk branch

  Scenario: Work branch naming convention
    When I run "git @ work feature add-user-authentication"
    Then it should create a branch named "feature-add-user-authentication"
    And the branch name should follow the format "type-description"

  Scenario: Work branch workflow integration
    Given I have created a feature branch
    When I make changes to implement the feature
    And I run "git @ save 'Add user authentication system'"
    And I run "git @ pr 'Feature: Add user authentication system'"
    Then the feature should be ready for review and deployment

  Scenario: Conventional Commits integration
    Given I am on a feature branch
    When I run "git @ save 'Add login functionality'"
    Then the commit message should start with "[FEATURE]"

    Given I am on a hotfix branch
    When I run "git @ save 'Fix critical security bug'"
    Then the commit message should start with "[HOTFIX]"

    Given I am on a bugfix branch
    When I run "git @ save 'Fix crash on startup'"
    Then the commit message should start with "[BUGFIX]"

    Given I am on a docs branch
    When I run "git @ save 'Update API documentation'"
    Then the commit message should start with "[DOCS]"

  Scenario: Work type prefix detection
    Given I am on a branch named "feature-add-auth"
    When I run "git @ save 'Add authentication'"
    Then the commit message should be "[FEATURE] <label> Add authentication"

    Given I am on a branch named "hotfix-fix-login"
    When I run "git @ save 'Fix login bug'"
    Then the commit message should be "[HOTFIX] <label> Fix login bug"

    Given I am on a branch named "chore-update-deps"
    When I run "git @ save 'Update dependencies'"
    Then the commit message should be "[CHORE] <label> Update dependencies" 