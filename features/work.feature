Feature: Switch to Working Branch
  As a developer
  I want to switch to my working branch with stash management
  So that I can start working on my feature

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Switch to working branch
    Given I am on branch "develop"
    And the working branch is set to "feature-auth"
    When I run "git @ work"
    Then I should be switched to branch "feature-auth"
    And any changes should be stashed
    And the stash should be reapplied

  Scenario: Already on working branch
    Given I am on branch "feature-auth"
    And the working branch is set to "feature-auth"
    When I run "git @ work"
    Then I should see "You're already in the working branch"
    And I should remain on branch "feature-auth"

  Scenario: Create working branch if it doesn't exist
    Given I am on branch "develop"
    And the working branch is set to "feature-auth"
    And branch "feature-auth" doesn't exist
    When I run "git @ work"
    Then branch "feature-auth" should be created
    And I should be switched to branch "feature-auth"
    And the branch should be based on the latest develop

  Scenario: Handle stash management
    Given I have uncommitted changes
    And I am on branch "develop"
    And the working branch is set to "feature-auth"
    When I run "git @ work"
    Then the changes should be stashed with key "autostash-work-branch"
    And the stash should be reapplied on the working branch

  Scenario: Fetch latest changes
    Given I am on branch "develop"
    And the working branch is set to "feature-auth"
    When I run "git @ work"
    Then git fetch should be executed
    And the base branch should be updated

  Scenario: Show help
    When I run "git @ work --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ work"

  Scenario: Handle no working branch configured
    Given no working branch is configured
    When I run "git @ work"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Handle remote branch updates
    Given I am on branch "develop"
    And the working branch is set to "feature-auth"
    And there are remote updates
    When I run "git @ work"
    Then the remote changes should be pulled
    And the working branch should be updated 