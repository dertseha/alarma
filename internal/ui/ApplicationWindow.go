package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dertseha/alarma/internal/config"
)

// ApplicationWindow is the main window.
type ApplicationWindow struct {
	configUpdater func(config.Instance)
	configuration config.Instance

	fyneApp fyne.App
	fyneWin fyne.Window
	button  *widget.Button
}

// NewApplicationWindow returns a new instance.
func NewApplicationWindow(configUpdater func(config.Instance)) *ApplicationWindow {
	appWindow := &ApplicationWindow{
		configUpdater: configUpdater,

		fyneApp: app.New(),
	}
	appWindow.fyneApp.Settings().SetTheme(darkTheme{isActive: func() bool { return appWindow.configuration.TimeSpansActive }})

	appWindow.fyneWin = appWindow.fyneApp.NewWindow("alarma")
	appWindow.fyneWin.Resize(fyne.Size{
		Width:  800.0,
		Height: 450.0,
	})
	appWindow.button = widget.NewButton("__:__", func() {
		appWindow.toggleGlobalActive()
	})
	panel := container.NewBorder(nil, nil, nil, nil, appWindow.button)
	appWindow.fyneWin.SetContent(panel)

	return appWindow
}

// Show opens and runs the actual window.
func (appWindow *ApplicationWindow) Show() {
	appWindow.fyneWin.SetFullScreen(true)
	appWindow.fyneWin.ShowAndRun()
}

// Update sets the UI to the provided state.
func (appWindow *ApplicationWindow) Update(now time.Time, configuration config.Instance) {
	appWindow.configuration = configuration
	appWindow.button.SetText(now.Format("15:04"))
}

func (appWindow *ApplicationWindow) toggleGlobalActive() {
	newConfiguration := appWindow.configuration
	newConfiguration.TimeSpansActive = !appWindow.configuration.TimeSpansActive
	appWindow.configUpdater(newConfiguration)
}
