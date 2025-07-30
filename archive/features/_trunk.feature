Feature: Trunk Branch Management
  As a developer
  I want to manage the trunk/base branch
  So that I can track which branch is the main development branch

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show current trunk branch
    When I run "git @ _trunk"
    Then I should see the current trunk branch
    And the output should be from git config at.trunk

  Scenario: Set trunk branch
    When I run "git @ _trunk develop"
    Then the trunk branch should be set to "develop"
    And git config at.trunk should be "develop"
    And I should see "Base branch updated to: develop from"

  Scenario: Set trunk branch to master
    When I run "git @ _trunk master"
    Then the trunk branch should be set to "master"
    And git config at.trunk should be "master"

  Scenario: Auto-detect trunk branch
    Given no trunk branch is set
    And the remote HEAD points to "main"
    When I run "git @ _trunk"
    Then the trunk branch should be auto-detected as "main"
    And I should see "Auto-detected trunk branch: main"
    And git config at.trunk should be "main"

  Scenario: Auto-detect trunk branch with master
    Given no trunk branch is set
    And the remote HEAD points to "master"
    When I run "git @ _trunk"
    Then the trunk branch should be auto-detected as "master"
    And I should see "Auto-detected trunk branch: master"

  Scenario: Show help
    When I run "git @ _trunk --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ _trunk"

  Scenario: Handle empty trunk branch
    Given no trunk branch is set
    And remote HEAD detection fails
    When I run "git @ _trunk"
    Then the trunk branch should default to "develop"
    And I should see "Auto-detected trunk branch: develop"

  Scenario: Update existing trunk branch
    Given the trunk branch is set to "develop"
    When I run "git @ _trunk main"
    Then the trunk branch should be updated to "main"
    And I should see "Base branch updated to: main from develop" 