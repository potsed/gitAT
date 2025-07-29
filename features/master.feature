Feature: Switch to Master Branch
  As a developer
  I want to switch to master branch with stash management
  So that I can work on the main branch safely

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Switch to master branch
    Given I am on branch "feature-auth"
    And I have uncommitted changes
    When I run "git @ master"
    Then I should be switched to the master branch
    And the changes should be stashed
    And the latest changes should be pulled

  Scenario: Already on master branch
    Given I am on branch "master"
    When I run "git @ master"
    Then I should remain on branch "master"
    And no stash should be created

  Scenario: Switch without changes
    Given I am on branch "feature-auth"
    And I have no uncommitted changes
    When I run "git @ master"
    Then I should be switched to the master branch
    And no stash should be created
    And the latest changes should be pulled

  Scenario: Show help
    When I run "git @ master --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ master"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ master"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Handle master branch not existing
    Given the master branch doesn't exist
    When I run "git @ master"
    Then the command should fail
    And I should see an appropriate error message 