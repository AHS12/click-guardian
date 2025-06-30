package dialogs

import (
    "fmt"
    "image/color"
    "net/url"

    "click-guardian/internal/version"
    "click-guardian/internal/gui/resources"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
)

func ShowAboutDialog(window fyne.Window) {
    appInfo := version.GetAppInfo()

    icon := widget.NewIcon(resources.GetAppIcon())
    headerText := canvas.NewText("About", color.White)
    headerText.TextStyle = fyne.TextStyle{Bold: true}
    headerText.TextSize = 22
    headerText.Alignment = fyne.TextAlignCenter

    titleText := canvas.NewText("Click Guardian", color.White)
    titleText.TextStyle = fyne.TextStyle{Bold: true}
    titleText.TextSize = 20
    titleText.Alignment = fyne.TextAlignCenter

    versionText := canvas.NewText(fmt.Sprintf("Version %s", appInfo.Version), color.White)
    versionText.Alignment = fyne.TextAlignCenter

    var buildText *widget.Label
    if appInfo.BuildTime != "unknown" && appInfo.GitCommit != "unknown" {
        buildInfo := fmt.Sprintf("Built: %s\nCommit: %s\nBy: %s",
            appInfo.BuildTime,
            appInfo.GitCommit[:min(len(appInfo.GitCommit), 8)],
            appInfo.BuildBy)
        buildText = widget.NewLabel(buildInfo)
        buildText.Alignment = fyne.TextAlignCenter
    }

    systemInfo := fmt.Sprintf("Platform: %s %s\nGo: %s",
        appInfo.Platform,
        appInfo.Arch,
        appInfo.GoVersion)
    systemText := widget.NewLabel(systemInfo)
    systemText.Alignment = fyne.TextAlignCenter

    descText := widget.NewLabel(appInfo.Description)
    descText.Wrapping = fyne.TextWrapWord
    descText.Alignment = fyne.TextAlignCenter

    copyrightText := widget.NewLabel(appInfo.Copyright)
    copyrightText.Alignment = fyne.TextAlignCenter

    licenseText := widget.NewLabel("License: GNU General Public License v3.0")
    licenseText.Alignment = fyne.TextAlignCenter

    githubURL, _ := url.Parse("https://github.com/ahs12/click-guardian")
    githubLink := widget.NewHyperlink("GitHub Repository", githubURL)
    githubLink.Alignment = fyne.TextAlignCenter

    content := container.NewVBox(
        container.NewHBox(layout.NewSpacer(), icon, layout.NewSpacer()),
        container.NewCenter(headerText),
        widget.NewSeparator(),
        container.NewCenter(titleText),
        container.NewCenter(versionText),
    )
    if buildText != nil {
        content.Add(widget.NewSeparator())
        content.Add(buildText)
    }
    content.Add(widget.NewSeparator())
    content.Add(systemText)
    content.Add(widget.NewSeparator())
    content.Add(descText)
    content.Add(widget.NewSeparator())
    content.Add(copyrightText)
    content.Add(licenseText)
    content.Add(githubLink)

    aboutDialog := dialog.NewCustom("", "Close", content, window)
    aboutDialog.Resize(fyne.NewSize(400, 500))
    aboutDialog.Show()
}

// Helper for substring length
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}