package resources

import "fyne.io/fyne/v2"

// GetAppIcon returns the application icon as a Fyne resource
func GetAppIcon() fyne.Resource {
	return resourceIconModernShieldSolidSvg
}

// GetTrayIcon returns the application system tray icon as a Fyne resource(ico)
func GetTrayIcon() fyne.Resource {
	return resourceTrayIconIco
}
