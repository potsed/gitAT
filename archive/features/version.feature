Feature: Version Management
  As a developer
  I want to manage semantic versioning
  So that I can track version changes properly

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show current version
    When I run "git @ version"
    Then I should see the current version in format "major.minor.fix"
    And the output should be from git config at.major.at.minor.at.fix

  Scenario: Show version tag
    When I run "git @ version -t"
    Then I should see the version with "v" prefix
    And the output should be "v" + current version

  Scenario: Reset version to 0.0.0
    When I run "git @ version -r"
    Then the version should be reset to "0.0.0"
    And git config at.major should be "0"
    And git config at.minor should be "0"
    And git config at.fix should be "0"

  Scenario: Increment major version
    Given the current version is "1.2.3"
    When I run "git @ version -M"
    Then the major version should be incremented to "2"
    And the minor version should be reset to "0"
    And the fix version should be reset to "0"
    And the new version should be "2.0.0"

  Scenario: Increment minor version
    Given the current version is "1.2.3"
    When I run "git @ version -m"
    Then the minor version should be incremented to "3"
    And the fix version should be reset to "0"
    And the new version should be "1.3.0"

  Scenario: Increment fix version
    Given the current version is "1.2.3"
    When I run "git @ version -b"
    Then the fix version should be incremented to "4"
    And the new version should be "1.2.4"

  Scenario: Show help
    When I run "git @ version --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ version"

  Scenario: Handle unset version
    Given no version is set
    When I run "git @ version"
    Then I should see "0.0.0"
    And the exit code should be 0 