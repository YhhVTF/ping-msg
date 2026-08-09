package chat

import (
    "io"

    "github.com/YhhVTF/ping-msg/protocol"
)

// Cache and data for a chat
type Chat struct {
    // Readers for files to be attached to the next message the client sends in this chat
    Attachments     []io.Reader
    // Data for each message in the chat
    //  Key (int) - Message ID
    //  Val (*prot.MessageRaw) - Message data
    Messages        map[int]*prot.MessageRaw
    // Chat metadata
    Metadata        prot.ChatMetadata
    // Cache for replied messages
    RepliedMessages *ReplyCache
    // IDs of messages that the next message the client sends in this chat will reply to
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

type RepliedMessage struct {
    ID              int
    Message         *prot.MessageRaw
    TextForWidget   string
}

type ReplyCache struct {
    // Cache for all replied messages, accessed with the ID of the replied message
    //  Key (int) - Replied message ID
    //  Val (*RepliedMessage) - Replied message cache
    RepliedMessages map[int]*RepliedMessage
    // Cache for the replied messages of each message replying to them, accessed with the ID of a message which replies to the replied message (holy friggie we need a naming convention for thissie T^T)
    //  Key (int) - ID of a message replying to the replied message
    //  Val ([]*RepliedMessage) - Cache for replied messages that the message replies to
    RepliedMessagesByReplies map[int][]*RepliedMessage
}

func NewChatCache() *Chat {
    re := &ReplyCache{
        RepliedMessages:            make(map[int]*RepliedMessage),
        RepliedMessagesByReplies:   make(map[int][]*RepliedMessage),
    }

    return &Chat{
        Attachments:        make([]io.Reader, 0),
        Messages:           make(map[int]*prot.MessageRaw),
        Metadata:           prot.ChatMetadata{ ID: 1, },
        RepliedMessages:    re,
        ReplyingTo:         make([]int, 0),
    }
}

func (re *ReplyCache) NewReply(c *Chat, replyingID int, repliedIDs []int, repliedMsgsText []string) {
    re.RepliedMessagesByReplies[replyingID] = 
        make([]*RepliedMessage, len(repliedIDs))

    for i, repliedID := range repliedIDs {
        if _, exists := re.RepliedMessages[repliedID]; !exists {
            re.RepliedMessages[repliedID] = &RepliedMessage{
                ID:             repliedID,
                Message:        c.Messages[repliedID],
                TextForWidget:  repliedMsgsText[i],
            }
        }
        re.RepliedMessagesByReplies[replyingID][i] = re.RepliedMessages[repliedID]
    }
}
