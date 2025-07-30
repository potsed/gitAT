Feature: Work In Progress Management
  As a developer
  I want to manage my work-in-progress branch
  So that I can context switch between different features

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show current WIP branch
    When I run "git @ wip"
    Then I should see the current WIP branch
    And the output should be from git config at.wip

  Scenario: Set current branch as WIP
    Given I am on branch "feature-auth"
    When I run "git @ wip -s"
    Then the WIP branch should be set to "feature-auth"
    And git config at.wip should be "feature-auth"
    And I should see "WIP updated to feature-auth from"

  Scenario: Set current branch as WIP with dot
    Given I am on branch "feature-auth"
    When I run "git @ wip ."
    Then the WIP branch should be set to "feature-auth"
    And git config at.wip should be "feature-auth"

  Scenario: Checkout WIP branch
    Given the WIP branch is set to "feature-auth"
    When I run "git @ wip -c"
    Then I should be switched to branch "feature-auth"
    And I should see "Switched to WIP branch: feature-auth"

  Scenario: Restore WIP to working branch
    Given the WIP branch is set to "feature-auth"
    When I run "git @ wip -r"
    Then the working branch should be set to "feature-auth"
    And git @ work should be executed
    And I should see "Restored WIP branch: feature-auth"

  Scenario: Handle no WIP branch configured for checkout
    Given no WIP branch is configured
    When I run "git @ wip -c"
    Then the command should fail
    And I should see "Error: No WIP branch configured"

  Scenario: Handle no WIP branch configured for restore
    Given no WIP branch is configured
    When I run "git @ wip -r"
    Then the command should fail
    And I should see "Error: No WIP branch configured"

  Scenario: Show help
    When I run "git @ wip --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ wip"

  Scenario: Handle empty WIP branch
    Given no WIP branch is set
    When I run "git @ wip"
    Then I should see "Current WIP is: "
    And the exit code should be 0 