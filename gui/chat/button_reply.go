package schat

import (
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "slices"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/log"
)

func createButtonReply(g *ScreenChat, msgID int, chat *chat.Chat) *widget.Button {
    return widget.NewButtonWithIcon("", theme.MailReplyIcon(), func() {
        log.Info.Printf("Reply button on message %d pressed\n", msgID)
        // Give focus back to the message entry when done
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

        // Add the message ID to chat.Selections.ReplyingTo if it isn't already there or remove it from there if it is
        if slices.Contains(chat.Selections.ReplyingTo, msgID) {
            i := slices.Index(chat.Selections.ReplyingTo, msgID)
            chat.Selections.ReplyingTo = slices.Delete(chat.Selections.ReplyingTo, i, i+1)
        } else {
            chat.Selections.ReplyingTo = append(chat.Selections.ReplyingTo, msgID)
        }

        log.Info.Printf("The next message sent will reply to messages %d\n", chat.Selections.ReplyingTo)
    })
}
