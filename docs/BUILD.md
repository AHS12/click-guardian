# Double Click Fix - Cross-Platform Build Instructions

## Quick Start

### Windows

```cmd
scripts\build.bat
```

### Linux/macOS

```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

## Manual Build Commands

### Single Platform Build

```bash
# Windows GUI version
go build -ldflags "-s -w -H=windowsgui" -o dist/go-double-click-fix-gui.exe ./cmd/doubleclick-fix

# Windows Console version
go build -ldflags "-s -w" -o dist/go-double-click-fix.exe ./cmd/doubleclick-fix

# Linux/macOS
go build -ldflags "-s -w" -o dist/go-double-click-fix ./cmd/doubleclick-fix
```

### Cross-Platform Build

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/go-double-click-fix-windows-amd64.exe ./cmd/doubleclick-fix

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/go-double-click-fix-linux-amd64 ./cmd/doubleclick-fix

# macOS
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/go-double-click-fix-darwin-amd64 ./cmd/doubleclick-fix
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o dist/go-double-click-fix-darwin-arm64 ./cmd/doubleclick-fix
```

## Development

### Run in Development Mode

```bash
# Windows
scripts\dev.bat

# Linux/macOS
go run ./cmd/doubleclick-fix
```

### Clean Build Artifacts

```bash
# Windows
rmdir /s /q dist

# Linux/macOS
rm -rf dist
```

## Dependencies

### Prerequisites

- Go 1.24.1 or later
- C compiler (for CGO on Windows)

### Install Dependencies

```bash
go mod tidy
```
