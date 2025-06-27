//go:build windows

package hooks

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

type windowsHook struct {
	hook            C.HHOOK
	delay           time.Duration
	lastClickTime   time.Time
	lastClickButton C.WPARAM
	logChannel      chan string
	blockedCount    int
	isRunning       bool
}

func newPlatformHook() MouseHook {
	return &windowsHook{}
}

var globalHook *windowsHook

func (w *windowsHook) sendLog(msg string) {
	select {
	case w.logChannel <- msg:
	default:
		// Log channel is full, message is dropped to prevent blocking.
	}
}

//export LowLevelMouseProc
func LowLevelMouseProc(nCode C.int, wParam C.WPARAM, lParam C.LPARAM) C.LRESULT {
	if nCode < 0 || globalHook == nil {
		return C.CallNextHookEx(globalHook.hook, nCode, wParam, lParam)
	}

	if wParam == C.WM_LBUTTONDOWN || wParam == C.WM_RBUTTONDOWN {
		now := time.Now()
		buttonName := "Left"
		if wParam == C.WM_RBUTTONDOWN {
			buttonName = "Right"
		}

		if now.Sub(globalHook.lastClickTime) < globalHook.delay && wParam == globalHook.lastClickButton {
			globalHook.blockedCount++
			globalHook.sendLog(fmt.Sprintf("🚫 BLOCKED: %s button double-click (%.0fms interval) - Total blocked: %d",
				buttonName, float64(now.Sub(globalHook.lastClickTime).Nanoseconds())/1000000, globalHook.blockedCount))
			return 1 // Block the click
		}
		globalHook.lastClickTime = now
		globalHook.lastClickButton = wParam
		globalHook.sendLog(fmt.Sprintf("✅ ALLOWED: %s button click", buttonName))
	}

	return C.CallNextHookEx(globalHook.hook, nCode, wParam, lParam)
}

func (w *windowsHook) Start(delay time.Duration, logChan chan string) error {
	if w.isRunning {
		return fmt.Errorf("hook is already running")
	}

	w.delay = delay
	w.logChannel = logChan
	w.isRunning = true
	globalHook = w

	go func() {
		w.hook = C.SetWindowsHookExW(C.WH_MOUSE_LL, C.HOOKPROC(C.LowLevelMouseProc), nil, 0)
		if w.hook == nil {
			w.logChannel <- "❌ Failed to install mouse hook"
			w.isRunning = false
			return
		}
		w.logChannel <- "🎯 Mouse hook installed successfully - protection active!"
		var msg C.MSG
		for w.isRunning && C.GetMessage(&msg, nil, 0, 0) != 0 {
			C.TranslateMessage(&msg)
			C.DispatchMessage(&msg)
		}
	}()

	return nil
}

func (w *windowsHook) Stop() error {
	if !w.isRunning {
		return nil
	}

	w.isRunning = false
	if w.hook != nil {
		C.UnhookWindowsHookEx(w.hook)
		w.hook = nil
		if w.logChannel != nil {
			w.logChannel <- "🛑 Mouse hook removed - protection stopped"
		}
	}
	globalHook = nil
	return nil
}

func (w *windowsHook) GetBlockedCount() int {
	return w.blockedCount
}

func (w *windowsHook) ResetBlockedCount() {
	w.blockedCount = 0
}

func (w *windowsHook) IsSupported() bool {
	return true
}
