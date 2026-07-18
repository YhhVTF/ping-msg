package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

// Messages from chats that aren't currently shown are stored here
type ChatCache struct {
    // Key (int) - Chat ID
    // Val ([]string) - List of cached messages from that chat
    CachedMessages map[int][]string
}

// Wrapper struct containing all containers that compose the chat section of the screen
type Chat struct {
    // Base container of the chat section, is a stack to allow children to fill all avaliable space
    Base *fyne.Container
    // Makes child scrollable
    VScroll *container.Scroll
    // Makes child fill all avaliable space
    Stack2 *fyne.Container
    // Used to position messages at the bottom of the section when there aren't enough to fill the whole section
    Border *fyne.Container
    // Lists messages vertically
    VBox *fyne.Container
}

type Message struct {}

// NewChat: Creates and returns all containers needed for the section of the GUI
// Returns:
//  ChatContainer - Wrapper struct containing all containers that compose the chat section
func NewChat() Chat {
    c := Chat{}

    c.VBox = container.NewVBox()
    c.Border = container.NewBorder(nil, c.VBox, nil, nil)
    c.Stack2 = container.NewStack(c.Border)
    c.VScroll = container.NewVScroll(c.Stack2)
    c.Base = container.NewStack(c.VScroll)

    return c
}

func NewMessage(content string, username string) *widget.Card {
    card := widget.NewCard("", "", widget.NewLabel("buh"))
    return card
}
