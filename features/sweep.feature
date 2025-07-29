Feature: Clean Up Merged Branches
  As a developer
  I want to clean up merged branches
  So that I can keep my repository organized

  Background:
    Given I am in a git repository
    And GitAT is initialized
    And there are merged branches

  Scenario: Clean up merged branches
    When I run "git @ sweep"
    Then merged branches should be deleted
    And important branches should be preserved
    And I should see "🧹 Cleaning up merged branches..."
    And I should see "✅ Sweep completed"

  Scenario: Preserve important branches
    Given there are branches: master, develop, staging, feature-auth
    And feature-auth is merged
    When I run "git @ sweep"
    Then master should be preserved
    And develop should be preserved
    And staging should be preserved
    And feature-auth should be deleted

  Scenario: Preserve trunk branch
    Given the trunk branch is set to "main"
    And there are merged branches
    When I run "git @ sweep"
    Then the trunk branch "main" should be preserved
    And I should see "Preserving branches:" with the trunk branch

  Scenario: Handle no merged branches
    Given there are no merged branches to clean up
    When I run "git @ sweep"
    Then I should see "No merged branches to clean up"
    And the command should complete successfully

  Scenario: Handle failed branch deletions
    Given there are branches with unmerged changes
    When I run "git @ sweep"
    Then the command should show which branches failed to delete
    And I should see "(failed to delete - may have unmerged changes)"

  Scenario: Show help
    When I run "git @ sweep --help"
    Then I should see the usage information
    And the output should contain "Usage: git @ sweep"

  Scenario: Preserve all protected branches
    Given there are branches: master, main, dev, develop, staging, stage, qa
    When I run "git @ sweep"
    Then all protected branches should be preserved
    And I should see "Preserving branches: master|main|dev|develop|staging|stage|qa"

  Scenario: Show deletion progress
    Given there are merged branches to delete
    When I run "git @ sweep"
    Then I should see "Deleting merged branches:"
    And each branch should be listed before deletion
    And I should see deletion counts at the end

  Scenario: Handle dynamic trunk branch
    Given the trunk branch is set to "main"
    When I run "git @ sweep"
    Then I should see "Trunk branch: main"
    And "main" should be in the preserved branches list 