package sside

import (
    "fyne.io/fyne/v2/widget"

    "fmt"

    "github.com/YhhVTF/ping-msg/chat"
)

// A card with information about a chat. It will open said chat in the chat screen upon being clicked
type ChatCard struct {
    Base *widget.Card
}

func createChatCard(chatCache *chat.Chat) *ChatCard {
    chatCard := &ChatCard{}

    chatCard.Base = widget.NewCard("", fmt.Sprintf("%d", chatCache.Metadata.ID), widget.NewLabel("buh :P"))

    return chatCard
}
