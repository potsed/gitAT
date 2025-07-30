Feature: Project ID Generation
  As a developer
  I want to generate a unique project identifier
  So that I can identify the current product state

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show project ID
    When I run "git @ _id"
    Then I should see the project ID in format "product:major.minor.fix"
    And the output should be generated from git config values

  Scenario: Show project ID with all components
    Given the product is set to "GitAT"
    And the major version is "1"
    And the minor version is "2"
    And the fix version is "3"
    When I run "git @ _id"
    Then I should see "GitAT:123"

  Scenario: Show project ID with missing components
    Given the product is set to "GitAT"
    And no version components are set
    When I run "git @ _id"
    Then I should see "GitAT:"

  Scenario: Show project ID with no configuration
    Given no GitAT configuration is set
    When I run "git @ _id"
    Then I should see ":"

  Scenario: Show project ID with partial version
    Given the product is set to "GitAT"
    And the major version is "1"
    And the minor version is "0"
    And no fix version is set
    When I run "git @ _id"
    Then I should see "GitAT:10"

  Scenario: Show help
    When I run "git @ _id --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ _id"

  Scenario: Handle empty project name
    Given no product name is set
    And the version is "1.2.3"
    When I run "git @ _id"
    Then I should see ":123"

  Scenario: Handle zero versions
    Given the product is set to "GitAT"
    And all version components are "0"
    When I run "git @ _id"
    Then I should see "GitAT:000" 