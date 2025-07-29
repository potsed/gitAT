Feature: Initialize GitAT
  As a developer
  I want to initialize GitAT for a new repository
  So that I can set up all configurations for general use

  Background:
    Given I am in a git repository
    And GitAT is not initialized

  Scenario: Initialize GitAT
    When I run "git @ _go"
    Then the base branch should be set based on remote HEAD
    And the version should be reset to 0.0.0
    And the current working branch should be set
    And the current WIP branch should be set
    And the repository should be marked as initialized

  Scenario: Set base branch from remote HEAD
    Given the remote HEAD points to "main"
    When I run "git @ _go"
    Then the trunk branch should be set to "main"
    And git config at.trunk should be "main"

  Scenario: Set base branch from remote HEAD with master
    Given the remote HEAD points to "master"
    When I run "git @ _go"
    Then the trunk branch should be set to "master"
    And git config at.trunk should be "master"

  Scenario: Reset version to 0.0.0
    When I run "git @ _go"
    Then the version should be reset to 0.0.0
    And git config at.major should be "0"
    And git config at.minor should be "0"
    And git config at.fix should be "0"

  Scenario: Set current working branch
    Given I am on branch "develop"
    When I run "git @ _go"
    Then the working branch should be set to "develop"
    And git config at.branch should be "develop"

  Scenario: Set current WIP branch
    Given I am on branch "develop"
    When I run "git @ _go"
    Then the WIP branch should be set to "develop"
    And git config at.wip should be "develop"

  Scenario: Mark repository as initialized
    When I run "git @ _go"
    Then the repository should be marked as initialized
    And git config at.initialised should be "true"

  Scenario: Show help
    When I run "git @ _go --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ _go"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ _go"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Handle already initialized repository
    Given GitAT is already initialized
    When I run "git @ _go"
    Then the configurations should be updated
    And the repository should remain marked as initialized 