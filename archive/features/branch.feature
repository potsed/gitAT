Feature: Branch Configuration
  As a developer
  I want to manage my working branch
  So that I can track which branch I should be working on

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show current working branch
    When I run "git @ branch"
    Then I should see the configured working branch
    And the output should be from git config at.branch

  Scenario: Set working branch
    When I run "git @ branch feature-auth"
    Then the working branch should be set to "feature-auth"
    And git config at.branch should be "feature-auth"
    And I should see "Branch updated to: feature-auth from"

  Scenario: Set working branch to current branch
    Given I am on branch "feature-auth"
    When I run "git @ branch -s"
    Then the working branch should be set to "feature-auth"
    And git config at.branch should be "feature-auth"

  Scenario: Set working branch to current branch with dot
    Given I am on branch "feature-auth"
    When I run "git @ branch ."
    Then the working branch should be set to "feature-auth"
    And git config at.branch should be "feature-auth"

  Scenario: Show current git branch
    When I run "git @ branch -c"
    Then I should see the current git branch
    And the output should match "git rev-parse --abbrev-ref HEAD"

  Scenario: Create new working branch
    When I run "git @ branch -n"
    Then a new branch should be created with timestamp
    And the working branch should be set to the new branch
    And I should see "New working branch created and set:"

  Scenario: Show help
    When I run "git @ branch --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ branch"

  Scenario: Handle empty working branch
    Given no working branch is set
    When I run "git @ branch"
    Then I should see an empty line
    And the exit code should be 0 