//go:build windows

package main

/*
#cgo LDFLAGS: -luser32 -lkernel32
#include <windows.h>

extern LRESULT LowLevelMouseProc(int nCode, WPARAM wParam, LPARAM lParam);
*/
import "C"

import (
	"fmt"
	"time"
)

var (
	hook            C.HHOOK
	delay           time.Duration
	lastClickTime   time.Time
	lastClickButton C.WPARAM
	logChannel      chan string
	blockedCount    int
)

func sendLog(msg string) {
	select {
	case logChannel <- msg:
	default:
		// Log channel is full, message is dropped to prevent blocking.
	}
}

//export LowLevelMouseProc
func LowLevelMouseProc(nCode C.int, wParam C.WPARAM, lParam C.LPARAM) C.LRESULT {
	if nCode < 0 {
		return C.CallNextHookEx(hook, nCode, wParam, lParam)
	}

	if wParam == C.WM_LBUTTONDOWN || wParam == C.WM_RBUTTONDOWN {
		now := time.Now()
		buttonName := "Left"
		if wParam == C.WM_RBUTTONDOWN {
			buttonName = "Right"
		}

		if now.Sub(lastClickTime) < delay && wParam == lastClickButton {
			blockedCount++
			sendLog(fmt.Sprintf("🚫 BLOCKED: %s button double-click (%.0fms interval) - Total blocked: %d", buttonName, float64(now.Sub(lastClickTime).Nanoseconds())/1000000, blockedCount))
			return 1 // Block the click
		}
		lastClickTime = now
		lastClickButton = wParam
		sendLog(fmt.Sprintf("✅ ALLOWED: %s button click", buttonName))
	}

	return C.CallNextHookEx(hook, nCode, wParam, lParam)
}

func startHook(d time.Duration, logChan chan string) {
	delay = d
	logChannel = logChan
	go func() {
		hook = C.SetWindowsHookExW(C.WH_MOUSE_LL, C.HOOKPROC(C.LowLevelMouseProc), nil, 0)
		if hook == nil {
			logChannel <- "❌ Failed to install mouse hook"
			return
		}
		logChannel <- "🎯 Mouse hook installed successfully - protection active!"
		var msg C.MSG
		for C.GetMessage(&msg, nil, 0, 0) != 0 {
			C.TranslateMessage(&msg)
			C.DispatchMessage(&msg)
		}
	}()
}

func stopHook() {
	if hook != nil {
		C.UnhookWindowsHookEx(hook)
		hook = nil
		logChannel <- "🛑 Mouse hook removed - protection stopped"
	}
}

func getBlockedCount() int {
	return blockedCount
}

func resetBlockedCount() {
	blockedCount = 0
}
