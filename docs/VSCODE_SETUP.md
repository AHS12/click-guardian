# VSCode Go Extension Configuration

This guide helps you set up VSCode for optimal Go + Fyne development with this project.

## The Problem

VSCode's Go language server (gopls) tries to analyze all build tags and dependencies, including platform-specific ones. This can cause false positive errors like:

```
error while importing fyne.io/fyne/v2/internal/painter/gl: build constraints exclude all Go files in C:\Users\...\go-gl\gl@v0.0.0-...\v2.1\gl [darwin]
```

The `[darwin]` indicates VSCode is trying to process macOS-specific code on Windows.

## Solution: Use the Workspace File

This project includes a pre-configured VSCode workspace file that solves all these issues automatically.

### Setup Instructions

**Open the workspace file instead of the folder:**

```
File → Open Workspace from File → click-guardian.code-workspace
```

### What the Workspace File Provides

The workspace file (`click-guardian.code-workspace`) includes:

- **Go + CGO + Fyne optimized settings** - Forces CGO_ENABLED=1, sets GOOS=windows
- **Platform-specific build constraints** - Configures gopls with Windows-specific build flags
- **Recommended extensions** - Go and C++ extensions (needed for CGO)
- **Optimized gopls configuration** - Disables problematic static analysis
- **Build environment setup** - Platform-specific build environment

## Troubleshooting

### Automated Diagnosis

Run the automated troubleshooting script:

```bash
scripts\troubleshoot.bat
```

This script will:

- Check Go and CGO setup
- Test CGO compilation
- Verify project builds
- Clean and refresh dependencies

### Manual Steps

If issues persist after using the workspace file:

**1. Restart Go Language Server**

- `Ctrl+Shift+P` → "Go: Restart Language Server"

**2. Update Go Tools**

- `Ctrl+Shift+P` → "Go: Install/Update Tools"
- Select all tools and update

**3. Clear Go Module Cache**

```bash
go clean -modcache
go mod download
```

**4. Reload VSCode Window**

```bash
Ctrl+Shift+P → "Developer: Reload Window"
```

## Important Notes

- **Always use the workspace file** - Don't open the folder directly
- **These errors are cosmetic** - Your builds will work fine regardless
- **The workspace file handles everything** - No manual configuration needed
- **CGO is required** - For Windows API integration with Fyne
