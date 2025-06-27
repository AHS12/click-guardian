@echo off
echo Running Double-Click Fix in development mode...

REM Navigate to project root
cd /d "%~dp0.."
echo Running from: %CD%
echo.

go run .\cmd\doubleclick-fix\main.go
pause
