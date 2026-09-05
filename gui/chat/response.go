package schat

import (
    "fmt"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func (g *ScreenChat) RespAdd(
    resp *prot.ChatResponse, c *chat.ChatCache, u *user.UserCache, opt *options.Options,
) {
    chat := c.Chats[resp.ChatID]

    for _, msgRaw := range resp.Messages {
        log.Info.Printf("Updating message %d cache (%s)\n", resp.MessageID, resp.Type)

        // Cache the new message
        err := chat.CacheMessages(resp.Messages, u)
        if err != nil {
            log.Error.Printf("Failed to cache messages: %s\n", err)
            return
        }

        // Create new message widget if the chat involved is currently on screen
        if resp.ChatID == c.ThisChat.Metadata.ID {
            msgWidget := createMessage(
                g, chat.Messages[msgRaw.ID], chat.MessagesBind[msgRaw.ID], chat, u, opt,
            )
            g.Widgets.Messages[msgRaw.ID] = msgWidget
            g.Containers.Chat.VBox.Add(msgWidget.Base)

            g.Containers.Chat.VBox.Refresh()
            g.Containers.Chat.VScroll.ScrollToBottom()
        }
    }
}

func (g *ScreenChat) RespDel(resp *prot.ChatResponse, c *chat.ChatCache, opt *options.Options) {
    // Delete corresponding message widget if it exists
    if msgWidget, exists := g.Widgets.Messages[resp.MessageID]; exists &&
    resp.ChatID == c.ThisChat.Metadata.ID {
        // Replace replied message widget text
        if repliedMsg, exists := g.Widgets.RepliedMessages[resp.MessageID]; exists {
            repliedMsg.Text.Text = opt.GUIText.Greeting.PlaceholderDeleted
        }
        // Deallocate message widget
        msgWidget.Base.Hide()
        delete(g.Widgets.Messages, resp.MessageID)
        g.Containers.Chat.VScroll.Refresh()
    }
    log.Info.Printf("Updating message %d cache (%s)\n", resp.MessageID, resp.Type)
    // Delete message in cache
    chat := c.Chats[resp.ChatID]
    delete(chat.MessagesBind, resp.MessageID)
    delete(chat.Messages, resp.MessageID)
}

func (g *ScreenChat) RespEdit(resp *prot.ChatResponse, c *chat.ChatCache, u *user.UserCache) {
    log.Info.Printf("Updating message %d cache (%s)\n", resp.MessageID, resp.Type)
    // Edit message in cache
    chat := c.Chats[resp.ChatID]
    chat.MessagesBind[resp.MessageID].Set(resp.Messages[0].Content)

    // Update replied message widget for message if there is one
    if repliedMsgWidget, exists := g.Widgets.RepliedMessages[resp.MessageID]; exists &&
    resp.ChatID == c.ThisChat.Metadata.ID {
        repliedMsgWidget.Text.Text = fmt.Sprintf("%s: %s", u.Users[resp.Messages[0].Username].Username, chat.Messages[resp.MessageID].Content)
    }
}
