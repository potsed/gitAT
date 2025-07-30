Feature: Add to .gitignore
  As a developer
  I want to add patterns to .gitignore
  So that I can exclude files from version control

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Add pattern to .gitignore
    When I run "git @ ignore \"*.log\""
    Then the pattern should be added to .gitignore
    And the .gitignore file should contain "*.log"

  Scenario: Add multiple patterns
    When I run "git @ ignore \"node_modules/\""
    And I run "git @ ignore \"*.tmp\""
    Then the .gitignore file should contain "node_modules/"
    And the .gitignore file should contain "*.tmp"

  Scenario: Add pattern with spaces
    When I run "git @ ignore \"build files/\""
    Then the .gitignore file should contain "build files/"

  Scenario: Show help
    When I run "git @ ignore --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ ignore"

  Scenario: Handle missing pattern argument
    When I run "git @ ignore"
    Then the command should show usage information
    And the exit code should be 0

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ ignore \"*.log\""
    Then the command should fail
    And I should see an appropriate error message

  Scenario: Create .gitignore if it doesn't exist
    Given .gitignore doesn't exist
    When I run "git @ ignore \"*.log\""
    Then .gitignore should be created
    And it should contain "*.log"

  Scenario: Append to existing .gitignore
    Given .gitignore exists with content
    When I run "git @ ignore \"*.log\""
    Then the pattern should be appended to .gitignore
    And existing content should be preserved 