package sside

import (
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/user"
)

func createButtonCreate(g *ScreenSidebar, u *user.UserCache, opt *options.Options) *widget.Button {
    return widget.NewButton(opt.GUIText.ButtonCreateChat.Label, func() {
        log.Info.Printf("Widget ButtonCreateChat pressed\n")
        g.ChatMetadataRequestCreate(u)
    })
}
