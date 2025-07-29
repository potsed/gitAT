Feature: Squash Commits
  As a developer
  I want to squash commits to clean up history
  So that I can create clean commit history before PR

  Background:
    Given I am in a git repository
    And GitAT is initialized
    And I have multiple commits on my branch

  Scenario: Squash to develop branch
    When I run "git @ squash develop"
    Then the commits should be soft reset to develop HEAD
    And all changes should remain staged
    And I should see "Squashed branch" with the current branch name

  Scenario: Squash to master branch
    When I run "git @ squash master"
    Then the commits should be soft reset to master HEAD
    And all changes should remain staged

  Scenario: Squash and save
    When I run "git @ squash develop -s"
    Then the commits should be soft reset to develop HEAD
    And git @ save should be executed
    And the changes should be committed

  Scenario: Handle non-existent target branch
    When I run "git @ squash non-existent-branch"
    Then the command should fail
    And I should see "ERROR: Branch \"non-existent-branch\" does not exist locally"

  Scenario: Handle missing target branch argument
    When I run "git @ squash"
    Then the command should show usage information
    And the exit code should be 0

  Scenario: Show help
    When I run "git @ squash --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ squash"

  Scenario: Squash with custom target branch
    Given there is a branch called "staging"
    When I run "git @ squash staging"
    Then the commits should be soft reset to staging HEAD
    And all changes should remain staged

  Scenario: Handle squash on protected branch
    Given I am on branch "master"
    When I run "git @ squash develop"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Verify staged changes after squash
    When I run "git @ squash develop"
    Then git status should show staged changes
    And the commit history should be clean 