//go:build !windows

package hooks

import (
	"fmt"
	"time"
)

type unsupportedHook struct{}

func newPlatformHook() MouseHook {
	return &unsupportedHook{}
}

func (u *unsupportedHook) Start(delay time.Duration, logChan chan string) error {
	logChan <- "❌ Mouse hooking is only supported on Windows."
	return fmt.Errorf("mouse hooking not supported on this platform")
}

func (u *unsupportedHook) Stop() error {
	return nil
}

func (u *unsupportedHook) GetBlockedCount() int {
	return 0
}

func (u *unsupportedHook) ResetBlockedCount() {
	// No-op
}

func (u *unsupportedHook) IsSupported() bool {
	return false
}
