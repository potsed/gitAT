Feature: Commit Label Management
  As a developer
  I want to manage commit labels
  So that I can create consistent commit messages

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show current label
    When I run "git @ _label"
    Then I should see the current label in format [product.feature.issue]
    And the output should be generated from git config values

  Scenario: Show label with all components
    Given the product is set to "GitAT"
    And the feature is set to "user-auth"
    And the issue is set to "PROJ-123"
    When I run "git @ _label"
    Then I should see "[GitAT.user-auth.PROJ-123]"

  Scenario: Show label with missing components
    Given the product is set to "GitAT"
    And no feature is set
    And no issue is set
    When I run "git @ _label"
    Then I should see "[GitAT..]"

  Scenario: Show label with no configuration
    Given no GitAT configuration is set
    When I run "git @ _label"
    Then I should see "[Update]"

  Scenario: Set custom label
    When I run "git @ _label \"Custom Label\""
    Then the custom label should be set
    And git config at.label should be "Custom Label"
    And I should see "Label updated to: Custom Label"

  Scenario: Set custom label with special characters
    When I run "git @ _label \"[HOTFIX] Security patch\""
    Then the custom label should be set
    And git config at.label should be "[HOTFIX] Security patch"

  Scenario: Show help
    When I run "git @ _label --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ _label"

  Scenario: Handle empty label
    Given no label is set
    When I run "git @ _label"
    Then I should see the default format or "[Update]"
    And the exit code should be 0

  Scenario: Update existing label
    Given a custom label is set
    When I run "git @ _label \"New Label\""
    Then the label should be updated to "New Label"
    And I should see "Label updated to: New Label" 