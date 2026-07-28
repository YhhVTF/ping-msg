package gui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
)

func createButtonOptions(g *GUI, opt *options.Options) *widget.Button {
    // Initialize options button
    button := widget.NewButton(opt.GUIText.ButtonOptions.Label, func(){
        log.Info.Printf("Widget ButtonOptions pressed\n")
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

        w := container.NewInnerWindow("Options", widget.NewButton("e", func(){}))
        mw := container.NewMultipleWindows(w)

        w.Resize(fyne.NewSize(600, 500))
        w.SetMaximized(false)
        w.CloseIntercept = func() {
            g.Window.Canvas().Overlays().Remove(mw)
            g.Window.Canvas().Content().Refresh()
        }
        g.Window.Canvas().Overlays().Add(mw)
    })
    return button
}
