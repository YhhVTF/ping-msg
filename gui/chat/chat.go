package schat

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

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

func (g *ScreenChat) RespAdd(r *prot.ChatResponse, u *user.UserCache) {
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
