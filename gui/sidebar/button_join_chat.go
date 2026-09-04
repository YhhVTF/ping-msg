package sside

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/gui/dialogs"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func createButtonJoin(
    g *ScreenSidebar, w fyne.Window, u *user.UserCache, opt *options.Options,
) *widget.Button {
    return widget.NewButton(opt.GUIText.ButtonJoinChat.Label, func() {
        log.Info.Printf("Widget ButtonJoinChat pressed\n")
        defer g.ScreenManager.ScreenChatFocusDefault()

        dialogDone := make(chan int)
        dialogs.InitDialogJoinChat(w, dialogDone, opt)

        go func() {
            chatID := <-dialogDone
            if chatID != prot.NONE_INT {
                g.ChatMetadataRequestJoin(chatID, u)
            }
        }()
    })
}
