# Multi-Team Version Management Example

## 🎯 **Scenario: E-commerce Platform**

We have a multi-team e-commerce platform with the following teams:

- **Frontend Team**: React web application
- **Backend Team**: Node.js API services
- **Mobile Team**: React Native mobile app
- **DevOps Team**: Infrastructure and deployment
- **QA Team**: Testing and quality assurance

## 📋 **Current Situation**

```
Current Versions:
- Frontend: 2.0.5
- Backend: 2.0.5
- Mobile: 2.0.5
- DevOps: 2.0.5
- QA: 2.0.5

All teams are synchronized at version 2.0.5
```

## 🚀 **Feature Development: Payment System Enhancement**

### **Step 1: Backend Team Proposes Version Change**

The backend team has implemented new payment API endpoints and needs to bump the version.

```bash
# Backend team sets their identity
git config at.team backend

# Propose version change
git @ version multi-team propose 2.1.0 "Add new payment API endpoints for enhanced checkout flow"

# Output:
# ✅ Version proposal created: prop_1703123456
# Proposed version: 2.1.0
# Reason: Add new payment API endpoints for enhanced checkout flow
# Teams to approve: frontend, mobile, devops, qa
```

### **Step 2: Teams Review and Approve**

#### **Frontend Team Approval**

```bash
# Frontend team sets their identity
git config at.team frontend

# Check pending proposals
git @ version multi-team

# Output:
# 📋 Pending Proposals
# Proposal prop_1703123456: 2.0.5 -> 2.1.0 by backend
# Reason: Add new payment API endpoints for enhanced checkout flow

# Approve the proposal
git @ version multi-team approve prop_1703123456

# Output:
# ✅ Approval recorded. Waiting for other teams...
```

#### **Mobile Team Approval**

```bash
# Mobile team sets their identity
git config at.team mobile

# Approve the proposal
git @ version multi-team approve prop_1703123456

# Output:
# ✅ Approval recorded. Waiting for other teams...
```

#### **DevOps Team Approval**

```bash
# DevOps team sets their identity
git config at.team devops

# Approve the proposal
git @ version multi-team approve prop_1703123456

# Output:
# ✅ Approval recorded. Waiting for other teams...
```

#### **QA Team Approval**

```bash
# QA team sets their identity
git config at.team qa

# Approve the proposal
git @ version multi-team approve prop_1703123456

# Output:
# ✅ Version change approved and applied: 2.1.0
```

### **Step 3: Version Synchronization**

Now all teams need to sync to the new version.

```bash
# Sync all teams to version 2.1.0
git @ version multi-team sync

# Output:
# 🔄 Synchronizing Versions Across Teams
# ✅ Successfully synchronized all teams to version 2.1.0
```

### **Step 4: Frontend Team Implements UI Changes**

The frontend team now implements the UI for the new payment features.

```bash
# Frontend team creates a work branch
git @ work feature enhanced-payment-ui

# After implementing changes, they want to propose a minor version bump
git @ version multi-team propose 2.1.1 "Add enhanced payment UI with new payment methods"

# Output:
# ✅ Version proposal created: prop_1703123789
# Proposed version: 2.1.1
# Reason: Add enhanced payment UI with new payment methods
# Teams to approve: backend, mobile, devops, qa
```

### **Step 5: Release Window Coordination**

Before the release, the DevOps team locks the version to prevent conflicts.

```bash
# DevOps team locks version for release window
git @ version multi-team lock

# Output:
# 🔒 Lock Version
# ✅ Version locked by devops
# Reason: Release window - deploying payment system
# Expires: 2024-01-15 18:00:00
```

### **Step 6: Release Execution**

During the release window:

```bash
# Check current status
git @ version multi-team

# Output:
# 🌐 Multi-Team Version Status
# ┌─────────────────┬─────────────────────────┐
# │ Property        │ Value                   │
# ├─────────────────┼─────────────────────────┤
# │ Current Version │ 2.1.0                   │
# │ Lock Status     │ Locked by devops        │
# │ Team Approvals  │ frontend, mobile, qa    │
# └─────────────────┴─────────────────────────┘
# 
# 📋 Pending Proposals
# Proposal prop_1703123789: 2.1.0 -> 2.1.1 by frontend
# Reason: Add enhanced payment UI with new payment methods
```

### **Step 7: Post-Release Unlock**

After successful deployment:

```bash
# DevOps team unlocks version
git @ version multi-team unlock

# Output:
# 🔓 Unlock Version
# ✅ Version unlocked successfully
```

## 🚨 **Emergency Scenario: Critical Bug Found**

### **Problem**

A critical security vulnerability is discovered in the payment system. All teams need to rollback to version 2.0.5 immediately.

### **Solution**

