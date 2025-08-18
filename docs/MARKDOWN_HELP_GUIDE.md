# Markdown Help Guide

## 🎯 **Overview**

All GitAT command help text now uses [Charmbracelet Glamour](https://github.com/charmbracelet/glamour) for beautiful markdown rendering in the terminal. This provides consistent, professional-looking help documentation across all commands.

## 🚀 **Implementation Status**

### ✅ **Completed Commands**

- **Work Command**: Full markdown help with usage, arguments, options, examples, workflow, and special cases
- **Version Command**: Complete markdown help with semantic versioning documentation
- **Save Command**: Basic markdown help (placeholder - ready for full implementation)

### 🔄 **In Progress**

- **Multi-Team Version Handler**: Complete markdown help with approval workflow documentation

### 📋 **Pending Commands**

- **Squash Command**: Needs markdown help implementation
- **Pull Request Command**: Needs markdown help implementation
- **Branch Command**: Needs markdown help implementation
- **Sweep Command**: Needs markdown help implementation
- **WIP Command**: Needs markdown help implementation
- **Hotfix Command**: Needs markdown help implementation
- **Utility Commands**: Needs markdown help implementation
- **Info Command**: Needs markdown help implementation
- **Release Command**: Needs markdown help implementation
- **Feature Command**: Needs markdown help implementation
- **Product Command**: Needs markdown help implementation
- **Issue Command**: Needs markdown help implementation
- **Label Command**: Needs markdown help implementation
- **Trunk Command**: Needs markdown help implementation
- **Ignore Command**: Needs markdown help implementation
- **Init Commands**: Needs markdown help implementation
- **Security Command**: Needs markdown help implementation
- **Go Command**: Needs markdown help implementation

## 🛠️ **Implementation Pattern**

### **1. Basic Structure**

```go
// Execute handles the command
func (h *Handler) Execute(args []string) error {
    if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
        return h.showUsage()
    }
    // Command implementation
    return nil
}

// showUsage displays the command usage
func (h *Handler) showUsage() error {
    return output.Markdown(`# Command Name

Brief description of what the command does.

## Usage

` + "```" + `bash
git @ command [options] [arguments]
` + "```" + `

## Arguments

- **argument**: Description of the argument

## Options

- **-o, --option**: Description of the option
- **-h, --help**: Show this help message

## Examples

` + "```" + `bash
# Example 1
git @ command example1

# Example 2
git @ command example2
` + "```" + `

## Workflow

1. Step one
2. Step two
3. Step three

## Special Cases

- **Special case 1**: Description
- **Special case 2**: Description
`)
}
```

### **2. Advanced Structure with Sections**

```go
func (h *Handler) showUsage() error {
    return output.Markdown(`# Command Name

Comprehensive description with multiple sections.

## Usage

` + "```" + `bash
git @ command <required> [optional]
git @ command --flag value
` + "```" + `

## Arguments

- **required**: Required argument description
- **optional**: Optional argument description

## Options

### Basic Options
- **-f, --flag**: Basic flag description
- **-v, --verbose**: Enable verbose output

### Advanced Options
- **--config**: Path to configuration file
- **--dry-run**: Preview changes without applying

## Examples

### Basic Usage
` + "```" + `bash
# Simple example
git @ command basic-arg
` + "```" + `

### Advanced Usage
` + "```" + `bash
# Complex example
git @ command --config config.yaml --dry-run
` + "```" + `

## Workflow

1. **Preparation**: Initial setup steps
2. **Execution**: Main command execution
3. **Verification**: Post-execution checks

## Configuration

### Environment Variables
- `ENV_VAR`: Description of environment variable

### Configuration Files
- `config.yaml`: Description of config file

## Troubleshooting

### Common Issues
- **Issue 1**: Solution description
- **Issue 2**: Solution description

### Error Messages
- `Error: message`: What it means and how to fix
`)
}
```

## 📝 **Markdown Features Used**

### **1. Headers**

```markdown
# Main Title
## Section Title
### Subsection Title
```

### **2. Code Blocks**

```markdown
` + "```" + `bash
git @ command example
` + "```" + `
```

### **3. Lists**

```markdown
## Unordered List
- Item 1
- Item 2
- Item 3

## Ordered List
1. First step
2. Second step
3. Third step

## Definition List
- **Term**: Definition
- **Another Term**: Another definition
```

### **4. Emphasis**

```markdown
**Bold text** for important information
*Italic text* for emphasis
`code` for inline code
```

### **5. Tables**

```markdown
| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Data 1   | Data 2   | Data 3   |
| Data 4   | Data 5   | Data 6   |
```

## 🎨 **Styling Guidelines**

### **1. Consistent Structure**

- Always start with a main title (`# Command Name`)
- Include a brief description
- Use consistent section headers
- End with troubleshooting or special cases

### **2. Code Examples**

- Use `bash` syntax highlighting for shell commands
- Include realistic examples
- Show both simple and complex usage
- Add comments to explain examples

### **3. Argument Documentation**

- Use bold for argument names
- Provide clear descriptions
- Indicate required vs optional
- Show data types when relevant

### **4. Option Documentation**

- Group related options together
- Use consistent formatting
- Include both short and long forms
- Explain default values

## 🔧 **Implementation Steps**

### **Step 1: Add Help Detection**

```go
func (h *Handler) Execute(args []string) error {
    if len(args) > 0 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
        return h.showUsage()
    }
    // Command implementation
    return nil
}
```

### **Step 2: Create showUsage Method**

```go
func (h *Handler) showUsage() error {
    return output.Markdown(`# Command Documentation`)
}
```

### **Step 3: Add Comprehensive Documentation**

- Usage examples
- Argument descriptions
- Option explanations
- Workflow steps
- Special cases
- Troubleshooting

### **Step 4: Test the Help**

```bash
git @ command --help
git @ command -h
git @ command help
```

## 📊 **Command Priority List**

### **High Priority (Core Commands)**

1. **Sweep Command** - Critical for branch management
2. **Save Command** - Essential for workflow
3. **Pull Request Command** - Important for collaboration
4. **Branch Command** - Core functionality

### **Medium Priority (Workflow Commands)**

5. **WIP Command** - Work in progress management
6. **Hotfix Command** - Emergency fixes
7. **Release Command** - Release management
8. **Feature Command** - Feature development

### **Low Priority (Utility Commands)**

9. **Info Command** - Information display
10. **Utility Commands** - Various utilities
11. **Product Command** - Product management
12. **Issue Command** - Issue tracking
13. **Label Command** - Label management
14. **Trunk Command** - Trunk operations
15. **Ignore Command** - Ignore management
16. **Init Commands** - Initialization
17. **Security Command** - Security features
18. **Go Command** - Go-specific features

## 🧪 **Testing Guidelines**

### **1. Visual Testing**

```bash
# Test help display
git @ command --help

