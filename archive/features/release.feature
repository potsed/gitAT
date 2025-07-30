Feature: Release Management
  As a developer
  I want to create releases with version bumping
  So that I can properly version and tag releases

  Background:
    Given I am in a git repository
    And GitAT is initialized
    And I am on the master branch

  Scenario: Create minor release
    When I run "git @ release -m"
    Then the minor version should be incremented
    And a release commit should be created
    And a version tag should be created
    And I should see "Release created successfully"

  Scenario: Create major release
    When I run "git @ release -M"
    Then the major version should be incremented
    And the minor version should be reset to 0
    And a release commit should be created
    And a version tag should be created

  Scenario: Create fix release
    When I run "git @ release -b"
    Then the fix version should be incremented
    And a release commit should be created
    And a version tag should be created

  Scenario: Create release with custom message
    When I run "git @ release -m \"Add new features\""
    Then the minor version should be incremented
    And the commit message should contain "Add new features"
    And a version tag should be created

  Scenario: Prevent release on non-master branch
    Given I am on branch "develop"
    When I run "git @ release -m"
    Then the command should fail
    And I should see "Error: Releases can only be created from master branch"

  Scenario: Prevent release with uncommitted changes
    Given I have uncommitted changes
    When I run "git @ release -m"
    Then the command should fail
    And I should see "Error: Cannot create release with uncommitted changes"

  Scenario: Show help
    When I run "git @ release --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ release"

  Scenario: Handle no version configured
    Given no version is configured
    When I run "git @ release -m"
    Then the version should be initialized to 0.1.0
    And a release should be created

  Scenario: Create release with all components
    Given the current version is "1.2.3"
    When I run "git @ release -M -m -b"
    Then the version should be "2.1.1"
    And a release should be created 