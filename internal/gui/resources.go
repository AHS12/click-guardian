// Icon resources for Click Guardian application

package gui

import "fyne.io/fyne/v2"

var appIconResource = &fyne.StaticResource{
	StaticName: "icon-modern-shield.svg",
	StaticContent: []byte(
		"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"256\" height=\"256\" viewBox=\"0 0 256 256\">\r\n    <!-- Shield outline -->\r\n    <path d=\"M128 30L70 65v75c0 35 58 86 58 86s58-51 58-86V65z\"\r\n        fill=\"#E3F2FD\" stroke=\"#1976D2\" stroke-width=\"4\" />\r\n\r\n    <!-- Inner shield -->\r\n    <path d=\"M128 45L85 75v60c0 25 43 65 43 65s43-40 43-65V75z\"\r\n        fill=\"#2196F3\" />\r\n\r\n    <!-- Mouse cursor with modern design -->\r\n    <g transform=\"translate(118, 105)\">\r\n        <path d=\"M0 0L0 28L6 22L10 32L14 30L10 20L20 20z\"\r\n            fill=\"#FFFFFF\" stroke=\"#0D47A1\" stroke-width=\"1.2\" />\r\n        <!-- Click indicator -->\r\n        <circle cx=\"2\" cy=\"4\" r=\"4\" fill=\"none\" stroke=\"#FF5722\" stroke-width=\"2\" opacity=\"0.9\" />\r\n        <circle cx=\"2\" cy=\"4\" r=\"8\" fill=\"none\" stroke=\"#FF5722\" stroke-width=\"1\" opacity=\"0.6\" />\r\n    </g>\r\n</svg>"),
}

// GetAppIcon returns the application icon as a Fyne resource
func GetAppIcon() fyne.Resource {
	return appIconResource
}
