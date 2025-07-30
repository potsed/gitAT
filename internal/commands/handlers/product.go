package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/potsed/gitAT/internal/config"
	"github.com/potsed/gitAT/internal/git"
	"github.com/potsed/gitAT/pkg/output"
)

// ProductHandler handles product-related commands
type ProductHandler struct {
	BaseHandler
}

// NewProductHandler creates a new product handler
func NewProductHandler(cfg *config.Config, gitRepo *git.Repository) *ProductHandler {
	return &ProductHandler{
		BaseHandler: NewBaseHandler(cfg, gitRepo),
	}
}

// Execute handles the product command
func (p *ProductHandler) Execute(args []string) error {
	if len(args) == 0 {
		return p.showProduct()
	}

	switch args[0] {
	case "set":
		if len(args) < 2 {
			return fmt.Errorf("product name required")
		}
		return p.setProduct(strings.Join(args[1:], " "))
	case "get":
		return p.showProduct()
	case "clear", "unset":
		return p.clearProduct()
	case "-h", "--help":
		return p.showUsage()
	default:
		// If no subcommand, treat as set
		return p.setProduct(strings.Join(args, " "))
	}
}

// setProduct sets the product name
func (p *ProductHandler) setProduct(name string) error {
	if name == "" {
		return fmt.Errorf("product name cannot be empty")
	}

	// Validate product name (alphanumeric, hyphens, underscores)
	if !p.isValidProductName(name) {
		return fmt.Errorf("invalid product name: %s (use alphanumeric, hyphens, underscores only)", name)
	}

	// Get current product
	currentProduct, _ := p.git.GetConfig("at.product")

	// If same product, just confirm
	if currentProduct == name {
		output.Success("Product already set to: %s", name)
		return nil
	}

	// Show confirmation if changing
	if currentProduct != "" {
		var proceed bool
		err := huh.NewConfirm().
			Title("Change Product").
			Description(fmt.Sprintf("Change product from '%s' to '%s'?", currentProduct, name)).
			Value(&proceed).
			Run()

		if err != nil {
			return fmt.Errorf("failed to get user input: %v", err)
		}

		if !proceed {
			output.Info("Product change cancelled")
			return nil
		}
	}

	// Set the product
	if err := p.git.SetConfig("at.product", name); err != nil {
		return fmt.Errorf("failed to set product: %v", err)
	}

	output.Success("Product set to: %s", name)
	return nil
}

// showProduct shows the current product
func (p *ProductHandler) showProduct() error {
	product, err := p.git.GetConfig("at.product")
	if err != nil {
		output.Info("No product configured")
		return nil
	}

	if product == "" {
		output.Info("No product configured")
		return nil
	}

	output.Success("Current product: %s", product)
	return nil
}

// clearProduct clears the product setting
func (p *ProductHandler) clearProduct() error {
	currentProduct, _ := p.git.GetConfig("at.product")
	if currentProduct == "" {
		output.Info("No product configured")
		return nil
	}

	var proceed bool
	err := huh.NewConfirm().
		Title("Clear Product").
		Description(fmt.Sprintf("Clear product setting '%s'?", currentProduct)).
		Value(&proceed).
		Run()

	if err != nil {
		return fmt.Errorf("failed to get user input: %v", err)
	}

	if !proceed {
		output.Info("Product clear cancelled")
		return nil
	}

	if err := p.git.SetConfig("at.product", ""); err != nil {
		return fmt.Errorf("failed to clear product: %v", err)
	}

	output.Success("Product cleared")
	return nil
}

// isValidProductName validates product name format
func (p *ProductHandler) isValidProductName(name string) bool {
	// Allow alphanumeric, hyphens, underscores, spaces
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' || char == ' ') {
			return false
		}
	}
	return len(strings.TrimSpace(name)) > 0
}

// showUsage displays the usage information
func (p *ProductHandler) showUsage() error {
	usage := `# Product Command

Manages the product name configuration for the repository.

## Usage

  git @ product [<name>]
  git @ product set <name>
  git @ product get
  git @ product clear

## Commands

• **set <name>**: Set the product name
• **get**: Show current product name (default)
• **clear, unset**: Clear the product setting

## Options

• **-h, --help**: Show this help message

## Examples

  # Set product name
  git @ product "My Awesome Product"
  git @ product set "My Awesome Product"

  # Show current product
  git @ product
  git @ product get

  # Clear product setting
  git @ product clear

## Features

• **Validation**: Ensures product names are valid
• **Confirmation**: Confirms changes when updating existing product
• **Persistence**: Stores in Git configuration
• **Integration**: Works with other GitAT commands

## Product Name Rules

• Alphanumeric characters (a-z, A-Z, 0-9)
• Hyphens (-) and underscores (_)
• Spaces allowed
• Cannot be empty

## Use Cases

• **Project Identification**: Identify which product this repository belongs to
• **Multi-Product Organizations**: Manage multiple products in one organization
• **CI/CD Integration**: Use in build and deployment scripts
• **Documentation**: Auto-generate product-specific documentation

## Configuration

Product names are stored in Git configuration:
  git config at.product "Product Name"

## Notes

• Product names are case-sensitive
• Changes are stored in repository configuration
• Can be used by other GitAT commands for context
`

	return output.Markdown(usage)
}
