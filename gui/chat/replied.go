package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "fmt"
    "image/color"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
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
    repliedMsg *chat.Message, repliedID int, u *user.UserCache, opt *options.Options,
) *RepliedMessage {
    log.Info.Printf("Creating replied message widget for message %d\n", repliedID)

    repliedMsgWidget := &RepliedMessage{}
    repliedMsgWidget.MessageID = repliedID

    var repliedText string
    if repliedMsg == nil {
        repliedText = opt.GUIText.Greeting.PlaceholderUnloaded
    } else {
        repliedText = fmt.Sprintf("%s: %s", 
            *repliedMsg.Username,
            repliedMsg.Content,
        )
    }

    repliedMsgWidget.Icon = widget.NewIcon(theme.Current().Icon(theme.IconNameMailReply))
    repliedMsgWidget.Text = canvas.NewText(repliedText, color.NRGBA{ 255, 255, 255, 255 })

    repliedMsgWidget.Base = container.NewHBox(repliedMsgWidget.Icon, repliedMsgWidget.Text)

    return repliedMsgWidget
}
