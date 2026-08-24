package chat

import (
    "fyne.io/fyne/v2/data/binding"

    "io"

    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

// Cache and data for a chat
type Chat struct {
    // Readers for files to be attached to the next message the client sends in this chat
    Attachments     []io.Reader
    // Data for each message in the chat
    //  Key (int) - Message ID
    //  Val (Message) - Message data
    Messages        map[int]*Message
    MessagesBind    map[int]binding.String
    // Chat metadata
    Metadata        prot.ChatMetadata
    ReplyingTo      []int
}

type ChatCache struct {
    // Cache and data for each loaded chat
    //  Key (int) - Chat ID
    //  Val (*Chat) - Chat cache and data
    Chats       map[int]*Chat
    // Cache and data for chat currently shown on screen, is null if there is no chat on screen
    ThisChat    *Chat
}

type Message struct {
    Content     string
    ID          int
    RepliedIDs  []int
    Time        int64
    Username    *string
}

func (chat *Chat) CacheMessages(messagesRaw []prot.MessageRaw, u *user.UserCache) {
    // Go through each raw message provided
    for _, msgRaw := range messagesRaw {
        // Initialize the cache for that message if it isn't already and fill in the data
        if _, exists := chat.Messages[msgRaw.ID]; !exists {
            chat.Messages[msgRaw.ID] = &Message{
                Content:    msgRaw.Content,
                ID:         msgRaw.ID,
                RepliedIDs: msgRaw.RepliedIDs,
                Time:       msgRaw.Time,
                Username:   &u.Users[msgRaw.Username].Username,
            }
            chat.MessagesBind[msgRaw.ID] =
                binding.BindString(&chat.Messages[msgRaw.ID].Content)
        // If the message is already present in cache, assign the new Content and RepliedIDs to it
        } else {
            chat.MessagesBind[msgRaw.ID].Set(msgRaw.Content)
            chat.Messages[msgRaw.ID].RepliedIDs = msgRaw.RepliedIDs
        }
    }
}

func NewChat() *Chat {
    return &Chat{
        Attachments:        make([]io.Reader, 0),
        Messages:           make(map[int]*Message),
        MessagesBind:       make(map[int]binding.String),
        Metadata:           prot.ChatMetadata{ ID: 1, },
        ReplyingTo:         make([]int, 0),
    }
}

func NewChatCache() *ChatCache {
    return &ChatCache{
        Chats: make(map[int]*Chat),
    }
}