```bash
# DevOps team creates emergency rollback proposal
git @ version multi-team propose 2.0.5 "EMERGENCY: Critical security vulnerability in payment system - immediate rollback required"

# Output:
# ✅ Version proposal created: prop_1703124000
# Proposed version: 2.0.5
# Reason: EMERGENCY: Critical security vulnerability in payment system - immediate rollback required
# Teams to approve: frontend, backend, mobile, qa

# All teams approve within 1 hour (emergency procedure)
git @ version multi-team approve prop_1703124000

# Output:
# ✅ Version change approved and applied: 2.0.5

# Sync all teams to rollback version
git @ version multi-team sync

# Output:
# 🔄 Synchronizing Versions Across Teams
# ✅ Successfully synchronized all teams to version 2.0.5
```

## 📊 **Monitoring and Metrics**

### **Version Health Check**

```bash
# Check version drift
git @ version multi-team teams

# Output:
# 👥 Team Status
# ┌──────────┬─────────┬─────────────────────┬────────┐
# │ Team     │ Version │ Last Sync           │ Status │
# ├──────────┼─────────┼─────────────────────┼────────┤
# │ frontend │ 2.0.5   │ 2024-01-15 16:30:00 │ active │
# │ backend  │ 2.0.5   │ 2024-01-15 16:30:00 │ active │
# │ mobile   │ 2.0.5   │ 2024-01-15 16:30:00 │ active │
# │ devops   │ 2.0.5   │ 2024-01-15 16:30:00 │ active │
# │ qa       │ 2.0.5   │ 2024-01-15 16:30:00 │ active │
# └──────────┴─────────┴─────────────────────┴────────┘
```

### **Version History**

```bash
# View version change history
git @ version multi-team history

# Output:
# 📜 Version History
# Last 20 version changes:
# [2024-01-15 16:30:00] Version change approved and applied: 2.0.5
# [2024-01-15 16:25:00] Version proposal prop_1703124000: 2.1.0 -> 2.0.5 by devops
# [2024-01-15 16:20:00] Version unlocked by devops
# [2024-01-15 15:45:00] Version locked by devops: Release window - deploying payment system
# [2024-01-15 15:30:00] Version change approved and applied: 2.1.0
# [2024-01-15 15:15:00] Version proposal prop_1703123456: 2.0.5 -> 2.1.0 by backend
```

## 🔧 **Configuration Files**

### **Team Configuration**

```bash
# .git/config (for each team member)
[at]
    team = frontend  # or backend, mobile, devops, qa
    major = 2
    minor = 0
    fix = 5
```

### **Version Lock File**

```bash
# .git/gitat-logs/version-lock.json
{
  "locked": true,
  "team": "devops",
  "reason": "Release window - deploying payment system",
  "expiresAt": "2024-01-15T18:00:00Z",
  "createdAt": "2024-01-15T15:45:00Z"
}
```

### **Proposal File**

```bash
# .git/gitat-logs/proposals/prop_1703123456.json
{
  "id": "prop_1703123456",
  "proposedVersion": "2.1.0",
  "currentVersion": "2.0.5",
  "reason": "Add new payment API endpoints for enhanced checkout flow",
  "proposingTeam": "backend",
  "teamsToApprove": ["frontend", "mobile", "devops", "qa"],
  "approvals": {
    "frontend": "2024-01-15T15:15:00Z",
    "mobile": "2024-01-15T15:20:00Z",
    "devops": "2024-01-15T15:25:00Z",
    "qa": "2024-01-15T15:30:00Z"
  },
  "status": "approved",
  "createdAt": "2024-01-15T15:10:00Z",
  "appliedAt": "2024-01-15T15:30:00Z"
}
```

## 📈 **Best Practices Demonstrated**

### **1. Clear Communication**

- All version changes have detailed reasons
- Teams are notified of pending proposals
- Emergency procedures are clearly defined

### **2. Coordinated Releases**

- Release windows with version locking
- Phased deployment (backend → frontend → mobile)
- All teams synchronized after changes

### **3. Emergency Procedures**

- Quick rollback capability
- Emergency override mechanisms
- Clear audit trail for all changes

### **4. Team Autonomy**

- Teams can propose changes independently
- Approval workflow ensures consensus
- No single point of failure

### **5. Audit and Compliance**

- Complete history of all version changes
- Team accountability for approvals
- Timestamped logs for compliance

## 🎯 **Key Takeaways**

1. **Proposal Workflow**: All version changes go through a proposal and approval process
2. **Team Coordination**: Clear communication and coordination between teams
3. **Emergency Procedures**: Quick response to critical issues with rollback capability
4. **Audit Trail**: Complete logging of all version-related activities
5. **Flexibility**: Support for both coordinated and independent team releases

This example demonstrates how GitAT's multi-team version management provides the structure and tooling needed for successful coordination across multiple development teams while maintaining team autonomy and ensuring release quality.
