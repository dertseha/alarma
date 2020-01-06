package ui

import (
	"time"

	"github.com/dertseha/jellui"
	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/controls"
	"github.com/dertseha/jellui/env/native"
	"github.com/dertseha/jellui/font"
	"github.com/dertseha/jellui/graphics"

	"github.com/dertseha/alarma/internal/config"
)

// ApplicationWindow is the main window.
type ApplicationWindow struct {
	app *jellui.StandardApplication

	activeToggle *area.Area
	timeWindow   *area.Area
	timeLabel    *controls.Label

	timeFontPainter graphics.TextPainter

	configUpdater func(config.Instance)
	configuration config.Instance
}

// NewApplicationWindow returns a new instance.
func NewApplicationWindow(configUpdater func(config.Instance)) *ApplicationWindow {
	return &ApplicationWindow{
		timeFontPainter: graphics.NewBitmapTextPainter(font.ColorHeadingShock, 0x00),
		configUpdater:   configUpdater}
}

// Show opens and runs the actual window.
func (appWindow *ApplicationWindow) Show(deferrer <-chan func()) {
	appWindow.app = jellui.NewStandardApplication(appWindow.initInterface)

	native.Run(appWindow.app, "Alarma", 15.0, deferrer)
}

func (appWindow *ApplicationWindow) initInterface(app *jellui.StandardApplication, rootArea *area.Area) {
	app.SetFullScreen(true)
	app.SetCursorVisible(false)
	appWindow.onGlobalActiveChanged(false)

	{
		activeToggleBuilder := area.NewAreaBuilder()
		activeToggleBuilder.SetParent(rootArea)
		activeToggleBuilder.SetLeft(area.NewOffsetAnchor(rootArea.Left(), 0))
		activeToggleBuilder.SetTop(area.NewOffsetAnchor(rootArea.Top(), 0))
		activeToggleBuilder.SetRight(area.NewOffsetAnchor(rootArea.Right(), 0))
		activeToggleBuilder.SetBottom(area.NewOffsetAnchor(rootArea.Bottom(), 0))
		activeToggleBuilder.OnEvent(events.MouseButtonUpEventType, func(area *area.Area, event events.Event) bool {
			appWindow.toggleGlobalActive()
			return true
		})
		appWindow.activeToggle = activeToggleBuilder.Build()
	}
	{
		windowBuilder := area.NewAreaBuilder()
		windowBuilder.SetParent(rootArea)

		windowHorizontalCenter := area.NewRelativeAnchor(rootArea.Left(), rootArea.Right(), 0.5)
		windowVerticalCenter := area.NewRelativeAnchor(rootArea.Top(), rootArea.Bottom(), 0.5)

		windowBuilder.SetLeft(area.NewOffsetAnchor(windowHorizontalCenter, -250.0))
		windowBuilder.SetRight(area.NewOffsetAnchor(windowHorizontalCenter, 250.0))
		windowBuilder.SetTop(area.NewOffsetAnchor(windowVerticalCenter, -70.0))
		windowBuilder.SetBottom(area.NewOffsetAnchor(windowVerticalCenter, 70.0))

		appWindow.timeWindow = windowBuilder.Build()
	}
	{
		labelBuilder := app.ForLabel()
		labelBuilder.SetParent(appWindow.timeWindow)
		labelBuilder.SetLeft(area.NewOffsetAnchor(appWindow.timeWindow.Left(), 0))
		labelBuilder.SetRight(area.NewOffsetAnchor(appWindow.timeWindow.Right(), 0))
		labelBuilder.SetTop(area.NewOffsetAnchor(appWindow.timeWindow.Top(), 0))
		labelBuilder.SetBottom(area.NewOffsetAnchor(appWindow.timeWindow.Bottom(), 0))
		labelBuilder.WithTextPainter(appWindow.timeFontPainter)
		labelBuilder.SetScale(8.0)
		appWindow.timeLabel = labelBuilder.Build()
	}
}

// Update sets the UI to the provided state.
func (appWindow *ApplicationWindow) Update(now time.Time, configuration config.Instance) {
	if appWindow.configuration.TimeSpansActive != configuration.TimeSpansActive {
		appWindow.onGlobalActiveChanged(configuration.TimeSpansActive)
	}
	appWindow.configuration = configuration
	appWindow.timeLabel.SetText(now.Format("15:04"))
}

func (appWindow *ApplicationWindow) toggleGlobalActive() {
	newConfiguration := appWindow.configuration
	newConfiguration.TimeSpansActive = !appWindow.configuration.TimeSpansActive
	appWindow.configUpdater(newConfiguration)
}

func (appWindow *ApplicationWindow) onGlobalActiveChanged(active bool) {
	var uiTextPalette map[int][4]byte

	if active {
		uiTextPalette = map[int][4]byte{
			0: {0x00, 0x00, 0x00, 0x00},
			1: {0x80, 0x94, 0x54, 0xFF},
			2: {0x00, 0x00, 0x00, 0xC0},

			90: {0x80, 0x54, 0x94, 0xFF},
			92: {0x70, 0x44, 0x84, 0xFF},
			94: {0x60, 0x34, 0x74, 0xFF},
			95: {0x50, 0x24, 0x64, 0xFF},
			98: {0x40, 0x14, 0x54, 0x20}}
	} else {
		uiTextPalette = map[int][4]byte{
			0: {0x00, 0x00, 0x00, 0x00},
			1: {0x80, 0x94, 0x54, 0xFF},
			2: {0x00, 0x00, 0x00, 0xC0},

			90: {0x80, 0x94, 0x54, 0xFF},
			92: {0x70, 0x84, 0x44, 0xFF},
			94: {0x60, 0x74, 0x34, 0xFF},
			95: {0x50, 0x64, 0x24, 0xFF},
			98: {0x40, 0x54, 0x14, 0x20}}
	}

	appWindow.app.SetUITextPalette(uiTextPalette)
}
