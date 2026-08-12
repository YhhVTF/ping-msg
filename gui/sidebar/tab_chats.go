package sside

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "fmt"

    "github.com/YhhVTF/ping-msg/chat"
)

// Chat card - A card with information about a chat. It will open said chat in the chat screen upon being clicked

func createChatCardTemplate() *widget.Card {
    c := container.NewVBox(widget.NewLabel(""))
    return widget.NewCard("", "", c)
}

func updateChatCard(chatCard *widget.Card, chatCache *chat.Chat) {
    chatCard.Content.(*fyne.Container).Objects[0].(*widget.Label).SetText(
        fmt.Sprintf("%d", chatCache.Metadata.ID),
    )
}
