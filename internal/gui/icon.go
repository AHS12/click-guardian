package gui

import (
	_ "embed"
)

//go:embed tray-icon.ico
var trayIconData []byte

// getTrayIcon returns the modern shield icon for the system tray
func getTrayIcon() []byte {
	return trayIconData
}
