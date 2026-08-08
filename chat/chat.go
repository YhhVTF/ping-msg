package chat

import (
    "io"

    "github.com/YhhVTF/ping-msg/protocol"
)

// Cache and data for a chat
type Chat struct {
    // Readers for files to be attached to the next message the client sends in this chat
    Attachments []io.Reader
    // Data for each message in the chat
    //  Key (int) - Message ID
    //  Val (*prot.MessageRaw) - Message data
    Messages    map[int]*prot.MessageRaw
    // Chat metadata
    Metadata    prot.ChatMetadata
    // 
    // IDs of messages that the next message the client sends in this chat will reply to
    ReplyingTo  []int
}

type ChatCache struct {
    // Cache and data for each loaded chat
    //  Key (int) - Chat ID
    //  Val (*Chat) - Chat cache and data
    Chats       map[int]*Chat
    // Cache and data for chat currently shown on screen, is null if there is no chat on screen
    ThisChat    *Chat
}

func NewChatCache() *Chat {
    return &Chat{
        Attachments:    make([]io.Reader, 0),
        Messages:       make(map[int]*prot.MessageRaw),
        ReplyingTo:     make([]int, 0),
    }
}
