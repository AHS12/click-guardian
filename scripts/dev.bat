@echo off
echo Running Click Guardian in development mode...

REM Navigate to project root
cd /d "%~dp0.."
echo Running from: %CD%
echo.

go run .\cmd\click-guardian\main.go
pause
