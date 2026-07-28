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
        w.Resize(fyne.NewSize(500, 600))
        w.SetMaximized(false)
    })
    return button
}
