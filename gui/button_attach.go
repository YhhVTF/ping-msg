package gui

import (
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
)

func createButtonAttach(g *GUI, opt *options.Options) *widget.Button {
    // Initialize attach button
    button := widget.NewButton(opt.GUIText.ButtonAttach.Label, func(){
        log.Info.Printf("Widget ButtonAttach pressed\n")
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)
    })
    return button
}
