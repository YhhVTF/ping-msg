package sside

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "fmt"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/gui/screen"
)

// Chat card - A card with information about a chat. It will open said chat in the chat screen upon being clicked

func createChatCardTemplate() *widget.Card {
    c := container.NewVBox(widget.NewLabel(""), widget.NewButton("->", func(){}))
    return widget.NewCard("", "", c)
}

func updateChatCard(
    chatCard *widget.Card, chatCache *chat.Chat,
    c *chat.ChatCache, s *screen.ScreenManager,
) {
    chatCard.Content.(*fyne.Container).Objects[0].(*widget.Label).SetText(
        fmt.Sprintf("%d", chatCache.Metadata.ID),
    )
    chatCard.Content.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
        defer s.ScreenChatFocusDefault()

        if c.ThisChat == nil {
            c.ThisChat = chatCache
            s.ScreenChatFull()
            return
        }
        if chatCache.Metadata.ID == c.ThisChat.Metadata.ID { return }
        c.ThisChat = chatCache
        s.ScreenChatFull()
    }
}
