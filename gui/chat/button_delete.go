package schat

import (
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/user"
)

func createButtonDelete(g *ScreenChat, msgID int, c *chat.ChatCache, u *user.UserCache) *widget.Button {
    return widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
        log.Info.Printf("Delete button on message %d pressed\n", msgID)

        // Give focus back to the message entry when done
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

        // Send new DEL chat request to net.serverSend
        g.ChatRequestDelete(msgID, c, u)
    })
}
