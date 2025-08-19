package main

import (
	"fmt"
	"os"

	"github.com/potsed/gitAT/internal/commands/handlers"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
)

func main() {
	// Create a basic config
	cfg := &config.Config{}
	
	// Get current directory as repo path
	currentDir, _ := os.Getwd()
	cfg.RepoPath = currentDir
	
	// Create git repo instance
	gitRepo := git.NewRepository(currentDir)
	
	// Test creating the new handlers
	fmt.Println("Testing new handlers...")
	
	changelogHandler := handlers.NewChangelogHandler(cfg, gitRepo)
	if changelogHandler == nil {
		fmt.Println("ERROR: Failed to create ChangelogHandler")
		return
	}
	fmt.Println("✓ ChangelogHandler created successfully")
	
	rebaseHandler := handlers.NewRebaseHandler(cfg, gitRepo)
	if rebaseHandler == nil {
		fmt.Println("ERROR: Failed to create RebaseHandler")
		return
	}
	fmt.Println("✓ RebaseHandler created successfully")
	
	commitizenHandler := handlers.NewCommitizenHandler(cfg, gitRepo)
	if commitizenHandler == nil {
		fmt.Println("ERROR: Failed to create CommitizenHandler")
		return
	}
	fmt.Println("✓ CommitizenHandler created successfully")
	
	// Test showing usage for each handler
	fmt.Println("\nTesting handler usage...")
	
	fmt.Println("\n--- Changelog Handler Usage ---")
	changelogHandler.Execute([]string{"--help"})
	
	fmt.Println("\n--- Rebase Handler Usage ---")
	rebaseHandler.Execute([]string{"--help"})
	
	fmt.Println("\n--- Commitizen Handler Usage ---")
	commitizenHandler.Execute([]string{"--help"})
	
	fmt.Println("\nAll handlers created and tested successfully!")
}