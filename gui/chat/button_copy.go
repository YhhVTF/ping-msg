package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/log"
)

func createButtonCopy(g *ScreenChat, listItemID widget.ListItemID, chat *chat.Chat) *widget.Button {
    msgID := chat.Islands[chat.RenderedIsland].MessageIDs[listItemID]
    msg := chat.Messages[msgID]

	return widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		log.Info.Printf("Copied message %d\n", msgID)

		// Give focus back to the message entry when done
		defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)
		// Copy the message contents
		fyne.CurrentApp().Clipboard().SetContent(msg.Content)
	})
}
