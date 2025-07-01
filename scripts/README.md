# Build Scripts

This directory contains build and development scripts for the Click Guardian project.

## Quick Reference

```cmd
# Development (quick build and run)
scripts\dev.bat

# Production build
scripts\build.bat

# Troubleshoot setup issues
scripts\troubleshoot.bat

# Cross-platform build (Linux/macOS)
scripts/build.sh

# release build(windows)
scripts/release-build.bat
```

## Documentation

For complete build instructions, script details, manual build commands, and troubleshooting, see:

**📖 [Build Instructions](../docs/BUILD.md)**

This document covers:

- Detailed script descriptions and usage
- Manual build commands for all platforms
- Development workflow
- Prerequisites and setup
- Troubleshooting common build issues
- Platform-specific notes

## Scripts Overview

- **`build.bat`** - Windows production build (GUI + console versions)
- **`build.sh`** - Cross-platform build script
- **`dev.bat`** - Quick development testing
- **`troubleshoot.bat`** - Go/CGO setup diagnosis
- **`release-build.bat`** - Production release build (windows)

All scripts can be run from anywhere in the project directory.
