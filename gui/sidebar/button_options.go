package sside

import (
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/log"
)

func createButtonOptions(s *screen.ScreenManager) *widget.Button {
    // Initialize options button
    button := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameSettings), func(){
        log.Info.Printf("Widget ButtonOptions pressed\n")
        defer s.ScreenChatFocusDefault()
        // Create options screen as an inner/floating window
        s.ScreenOptionsFloat()
    })
    return button
}
