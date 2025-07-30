package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

// createTestManager creates a test manager with a temporary repository
func createTestManager(t *testing.T) *Manager {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gitat-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repository
	repo := git.NewRepository(tempDir)
	if _, err := repo.Run("init"); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit the file
	if _, err := repo.Run("add", "test.txt"); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}
	if _, err := repo.Run("commit", "-m", "Initial commit"); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Create config
	cfg := &config.Config{
		RepoPath: tempDir,
	}

	return &Manager{
		config: cfg,
		git:    repo,
	}
}

// TestPath tests the _path command
func TestPath(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test _path command
	err := manager.Path([]string{})
	if err != nil {
		t.Errorf("Path command failed: %v", err)
	}
}

// TestChanges tests the changes command
func TestChanges(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Create a modified file
	testFile := filepath.Join(manager.config.RepoPath, "modified.txt")
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to create modified file: %v", err)
	}

	// Test changes command
	err := manager.Changes([]string{})
	if err != nil {
		t.Errorf("Changes command failed: %v", err)
	}
}

// TestLogs tests the logs command
func TestLogs(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test logs command
	err := manager.Logs([]string{})
	if err != nil {
		t.Errorf("Logs command failed: %v", err)
	}
}

// TestProduct tests the product command
func TestProduct(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test setting product
	err := manager.Product([]string{"test-product"})
	if err != nil {
		t.Errorf("Product set command failed: %v", err)
	}

	// Test getting product
	err = manager.Product([]string{})
	if err != nil {
		t.Errorf("Product get command failed: %v", err)
	}
}

// TestFeature tests the feature command
func TestFeature(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test setting feature
	err := manager.Feature([]string{"test-feature"})
	if err != nil {
		t.Errorf("Feature set command failed: %v", err)
	}

	// Test getting feature
	err = manager.Feature([]string{})
	if err != nil {
		t.Errorf("Feature get command failed: %v", err)
	}
}

// TestIssue tests the issue command
func TestIssue(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test setting issue
	err := manager.Issue([]string{"TEST-123"})
	if err != nil {
		t.Errorf("Issue set command failed: %v", err)
	}

	// Test getting issue
	err = manager.Issue([]string{})
	if err != nil {
		t.Errorf("Issue get command failed: %v", err)
	}
}

// TestVersion tests the version command
func TestVersion(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test getting version
	err := manager.Version([]string{})
	if err != nil {
		t.Errorf("Version get command failed: %v", err)
	}

	// Test version tag
	err = manager.Version([]string{"-t"})
	if err != nil {
		t.Errorf("Version tag command failed: %v", err)
	}

	// Test version reset
	err = manager.Version([]string{"-r"})
	if err != nil {
		t.Errorf("Version reset command failed: %v", err)
	}
}

// TestTrunk tests the _trunk command
func TestTrunk(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test setting trunk
	err := manager.Trunk([]string{"master"})
	if err != nil {
		t.Errorf("Trunk set command failed: %v", err)
	}

	// Test getting trunk
	err = manager.Trunk([]string{})
	if err != nil {
		t.Errorf("Trunk get command failed: %v", err)
	}
}

// TestLabel tests the _label command
func TestLabel(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Set up product, feature, and issue
	manager.Product([]string{"test-product"})
	manager.Feature([]string{"test-feature"})
	manager.Issue([]string{"TEST-123"})

	// Test getting label
	err := manager.Label([]string{})
	if err != nil {
		t.Errorf("Label get command failed: %v", err)
	}

	// Test setting custom label
	err = manager.Label([]string{"custom-label"})
	if err != nil {
		t.Errorf("Label set command failed: %v", err)
	}
}

// TestID tests the _id command
func TestID(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Set up product and version
	manager.Product([]string{"test-product"})
	manager.Version([]string{"-r"}) // Reset to 0.0.0

	// Test getting ID
	err := manager.ID([]string{})
	if err != nil {
		t.Errorf("ID command failed: %v", err)
	}
}

