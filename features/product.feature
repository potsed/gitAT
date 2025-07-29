Feature: Product Configuration
  As a developer
  I want to set and get the product name
  So that I can track which product I'm working on

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show current product name
    When I run "git @ product"
    Then I should see the current product name
    And the output should be from git config at.product

  Scenario: Set product name
    When I run "git @ product GitAT"
    Then the product name should be set to "GitAT"
    And git config at.product should be "GitAT"
    And I should see "Project updated to: GitAT from"

  Scenario: Set product name with special characters
    When I run "git @ product my-app-v2"
    Then the product name should be set to "my-app-v2"
    And git config at.product should be "my-app-v2"

  Scenario: Show help
    When I run "git @ product --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ product"

  Scenario: Handle empty product name
    Given no product name is set
    When I run "git @ product"
    Then I should see an empty line
    And the exit code should be 0 