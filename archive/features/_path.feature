Feature: Repository Path
  As a developer
  I want to get the repository root path
  So that I can navigate to the correct directory

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show repository path
    When I run "git @ _path"
    Then I should see the repository root path
    And the output should be from "git rev-parse --show-toplevel"

  Scenario: Show absolute path
    Given I am in a subdirectory of the repository
    When I run "git @ _path"
    Then I should see the absolute path to the repository root
    And the path should be absolute

  Scenario: Show help
    When I run "git @ _path --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ _path"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ _path"
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Handle deep subdirectory
    Given I am in a deep subdirectory of the repository
    When I run "git @ _path"
    Then I should see the repository root path
    And the path should not include the subdirectory 