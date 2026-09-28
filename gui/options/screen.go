package sopt

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	options "github.com/YhhVTF/ping-msg/opt"
)

type ContainerTableOptions struct {
	// Base container of the options sceen, makes the child container scrollable
	Base *container.Scroll
	// Contains widgets for changing options
	VBox *fyne.Container
}

type ScreenOptions struct {
	Containers ContainerTableOptions
	Widgets    WidgetTableOptions
	Window     fyne.Window
}

type WidgetTableOptions struct {
	CardLanguage   *widget.Card
	CardAudio      *widget.Card
	SelectLanguage *widget.Select
	CheckAudioEn   *widget.Check
	SliderVolume   *widget.Slider
}

func (g *ScreenOptions) ExitFloat(
	wOptions *container.InnerWindow,
	wMain fyne.Window, innerWindows *container.MultipleWindows,
	opt *options.Options,
) {
	// Set options window size in options to its current size
	opt.GUI.WindowOptions.Size[0] = wOptions.Size().Width
	opt.GUI.WindowOptions.Size[1] = wOptions.Size().Height

	// Close the options window
	wMain.Canvas().Overlays().Remove(innerWindows)
	wOptions.Hide()
	wMain.Canvas().Content().Refresh()

	// Save options
	opt.SaveGUI(".ping/options")
	opt.SaveNet(".ping/options")
}

func (g *ScreenOptions) Float(
	wMain fyne.Window, innerWindows *container.MultipleWindows, opt *options.Options,
) {
	w := container.NewInnerWindow(
		opt.GUIText.WindowOptions.Title, g.Containers.Base,
	)
	innerWindows.Add(w)

	w.SetPadded(true)
	w.Resize(fyne.NewSize(
		opt.GUI.WindowOptions.Size[0], opt.GUI.WindowOptions.Size[1],
	))
	w.SetMaximized(false)

	// Center options floating window
	posX := (wMain.Canvas().Size().Width / 2) - (w.Size().Width / 2)
	posY := (wMain.Canvas().Size().Height / 2) - (w.Size().Height / 2)
	w.Move(fyne.NewPos(posX, posY))

	w.CloseIntercept = func() {
		g.ExitFloat(w, wMain, innerWindows, opt)
	}
	g.Window.Canvas().Overlays().Add(innerWindows)
}

func InitScreenOptions(
	w fyne.Window, opt *options.Options,
) *ScreenOptions {
	g := &ScreenOptions{}
	g.Window = w

	g.Widgets.SelectLanguage = createSelectLanguage(opt)
	g.Widgets.CardLanguage =
		widget.NewCard("", opt.GUIText.OptionCardLanguage.Subtitle, g.Widgets.SelectLanguage)

	g.Widgets.CheckAudioEn = widget.NewCheck("Enable Sound Effects", func(checked bool) {
		opt.GUI.Audio.Enabled = checked
		opt.SaveGUI(".ping/options")
	})
	g.Widgets.CheckAudioEn.SetChecked(opt.GUI.Audio.Enabled)

	g.Widgets.SliderVolume = widget.NewSlider(0, 1)
	g.Widgets.SliderVolume.Value = opt.GUI.Audio.Volume
	g.Widgets.SliderVolume.Step = 0.05
	g.Widgets.SliderVolume.OnChanged = func(val float64) {
		opt.GUI.Audio.Volume = val
		opt.SaveGUI(".ping/options")
	}

	audioContainer := container.NewVBox(
		g.Widgets.CheckAudioEn,
		widget.NewLabel("Volume"),
		g.Widgets.SliderVolume,
	)
	g.Widgets.CardAudio = widget.NewCard("Audio Settings", "Configure client sound effects", audioContainer)

	g.Containers.VBox = container.NewVBox(g.Widgets.CardLanguage, g.Widgets.CardAudio)
	g.Containers.Base = container.NewVScroll(g.Containers.VBox)

	return g
}
