# Build Scripts

This directory contains build and development scripts for the Click Guardian project.

## Scripts

### `build.bat` (Windows)

Main build script that creates both GUI and console versions of the application.

- Creates executables in `dist/` directory
- Builds GUI version (no console window) and console version (for debugging)
- Handles build failures gracefully

### `build.sh` (Linux/macOS)

Cross-platform build script that creates binaries for multiple platforms.

- Requires bash shell
- Builds for Windows, Linux, and macOS
- Creates both GUI and console versions for Windows

### `dev.bat` (Windows)

Development script for quick testing.

- Runs the application directly without building
- Useful for rapid development iteration

### `troubleshoot.bat` (Windows)

Troubleshooting script for Go/CGO setup issues.

- Tests CGO compilation
- Verifies project builds
- Cleans and refreshes dependencies

## Usage

### Windows Development

```cmd
# From anywhere in the project
scripts\dev.bat

# Production build
scripts\build.bat

# Troubleshoot Go/CGO setup
scripts\troubleshoot.bat
```

### Cross-platform Build

```bash
# Make executable (first time only)
chmod +x scripts/build.sh

# Build for all platforms (from anywhere in project)
scripts/build.sh
```

**Note:** All scripts automatically navigate to the project root directory, so they can be run from anywhere within the project.

### Windows Development

```cmd
# Quick development testing
scripts\dev.bat

# Production build
scripts\build.bat
```

### Cross-platform Build

```bash
# Make executable (first time only)
chmod +x scripts/build.sh

# Build for all platforms
./scripts/build.sh
```

## Output

All builds output to the `dist/` directory, which is git-ignored.

### Naming Convention

- `click-guardian-gui.exe` - Windows GUI version (recommended)
- `click-guardian.exe` - Windows console version (debugging)
- `click-guardian-{os}-{arch}` - Cross-platform builds
