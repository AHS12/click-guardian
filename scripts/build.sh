#!/bin/bash

set -e

echo "Building Click Guardian for multiple platforms..."

# Get script directory and navigate to project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."
echo "Building from: $(pwd)"
echo

# Create dist directory if it doesn't exist
mkdir -p dist

# Clean previous builds
rm -f dist/*

# Define build targets
declare -A targets=(
    ["windows/amd64"]="dist/click-guardian-windows-amd64.exe"
    ["windows/386"]="dist/click-guardian-windows-386.exe"
    ["linux/amd64"]="dist/click-guardian-linux-amd64"
    ["darwin/amd64"]="dist/click-guardian-darwin-amd64"
    ["darwin/arm64"]="dist/click-guardian-darwin-arm64"
)

# Build for each target
for target in "${!targets[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$target"
    output="${targets[$target]}"
    
    echo "Building for $GOOS/$GOARCH..."
    
    if [ "$GOOS" = "windows" ]; then
        # GUI version for Windows
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -H=windowsgui" -o "${output%.*}-gui.exe" ./cmd/click-guardian
        # Console version for Windows
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -o "$output" ./cmd/click-guardian
    else
        # Regular build for other platforms
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -o "$output" ./cmd/click-guardian
    fi
    
    if [ $? -eq 0 ]; then
        echo "✅ Build successful for $GOOS/$GOARCH"
    else
        echo "❌ Build failed for $GOOS/$GOARCH"
    fi
done

echo ""
echo "Cross-platform build complete! Executables are in the dist folder:"
ls -la dist/

echo ""
echo "Note: Only Windows builds will have full mouse hooking functionality."
echo "Other platforms will show an unsupported message."
