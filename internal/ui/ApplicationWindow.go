package ui

import (
	"image"
	"image/color"
	"log"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/dertseha/alarma/internal/config"
)

// ApplicationWindow is the main window.
type ApplicationWindow struct {
	configUpdater func(config.Instance)
	configuration config.Instance

	window *app.Window
}

// NewApplicationWindow returns a new instance.
func NewApplicationWindow(configUpdater func(config.Instance)) *ApplicationWindow {
	return &ApplicationWindow{
		configUpdater: configUpdater,
	}
}

// Show opens and runs the actual window.
func (appWindow *ApplicationWindow) Show() {
	go func() {
		appWindow.window = app.NewWindow(app.Size(unit.Px(800), unit.Px(450)), app.Title("alarma"))
		if err := appWindow.loop(); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

// Update sets the UI to the provided state.
func (appWindow *ApplicationWindow) Update(now time.Time, configuration config.Instance) {
	appWindow.configuration = configuration
	appWindow.window.Invalidate()
}

func (appWindow *ApplicationWindow) toggleGlobalActive() {
	newConfiguration := appWindow.configuration
	newConfiguration.TimeSpansActive = !appWindow.configuration.TimeSpansActive
	appWindow.configUpdater(newConfiguration)
}

func (appWindow *ApplicationWindow) loop() error {
	gofont.Register()
	th := material.NewTheme()
	gtx := layout.NewContext(appWindow.window.Queue())
	button := new(widget.Button)
	for {
		e := <-appWindow.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)
			fill(gtx, color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF})

			globalActive := appWindow.configuration.TimeSpansActive
			for button.Clicked(gtx) {
				globalActive = !globalActive
			}
			if globalActive != appWindow.configuration.TimeSpansActive {
				appWindow.toggleGlobalActive()
			}

			themeButton := th.Button(time.Now().Format("15:04"))
			themeButton.Font.Size = unit.Dp(100)
			themeButton.Background = color.RGBA{}
			if appWindow.configuration.TimeSpansActive {
				themeButton.Color = color.RGBA{R: 0x40, G: 0x14, B: 0x54, A: 0xFF}
			} else {
				themeButton.Color = color.RGBA{R: 0x40, G: 0x54, B: 0x14, A: 0xFF}
			}
			themeButton.Layout(gtx, button)

			e.Frame(gtx.Ops)
		}
	}
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
