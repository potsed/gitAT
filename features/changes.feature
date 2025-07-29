Feature: Show Uncommitted Changes
  As a developer
  I want to see uncommitted changes
  So that I can review what files have been modified

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show uncommitted changes
    Given I have modified files
    When I run "git @ changes"
    Then I should see a list of modified file names
    And the output should be from "git diff --name-only"

  Scenario: Show no changes
    Given I have no uncommitted changes
    When I run "git @ changes"
    Then I should see no output
    And the exit code should be 0

  Scenario: Show help
    When I run "git @ changes --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ changes"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ changes"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Show staged and unstaged changes
    Given I have both staged and unstaged changes
    When I run "git @ changes"
    Then I should see all modified file names
    And the output should not show file status indicators 