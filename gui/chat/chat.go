package schat

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

// Messages from chats that aren't currently shown are stored here
type ChatCache struct {
	// Key (int) - Chat ID
	// Val ([]MessageRaw) - List of cached messages from that chat
	CachedMessages map[int][]prot.MessageRaw
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

type Message struct {
    // Surrounds message in a card
    Base *widget.Card
    // Base VBox container for card content
    VBox *fyne.Container
    // Border container for message metadata, added to VBox first
    Border *fyne.Container
    // Label for message content, added to VBox second
    Content *widget.Label
    // Label for username, added to left side of Border
    Username *widget.Label
    // Label for time, added to right side of Border
    Time *widget.Label
}

// NewChat: Creates and returns all containers needed for the section of the GUI
// Returns:
//
//	ChatContainer - Wrapper struct containing all containers that compose the chat section
func NewChat() Chat {
	c := Chat{}

	c.VBox = container.NewVBox()
	c.Border = container.NewBorder(nil, c.VBox, nil, nil)
	c.Stack2 = container.NewStack(c.Border)
	c.VScroll = container.NewVScroll(c.Stack2)
	c.Base = container.NewStack(c.VScroll)

	return c
}

func (g *ScreenChat) RespAdd(r *prot.ChatResponse, u *user.UserData) {
    for _, msg := range r.Messages {
        msgWidget := createMessage(g, msg, u)
        g.Widgets.Messages[msg.ID] = msgWidget
        g.Containers.Chat.VBox.Add(msgWidget.Base)
    }
    g.Containers.Chat.VBox.Refresh()
    g.Containers.Chat.VScroll.ScrollToBottom()
}

func (g *ScreenChat) RespDel(r *prot.ChatResponse) {
    if _, exists := g.Widgets.Messages[r.MessageID]; exists {
        g.Widgets.Messages[r.MessageID].Base.Hide()
        delete(g.Widgets.Messages, r.MessageID)
        g.Containers.Chat.VScroll.Refresh()
    }
}

func (g *ScreenChat) RespEdit(r *prot.ChatResponse) {
    if _, exists := g.Widgets.Messages[r.MessageID]; exists {
        g.Widgets.Messages[r.MessageID].Content.Text =
        r.Messages[0].Content
        g.Widgets.Messages[r.MessageID].Content.Refresh()
    }
}
