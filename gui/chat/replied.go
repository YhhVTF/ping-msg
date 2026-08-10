package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "fmt"
    "image/color"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

// A widget representing a reply to a message. Is a component of Message widgets that reply to other messages
type RepliedMessage struct {
    Base        *fyne.Container
    Icon        *widget.Icon
    MessageID   int
    Text        *canvas.Text
}

func createRepliedMessage(
    g *ScreenChat, repliedMsg *prot.MessageRaw, u *user.UserCache,
) *RepliedMessage {
    log.Info.Printf("Creating replied message widget for message %d\n", repliedMsg.ID)

    repliedMsgWidget := &RepliedMessage{}
    repliedMsgWidget.MessageID = repliedMsg.ID

    var repliedText string
    if repliedMsg == nil {
        repliedText = "Message could not be loaded"
    } else {
        repliedText = fmt.Sprintf("%s: %s", 
            u.Users[repliedMsg.UserID].Username,
            repliedMsg.Content,
        )
    }

    repliedMsgWidget.Icon = widget.NewIcon(theme.Current().Icon(theme.IconNameMailReply))
    repliedMsgWidget.Text = canvas.NewText(repliedText, color.NRGBA{ 255, 255, 255, 255 })

    repliedMsgWidget.Base = container.NewHBox(repliedMsgWidget.Icon, repliedMsgWidget.Text)

    return repliedMsgWidget
}
