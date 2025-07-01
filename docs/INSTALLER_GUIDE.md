# 🛠️ Click Guardian — MSI Installer Guide

This guide shows how to create a **Windows MSI installer** for your `click-guardian.exe` using `go-msi`. This includes setting shortcuts, upgrade logic, and system integration.

---

## 📁 Project Structure

```
click-guardian-installer/
├── dist/
│   └── click-guardian.exe         # Compiled executable
├── assets/
│   └── icon.ico                   # Optional icon for shortcuts
├── wix.json                       # MSI configuration
```

---

## 🔧 Prerequisites

- ✅ Go installed and in `PATH`
- ✅ Windows machine
- ✅ **WiX Toolset** installed:
  ```bash
  # Option 1: Download from official site
  # https://wixtoolset.org/releases/
  
  # Option 2: Using Chocolatey
  choco install wixtoolset
  
  # Option 3: Using winget
  winget install WiXToolset.WiXToolset
  ```
  > ⚠️ **Important**: After installation, restart your terminal/command prompt to ensure WiX tools (`candle.exe` and `light.exe`) are in your PATH.

- ✅ `go-msi` installed:
  ```bash
  go install github.com/mh-cbon/go-msi@latest
  ```

---

## 📄 Step 1: Create `wix.json`

Create `wix.json` in the root of your installer folder:

```json
{
  "product": "Click Guardian",
  "company": "Click Guardian Project",
  "version": "1.0.1",
  "license": "GPL-3.0",
  "upgrade-code": "d900466f-5e04-43c5-b7f4-cbeb7ce26bce",
  "product-code": "d900466f-5e04-43c5-b7f4-cbeb7ce26rce",
  "base": "dist",
  "items": ["click-guardian.exe"],
  "files": {
    "guid": "ea91758f-b81c-4384-8bbc-ddb21be99e96",
    "items": ["click-guardian.exe"]
  },
  "env": {
    "guid": "9779be66-6aea-41d9-a1e5-6657cdb98ea0",
    "vars": []
  },
  "shortcuts": {
    "guid": "8a69c043-2d46-4604-954c-648f667b2233",
    "items": [
      {
        "name": "Click Guardian",
        "description": "Click Guardian - Double-Click Protection",
        "target": "[INSTALLDIR]\\click-guardian.exe",
        "wdir": "DesktopFolder",
        "arguments": "",
        "icon": "../assets/icon.ico"
      },
      {
        "name": "Click Guardian",
        "description": "Click Guardian - Double-Click Protection",
        "target": "[INSTALLDIR]\\click-guardian.exe",
        "wdir": "ProgramMenuFolder",
        "arguments": "",
        "icon": "../assets/icon.ico"
      }
    ]
  },
  "choco": {}
}
```

> 🔐 Generate GUIDs with:

```bash
go-msi set-guid
```

or use online tool like this: [https://www.guidgenerator.com/](https://www.guidgenerator.com/)

---

## 🏢 Step 2: Build the Installer

### Option 1: Automated Build (Recommended)

Use the automated release build script:

```bash
scripts\release-build.bat
```

This will:
- ⚙️ Build the executable
- 📎 Create portable ZIP package
- 🔄 Copy icon from `assets/` folder
- 📦 Generate MSI installer with shortcuts
- 🧽 Clean up temporary files
- 📁 List all created distribution files

### Option 2: Manual Build

From the project root, run:

```bash
# Copy icon from assets (required for WiX)
copy assets\icon.ico icon.ico

# Build MSI installer
go-msi make --msi dist/click-guardian-installer.msi --version 1.0.1 --src templates

# Clean up
del icon.ico
```

### What gets created:

- 💻 **click-guardian.exe** - Standalone executable
- 📎 **click-guardian-v1.0.1-windows-portable.zip** - Portable package
- 📦 **click-guardian-installer.msi** - Windows installer with:
  - Desktop + Start Menu shortcuts
  - Add/Remove Programs entry
  - Start Menu search visibility
  - Custom icon from assets folder

---

## 🔁 Future Updates

To update the app:

- 🔁 Keep the same `upgrade-code`
- 🆕 Change the `product-code`
- ⬆️ Bump `version`
- ✅ Rebuild `.msi`

Windows will automatically remove the previous version and install the new one.

---

## 🐛 Troubleshooting

### WiX Tools Not Found
```
CreateFile C:\Users\...\go\bin\templates: The system cannot find the file specified.
```
**Solution**: Install WiX Toolset and restart your terminal.

### Icon File Not Found
```
LGHT0103: The system cannot find the file 'icon.ico'.
```
**Solution**: Ensure `assets/icon.ico` exists and the build script copies it correctly.

### Relative Path Errors
```
Rel: can't make K:\...\file.exe relative to C:\...\temp
```
**Solution**: This usually occurs with cross-drive paths. The automated build script handles this by using proper temp directory management.

### MSI Build Fails
**Check**:
1. ✅ WiX Toolset installed (`candle.exe` and `light.exe` in PATH)
2. ✅ `go-msi` installed (`go-msi --version` works)
3. ✅ `templates/main.wxs` exists
4. ✅ `wix.json` syntax is valid
5. ✅ Icon file exists in assets folder

---

## ✅ Bonus Tips

- Add autostart, custom install paths, registry entries, or tray options by extending `wix.json`.
- For more customization, explore [`WiX Toolset`](https://wixtoolset.org/).
- Use `scripts\release-build.bat` for one-command building of all distribution formats.

---

Made with 💡 by Azizul (and Copilot 😄)
