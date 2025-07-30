Feature: GitAT Information Display
  As a developer
  I want to see comprehensive GitAT status
  So that I can understand my current configuration and repository state

  Background:
    Given I am in a git repository
    And GitAT is initialized

  Scenario: Show comprehensive information
    When I run "git @ info"
    Then I should see a formatted status report
    And the output should contain "GitAT Status Report"
    And the output should contain "Configuration & Information"
    And the output should contain "Git Repository Status"
    And the output should contain "Branch Information"
    And the output should contain "Available Commands"
    And the output should contain "Quick Actions"

  Scenario: Display configuration information
    Given the product is set to "GitAT"
    And the feature is set to "user-auth"
    And the issue is set to "PROJ-123"
    When I run "git @ info"
    Then I should see "Product Name │ GitAT"
    And I should see "Feature Name │ user-auth"
    And I should see "Issue/Task ID │ PROJ-123"

  Scenario: Display git repository status
    Given I am on branch "feature-auth"
    And I have uncommitted changes
    When I run "git @ info"
    Then I should see "Current Branch │ feature-auth"
    And I should see "Uncommitted Changes │" with a number
    And I should see "Repository Root │" with a path

  Scenario: Display branch information
    Given I am on branch "feature-auth"
    And the working branch is set to "feature-auth"
    When I run "git @ info"
    Then I should see "Branch Status │ ✅ On correct branch"

  Scenario: Display branch warning
    Given I am on branch "develop"
    And the working branch is set to "feature-auth"
    When I run "git @ info"
    Then I should see "Branch Status │ ⚠️ On develop (should be on feature-auth)"

  Scenario: Display branch protection
    Given I am on branch "master"
    When I run "git @ info"
    Then I should see "Branch Protection │ 🛡️ Protected branch"

  Scenario: Display available commands
    When I run "git @ info"
    Then I should see "Workflow │ git @ save, git @ work, git @ wip"
    And I should see "Version Management │ git @ version, git @ release"
    And I should see "Branch Management │ git @ branch, git @ master, git @ root, git @ sweep"

  Scenario: Display quick actions
    When I run "git @ info"
    Then I should see "Switch to Work │ git @ work"
    And I should see "Save Changes │ git @ save \"message\""
    And I should see "Check Changes │ git @ changes"

  Scenario: Handle not in git repository
    Given I am not in a git repository
    When I run "git @ info"
    Then I should see "Current Branch │ <not in git repo>"
    And I should see "Repository Root │ <not in git repo>"

  Scenario: Show help
    When I run "git @ info --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ info"

  Scenario: Handle empty configurations
    Given no GitAT configurations are set
    When I run "git @ info"
    Then I should see "<not set>" for empty configurations
    And the command should complete successfully 