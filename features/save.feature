Feature: Save Changes
  As a developer
  I want to save my changes with proper commit messages
  So that I can track my work with meaningful commits

  Background:
    Given I am in a git repository
    And GitAT is initialized
    And I have uncommitted changes

  Scenario: Save with custom message
    When I run "git @ save \"Add user authentication\""
    Then the changes should be committed
    And the commit message should contain the label
    And the commit message should contain "Add user authentication"

  Scenario: Save without message
    When I run "git @ save"
    Then the changes should be committed
    And the commit message should be the default label or "Update"

  Scenario: Save with product and feature configured
    Given the product is set to "GitAT"
    And the feature is set to "user-auth"
    And the issue is set to "PROJ-123"
    When I run "git @ save \"Add login functionality\""
    Then the commit message should be "[GitAT.user-auth.PROJ-123] Add login functionality"

  Scenario: Prevent save on master branch
    Given I am on branch "master"
    When I run "git @ save \"Test commit\""
    Then the command should fail
    And I should see "Error: Cannot save changes on master"

  Scenario: Prevent save on develop branch
    Given I am on branch "develop"
    When I run "git @ save \"Test commit\""
    Then the command should fail
    And I should see "Error: Cannot save changes on develop"

  Scenario: Auto-set working branch if not configured
    Given no working branch is configured
    And I am on branch "feature-auth"
    When I run "git @ save \"Test commit\""
    Then the working branch should be set to "feature-auth"
    And the changes should be committed

  Scenario: Confirm save on production branch
    Given I am on branch "prod"
    When I run "git @ save \"Production fix\""
    Then I should be prompted for confirmation
    And if I confirm, the changes should be committed
    And a version tag should be created

  Scenario: Show help
    When I run "git @ save --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ save"

  Scenario: Handle invalid message characters
    When I run "git @ save \"Invalid; message\""
    Then the command should fail
    And I should see "Error: Invalid message"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ save \"Test commit\""
    Then the command should fail
    And I should see "Error: Not in a git repository" 