// TestWIP tests the wip command
func TestWIP(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test getting WIP (should be empty initially)
	err := manager.WIP([]string{})
	if err != nil {
		t.Errorf("WIP get command failed: %v", err)
	}

	// Test setting WIP
	err = manager.WIP([]string{"-s"})
	if err != nil {
		t.Errorf("WIP set command failed: %v", err)
	}
}

// TestHelpCommands tests help functionality for all commands
func TestHelpCommands(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	helpTests := []struct {
		name    string
		command func([]string) error
	}{
		{"Path", manager.Path},
		{"Changes", manager.Changes},
		{"Logs", manager.Logs},
		{"Product", manager.Product},
		{"Feature", manager.Feature},
		{"Issue", manager.Issue},
		{"Version", manager.Version},
		{"Trunk", manager.Trunk},
		{"Label", manager.Label},
		{"ID", manager.ID},
		{"WIP", manager.WIP},
	}

	for _, tt := range helpTests {
		t.Run(tt.name+"Help", func(t *testing.T) {
			err := tt.command([]string{"--help"})
			if err != nil {
				t.Errorf("%s help command failed: %v", tt.name, err)
			}
		})
	}
}

// TestErrorHandling tests error conditions
func TestErrorHandling(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test invalid arguments - these should show help instead of error
	err := manager.Product([]string{"arg1", "arg2"})
	if err != nil {
		t.Logf("Product with invalid args returned: %v", err)
	}

	err = manager.Feature([]string{"arg1", "arg2"})
	if err != nil {
		t.Logf("Feature with invalid args returned: %v", err)
	}

	err = manager.Issue([]string{"arg1", "arg2"})
	if err != nil {
		t.Logf("Issue with invalid args returned: %v", err)
	}
}

// cleanupTest cleans up the test environment
func cleanupTest(t *testing.T, manager *Manager) {
	if manager.config.RepoPath != "" {
		if err := os.RemoveAll(manager.config.RepoPath); err != nil {
			t.Logf("Failed to cleanup test directory: %v", err)
		}
	}
}

// TestManagerIntegration tests integration between commands
func TestManagerIntegration(t *testing.T) {
	manager := createTestManager(t)
	defer cleanupTest(t, manager)

	// Test full workflow
	// 1. Set product
	err := manager.Product([]string{"integration-test"})
	if err != nil {
		t.Fatalf("Failed to set product: %v", err)
	}

	// 2. Set feature
	err = manager.Feature([]string{"test-feature"})
	if err != nil {
		t.Fatalf("Failed to set feature: %v", err)
	}

	// 3. Set issue
	err = manager.Issue([]string{"INT-001"})
	if err != nil {
		t.Fatalf("Failed to set issue: %v", err)
	}

	// 4. Set trunk
	err = manager.Trunk([]string{"main"})
	if err != nil {
		t.Fatalf("Failed to set trunk: %v", err)
	}

	// 5. Test label generation
	err = manager.Label([]string{})
	if err != nil {
		t.Fatalf("Failed to generate label: %v", err)
	}

	// 6. Test ID generation
	err = manager.ID([]string{})
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}
}

// BenchmarkCommands benchmarks the command performance
func BenchmarkPath(b *testing.B) {
	manager := createTestManager(&testing.T{})
	defer cleanupTest(&testing.T{}, manager)

	for i := 0; i < b.N; i++ {
		manager.Path([]string{})
	}
}

func BenchmarkChanges(b *testing.B) {
	manager := createTestManager(&testing.T{})
	defer cleanupTest(&testing.T{}, manager)

	for i := 0; i < b.N; i++ {
		manager.Changes([]string{})
	}
}

func BenchmarkLogs(b *testing.B) {
	manager := createTestManager(&testing.T{})
	defer cleanupTest(&testing.T{}, manager)

	for i := 0; i < b.N; i++ {
		manager.Logs([]string{})
	}
}
