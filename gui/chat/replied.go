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
)

// A widget representing a reply to a message. Is a component of Message widgets that reply to other messages
type RepliedMessage struct {
    Base        *fyne.Container
    Icon        *widget.Icon
    MessageID   int
    Text        *canvas.Text
}

func createRepliedMessage(g *ScreenChat, msgID int) *RepliedMessage {
    log.Info.Printf("Creating replied message widget for message %d\n", msgID)

    repliedMsg := &RepliedMessage{}
    repliedMsg.MessageID = msgID

    msg, exists := g.Widgets.Messages[msgID]
    var repliedText string
    if !exists {
        repliedText = "Message could not be loaded"
    } else {
        repliedText = 
            fmt.Sprintf("%s: %s", msg.Username.Text, msg.Content.Text)
    }

    repliedMsg.Icon = widget.NewIcon(theme.Current().Icon(theme.IconNameMailReply))
    repliedMsg.Text = canvas.NewText(repliedText, color.NRGBA{ 255, 255, 255, 255 })

    repliedMsg.Base = container.NewHBox(repliedMsg.Icon, repliedMsg.Text)

    return repliedMsg
}
