# VSCode Go Extension Configuration

This file helps resolve platform-specific build constraint issues in VSCode.

## The Problem

VSCode's Go language server (gopls) tries to analyze all build tags and dependencies, including platform-specific ones. This can cause false positive errors like:

```
error while importing fyne.io/fyne/v2/internal/painter/gl: build constraints exclude all Go files in C:\Users\...\go-gl\gl@v0.0.0-...\v2.1\gl [darwin]
```

The `[darwin]` indicates VSCode is trying to process macOS-specific code on Windows.

## Solutions Applied

### 1. Workspace Settings (`.vscode/settings.json`)

- Forces CGO_ENABLED=1 for Fyne
- Sets GOOS=windows explicitly
- Configures gopls with Windows-specific build flags
- Disables problematic static analysis

### 2. Gopls Configuration (`gopls.yaml`)

- Root-level gopls configuration file
- Platform-specific build environment
- Optimized for Windows CGO compilation

### 3. VSCode Workspace File (`go-double-click-fix.code-workspace`)

- Provides workspace-level configuration
- Recommended extensions for Go and C++ (needed for CGO)

### 4. Build Constraints

- Platform-specific build tags in hook files
- Prevents cross-platform compilation issues

### 3. Usage

Open the workspace file instead of the folder:

```
File → Open Workspace from File → go-double-click-fix.code-workspace
```

Or reload VSCode window after the settings are applied:

```
Ctrl+Shift+P → "Developer: Reload Window"
```

### 4. Troubleshooting Script

Run the automated troubleshooting script:

```bash
scripts\troubleshoot.bat
```

This script will:

- Check Go and CGO setup
- Test CGO compilation
- Verify project builds
- Clean and refresh dependencies

## Alternative Solutions

If the error persists, try these steps:

### 1. Clear Go Module Cache

```bash
go clean -modcache
go mod download
```

### 2. Restart Go Language Server

- `Ctrl+Shift+P`
- Type "Go: Restart Language Server"

### 3. Update Go Tools

- `Ctrl+Shift+P`
- Type "Go: Install/Update Tools"
- Select all tools and update

### 4. Environment Variables

Add to your system environment or `.bashrc`:

```bash
export CGO_ENABLED=1
export GOOS=windows
export GOARCH=amd64
```

## Note

These errors are cosmetic in VSCode - your builds will work fine regardless. The configuration files help VSCode's language server understand your platform-specific setup better.
