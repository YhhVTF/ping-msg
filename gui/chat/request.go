package schat

import (
    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func (g *ScreenChat) ChatRequestAdd(msgContent string, c *chat.ChatCache, u *user.UserCache) {
    g.OutgoingReqChat <- prot.ChatRequest{
        ChatID:         c.ThisChat.Metadata.ID,
        MessageContent: msgContent,
        MessageID:      prot.NONE_INT,
        RepliedIDs:     c.ThisChat.ReplyingTo,
        Type:           prot.REQ_ADD,
        Username:       u.ThisUser.Username,
    }
}

func (g *ScreenChat) ChatRequestDelete(msgID int, c *chat.ChatCache, u *user.UserCache) {
    g.OutgoingReqChat <- prot.ChatRequest{
        ChatID:         c.ThisChat.Metadata.ID,
        MessageID:      msgID,
        Type:           prot.REQ_DEL,
        Username:       u.ThisUser.Username,
    }
}

func (g *ScreenChat) ChatRequestEdit(msgID int, msgContent string, c *chat.ChatCache, u *user.UserCache) {
    g.OutgoingReqChat <- prot.ChatRequest{
        ChatID:         c.ThisChat.Metadata.ID,
        MessageContent: msgContent,
        MessageID:      msgID,
        Type:           prot.REQ_EDIT,
        Username:       u.ThisUser.Username,
    }
}
