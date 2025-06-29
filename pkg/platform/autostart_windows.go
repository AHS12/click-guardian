//go:build windows

package platform

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

const (
	registryKey   = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	appName       = "Click Guardian"
	registryValue = "ClickGuardian"
)

// EnableAutoStart enables the application to start with Windows
func EnableAutoStart() error {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Add --minimized flag for auto-start (protection will auto-start if enabled in settings)
	autoStartCommand := fmt.Sprintf(`"%s" --minimized`, exePath)

	// Open the registry key for writing
	key, err := registry.OpenKey(registry.CURRENT_USER, registryKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer key.Close()

	// Set the registry value
	err = key.SetStringValue(registryValue, autoStartCommand)
	if err != nil {
		return fmt.Errorf("failed to set registry value: %v", err)
	}

	return nil
}

// DisableAutoStart disables the application from starting with Windows
func DisableAutoStart() error {
	// Open the registry key for writing
	key, err := registry.OpenKey(registry.CURRENT_USER, registryKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer key.Close()

	// Delete the registry value
	err = key.DeleteValue(registryValue)
	if err != nil {
		// If the value doesn't exist, that's fine
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("failed to delete registry value: %v", err)
	}

	return nil
}

// IsAutoStartEnabled checks if auto-start is currently enabled
func IsAutoStartEnabled() bool {
	// Open the registry key for reading
	key, err := registry.OpenKey(registry.CURRENT_USER, registryKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	// Try to read the registry value
	_, _, err = key.GetStringValue(registryValue)
	return err == nil
}
