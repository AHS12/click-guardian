//go:build windows

// This is the pure Go implementation of the Windows mouse hook.
// The old CGO implementation has been removed - this is now the only implementation.

package hooks

import (
	"fmt"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procSetWindowsHookEx         = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx           = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx      = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage               = user32.NewProc("GetMessageW")
	procTranslateMessage         = user32.NewProc("TranslateMessage")
	procDispatchMessage          = user32.NewProc("DispatchMessageW")
	procGetSystemMetrics         = user32.NewProc("GetSystemMetrics")
	procMouseEvent               = user32.NewProc("mouse_event")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess              = kernel32.NewProc("OpenProcess")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
)

const (
	PROCESS_QUERY_INFORMATION = 0x0400
	WH_MOUSE_LL               = 14
	WM_LBUTTONDOWN            = 0x0201
	WM_LBUTTONUP              = 0x0202
	WM_RBUTTONDOWN            = 0x0204
	WM_RBUTTONUP              = 0x0205
	WM_MOUSEMOVE              = 0x0200
	WM_MBUTTONDOWN            = 0x0207
	WM_MBUTTONUP              = 0x0208
	WM_XBUTTONDOWN            = 0x020B
	WM_XBUTTONUP              = 0x020C

	MOUSEEVENTF_LEFTDOWN   = 0x0002
	MOUSEEVENTF_LEFTUP     = 0x0004
	MOUSEEVENTF_RIGHTDOWN  = 0x0008
	MOUSEEVENTF_RIGHTUP    = 0x0010
	MOUSEEVENTF_MIDDLEDOWN = 0x0020
	MOUSEEVENTF_MIDDLEUP   = 0x0040
	MOUSEEVENTF_XDOWN      = 0x0080
	MOUSEEVENTF_XUP        = 0x0100
)

type POINT struct {
	X, Y int32
}

type MSLLHOOKSTRUCT struct {
	Pt          POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type WindowsHook struct {
	hook              syscall.Handle
	delay             time.Duration
	lastCompleteClick time.Time
	lastClickButton   uintptr
	buttonPressed     map[uintptr]bool
	buttonPressTime   map[uintptr]time.Time
	dragDetected      map[uintptr]bool

	// Faulty hardware detection
	faultyClickPattern map[uintptr][]time.Time   // Track recent click times
	adaptiveDelay      map[uintptr]time.Duration // Per-button adaptive delay
	shortClickCount    map[uintptr]int           // Count of very short clicks (likely low pressure)

	logChannel   chan string
	blockedCount int
	isRunning    bool

	// Track last DOWN and UP event times
	lastDownTime    map[uintptr]time.Time // Track last DOWN event time
	lastDownBlocked map[uintptr]bool      // Track if last DOWN was blocked
	lastUpTime      map[uintptr]time.Time // Track last UP event time

	// Protected buttons configuration
	protectedButtons map[string]bool

	// Drag Fix
	dragFixEnabled   bool
	dragFixThreshold time.Duration
	pendingUpTimes   map[uintptr]time.Time
	mu               sync.Mutex

	// Pause hook logic
	pauseUntil    time.Time
	pauseDuration time.Duration

	// Cache for process permission checks
	permCache   map[uint32]bool      // PID -> canInject
	permCacheTs map[uint32]time.Time // PID -> timestamp
}

func newPlatformHook() MouseHook {
	return &WindowsHook{
		permCache:     make(map[uint32]bool),
		permCacheTs:   make(map[uint32]time.Time),
		pauseDuration: 3 * time.Second, // Default pause duration
	}
}

var globalHook *WindowsHook

func (w *WindowsHook) sendLog(msg string) {
	select {
	case w.logChannel <- msg:
	default:
		// Log channel is full, message is dropped to prevent blocking.
	}
}

func (w *WindowsHook) canInjectToForeground() bool {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return true // No foreground window, assume safe
	}

	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return true
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// Check cache (valid for 5 seconds)
	if ts, ok := w.permCacheTs[pid]; ok && time.Since(ts) < 5*time.Second {
		return w.permCache[pid]
	}

	// Heuristic: Try to open the process with QUERY_INFORMATION.
	// If we are Standard User and target is Admin, this usually fails with Access Denied.
	hProcess, _, _ := procOpenProcess.Call(
		uintptr(PROCESS_QUERY_INFORMATION),
		0,
		uintptr(pid),
	)

	canInject := true
	if hProcess == 0 {
		// Failed to open process. Assume high privilege or protected.
		canInject = false
	} else {
		procCloseHandle.Call(hProcess)
	}

	// Update cache
	w.permCache[pid] = canInject
	w.permCacheTs[pid] = time.Now()

	return canInject
}

