package schat

import (
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
)

func createButtonOptions(g *ScreenChat, s *screen.ScreenManager, opt *options.Options) *widget.Button {
    // Initialize options button
    button := widget.NewButton(opt.GUIText.ButtonOptions.Label, func(){
        log.Info.Printf("Widget ButtonOptions pressed\n")
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)
        // Create options screen as an inner/floating window
        s.ScreenOptionsFloat()
    })
    return button
}
