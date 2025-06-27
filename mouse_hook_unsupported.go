//go:build !windows

package main

import (
	"time"
)

func startHook(d time.Duration, logChan chan string) {
	logChan <- "Mouse hooking is only supported on Windows."
}

func stopHook() {}
