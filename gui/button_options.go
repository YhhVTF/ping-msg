package gui

import (
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
)

func createButtonOptions(g *GUI, opt *options.Options) *widget.Button {
    // Initialize options button
    button := widget.NewButton(opt.GUIText.ButtonOptions.Label, func(){
        log.Info.Printf("Widget ButtonOptions pressed\n")
        g.Window.Canvas().Focus(g.Widgets.EntryMessage)
    })
    return button
}
