#!/bin/bash

# GitAT Installation Script
# Automates the installation of the GitAT plugin

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script information
SCRIPT_NAME="GitAT Installer"
VERSION="1.0.0"

# Installation options
INSTALL_METHOD=""
TARGET_DIR=""
BACKUP_EXISTING=false

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_banner() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    GitAT Installation Script                  ║"
    echo "║                        Version $VERSION                        ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Options:
    -m, --method METHOD    Installation method (profile|link|copy)
    -d, --directory DIR    Target directory for copy method
    -b, --backup          Backup existing installation
    -h, --help            Show this help message

Installation Methods:
    profile    Add to shell profile (recommended)
    link       Create symbolic link in /usr/local/bin
    copy       Copy to specified directory

Examples:
    $0 --method profile
    $0 --method link
    $0 --method copy --directory ~/.local/bin
    $0 --method copy --directory /usr/local/bin --backup

EOF
}

detect_shell() {
    local shell_name
    shell_name=$(basename "$SHELL")
    echo "$shell_name"
}

get_shell_profile() {
    local shell_name="$1"
    local profile=""
    
    case "$shell_name" in
        "bash")
            if [ -f "$HOME/.bashrc" ]; then
                profile="$HOME/.bashrc"
            elif [ -f "$HOME/.bash_profile" ]; then
                profile="$HOME/.bash_profile"
            fi
            ;;
        "zsh")
            if [ -f "$HOME/.zshrc" ]; then
                profile="$HOME/.zshrc"
            fi
            ;;
        "fish")
            if [ -f "$HOME/.config/fish/config.fish" ]; then
                profile="$HOME/.config/fish/config.fish"
            fi
            ;;
    esac
    
    echo "$profile"
}

check_dependencies() {
    log_info "Checking dependencies..."
    
    # Check if git is available
    if ! command -v git >/dev/null 2>&1; then
        log_error "Git is not installed. Please install Git first."
        exit 1
    fi
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        log_error "Not in a git repository. Please run this script from the GitAT repository."
        exit 1
    fi
    
    log_success "Dependencies check passed"
}

check_permissions() {
    log_info "Checking script permissions..."
    
    # Make git-@ executable
    if [ ! -x "git-@" ]; then
        log_info "Making git-@ executable..."
        chmod +x git-@
    fi
    
    log_success "Permissions check passed"
}

install_profile_method() {
    log_info "Installing using profile method..."
    
    local shell_name
    local profile
    local current_dir
    
    shell_name=$(detect_shell)
    profile=$(get_shell_profile "$shell_name")
    current_dir=$(pwd)
    
    if [ -z "$profile" ]; then
        log_error "Could not detect shell profile for $shell_name"
        return 1
    fi
    
    # Check if already in PATH
    if echo "$PATH" | grep -q "$current_dir"; then
        log_warning "GitAT directory is already in PATH"
        return 0
    fi
    
    # Add to profile
    log_info "Adding to $profile..."
    
    # Create backup if requested
    if [ "$BACKUP_EXISTING" = true ]; then
        cp "$profile" "${profile}.backup.$(date +%Y%m%d_%H%M%S)"
        log_info "Backup created: ${profile}.backup.$(date +%Y%m%d_%H%M%S)"
    fi
    
    # Add PATH export
    echo "" >> "$profile"
    echo "# GitAT plugin" >> "$profile"
    echo "export PATH=\"\$PATH:$current_dir\"" >> "$profile"
    
    log_success "Added to $profile"
    log_info "Please restart your terminal or run: source $profile"
}

install_link_method() {
    log_info "Installing using symbolic link method..."
    
    local current_dir
    local target_path
    
    current_dir=$(pwd)
    target_path="/usr/local/bin/git-@"
    
    # Check if already exists
    if [ -L "$target_path" ]; then
        if [ "$BACKUP_EXISTING" = true ]; then
            log_info "Backing up existing link..."
            mv "$target_path" "${target_path}.backup.$(date +%Y%m%d_%H%M%S)"
        else
            log_warning "Symbolic link already exists at $target_path"
            return 0
        fi
    fi
    
    # Create symbolic link
    log_info "Creating symbolic link..."
    sudo ln -s "$current_dir/git-@" "$target_path"
    
    log_success "Created symbolic link: $target_path"
}

