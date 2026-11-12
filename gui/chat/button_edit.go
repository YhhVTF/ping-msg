package schat

import (
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/log"
)

func createButtonEdit(g *ScreenChat, listItemID widget.ListItemID, chat *chat.Chat) *widget.Button {
    msgID := chat.Islands[chat.RenderedIsland].MessageIDs[listItemID]

    return widget.NewButtonWithIcon("", theme.MailComposeIcon(), func() {
        log.Info.Printf("Edit button on message %d pressed\n", msgID)

        // Mark this message as being edited and refresh the messages list
        chat.Editing = msgID
        g.Widgets.MessagesList.Refresh()
    })
}
