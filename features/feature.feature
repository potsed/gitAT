Feature: Feature Configuration
  As a developer
  I want to set and get the feature name
  So that I can track which feature I'm working on

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show current feature name
    When I run "git @ feature"
    Then I should see the current feature name
    And the output should be from git config at.feature

  Scenario: Set feature name
    When I run "git @ feature user-auth"
    Then the feature name should be set to "user-auth"
    And git config at.feature should be "user-auth"
    And I should see "Feature updated to: user-auth from"

  Scenario: Set feature name with hyphens
    When I run "git @ feature payment-integration"
    Then the feature name should be set to "payment-integration"
    And git config at.feature should be "payment-integration"

  Scenario: Set feature name with underscores
    When I run "git @ feature user_management"
    Then the feature name should be set to "user_management"
    And git config at.feature should be "user_management"

  Scenario: Show help
    When I run "git @ feature --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ feature"

  Scenario: Handle empty feature name
    Given no feature name is set
    When I run "git @ feature"
    Then I should see an empty line
    And the exit code should be 0 