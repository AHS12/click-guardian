@echo off
setlocal EnableDelayedExpansion

echo =====================================
echo  Click Guardian - Release Builder
echo =====================================
echo.

REM Navigate to project root
cd /d "%~dp0\.."
echo Building from: %CD%
echo.

REM Load version from config or use default
set VERSION=1.0.0
if exist "build\build.conf" (
    for /f "usebackq tokens=1,2 delims==" %%A in ("build\build.conf") do (
        if "%%A"=="VERSION" set VERSION=%%B
    )
)

REM Get build info
for /f %%A in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%A
if not defined GIT_COMMIT set GIT_COMMIT=unknown

for /f %%A in ('git config user.name 2^>nul') do set BUILD_BY=%%A
if not defined BUILD_BY set BUILD_BY=%USERNAME%

REM Get current date/time for build timestamp
for /f "tokens=2-4 delims=/ " %%A in ('date /t') do set BUILD_DATE=%%C-%%A-%%B
for /f "tokens=1-2 delims=: " %%A in ('time /t') do set BUILD_TIME=%%A:%%B
set BUILD_TIMESTAMP=!BUILD_DATE! !BUILD_TIME!

echo Building Click Guardian v%VERSION%
echo Git Commit: %GIT_COMMIT%
echo Build Time: %BUILD_TIMESTAMP%
echo Built By: %BUILD_BY%
echo.

REM Create dist directory if it doesn't exist
if not exist "dist" mkdir dist

REM Clean previous builds
del /Q dist\*.exe 2>nul
del /Q dist\*.zip 2>nul

REM --- Update version in app.rc automatically ---
set VERSION_RC=%VERSION:.=,%,0
powershell -Command "(Get-Content build\windows\app.rc) -replace 'FILEVERSION [0-9,]+', 'FILEVERSION %VERSION_RC%' | Set-Content build\windows\app.rc"
powershell -Command "(Get-Content build\windows\app.rc) -replace 'PRODUCTVERSION [0-9,]+', 'PRODUCTVERSION %VERSION_RC%' | Set-Content build\windows\app.rc"
powershell -Command "(Get-Content build\windows\app.rc) -replace 'VALUE \"FileVersion\", \".*\"', 'VALUE \"FileVersion\", \"%VERSION%\"' | Set-Content build\windows\app.rc"
powershell -Command "(Get-Content build\windows\app.rc) -replace 'VALUE \"ProductVersion\", \".*\"', 'VALUE \"ProductVersion\", \"%VERSION%\"' | Set-Content build\windows\app.rc"

REM --- Update version in app-manifest.xml automatically ---
powershell -NoProfile -ExecutionPolicy Bypass -Command "[xml]$xml = Get-Content 'build\windows\app-manifest.xml'; $xml.assembly.assemblyIdentity.version = '%VERSION%.0'; $xml.Save((Resolve-Path 'build\windows\app-manifest.xml').Path)"


REM Generate Windows resource file (icon, manifest, version info)
echo Generating Windows resource file...
windres build\windows\app.rc -O coff -o cmd\click-guardian\click-guardian.syso

echo Building release version...
go build -ldflags "-s -w -H=windowsgui %VERSION_FLAGS%" -o dist\click-guardian.exe .\cmd\click-guardian

REM Set version flags (keep it simple)
set "VERSION_FLAGS=-X click-guardian/internal/version.Version=%VERSION% -X click-guardian/internal/version.GitCommit=%GIT_COMMIT% -X click-guardian/internal/version.BuildBy=%BUILD_BY%"

REM Build GUI version (no console window) - now just click-guardian.exe

echo Building release version...
go build -ldflags "-s -w -H=windowsgui %VERSION_FLAGS%" -o dist\click-guardian.exe .\cmd\click-guardian

if %ERRORLEVEL% EQU 0 (
    echo ✅ Build successful! Created dist\click-guardian.exe
) else (
    echo ❌ Build failed!
    goto end
)

REM Test the build
echo.
echo Testing build...
if exist "dist\click-guardian.exe" (
    echo Release version build completed successfully
    echo.
)

REM Create release package
echo Creating release package...
set RELEASE_DIR=build\temp\click-guardian-v%VERSION%-windows
if exist "%RELEASE_DIR%" rmdir /s /q "%RELEASE_DIR%"
if not exist "build\temp" mkdir "build\temp"
mkdir "%RELEASE_DIR%"

REM Copy files to release directory
copy "dist\click-guardian.exe" "%RELEASE_DIR%\" >nul

for /f "tokens=2 delims==" %%Y in ('"wmic os get localdatetime /value | findstr ="') do set CURYEAR=%%Y
set CURYEAR=!CURYEAR:~0,4!

REM Create README for release
(
echo Click Guardian v%VERSION%
echo.
echo Prevents accidental double-clicks with configurable delay protection
echo.
echo © %CURYEAR% Azizul Hakim
echo.
echo Build Information:
echo - Version: %VERSION%
echo - Built: %BUILD_TIMESTAMP%
echo - Commit: %GIT_COMMIT%
echo - Built by: %BUILD_BY%
echo.
echo Files:
echo - click-guardian.exe  - Main application
echo.
echo Installation:
echo 1. Run click-guardian.exe
echo 2. Configure your preferred delay
echo 3. Click "Start Protection"
echo.
echo For auto-start with Windows:
echo - Check "Start with Windows and auto-enable protection"
) > "%RELEASE_DIR%\README.txt"

REM Copy license if it exists
if exist "LICENSE" copy "LICENSE" "%RELEASE_DIR%\" >nul
if exist "LICENSE.txt" copy "LICENSE.txt" "%RELEASE_DIR%\" >nul
if exist "LICENSE.md" copy "LICENSE.md" "%RELEASE_DIR%\" >nul

REM Create ZIP using PowerShell
echo Creating ZIP package...
powershell -Command "Compress-Archive -Path '%RELEASE_DIR%\*' -DestinationPath 'dist\click-guardian-v%VERSION%-windows-portable.zip' -Force"

if exist "dist\click-guardian-v%VERSION%-windows-portable.zip" (
    echo ✅ Release package created: dist\click-guardian-v%VERSION%-windows-portable.zip
) else (
    echo ❌ Failed to create release package
)

:end
echo.
echo =====================================
echo   🎉 RELEASE BUILD COMPLETE!
echo =====================================
echo.
echo Files created:
dir dist\*.exe dist\*.zip 2>nul
echo.
echo Version: %VERSION%
echo Ready for distribution! 🚀
echo.
pause
