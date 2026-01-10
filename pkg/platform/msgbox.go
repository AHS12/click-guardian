package platform

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

// ShowMessageBox shows a native message box on Windows
func ShowMessageBox(title, message string) {
	if runtime.GOOS == "windows" {
		showWindowsMessageBox(title, message)
	} else {
		// For other platforms, just print to console for now
		fmt.Printf("[%s] %s\n", title, message)
	}
}

// showWindowsMessageBox shows a Windows message box using syscall
func showWindowsMessageBox(title, message string) {
	// Load the user32.dll library
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBox := user32.NewProc("MessageBoxW")

	// Convert Go strings to Windows-compatible UTF16 pointers
	titlePtr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return
	}

	messagePtr, err := syscall.UTF16PtrFromString(message)
	if err != nil {
		return
	}

	// Show the message box
	// Parameters: HWND (0 for no parent), message, title, type (0 for OK button)
	MessageBox.Call(
		0, // No parent window
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		0, // OK button only
	)
}
