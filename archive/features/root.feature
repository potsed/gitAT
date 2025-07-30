Feature: Switch to Root/Trunk Branch
  As a developer
  I want to switch to the root/trunk branch with stash management
  So that I can work on the main development branch safely

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Switch to trunk branch
    Given I am on branch "feature-auth"
    And I have uncommitted changes
    And the trunk branch is set to "develop"
    When I run "git @ root"
    Then I should be switched to the develop branch
    And the changes should be stashed
    And the latest changes should be pulled with rebase

  Scenario: Already on trunk branch
    Given I am on branch "develop"
    And the trunk branch is set to "develop"
    When I run "git @ root"
    Then I should remain on branch "develop"
    And no stash should be created

  Scenario: Switch without changes
    Given I am on branch "feature-auth"
    And I have no uncommitted changes
    And the trunk branch is set to "develop"
    When I run "git @ root"
    Then I should be switched to the develop branch
    And no stash should be created
    And the latest changes should be pulled with rebase

  Scenario: Show help
    When I run "git @ root --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ root"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ root"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Handle trunk branch not existing
    Given the trunk branch doesn't exist
    When I run "git @ root"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Use dynamic trunk branch
    Given the trunk branch is set to "main"
    When I run "git @ root"
    Then I should be switched to the main branch 