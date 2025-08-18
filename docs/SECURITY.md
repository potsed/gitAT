# Security Guidelines

## Overview

GitAT takes security seriously and follows best practices for dependency management, code security, and safe execution of Git commands.

## Dependency Security

### Regular Updates

Dependencies should be regularly updated to address security vulnerabilities:

```bash
# Update security-critical dependencies
make security-update

# Update all dependencies
make update-deps

# Check for vulnerabilities
go list -json -deps ./... | nancy sleuth
```

### Key Security Dependencies

The following dependencies are critical for security and should be kept up to date:

- `golang.org/x/net` - Network libraries (HTTP, proxy handling)
- `golang.org/x/sys` - System calls and OS interfaces
- `golang.org/x/text` - Text processing and encoding
- `golang.org/x/term` - Terminal handling

### Vulnerability Monitoring

1. **GitHub Security Alerts**: Monitor repository security alerts
2. **Go Vulnerability Database**: Use `govulncheck` for scanning
3. **Dependency Updates**: Regular updates via `make security-update`

## Code Security

### Git Command Execution

GitAT executes Git commands through a secure fallthrough mechanism:

#### Argument Validation
- All arguments are validated for command injection patterns
- Dangerous patterns like `$(...)` and `` `...` `` are blocked
- Arguments are passed safely to subprocess execution

#### Command Blacklisting
- Dangerous commands can be blacklisted via configuration
- Reserved GitAT commands are protected from fallthrough
- Validation prevents bypassing security controls

#### Process Execution
- Commands are executed with proper argument separation
- No shell interpretation of arguments
- Exit codes and errors are properly handled

### Input Sanitization

```go
// Example: Safe argument validation
func ValidateArguments(args []string) error {
    for i, arg := range args {
        if strings.Contains(arg, "$(") || strings.Contains(arg, "`") {
            return fmt.Errorf("potentially unsafe argument at position %d: %s", i, arg)
        }
    }
    return nil
}
```

### Configuration Security

- Configuration files are validated for malicious content
- Blacklist entries are checked for duplicates and empty values
- Configuration tampering is prevented through validation

## Build Security

### Binary Naming Protection

Multiple safeguards prevent building with deprecated binary names:

1. **Makefile Validation**: Prevents building `gitat` instead of `git-@`
2. **Pre-commit Hooks**: Blocks committing wrong binary names
3. **CI/CD Validation**: Automated checks in build pipeline

### Secure Build Process

```bash
# Validate before building
make validate

# Clean build
make clean && make build-local

# Verify binary
./git-@ version
```

## Runtime Security

### Command Injection Prevention

GitAT prevents command injection through multiple layers:

1. **Argument Validation**: Blocks dangerous patterns
2. **Process Isolation**: Uses Go's `exec.Command` safely
3. **No Shell Execution**: Arguments passed directly to Git

### Path Traversal Protection

- File paths are validated by Git itself
- No custom path resolution that could be exploited
- Git handles path security according to its own rules

### Error Information Disclosure

- Error messages provide helpful information without exposing sensitive data
- Stack traces are not exposed to end users
- Debug information is only available in verbose mode

## Security Best Practices

### For Developers

1. **Validate All Input**: Never trust user input
2. **Use Safe APIs**: Prefer `exec.Command` over shell execution
3. **Minimize Privileges**: Run with least necessary permissions
4. **Regular Updates**: Keep dependencies current
5. **Security Testing**: Include security tests in test suite

### For Users

1. **Keep Updated**: Regularly update GitAT to latest version
2. **Review Configuration**: Audit blacklist and settings
3. **Monitor Commands**: Use verbose mode to see executed commands
4. **Report Issues**: Report security concerns promptly

### Configuration Security

```bash
# Review current configuration
git config --list | grep "at\."

# Validate fallthrough settings
git config at.fallthrough.enabled
git config at.fallthrough.blacklist

# Enable verbose mode for auditing
git config at.fallthrough.verbose true
```

## Incident Response

### If Security Issue is Discovered

1. **Immediate Action**:
   - Stop using affected functionality
   - Update to latest version if fix is available
   - Review logs for potential exploitation

2. **Reporting**:
   - Report security issues via GitHub Security Advisories
   - Include reproduction steps and impact assessment
   - Do not disclose publicly until fix is available

3. **Mitigation**:
   - Disable fallthrough if needed: `git config at.fallthrough.enabled false`
   - Use blacklist to block problematic commands
   - Monitor for suspicious activity

### Security Checklist

- [ ] Dependencies are up to date
- [ ] No deprecated binary names in use
- [ ] Fallthrough configuration is reviewed
- [ ] Blacklist includes necessary commands
- [ ] Verbose logging is enabled for auditing
- [ ] Build process uses correct binary name
- [ ] Pre-commit hooks are installed
- [ ] Security tests are passing

## Security Updates

### Recent Security Fixes

#### August 2025
- **golang.org/x/net**: Updated from v0.33.0 to v0.43.0
  - Fixed Cross-site Scripting vulnerability
  - Fixed HTTP Proxy bypass using IPv6 Zone IDs
- **All dependencies**: Updated to latest secure versions

### Update Process

```bash
# Check for security updates
go list -m -u all | grep golang.org/x

# Apply security updates
make security-update

# Verify updates
go mod why golang.org/x/net
```

## Contact

For security-related questions or to report vulnerabilities:
- GitHub Security Advisories (preferred)
- Create an issue with "Security" label
- Email maintainers for critical issues

## References

- [Go Security Policy](https://golang.org/security)
- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)
- [Git Security Documentation](https://git-scm.com/docs/git-config#Documentation/git-config.txt-safedirectory)