# Release Build Guide

This guide explains how to create professional release builds of Click Guardian with proper versioning and packaging.

## 🚀 Quick Start

For a simple release build:

```cmd
scripts\release-build.bat
```

This will create:

- `dist\click-guardian.exe` - Main application
- `dist\click-guardian-v1.0.4-windows-portable.zip` - Complete release package
```

## 📋 Prerequisites

### Required

- **Go 1.24.1+** with CGO enabled
- **Git** (for version information)

### Optional

- **Windows SDK** (for code signing with `signtool`)
- **Code signing certificate** (for production releases)

## 📝 Version Management

### Change Version Number

Edit **ONE** file to change the version:

**File: `build\build.conf`**

```ini
VERSION=1.0.4
```

That's it! The build script automatically uses this version for:

- Executable metadata
- Package naming
- Release documentation

The release build script automatically updates the following files with the new version:

- `wix.json` - Windows installer configuration
- `build/windows/app.rc` - Windows resource file (contains version info embedded in executable)
- `build/windows/app-manifest.xml` - Application manifest file

### Version Auto-Detection

The build script automatically includes:

- **Git commit hash** - from `git rev-parse --short HEAD`
- **Build timestamp** - current date/time
- **Builder name** - from `git config user.name` or Windows username

## 🏗️ Build Process

The simplified build process:

1. **Load version** from `build\build.conf`
2. **Get git info** (commit, builder name)
3. **Generate Windows resource file** (icon, manifest, version info)
   ```cmd
   windres build\windows\app.rc -O coff -o cmd\click-guardian\click-guardian.syso
   ```
4. **Build executables** with embedded version info
5. **Create release package** with documentation
6. **Generate ZIP file** ready for distribution

### What Gets Built

- **GUI Version** (`click-guardian.exe`) - Main application for end users  
  _(Windows icon, manifest, and version info are embedded via `click-guardian.syso`)_
- **Release Package** (`click-guardian-v1.0.4-windows-portable.zip`) - Complete distribution package

## 🔧 Manual Build Commands

If you prefer to build manually:

### Basic Build (Same as your original build.bat)

```cmd
# GUI version (recommended for users)
go build -ldflags "-s -w -H=windowsgui" -o dist\click-guardian-gui.exe .\cmd\click-guardian

# Console version (for debugging)
go build -ldflags "-s -w" -o dist\click-guardian.exe .\cmd\click-guardian
```

### Build with Version Information

```cmd
# Set your version
set VERSION=1.0.4
set GIT_COMMIT=abc1234
set BUILD_BY=YourName

# Build with version info
go build -ldflags "-s -w -H=windowsgui -X click-guardian/internal/version.Version=%VERSION% -X click-guardian/internal/version.GitCommit=%GIT_COMMIT% -X click-guardian/internal/version.BuildBy=%BUILD_BY%" -o dist\click-guardian-gui.exe .\cmd\click-guardian
```

## 📝 Windows Resource File (`.syso`)

The build process uses a Windows resource file to embed the application icon, manifest, and version info into the executable.

- **Resource script:** `build/windows/app.rc`
- **Icon:** `build/windows/app-icon.ico`
- **Manifest:** `build/windows/app-manifest.xml`
- **How to generate:**
  ```cmd
  windres build\windows\app.rc -O coff -o cmd\click-guardian\click-guardian.syso
  ```
- The `.syso` file must be in the same directory as `main.go` (`cmd\click-guardian\`).

## 🔐 Code Signing (Optional)

### For Production Releases

1. **Get a certificate** from a trusted CA (DigiCert, Sectigo, etc.)
2. **Use the signing script**:
   ```cmd
   build\scripts\sign-code.bat dist\click-guardian-gui.exe path\to\cert.pfx password
   ```

### Self-Signed Certificate (Testing Only)

```cmd
# Create self-signed certificate (Windows will show warnings)
makecert -sv mykey.pvk -n "CN=YourName" mycert.cer
pvk2pfx -pvk mykey.pvk -spc mycert.cer -pfx mycert.pfx
```

## 📦 Release Package Contents

The ZIP package includes:

```
click-guardian-v1.0.4-windows/
├── click-guardian-gui.exe      # Main application
├── README.txt                  # Usage instructions
└── LICENSE                     # License (if present)
```

## 🚀 Distribution

### GitHub Releases (Recommended)\n\n1. **Update version** in `build\\build.conf`\n2. **Build release**:\n   ```cmd\n   scripts\
elease-build.bat\n   ```\n3. **Create git tag**:\n   ```cmd\n   git tag -a v1.0.4 -m \"Release version 1.0.4\"\n   git push origin v1.0.4\n   ```\n4. **Upload to GitHub**:\n   - Go to GitHub → Releases → Create new release\n   - Upload `dist\\click-guardian-v1.0.4-windows-portable.zip`

### Direct Distribution

- Share the ZIP file directly
- Upload to your website
- Distribute the signed executables

## 🔍 Testing Your Build

### Check Version Information

```cmd
# The GUI version doesn't show console output, but you can verify it was built correctly
# by checking the file exists and running it to see the About dialog
```

### Verify Functionality

1. **GUI Version** - Should start without console window and show the main interface
2. **About Dialog** - Check "About" button to verify version info is embedded
3. **Protection** - Test double-click blocking works

## 🐛 Simple Troubleshooting

### Build Fails

```cmd
# Check Go environment
go version
go env

# Clean and retry
rmdir /s /q dist
scripts\release-build.bat
```

### No Version Info

- Make sure `build\build.conf` exists with `VERSION=1.0.4`
- Check that Git is installed and working

### Large File Size

- Executable is ~24MB (normal for Go with CGO and GUI)
- Use the original `scripts\build.bat` for smaller builds without version info if needed

## 📂 Files You Might Need to Edit

### For Version Changes

- **`build\build.conf`** - Change `VERSION=1.0.4` to your new version

### For Build Customization

- **`scripts\release-build.bat`** - Modify build process
- **`scripts\build.bat`** - Your original working build script

### For App Changes

- **`cmd\click-guardian\main.go`** - Main application entry point
- **`internal\version\version.go`** - Version handling code

---

## 🎯 Quick Reference

**To release a new version:**

1. Edit `build\build.conf` → change VERSION
2. Run `scripts\release-build.bat`
3. Upload `dist\click-guardian-v1.0.4-windows-portable.zip`

**Simple build (no packaging):**

```cmd
scripts\build.bat
```

**Release build (with packaging):**

```cmd
scripts\release-build.bat
```

That's it! 🎉