install_copy_method() {
    log_info "Installing using copy method..."
    
    local current_dir
    local target_dir
    
    current_dir=$(pwd)
    target_dir="${TARGET_DIR:-/usr/local/bin}"
    
    # Check if target directory exists
    if [ ! -d "$target_dir" ]; then
        log_error "Target directory does not exist: $target_dir"
        return 1
    fi
    
    # Check if already exists
    if [ -f "$target_dir/git-@" ]; then
        if [ "$BACKUP_EXISTING" = true ]; then
            log_info "Backing up existing file..."
            mv "$target_dir/git-@" "$target_dir/git-@.backup.$(date +%Y%m%d_%H%M%S)"
        else
            log_warning "GitAT already exists at $target_dir/git-@"
            return 0
        fi
    fi
    
    # Copy file
    log_info "Copying to $target_dir..."
    sudo cp "git-@" "$target_dir/"
    
    log_success "Copied to $target_dir/git-@"
}

detect_install_method() {
    log_info "Auto-detecting best installation method..."
    
    # Check if /usr/local/bin is writable
    if [ -w "/usr/local/bin" ] || sudo -n true 2>/dev/null; then
        INSTALL_METHOD="link"
        log_info "Detected: symbolic link method (recommended)"
    else
        INSTALL_METHOD="profile"
        log_info "Detected: profile method"
    fi
}

verify_installation() {
    log_info "Verifying installation..."
    
    # Check if git-@ script exists and is executable
    if [ -x "git-@" ]; then
        log_success "GitAT script is executable"
    else
        log_error "GitAT script is not executable"
        return 1
    fi
    
    # Check if git-@ is in PATH (may not be available in current session)
    if command -v git-@ >/dev/null 2>&1; then
        log_success "GitAT is available in PATH"
        log_success "GitAT installation verified successfully"
        return 0
    else
        # If not in PATH, check if we just added it to profile
        log_info "GitAT not in current PATH (this is expected after profile installation)"
        log_info "Please restart your terminal or run: source ~/.bashrc (or your shell profile)"
        return 0
    fi
}

run_tests() {
    log_info "Running installation tests..."
    
    if [ -f "package.json" ] && command -v npm >/dev/null 2>&1; then
        if npm run test:unit >/dev/null 2>&1; then
            log_success "Unit tests passed"
        else
            log_warning "Unit tests failed (this may be expected)"
        fi
    else
        log_info "Skipping tests (npm not available)"
    fi
}

print_post_install_info() {
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    Installation Complete!                     ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo
    echo "Next steps:"
    echo "1. Restart your terminal or run: source ~/.bashrc (or your shell profile)"
    echo "2. Navigate to a git repository"
    echo "3. Run: git @"
    echo "4. Explore the available commands"
    echo
    echo "For help:"
    echo "  git @ help"
    echo "  git @ <command> help"
    echo
    echo "Documentation: https://github.com/potsed/gitAT"
    echo
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -m|--method)
            INSTALL_METHOD="$2"
            shift 2
            ;;
        -d|--directory)
            TARGET_DIR="$2"
            shift 2
            ;;
        -b|--backup)
            BACKUP_EXISTING=true
            shift
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Main installation process
main() {
    print_banner
    
    # Check dependencies
    check_dependencies
    
    # Check permissions
    check_permissions
    
    # Detect installation method if not specified
    if [ -z "$INSTALL_METHOD" ]; then
        detect_install_method
    fi
    
    # Perform installation
    case "$INSTALL_METHOD" in
        "profile")
            install_profile_method
            ;;
        "link")
            install_link_method
            ;;
        "copy")
            install_copy_method
            ;;
        *)
            log_error "Invalid installation method: $INSTALL_METHOD"
            print_usage
            exit 1
            ;;
    esac
    
    # Verify installation
    if verify_installation; then
        log_success "Installation completed successfully!"
        
        # Run tests if available
        run_tests
        
        # Print post-install information
        print_post_install_info
    else
        log_error "Installation verification failed"
        echo
        echo "Troubleshooting:"
        echo "1. Make sure you're in the GitAT repository directory"
        echo "2. Check that git-@ is executable: ls -la git-@"
        echo "3. Verify your PATH includes the GitAT directory"
        echo "4. Try restarting your terminal"
        exit 1
    fi
}

# Run main function
main "$@" 