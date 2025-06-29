//go:build !windows

package platform

import (
	"fmt"
)

// EnableAutoStart enables the application to start with the system (not implemented for this platform)
func EnableAutoStart() error {
	return fmt.Errorf("auto-start not supported on this platform")
}

// DisableAutoStart disables the application from starting with the system (not implemented for this platform)
func DisableAutoStart() error {
	return fmt.Errorf("auto-start not supported on this platform")
}

// IsAutoStartEnabled checks if auto-start is currently enabled (not implemented for this platform)
func IsAutoStartEnabled() bool {
	return false
}
