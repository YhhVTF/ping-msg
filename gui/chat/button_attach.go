package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
)

func buttonAttachOnPressed(g *ScreenChat) {
    log.Info.Printf("Creating fyne default open file picker dialog\n")

    filePicker := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
        if err != nil {
            log.Error.Printf("Failed to open file: %s\n", err)
            return
        }
        g.Attachments = append(g.Attachments, reader)

        log.Info.Printf("%d files attached\n", len(g.Attachments))
    }, g.Window)
    filePicker.Show()
}

func createButtonAttach(g *ScreenChat, opt *options.Options) *widget.Button {
    // Initialize attach button
    button := widget.NewButton(opt.GUIText.ButtonAttach.Label, func(){
        log.Info.Printf("Widget ButtonAttach pressed\n")
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

        buttonAttachOnPressed(g)
    })
    return button
}