// LowLevelMouseProc is the mouse hook callback function
func LowLevelMouseProc(nCode int32, wParam uintptr, lParam unsafe.Pointer) uintptr {
	if nCode < 0 || globalHook == nil {
		ret, _, _ := procCallNextHookEx.Call(
			uintptr(globalHook.hook),
			uintptr(nCode),
			wParam,
			uintptr(lParam),
		)
		return ret
	}

	// Check if hook is temporarily paused
	if !globalHook.pauseUntil.IsZero() {
		if time.Now().Before(globalHook.pauseUntil) {
			// Allow everything through during pause
			ret, _, _ := procCallNextHookEx.Call(
				uintptr(globalHook.hook),
				uintptr(nCode),
				wParam,
				uintptr(lParam),
			)
			return ret
		} else {
			// Pause expired
			globalHook.mu.Lock()
			globalHook.pauseUntil = time.Time{}
			globalHook.mu.Unlock()
			globalHook.sendLog("▶️ Protection resumed after temporary pause")
		}
	}

	// Check for injected events (LLMHF_INJECTED is bit 1 of Flags)
	mouseStruct := (*MSLLHOOKSTRUCT)(lParam)
	isInjected := (mouseStruct.Flags & 0x01) != 0

	switch wParam {
	case WM_LBUTTONDOWN, WM_RBUTTONDOWN, WM_MBUTTONDOWN, WM_XBUTTONDOWN:
		buttonName := "Left"
		isProtected := globalHook.protectedButtons["left"]

		switch wParam {
		case WM_LBUTTONDOWN:
			buttonName = "Left"
			isProtected = globalHook.protectedButtons["left"]
		case WM_RBUTTONDOWN:
			buttonName = "Right"
			isProtected = globalHook.protectedButtons["right"]
		case WM_MBUTTONDOWN:
			buttonName = "Middle"
			isProtected = globalHook.protectedButtons["middle"]
		case WM_XBUTTONDOWN:
			// For XBUTTON events in low-level mouse hook, button info is in high word of mouseData
			mouseStruct := (*MSLLHOOKSTRUCT)(lParam)
			buttonFlag := uint32(mouseStruct.MouseData>>16) & 0xFFFF

			// Check which X button is pressed based on the high word of mouseData
			if buttonFlag == 1 { // XBUTTON1
				buttonName = "XBUTTON1"
				isProtected = globalHook.protectedButtons["xbutton1"]
			} else if buttonFlag == 2 { // XBUTTON2
				buttonName = "XBUTTON2"
				isProtected = globalHook.protectedButtons["xbutton2"]
			} else {
				// Unknown XBUTTON, allow through
				ret, _, _ := procCallNextHookEx.Call(
					uintptr(globalHook.hook),
					uintptr(nCode),
					wParam,
					uintptr(lParam),
				)
				return ret
			}
		}

		// If this button is not protected, allow the event through
		if !isProtected {
			ret, _, _ := procCallNextHookEx.Call(
				uintptr(globalHook.hook),
				uintptr(nCode),
				wParam,
				uintptr(lParam),
			)
			return ret
		}

		// Drag Fix: Check if this DOWN is a bounce after a blocked UP
		if globalHook.dragFixEnabled {
			var upEvent uintptr
			switch wParam {
			case WM_LBUTTONDOWN:
				upEvent = WM_LBUTTONUP
			case WM_RBUTTONDOWN:
				upEvent = WM_RBUTTONUP
			case WM_MBUTTONDOWN:
				upEvent = WM_MBUTTONUP
			}

			if upEvent != 0 {
				globalHook.mu.Lock()
				pendingTime, hasPending := globalHook.pendingUpTimes[upEvent]
				if hasPending {
					// We have a pending UP event. This DOWN is likely a bounce.
					// Check if it's within a reasonable window of the UP (timer hasn't fired yet)
					if time.Since(pendingTime) < globalHook.dragFixThreshold+50*time.Millisecond {
						// It is a bounce! Clear the pending UP so it won't be injected.
						delete(globalHook.pendingUpTimes, upEvent)
						globalHook.mu.Unlock()

						globalHook.blockedCount++
						globalHook.sendLog(fmt.Sprintf("🛡️ DRAG FIX: Blocked bounce for %s button (UP+DOWN suppressed)", buttonName))
						return 1
					}
				}
				globalHook.mu.Unlock()
			}
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

	case WM_LBUTTONUP, WM_RBUTTONUP, WM_MBUTTONUP, WM_XBUTTONUP:
		// Determine which button was released
		var downEvent uintptr
		buttonName := "Left"
		isProtected := globalHook.protectedButtons["left"]

		switch wParam {
		case WM_LBUTTONUP:
			downEvent = WM_LBUTTONDOWN
			buttonName = "Left"
			isProtected = globalHook.protectedButtons["left"]
		case WM_RBUTTONUP:
			downEvent = WM_RBUTTONDOWN
			buttonName = "Right"
			isProtected = globalHook.protectedButtons["right"]
		case WM_MBUTTONUP:
			downEvent = WM_MBUTTONDOWN
			buttonName = "Middle"
			isProtected = globalHook.protectedButtons["middle"]
		case WM_XBUTTONUP:
			// For XBUTTON events in low-level mouse hook, button info is in high word of mouseData
			mouseStruct := (*MSLLHOOKSTRUCT)(lParam)
			buttonFlag := uint32(mouseStruct.MouseData>>16) & 0xFFFF

			// Check which X button is pressed based on the high word of mouseData
			if buttonFlag == 1 { // XBUTTON1
				downEvent = WM_XBUTTONDOWN
				buttonName = "XBUTTON1"
				isProtected = globalHook.protectedButtons["xbutton1"]
			} else if buttonFlag == 2 { // XBUTTON2
				downEvent = WM_XBUTTONDOWN
				buttonName = "XBUTTON2"
				isProtected = globalHook.protectedButtons["xbutton2"]
			} else {
				// Unknown XBUTTON, allow through
				ret, _, _ := procCallNextHookEx.Call(
					uintptr(globalHook.hook),
					uintptr(nCode),
					wParam,
					uintptr(lParam),
				)
				return ret
			}
		}

		// If this button is not protected, allow the event through
		if !isProtected {
			ret, _, _ := procCallNextHookEx.Call(
				uintptr(globalHook.hook),
				uintptr(nCode),
				wParam,
				uintptr(lParam),
			)
			return ret
		}

		now := time.Now()

		// Drag Fix: Delay UP event to check for bounce
		if !isInjected && globalHook.dragFixEnabled && isProtected {
			if wParam == WM_LBUTTONUP || wParam == WM_RBUTTONUP || wParam == WM_MBUTTONUP {
				// Safety Check: Can we inject into the foreground window?
				// If not (e.g. Admin window), do NOT block the UP event, or we'll get stuck.
				if !globalHook.canInjectToForeground() {
					// We detected a high-privilege window where injection would fail.
					// Instead of just allowing this ONE event, let's pause protection for a few seconds.
					// This allows the user to interact with the Admin app/VM without interference.

					globalHook.mu.Lock()
					globalHook.pauseUntil = time.Now().Add(globalHook.pauseDuration)
					globalHook.mu.Unlock()

					globalHook.sendLog(fmt.Sprintf("⚠️ High privilege window detected. Pausing protection for %v", globalHook.pauseDuration))

					// Allow this event through immediately
				} else {
					globalHook.mu.Lock()
					globalHook.pendingUpTimes[wParam] = now
					globalHook.mu.Unlock()

					// Start timer to inject event if no bounce occurs
					go func(btn uintptr, triggerTime time.Time) {
						time.Sleep(globalHook.dragFixThreshold)

						globalHook.mu.Lock()
						if pendingTime, exists := globalHook.pendingUpTimes[btn]; exists && pendingTime.Equal(triggerTime) {
							// Still pending! No DOWN came to clear it.
							// It was a real release.
							delete(globalHook.pendingUpTimes, btn)
							globalHook.mu.Unlock()

							// Inject the UP event
							globalHook.injectUp(btn)
						} else {
							globalHook.mu.Unlock()
						}
					}(wParam, now)

					return 1 // Block the original UP
				}
			}
		}

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

	case WM_MOUSEMOVE:
		// Check if any button is currently pressed and mark as drag (but don't spam logs)
		if globalHook.buttonPressed[WM_LBUTTONDOWN] && globalHook.protectedButtons["left"] && !globalHook.dragDetected[WM_LBUTTONDOWN] {
			globalHook.dragDetected[WM_LBUTTONDOWN] = true
			globalHook.sendLog("🖱️  Left button drag operation detected")
		}
		if globalHook.buttonPressed[WM_RBUTTONDOWN] && globalHook.protectedButtons["right"] && !globalHook.dragDetected[WM_RBUTTONDOWN] {
			globalHook.dragDetected[WM_RBUTTONDOWN] = true
			globalHook.sendLog("🖱️  Right button drag operation detected")
		}
		if globalHook.buttonPressed[WM_MBUTTONDOWN] && globalHook.protectedButtons["middle"] && !globalHook.dragDetected[WM_MBUTTONDOWN] {
			globalHook.dragDetected[WM_MBUTTONDOWN] = true
			globalHook.sendLog("🖱️  Middle button drag operation detected")
		}
		// Note: XBUTTON drag detection would require additional tracking since we don't have specific DOWN/UP events
	}

	ret, _, _ := procCallNextHookEx.Call(
		uintptr(globalHook.hook),
		uintptr(nCode),
		wParam,
		uintptr(lParam),
	)
	return ret
}

// GetMouseButtonCount returns the number of mouse buttons available on the system
func GetMouseButtonCount() int {
	// Check if a mouse is present
	mousePresent, _, _ := procGetSystemMetrics.Call(19) // SM_MOUSEPRESENT
	if mousePresent == 0 {
		return 0
	}
	// Get the number of mouse buttons
	buttonCount, _, _ := procGetSystemMetrics.Call(43) // SM_CMOUSEBUTTONS
	return int(buttonCount)
}

// GetMouseButtons returns a list of available mouse buttons with their names
func GetMouseButtons() []struct {
	Name string
	ID   string
} {
	buttons := []struct {
		Name string
		ID   string
	}{
		{Name: "Left Mouse Button", ID: "left"},
		{Name: "Right Mouse Button", ID: "right"},
		{Name: "Middle Mouse Button", ID: "middle"},
	}

	// Check how many buttons the mouse has
	buttonCount := GetMouseButtonCount()

	// Add XBUTTON1 and XBUTTON2 if the mouse has 4 or more buttons
	if buttonCount >= 4 {
		buttons = append(buttons, struct {
			Name string
			ID   string
		}{Name: "Mouse Button 4 (XBUTTON1)", ID: "xbutton1"})
	}

	// Add XBUTTON2 if the mouse has 5 or more buttons
	if buttonCount >= 5 {
		buttons = append(buttons, struct {
			Name string
			ID   string
		}{Name: "Mouse Button 5 (XBUTTON2)", ID: "xbutton2"})
	}

	return buttons
}

// This is the callback function that Windows will call
var mouseProc = syscall.NewCallback(LowLevelMouseProc)

// SetProtectedButtons sets which mouse buttons should be protected
func (w *WindowsHook) SetProtectedButtons(buttons []string) {
	w.protectedButtons = make(map[string]bool)
	for _, button := range buttons {
		w.protectedButtons[button] = true
	}
}

func (w *WindowsHook) SetDragFix(enabled bool, threshold int, pauseDuration int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.dragFixEnabled = enabled
	w.dragFixThreshold = time.Duration(threshold) * time.Millisecond
	if pauseDuration < 1 {
		pauseDuration = 1
	} else if pauseDuration > 5 {
		pauseDuration = 5
	}
	w.pauseDuration = time.Duration(pauseDuration) * time.Second
	if !enabled {
		w.pendingUpTimes = make(map[uintptr]time.Time)
	}
}

func (w *WindowsHook) injectUp(button uintptr) {
	var flags uintptr
	switch button {
	case WM_LBUTTONUP:
		flags = MOUSEEVENTF_LEFTUP
	case WM_RBUTTONUP:
		flags = MOUSEEVENTF_RIGHTUP
	case WM_MBUTTONUP:
		flags = MOUSEEVENTF_MIDDLEUP
	}

	if flags != 0 {
		procMouseEvent.Call(flags, 0, 0, 0, 0)
	}
}

func (w *WindowsHook) Start(delay time.Duration, logChan chan string) error {
	if w.isRunning {
		return fmt.Errorf("hook is already running")
	}

	w.delay = delay
	w.logChannel = logChan
	w.isRunning = true
	w.buttonPressed = make(map[uintptr]bool)
	w.buttonPressTime = make(map[uintptr]time.Time)
	w.dragDetected = make(map[uintptr]bool)
	w.faultyClickPattern = make(map[uintptr][]time.Time)
	w.adaptiveDelay = make(map[uintptr]time.Duration)
	w.shortClickCount = make(map[uintptr]int)
	w.lastDownTime = make(map[uintptr]time.Time)
	w.lastDownBlocked = make(map[uintptr]bool)
	w.lastUpTime = make(map[uintptr]time.Time)
	w.pendingUpTimes = make(map[uintptr]time.Time)

	// Initialize protected buttons if not already set
	if w.protectedButtons == nil {
		w.protectedButtons = map[string]bool{"left": true} // Default to left button only
	}

	globalHook = w

	go func() {
		ret, _, err := procSetWindowsHookEx.Call(
			WH_MOUSE_LL,
			mouseProc,
			0,
			0,
		)

		if ret == 0 {
			w.logChannel <- fmt.Sprintf("❌ Failed to install mouse hook: %v", err)
			w.isRunning = false
			return
		}

		w.hook = syscall.Handle(ret)
		w.logChannel <- "🎯 Mouse hook installed successfully - protection active!"

		var msg MSG
		for w.isRunning {
			ret, _, _ := procGetMessage.Call(
				uintptr(unsafe.Pointer(&msg)),
				0,
				0,
				0,
			)

			if ret == 0 {
				break // WM_QUIT
			}

			procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}()

	return nil
}

func (w *WindowsHook) Stop() error {
	if !w.isRunning {
		return nil
	}

	w.isRunning = false
	if w.hook != 0 {
		procUnhookWindowsHookEx.Call(uintptr(w.hook))
		w.hook = 0
		if w.logChannel != nil {
			w.logChannel <- "🛑 Mouse hook removed - protection stopped"
		}
	}
	globalHook = nil
	return nil
}

func (w *WindowsHook) GetBlockedCount() int {
	return w.blockedCount
}

func (w *WindowsHook) ResetBlockedCount() {
	w.blockedCount = 0
}

func (w *WindowsHook) IsSupported() bool {
	return true
}

// detectFaultyHardware analyzes click patterns to detect faulty mouse behavior
func (w *WindowsHook) detectFaultyHardware(button uintptr, holdDuration time.Duration) {
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

		// Adaptive delay should never be less than user-selected delay
		// Only increase delay for detected faulty patterns, never decrease
		if float64(shortClicks)/5.0 > 0.6 {
			// Increase delay for faulty hardware (never reduce below user setting)
			newDelay := w.delay + (w.delay / 4) // Add 25% to user delay
			if newDelay > 200*time.Millisecond {
				newDelay = 200 * time.Millisecond // Maximum 200ms
			}

			buttonName := "Left"
			if button == WM_RBUTTONDOWN {
				buttonName = "Right"
			}

			if w.adaptiveDelay[button] != newDelay {
				w.adaptiveDelay[button] = newDelay
				w.sendLog(fmt.Sprintf("🔧 ADAPTIVE STRICT: %s button delay increased to %.0fms due to detected low-pressure pattern",
					buttonName, float64(newDelay.Nanoseconds())/1000000))
			}
		} else {
			// Reset to user-selected delay if pattern improves
			if w.adaptiveDelay[button] != w.delay {
				w.adaptiveDelay[button] = w.delay
				buttonName := "Left"
				if button == WM_RBUTTONDOWN {
					buttonName = "Right"
				}
				w.sendLog(fmt.Sprintf("🔧 ADAPTIVE STRICT: %s button delay reset to user setting (%.0fms)",
					buttonName, float64(w.delay.Nanoseconds())/1000000))
			}
		}
	}
}

// getEffectiveDelay returns the adaptive delay for a specific button
// Never returns a delay less than the user-selected delay
func (w *WindowsHook) getEffectiveDelay(button uintptr) time.Duration {
	if adaptiveDelay, exists := w.adaptiveDelay[button]; exists && adaptiveDelay >= w.delay {
		return adaptiveDelay
	}
	return w.delay
}
