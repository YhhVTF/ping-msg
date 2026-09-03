package chat

import (
    "fyne.io/fyne/v2/data/binding"

    "errors"
    "fmt"
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
    Metadata        ChatMetadata
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

type ChatMetadata struct {
	// A description of what the chat is about
	Description         string
	// ID of the chat
	ID                  int
	// Is the chat public?
	IsPublic            bool
	// ID of the last message in the chat. Should be set to -1 and not 0 when creating a chat to prevent and off by 1 error
	LastMessageID       int
	// Usernames of everyone who has access to the chat
	Members             []string
	// The name of the chat
	Name                string
	// Total number of messages in the chat. This should not be updated while the chat is loaded and instead should be calculated from the field `LastMessageID` when saving metadata
	NumberOfMessages    int
}

type Message struct {
    Content     string
    ID          int
    RepliedIDs  []int
    Time        int64
    Username    *string
}

func (chat *Chat) CacheMessages(messagesRaw []prot.MessageRaw, u *user.UserCache) error {
    // Go through each raw message provided
    for _, msgRaw := range messagesRaw {
        if _, exists := u.Users[msgRaw.Username]; !exists {
            return errors.New(
                fmt.Sprintf("User %s (sender of message %d) not present in user cache", msgRaw.Username, msgRaw.ID),
            )
        }
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
    return nil
}

func NewChat(metadataRaw *prot.ChatMetadataRaw) *Chat {
    metadata := ChatMetadata{
        Description:        metadataRaw.Description,
        ID:                 metadataRaw.ID,
        IsPublic:           metadataRaw.IsPublic,
        LastMessageID:      metadataRaw.LastMessageID,
        Members:            metadataRaw.Members,
        Name:               metadataRaw.Name,
        NumberOfMessages:   metadataRaw.NumberOfMessages,
    }

    return &Chat{
        Attachments:        make([]io.Reader, 0),
        Messages:           make(map[int]*Message),
        MessagesBind:       make(map[int]binding.String),
        Metadata:           metadata,
        ReplyingTo:         make([]int, 0),
    }
}

func NewChatCache() *ChatCache {
    return &ChatCache{
        Chats: make(map[int]*Chat),
    }
}
