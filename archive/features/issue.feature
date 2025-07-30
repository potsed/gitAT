Feature: Issue Configuration
  As a developer
  I want to set and get the issue/task identifier
  So that I can track which issue I'm working on

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show current issue ID
    When I run "git @ issue"
    Then I should see the current issue ID
    And the output should be from git config at.task

  Scenario: Set issue ID
    When I run "git @ issue PROJ-123"
    Then the issue ID should be set to "PROJ-123"
    And git config at.task should be "PROJ-123"
    And I should see "Task updated to: PROJ-123 from"

  Scenario: Set issue ID with different formats
    When I run "git @ issue BUG-456"
    Then the issue ID should be set to "BUG-456"
    And git config at.task should be "BUG-456"

  Scenario: Set issue ID with numbers only
    When I run "git @ issue 789"
    Then the issue ID should be set to "789"
    And git config at.task should be "789"

  Scenario: Show help
    When I run "git @ issue --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ issue"

  Scenario: Handle empty issue ID
    Given no issue ID is set
    When I run "git @ issue"
    Then I should see an empty line
    And the exit code should be 0 