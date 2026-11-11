package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "time"

    "github.com/YhhVTF/ping-msg/chat"
)

func createListMessages(chat *chat.Chat) *widget.List {
    return widget.NewList(
        func() int {
            if len(chat.Islands) == 0 {
                return 0
            }
            return len(chat.Islands[chat.RenderedIsland].MessageIDs)
        },
        func() fyne.CanvasObject {
            return createMessageCardTemplate()
        },
        func(itemID widget.ListItemID, o fyne.CanvasObject) {
            msgID := chat.Islands[chat.RenderedIsland].MessageIDs[itemID]
            msg := chat.Messages[msgID]

            card := o.(*fyne.Container).Objects[1].(*widget.Card)
            wTopOfMsg := card.Content.(*fyne.Container).Objects[0].(*fyne.Container)

            wUsername := wTopOfMsg.Objects[2].(*widget.Label)
            wUsername.SetText(*msg.Username)

            wTime := wTopOfMsg.Objects[3].(*widget.Label)
            wTime.SetText(time.Unix(msg.Time, 0).Format("3:04 PM"))

            wContent := card.Content.(*fyne.Container).Objects[1].(*widget.RichText)
            wContent.ParseMarkdown(msg.Content)
        },
    )
}

func createMessageCardTemplate() *fyne.Container {
    labelUsername := widget.NewLabel("")
    labelUsername.Wrapping = fyne.TextWrapWord
    labelUsername.TextStyle.Bold = true

    // Base VBox
    return container.NewVBox(
        // VBox for replied messages
        container.NewVBox(),

        // Main message card
        widget.NewCard("", "",
            // VBox for contents of message card
            container.NewVBox(
                // Border container for message metadata
                container.NewBorder(
                    nil, nil,
                    // Username label (left of border)
                    labelUsername,
                    // Time label (right of border)
                    widget.NewLabel(""),
                    nil,
                ),

                // Message content
                widget.NewRichTextWithText(""),
            ),
        ),
    )
}
