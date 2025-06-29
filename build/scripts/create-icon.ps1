# PowerShell script to convert SVG to ICO with multiple resolutions
param(
    [string]$SvgPath = "assets\icon-modern-shield-solid.svg",
    [string]$OutputPath = "build\windows\app-icon.ico"
)

Write-Host "Converting SVG to multi-resolution ICO file..." -ForegroundColor Green
Write-Host "Input: $SvgPath" -ForegroundColor Yellow
Write-Host "Output: $OutputPath" -ForegroundColor Yellow

# Check if SVG file exists
if (-not (Test-Path $SvgPath)) {
    Write-Error "SVG file not found: $SvgPath"
    exit 1
}

# Create output directory if it doesn't exist
$outputDir = Split-Path $OutputPath -Parent
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
}

# Try to use ImageMagick (if available)
try {
    $magickCmd = Get-Command "magick" -ErrorAction Stop
    Write-Host "Using ImageMagick for conversion..." -ForegroundColor Blue
    
    # Convert SVG to ICO with multiple resolutions
    & magick $SvgPath -background transparent -define icon:auto-resize=256, 128, 64, 48, 32, 16 $OutputPath
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ ICO file created successfully!" -ForegroundColor Green
        Write-Host "File size: $((Get-Item $OutputPath).Length) bytes" -ForegroundColor Cyan
        exit 0
    }
    else {
        throw "ImageMagick conversion failed"
    }
}
catch {
    Write-Warning "ImageMagick not found or failed. Trying alternative methods..."
}

# Fallback: Try using Inkscape (if available)
try {
    $inkscapeCmd = Get-Command "inkscape" -ErrorAction Stop
    Write-Host "Using Inkscape for conversion..." -ForegroundColor Blue
    
    # Create temporary PNG files for different sizes
    $tempDir = "build\temp"
    if (-not (Test-Path $tempDir)) {
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    }
    
    $sizes = @(16, 32, 48, 64, 128, 256)
    $pngFiles = @()
    
    foreach ($size in $sizes) {
        $pngFile = "$tempDir\icon-$size.png"
        & inkscape $SvgPath --export-type=png --export-filename=$pngFile --export-width=$size --export-height=$size
        if ($LASTEXITCODE -eq 0) {
            $pngFiles += $pngFile
        }
    }
    
    if ($pngFiles.Count -gt 0) {
        # Use ImageMagick to combine PNGs into ICO
        $pngArgs = $pngFiles -join " "
        Invoke-Expression "magick $pngArgs $OutputPath"
        
        # Clean up temporary files
        Remove-Item "$tempDir\*.png" -Force
        
        if (Test-Path $OutputPath) {
            Write-Host "✅ ICO file created successfully using Inkscape!" -ForegroundColor Green
            exit 0
        }
    }
}
catch {
    Write-Warning "Inkscape not found or failed."
}

# Final fallback: Copy existing tray icon if available
$existingIcon = "internal\gui\tray-icon.ico"
if (Test-Path $existingIcon) {
    Write-Host "Using existing tray icon as fallback..." -ForegroundColor Yellow
    Copy-Item $existingIcon $OutputPath -Force
    Write-Host "✅ ICO file copied from existing tray icon!" -ForegroundColor Green
    exit 0
}

# If all else fails, create a simple notice
Write-Host ""
Write-Host "⚠️  Manual ICO creation required!" -ForegroundColor Red
Write-Host ""
Write-Host "Neither ImageMagick nor Inkscape was found for automatic ICO creation."
Write-Host "Please manually create an ICO file with multiple resolutions:"
Write-Host ""
Write-Host "1. Use an online converter like https://convertio.co/svg-ico/"
Write-Host "2. Convert: $SvgPath"
Write-Host "3. Save as: $OutputPath"
Write-Host "4. Include resolutions: 16x16, 32x32, 48x48, 64x64, 128x128, 256x256"
Write-Host ""
Write-Host "Alternatively, install ImageMagick or Inkscape and run this script again."
Write-Host ""

exit 1
