//go:build windows

package hooks

/*
#cgo LDFLAGS: -luser32 -lkernel32
#include <windows.h>

// Define constants that might not be available
#ifndef WM_LBUTTONUP
#define WM_LBUTTONUP 0x0202
#endif
#ifndef WM_RBUTTONUP
#define WM_RBUTTONUP 0x0205
#endif

extern LRESULT LowLevelMouseProc(int nCode, WPARAM wParam, LPARAM lParam);
*/
import "C"

import (
	"fmt"
	"time"
)

// Windows message constants
const (
	WM_LBUTTONUP = 0x0202
	WM_RBUTTONUP = 0x0205
	WM_MOUSEMOVE = 0x0200
)

type windowsHook struct {
	hook              C.HHOOK
	delay             time.Duration
	lastCompleteClick time.Time
	lastClickButton   C.WPARAM
	buttonPressed     map[C.WPARAM]bool
	buttonPressTime   map[C.WPARAM]time.Time
	dragDetected      map[C.WPARAM]bool
	logChannel        chan string
	blockedCount      int
	isRunning         bool
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

	now := time.Now()

	switch wParam {
	case C.WM_LBUTTONDOWN, C.WM_RBUTTONDOWN:
		buttonName := "Left"
		if wParam == C.WM_RBUTTONDOWN {
			buttonName = "Right"
		}

		// Check if button is already pressed (hardware bounce detection)
		if globalHook.buttonPressed[wParam] {
			// This is likely a hardware bounce/double-click issue
			timeSincePress := now.Sub(globalHook.buttonPressTime[wParam])
			if timeSincePress < globalHook.delay {
				globalHook.blockedCount++
				globalHook.sendLog(fmt.Sprintf("🚫 BLOCKED: %s button hardware bounce (%.0fms since press) - Total blocked: %d",
					buttonName, float64(timeSincePress.Nanoseconds())/1000000, globalHook.blockedCount))
				return 1 // Block the bounce
			}
		}

		// Check for rapid successive complete clicks (main double-click protection)
		if !globalHook.lastCompleteClick.IsZero() &&
			now.Sub(globalHook.lastCompleteClick) < globalHook.delay &&
			wParam == globalHook.lastClickButton {
			globalHook.blockedCount++
			globalHook.sendLog(fmt.Sprintf("🚫 BLOCKED: %s button rapid double-click (%.0fms after complete click) - Total blocked: %d",
				buttonName, float64(now.Sub(globalHook.lastCompleteClick).Nanoseconds())/1000000, globalHook.blockedCount))
			return 1 // Block the rapid click
		}

		// Allow the click and mark button as pressed
		globalHook.buttonPressed[wParam] = true
		globalHook.buttonPressTime[wParam] = now
		globalHook.dragDetected[wParam] = false
		globalHook.sendLog(fmt.Sprintf("✅ ALLOWED: %s button press", buttonName))

	case C.WPARAM(WM_LBUTTONUP), C.WPARAM(WM_RBUTTONUP):
		// Determine which button was released
		var downEvent C.WPARAM
		buttonName := "Left"
		if wParam == C.WPARAM(WM_LBUTTONUP) {
			downEvent = C.WM_LBUTTONDOWN
		} else {
			downEvent = C.WM_RBUTTONDOWN
			buttonName = "Right"
		}

		// Only process if we have a corresponding button press
		if globalHook.buttonPressed[downEvent] {
			holdDuration := now.Sub(globalHook.buttonPressTime[downEvent])
			globalHook.lastCompleteClick = now
			globalHook.lastClickButton = downEvent
			globalHook.buttonPressed[downEvent] = false

			// Log based on operation type
			if globalHook.dragDetected[downEvent] {
				globalHook.sendLog(fmt.Sprintf("✅ ALLOWED: %s button release after drag operation (%.0fms hold)",
					buttonName, float64(holdDuration.Nanoseconds())/1000000))
			} else if holdDuration > 200*time.Millisecond {
				globalHook.sendLog(fmt.Sprintf("✅ ALLOWED: %s button release after long hold (%.0fms)",
					buttonName, float64(holdDuration.Nanoseconds())/1000000))
			} else {
				globalHook.sendLog(fmt.Sprintf("✅ ALLOWED: %s button release - quick click (%.0fms)",
					buttonName, float64(holdDuration.Nanoseconds())/1000000))
			}

			// Reset drag detection
			globalHook.dragDetected[downEvent] = false
		}

	case C.WPARAM(WM_MOUSEMOVE):
		// Check if any button is currently pressed and mark as drag (but don't spam logs)
		if globalHook.buttonPressed[C.WM_LBUTTONDOWN] && !globalHook.dragDetected[C.WM_LBUTTONDOWN] {
			globalHook.dragDetected[C.WM_LBUTTONDOWN] = true
			globalHook.sendLog("🖱️  Left button drag operation detected")
		}
		if globalHook.buttonPressed[C.WM_RBUTTONDOWN] && !globalHook.dragDetected[C.WM_RBUTTONDOWN] {
			globalHook.dragDetected[C.WM_RBUTTONDOWN] = true
			globalHook.sendLog("🖱️  Right button drag operation detected")
		}
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
	w.buttonPressed = make(map[C.WPARAM]bool)
	w.buttonPressTime = make(map[C.WPARAM]time.Time)
	w.dragDetected = make(map[C.WPARAM]bool)
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
