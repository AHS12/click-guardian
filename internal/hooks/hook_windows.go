//go:build windows

package hooks

/*
#cgo LDFLAGS: -luser32 -lkernel32
#include <windows.h>


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

	// Faulty hardware detection
	faultyClickPattern map[C.WPARAM][]time.Time   // Track recent click times
	adaptiveDelay      map[C.WPARAM]time.Duration // Per-button adaptive delay
	shortClickCount    map[C.WPARAM]int           // Count of very short clicks (likely low pressure)

	logChannel   chan string
	blockedCount int
	isRunning    bool

	// Track last DOWN and UP event times
	lastDownTime    map[C.WPARAM]time.Time // Track last DOWN event time
	lastDownBlocked map[C.WPARAM]bool      // Track if last DOWN was blocked
	lastUpTime      map[C.WPARAM]time.Time // Track last UP event time
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

	// now := time.Now()

	switch wParam {
	case C.WM_LBUTTONDOWN, C.WM_RBUTTONDOWN:
		buttonName := "Left"
		if wParam == C.WM_RBUTTONDOWN {
			buttonName = "Right"
		}

		now := time.Now()
		lastDown := globalHook.lastDownTime[wParam]
		interval := now.Sub(lastDown)
		if !lastDown.IsZero() && interval < globalHook.getEffectiveDelay(wParam) {
			// Always block rapid successive DOWN events (hardware bounce or double-click)
			globalHook.lastDownBlocked[wParam] = true
			globalHook.lastDownTime[wParam] = now
			globalHook.blockedCount++
			globalHook.sendLog(fmt.Sprintf("🛑 STRICT BLOCK: %s hardware bounce/double-click (%.0fms after previous DOWN, delay: %.0fms) - Total blocked: %d",
				buttonName, float64(interval.Nanoseconds())/1000000, float64(globalHook.getEffectiveDelay(wParam).Nanoseconds())/1000000, globalHook.blockedCount))
			return 1
		}
		globalHook.lastDownTime[wParam] = now
		globalHook.lastDownBlocked[wParam] = false

		// Enhance strictness of double-click blocking logic
		if globalHook.buttonPressed[wParam] {
			// Strictly block any DOWN events within the delay
			timeSincePress := now.Sub(globalHook.buttonPressTime[wParam])
			effectiveDelay := globalHook.getEffectiveDelay(wParam)
			if timeSincePress < effectiveDelay {
				globalHook.blockedCount++
				globalHook.sendLog(fmt.Sprintf("🚫 STRICT BLOCK: %s button hardware bounce (%.0fms since press, delay: %.0fms) - Total blocked: %d",
					buttonName, float64(timeSincePress.Nanoseconds())/1000000,
					float64(effectiveDelay.Nanoseconds())/1000000, globalHook.blockedCount))
				return 1 // Block the bounce
			}
		}

		// Strictly block rapid successive complete clicks
		if !globalHook.lastCompleteClick.IsZero() &&
			now.Sub(globalHook.lastCompleteClick) < globalHook.getEffectiveDelay(wParam) &&
			wParam == globalHook.lastClickButton {
			effectiveDelay := globalHook.getEffectiveDelay(wParam)
			globalHook.blockedCount++
			globalHook.sendLog(fmt.Sprintf("🚫 STRICT BLOCK: %s button rapid double-click (%.0fms after complete click, delay: %.0fms) - Total blocked: %d",
				buttonName, float64(now.Sub(globalHook.lastCompleteClick).Nanoseconds())/1000000,
				float64(effectiveDelay.Nanoseconds())/1000000, globalHook.blockedCount))
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

		now := time.Now()
		lastDown := globalHook.buttonPressTime[downEvent]
		lastUp := globalHook.lastUpTime[wParam]
		upInterval := now.Sub(lastUp)
		holdDuration := now.Sub(lastDown)

		// Block UP events that occur too quickly after a DOWN event
		if !globalHook.buttonPressed[downEvent] || (!lastUp.IsZero() && upInterval < globalHook.getEffectiveDelay(downEvent)) {
			globalHook.sendLog(fmt.Sprintf("🛑 STRICT BLOCK: %s spurious UP (%.0fms after previous UP)", buttonName, float64(upInterval.Nanoseconds())/1000000))
			return 1
		}
		globalHook.lastUpTime[wParam] = now

		// Only process if we have a corresponding button press
		if globalHook.buttonPressed[downEvent] {
			globalHook.lastCompleteClick = now
			globalHook.lastClickButton = downEvent
			globalHook.buttonPressed[downEvent] = false

			// Analyze click pattern for faulty hardware detection
			globalHook.detectFaultyHardware(downEvent, holdDuration)

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
	w.faultyClickPattern = make(map[C.WPARAM][]time.Time)
	w.adaptiveDelay = make(map[C.WPARAM]time.Duration)
	w.shortClickCount = make(map[C.WPARAM]int)
	w.lastDownTime = make(map[C.WPARAM]time.Time)
	w.lastDownBlocked = make(map[C.WPARAM]bool)
	w.lastUpTime = make(map[C.WPARAM]time.Time)
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

// detectFaultyHardware analyzes click patterns to detect faulty mouse behavior
func (w *windowsHook) detectFaultyHardware(button C.WPARAM, holdDuration time.Duration) {
	// Track recent click times (keep last 10 clicks)
	if w.faultyClickPattern[button] == nil {
		w.faultyClickPattern[button] = make([]time.Time, 0, 10)
	}

	now := time.Now()
	w.faultyClickPattern[button] = append(w.faultyClickPattern[button], now)
	if len(w.faultyClickPattern[button]) > 10 {
		w.faultyClickPattern[button] = w.faultyClickPattern[button][1:]
	}

	// Count very short clicks (likely insufficient pressure)
	if holdDuration < 30*time.Millisecond {
		w.shortClickCount[button]++
		w.sendLog(fmt.Sprintf("📊 Short click detected (%.0fms) - total short clicks: %d",
			float64(holdDuration.Nanoseconds())/1000000, w.shortClickCount[button]))
	}

	// Analyze pattern every 5 clicks
	if len(w.faultyClickPattern[button]) >= 5 {
		shortClicks := 0
		for i := len(w.faultyClickPattern[button]) - 5; i < len(w.faultyClickPattern[button]); i++ {
			// Check if this was a short click by looking at our count
			if w.shortClickCount[button] > 0 {
				shortClicks++
			}
		}

		// If more than 60% of recent clicks are short, reduce protection aggressively
		if float64(shortClicks)/5.0 > 0.6 {
			newDelay := w.delay / 2 // Reduce to half of original delay
			if newDelay < 20*time.Millisecond {
				newDelay = 20 * time.Millisecond // Minimum 20ms
			}

			buttonName := "Left"
			if button == C.WM_RBUTTONDOWN {
				buttonName = "Right"
			}

			if w.adaptiveDelay[button] != newDelay {
				w.adaptiveDelay[button] = newDelay
				w.sendLog(fmt.Sprintf("🔧 ADAPTIVE STRICT: %s button delay reduced to %.0fms due to detected low-pressure pattern",
					buttonName, float64(newDelay.Nanoseconds())/1000000))
			}
		} else {
			// Reset to normal delay if pattern improves
			if w.adaptiveDelay[button] != w.delay {
				w.adaptiveDelay[button] = w.delay
				buttonName := "Left"
				if button == C.WM_RBUTTONDOWN {
					buttonName = "Right"
				}
				w.sendLog(fmt.Sprintf("🔧 ADAPTIVE STRICT: %s button delay reset to normal (%.0fms)",
					buttonName, float64(w.delay.Nanoseconds())/1000000))
			}
		}
	}
}

// getEffectiveDelay returns the adaptive delay for a specific button
func (w *windowsHook) getEffectiveDelay(button C.WPARAM) time.Duration {
	if adaptiveDelay, exists := w.adaptiveDelay[button]; exists && adaptiveDelay > 0 {
		return adaptiveDelay
	}
	return w.delay
}
