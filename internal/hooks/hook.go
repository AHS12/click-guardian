//go:build windows || !windows

package hooks

import (
	"time"
)

// MouseHook defines the interface for mouse hooking functionality
type MouseHook interface {
	Start(delay time.Duration, logChan chan string) error
	Stop() error
	GetBlockedCount() int
	ResetBlockedCount()
	IsSupported() bool
}

// NewMouseHook creates a new mouse hook for the current platform
func NewMouseHook() MouseHook {
	return newPlatformHook()
}
