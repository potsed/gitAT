Feature: Show Branch Hash Information
  As a developer
  I want to see branch hash information and recent commits
  So that I can understand branch relationships, status, and recent history

  Background:
    Given I am in a git repository
    And GitAT is initialized
    And I have remote branches

  Scenario: Show branch hash information with recent commits
    When I run "git @ hash"
    Then I should see branch information
    And I should see current branch hash
    And I should see remote branch information
    And I should see ahead/behind information
    And I should see "RECENT COMMITS (last 5):"

  Scenario: Show current branch status with commits
    Given I am on branch "feature-auth"
    When I run "git @ hash"
    Then I should see "BRANCH feature-auth"
    And I should see current branch hash
    And I should see remote branch hash if it exists
    And I should see recent commits section

  Scenario: Show ahead/behind information
    Given I am ahead of remote branch
    When I run "git @ hash"
    Then I should see "CURRENT is X behind and Y ahead of origin/feature-auth"
    Where X and Y are numbers

  Scenario: Show develop branch relationship
    Given there is a develop branch
    When I run "git @ hash"
    Then I should see "CURRENT is X behind and Y ahead of origin/develop"
    Where X and Y are numbers

  Scenario: Show master branch relationship
    Given there is a master branch
    When I run "git @ hash"
    Then I should see "CURRENT is X behind and Y ahead of origin/master"
    Where X and Y are numbers

  Scenario: Display recent commits with hashes and messages
    Given there are at least 5 commits in history
    When I run "git @ hash"
    Then I should see "RECENT COMMITS (last 5):"
    And I should see 5 commit lines
    And each line should show hash, committer, and message in format "hash │ committer │ message"

  Scenario: Display fewer commits when history is limited
    Given there are only 3 commits in history
    When I run "git @ hash"
    Then I should see "RECENT COMMITS (last 5):"
    And I should see 3 commit lines
    And each line should show hash, committer, and message

  Scenario: Handle no commits in repository
    Given there are no commits in the repository
    When I run "git @ hash"
    Then I should see "RECENT COMMITS (last 5):"
    And I should see "No commits found"

  Scenario: Handle commits with unknown committers
    Given there are commits with missing committer information
    When I run "git @ hash"
    Then I should see "Unknown" for committer names
    And the command should complete successfully

  Scenario: Handle long committer names
    Given there are commits with very long committer names
    When I run "git @ hash"
    Then I should see truncated committer names (max 20 characters)
    And long names should end with "..."

  Scenario: Handle long commit messages
    Given there are commits with very long messages
    When I run "git @ hash"
    Then I should see truncated messages (max 50 characters)
    And long messages should end with "..."

  Scenario: Handle branch without remote
    Given I am on a local-only branch
    When I run "git @ hash"
    Then I should see "ORIGIN HASH: UNKNOWN"
    And I should see "CURRENT is 0 behind and 0 ahead of origin/branch"
    And I should see recent commits section

  Scenario: Show help
    When I run "git @ hash --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ hash"
    And the output should contain "Last 5 commits with hashes, committers, and messages"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ hash"
    Then the command should fail
    And I should see "Error: Not in a git repository"

  Scenario: Handle detached HEAD state
    Given I am in detached HEAD state
    When I run "git @ hash"
    Then the command should fail
    And I should see "Error: Could not determine current branch"

  Scenario: Handle missing remote branches
    Given there are no remote branches
    When I run "git @ hash"
    Then I should see "UNKNOWN" for remote hashes
    And the command should complete successfully
    And I should see recent commits section

  Scenario: Show detailed branch information with commits
    When I run "git @ hash"
    Then I should see branch name
    And I should see current hash
    And I should see origin hash
    And I should see ahead/behind counts
    And I should see recent commits with hashes, committers, and messages

  Scenario: Handle git fetch failures gracefully
    Given git fetch fails
    When I run "git @ hash"
    Then the command should continue
    And I should see "UNKNOWN" for remote information
    And I should see recent commits section

  Scenario: Format commit display correctly
    Given there are commits with long messages
    When I run "git @ hash"
    Then I should see commits formatted as "hash │ committer │ message"
    And the hash should be left-aligned in 8-character field
    And the committer should be left-aligned in 20-character field
    And the message should follow the second separator 