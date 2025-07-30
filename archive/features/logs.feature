Feature: Show Recent Logs
  As a developer
  I want to see recent commit history
  So that I can review recent changes

  Background:
    Given I am in a git repository
    And GitAT is initialized
    And there are commits in the history

  Scenario: Show recent logs
    When I run "git @ logs"
    Then I should see the last 10 commits
    And each line should show commit hash, author, date, and message
    And the output should be from "git log -10 --pretty=oneline --abbrev-commit"

  Scenario: Show logs with commits
    Given there are 5 commits in history
    When I run "git @ logs"
    Then I should see 5 lines of output
    And each line should contain a commit hash

  Scenario: Show logs with no commits
    Given there are no commits in history
    When I run "git @ logs"
    Then I should see no output
    And the exit code should be 0

  Scenario: Show help
    When I run "git @ logs --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ logs"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ logs"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Show logs in compact format
    When I run "git @ logs"
    Then the output should be in oneline format
    And commit hashes should be abbreviated
    And each line should be compact 