package version

import (
	"fmt"
	"runtime"
)

// Build information - set via ldflags during build
var (
	Version   = "dev"     // Version number
	GitCommit = "unknown" // Git commit hash
	BuildTime = "unknown" // Build timestamp
	BuildBy   = "Azizul Hakim" // Who built it
)

// AppInfo contains application metadata
type AppInfo struct {
	Name        string
	Version     string
	GitCommit   string
	BuildTime   string
	BuildBy     string
	GoVersion   string
	Platform    string
	Arch        string
	Description string
	Copyright   string
	Company     string
}

// GetAppInfo returns comprehensive application information
func GetAppInfo() AppInfo {
	return AppInfo{
		Name:        "Click Guardian",
		Version:     Version,
		GitCommit:   GitCommit,
		BuildTime:   BuildTime,
		BuildBy:     BuildBy,
		GoVersion:   runtime.Version(),
		Platform:    runtime.GOOS,
		Arch:        runtime.GOARCH,
		Description: "Prevents accidental double-clicks with configurable delay protection",
		Copyright:   "© 2025 Click Guardian Project",
		Company:     "Click Guardian Project",
	}
}

// GetVersionString returns a formatted version string
func GetVersionString() string {
	if Version == "dev" {
		return fmt.Sprintf("%s (development build)", Version)
	}
	return Version
}

// GetFullVersionString returns a complete version information string
func GetFullVersionString() string {
	info := GetAppInfo()
	return fmt.Sprintf("%s v%s (built %s, commit %s, %s %s/%s)",
		info.Name,
		info.Version,
		info.BuildTime,
		info.GitCommit[:min(len(info.GitCommit), 8)],
		info.GoVersion,
		info.Platform,
		info.Arch,
	)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
