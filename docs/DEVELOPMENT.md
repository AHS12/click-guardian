# Development Guide

## Project Structure

```
click-guardian/
├── assets/                      # Static assets
│   ├── icon-modern-shield.svg
│   ├── icon-modern-shield-solid.svg
│   ├── icon-modern-shield-simplified.svg
│   └── tray-icon.ico            # Source tray icon (ICO, for bundling)
├── build/
│   ├── build.conf               # Build configuration
│   ├── scripts/
│   │   ├── create-icon.ps1
│   │   └── sign-code.bat
│   ├── temp/
│   │   ├── app.syso
│   │   ├── click-guardian-v1.0.0-windows/
│   │   │   ├── LICENSE.txt
│   │   │   └── README.txt
│   │   └── click-guardian-v1.0.1-windows/
│   └── windows/
│       ├── app-icon.ico
│       ├── app-manifest.xml
│       └── app.rc
├── cmd/
│   └── click-guardian/          # Main application entry point
│       └── main.go
├── dist/                        # Build outputs (git ignored)
├── docs/                        # Documentation
│   ├── BUILD.md
│   ├── DEVELOPMENT.md
│   └── VSCODE_SETUP.md
├── internal/                    # Private application code
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── gui/                     # GUI application logic
│   │   ├── app.go
│   │   ├── icon.go
│   │   ├── resources.go         # All Fyne resources in one place (icon, SVG icon)
│   │   ├── icon_resource.go     # Auto-generated: main app icon (SVG)
│   │   ├── trayicon_resource.go # Auto-generated: tray icon (ICO)
│   │
│   ├── hooks/                   # Platform-specific mouse hooks
│   │   ├── hook.go
│   │   ├── hook_windows.go      # Windows implementation
│   │   └── hook_unsupported.go  # Fallback for other platforms
│   └── logger/                  # Logging functionality
│       └── logger.go
├── pkg/                         # Public packages (for future use)
│   └── platform/                # Platform detection utilities
│       ├── autostart_other.go
│       ├── autostart_windows.go
│       └── platform.go
├── scripts/                     # Build and development scripts
│   ├── build.bat
│   ├── build.sh
│   ├── dev.bat
│   ├── release-build.bat
│   ├── troubleshoot.bat
│   └── README.md
├── .vscode/                     # VSCode configuration
├── click-guardian.code-workspace # VSCode workspace file
├── go.mod
├── go.sum
└── README.md
```

---

## Installation & Building

### Prerequisites

