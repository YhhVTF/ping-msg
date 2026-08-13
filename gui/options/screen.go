package sopt

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/opt"
)

type ContainerTableOptions struct {
    // Base container of the options sceen, makes the child container scrollable
    Base *container.Scroll
    // Contains widgets for changing options
    VBox *fyne.Container
}

type ScreenOptions struct {
    Containers      ContainerTableOptions
    Widgets         WidgetTableOptions
    Window          fyne.Window
}

type WidgetTableOptions struct {
    CardLanguage    *widget.Card
    SelectLanguage  *widget.Select
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

    g.Containers.VBox = container.NewVBox(g.Widgets.CardLanguage)
    g.Containers.Base = container.NewVScroll(g.Containers.VBox)

    return g
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