# Verify formatting
git @ command -h

# Check consistency
git @ command help
```

### **2. Content Testing**

- Verify all sections are present
- Check code examples are correct
- Ensure argument descriptions are accurate
- Validate option documentation

### **3. Integration Testing**

- Test with different terminal sizes
- Verify markdown rendering in various environments
- Check for any rendering issues

## 🎯 **Best Practices**

### **1. Content Guidelines**

- Keep descriptions concise but informative
- Use consistent terminology
- Include practical examples
- Document edge cases and special scenarios

### **2. Formatting Guidelines**

- Use consistent header levels
- Maintain proper spacing
- Ensure code blocks are properly formatted
- Use bold for emphasis, not overuse

### **3. User Experience**

- Start with simple examples
- Progress to complex usage
- Include troubleshooting section
- Provide clear next steps

## 🔄 **Migration Checklist**

For each command that needs markdown help:

- [ ] Add help detection to Execute method
- [ ] Create showUsage method
- [ ] Write comprehensive markdown documentation
- [ ] Include usage examples
- [ ] Document all arguments and options
- [ ] Add workflow description
- [ ] Include special cases
- [ ] Add troubleshooting section
- [ ] Test help display
- [ ] Verify formatting
- [ ] Update any related documentation

## 📚 **Resources**

- [Charmbracelet Glamour Documentation](https://github.com/charmbracelet/glamour)
- [Markdown Syntax Guide](https://www.markdownguide.org/)
- [GitAT Output Package](../pkg/output/output.go)
- [Existing Command Examples](../internal/commands/handlers/)

---

**The markdown help system provides a consistent, professional, and user-friendly experience across all GitAT commands. Each command should follow the established patterns to maintain consistency and provide the best possible user experience.**
