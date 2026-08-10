package schat

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"

    "github.com/YhhVTF/ping-msg/chat"
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

func (g *ScreenChat) RespAdd(r *prot.ChatResponse, c *chat.ChatCache, u *user.UserCache) {
    chatCache := c.Chats[r.ChatID]

    for _, msg := range r.Messages {
        // Update chat cache with the data from the new message
        chatCache.Messages[msg.ID] = &msg
        msgCache := chatCache.Messages[msg.ID]

        chatCache.MessagesBind[msg.ID] = binding.BindString(&msgCache.Content)
        if len(msg.RepliedIDs) != 0 {
            chatCache.NewReply(msg.ID, msg.RepliedIDs)
        }

        // Create new message widget if the chat involved is currently on screen
        if r.ChatID == c.ThisChat.Metadata.ID {
            msgWidget := createMessage(g, msgCache, chatCache.MessagesBind[msg.ID], u)
            g.Widgets.Messages[msg.ID] = msgWidget
            g.Containers.Chat.VBox.Add(msgWidget.Base)

            g.Containers.Chat.VBox.Refresh()
            g.Containers.Chat.VScroll.ScrollToBottom()
        }
    }
}

func (g *ScreenChat) RespDel(r *prot.ChatResponse, c *chat.ChatCache) {
    // Delete corresponding message widget if it exists
    if msgWidget, exists := g.Widgets.Messages[r.MessageID]; exists &&
    r.ChatID == c.ThisChat.Metadata.ID {
        // Replace replied message widget text
        if repliedMsg, exists := g.Widgets.RepliedMessages[r.MessageID]; exists {
            repliedMsg.Text.Text = "Message deleted"
        }
        // Deallocate message widget
        msgWidget.Base.Hide()
        delete(g.Widgets.Messages, r.MessageID)
        g.Containers.Chat.VScroll.Refresh()
    }
    // Delete message in cache
    chatCache := c.Chats[r.ChatID]
    delete(chatCache.MessagesBind, r.MessageID)
    delete(chatCache.Messages, r.MessageID)
    //// If message is not cached...
    //if msgWidget, exists := g.Widgets.Messages[r.MessageID]; exists &&
    //r.ChatID == c.ThisChat.Metadata.ID {
    //    // Replace replied message widget text
    //    if repliedMsg, exists := g.Widgets.RepliedMessages[r.MessageID]; exists {
    //        repliedMsg.Text.Text = "Message deleted"
    //    }
    //    // Deallocate message widget
    //    msgWidget.Base.Hide()
    //    delete(g.Widgets.Messages, r.MessageID)
    //    g.Containers.Chat.VScroll.Refresh()
    //// If message is cached...
    //} else if _, exists := c.Chats[r.ChatID].Messages[r.MessageID]; exists {
    //    // Delete message in cache
    //    delete(c.Chats[r.ChatID].Messages, r.MessageID)
    //}
}

func (g *ScreenChat) RespEdit(r *prot.ChatResponse, c *chat.ChatCache, u *user.UserCache) {
    // Edit message in cache
    chatCache := c.Chats[r.ChatID]
    chatCache.MessagesBind[r.MessageID].Set(r.Messages[0].Content)

    // Update replied message widget for message if there is one
    if repliedMsgWidget, exists := g.Widgets.RepliedMessages[r.MessageID]; exists &&
    r.ChatID == c.ThisChat.Metadata.ID {
        repliedMsgWidget.Text.Text = "Message deleted"
    }

    //// If the message to be edited has an initialized corresponding widget, edit the widget, if the message widget is not initialized, edit the message cache, if the message is not loaded in any form, do nothing
    //if msgWidget, exists := g.Widgets.Messages[r.MessageID]; exists &&
    //r.ChatID == c.ThisChat.Metadata.ID {
    //    msgWidget.Content.SetText(r.Messages[0].Content)
    //    msgWidget.Content.Refresh()
//
    //    if repliedMsgWidget, exists := g.Widgets.RepliedMessages[r.MessageID]; exists {
    //        repliedMsgWidget.Text.Text = fmt.Sprintf("%s: %s", u.Users[
    //    }
    //} else if msg, exists := c.Chats[r.ChatID].Messages[r.MessageID]; exists {
    //    // Set the content of the message in cache to the new content
    //    msg.Content = r.Messages[0].Content
    //}
}
