package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
)

func createListMessages(chatCache *chat.Chat) *widget.List {
    return widget.NewList(
        func() int {
            return len(chatCache.Messages)
        },
        func() fyne.CanvasObject {
            return createMessageCardTemplate()
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