- **Windows operating system** (primary supported platform)
- **Go 1.24.1 or later** ([download](https://go.dev/dl/))
- **Fyne 2.6.1 or later** ([docs](https://fyne.io/))
- **CGO enabled** (required for Windows API integration)

#### Recommended Development Environment: MSYS2

For the best experience building Fyne apps with CGO on Windows, use the [MSYS2](https://www.msys2.org/) environment:

1. **Install MSYS2**  
   Download and install from [msys2.org](https://www.msys2.org/).

2. **Open the correct terminal**  
   After installation, **do not use the default MSYS terminal**.  
   Instead, open **“MSYS2 MinGW 64-bit”** from the Start menu.

3. **Update MSYS2 and install required packages:**  
   Run the following commands in the MinGW 64-bit terminal:

   ```sh
   pacman -Syu
   pacman -S git mingw-w64-x86_64-toolchain mingw-w64-x86_64-go
   ```

4. **Add Go to your PATH:**  
   To ensure Go binaries are available, add this to your shell profile:

   ```sh
   echo "export PATH=\$PATH:~/Go/bin" >> ~/.bashrc
   ```

5. **(Optional) Add MSYS2 tools to Windows PATH:**  
   To use the compiler from other terminals or editors (like VSCode), add the following to your Windows system PATH:

   ```
   C:\msys64\mingw64\bin
   ```

   - Open “Edit the system environment variables” → Advanced → Environment Variables → Edit `Path`.

6. **Install Fyne dependencies:**  
   Follow the [Fyne Getting Started guide](https://docs.fyne.io/started/) for any additional setup.

---

**Tip:**  
If you encounter build issues, ensure you are using the **MinGW 64-bit** terminal and that both Go and the C toolchain are available in your PATH.

### Quick Start

1. **Clone the repository:**

   ```bash
   git clone https://github.com/AHS12/click-guardian
   cd click-guardian
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Build and run:**

   ```bash
   # Quick development build and run
   scripts\dev.bat

   #dev build
   scripts\build.bat

   # Production build
   scripts\release-build.bat
   ```

### Manual Build Options

you can read the `scripts\build.bat` to get all the build command

### Running the Application

After building, run the executable:

```bash
go run .\cmd\click-guardian\main.go
```

**For Development:**

```bash
scripts\dev.bat
```

### VSCode Setup

**Recommended:** Open the workspace file for optimal configuration:

```
File → Open Workspace from File → click-guardian.code-workspace
```

This workspace file includes:

- Go-specific settings optimized for CGO and Fyne
- Platform-specific build constraints
- Recommended extensions for Go and C++ development
- Troubleshooting configurations for cross-platform dependencies

**Troubleshooting VSCode Issues:**

- Run `scripts\troubleshoot.bat` for automated diagnosis
- See `docs/VSCODE_SETUP.md` for detailed solutions
- Reload window: `Ctrl+Shift+P` → "Developer: Reload Window"

### Cross-platform Build

```bash
# Make script executable (Linux/macOS)
chmod +x scripts/build.sh

# Build for all platforms
scripts/build.sh
```

## Architecture

### Core Components

1. **GUI Layer** (`internal/gui/`)

   - Fyne-based user interface
   - Application state management
   - User interaction handling

2. **Hook Layer** (`internal/hooks/`)

   - Platform-specific mouse hook implementations
   - Windows: Low-level mouse hook using Windows API
   - Other platforms: Placeholder implementation

3. **Configuration** (`internal/config/`)

   - Application settings
   - Validation logic
   - Default values

4. **Logging** (`internal/logger/`)
   - Real-time log display
   - Message formatting
   - Thread-safe logging

### Platform Support

Currently fully supported:

- Windows (x86, x64)

Planned support:

- Linux (X11/Wayland)
- macOS

## Adding New Platforms

1. Create platform-specific hook implementation in `internal/hooks/`
2. Use build tags to conditionally compile: `//go:build yourplatform`
3. Implement the `MouseHook` interface
4. Update build scripts to include the new platform

## Dependencies

- **Fyne v2**: GUI framework
- **Windows API**: For mouse hooking on Windows (CGO)

## Development Tips

- **Use `scripts/dev.bat`** for quick testing during development
- **Check `internal/hooks/hook_unsupported.go`** for reference implementation
- **All builds output to `dist/` directory** (git ignored)
- **Open the `.code-workspace` file** instead of the folder for best VSCode experience
- **Use the console version** for debugging output
- **Icon resources** are auto-generated in `internal/gui/resources.go` using `fyne bundle`

## Project Configuration

## Icon Management

The application uses SVG and ICO icons that are converted to Go resources using `fyne bundle`:

```bash
# Regenerate main app icon resource (SVG)
fyne bundle -pkg resources -o internal/gui/resources/icon_resource.go assets/icon-modern-shield-solid.svg

# Regenerate tray icon resource (ICO)
fyne bundle -pkg resources -o internal/gui/resources/trayicon_resource.go assets/tray-icon.ico
```

- **Note:** The generated files `icon_resource.go` and `trayicon_resource.go` are used directly by the application for icon resources.
- If you update the icon files, rerun the above commands to regenerate the resources.

### Build Scripts

- `scripts/build.bat` - Windows production build
- `scripts/dev.bat` - Development build and run
- `scripts/troubleshoot.bat` - VSCode/Go environment diagnosis
- `scripts/build.sh` - Cross-platform build script (future)
