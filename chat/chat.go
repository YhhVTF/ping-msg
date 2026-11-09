package chat

import (
    "fyne.io/fyne/v2/data/binding"

    "errors"
    "fmt"
    "io"
    "slices"

    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

// Cache and data for a chat
type Chat struct {
    // Readers for files to be attached to the next message the client sends in this chat
    Attachments     []io.Reader
    // List of all islands. Earlier islands are earlier in the last and so for later ones
    Islands         []*Island
    // Data for each message in the chat
    //  Key (int) - Message ID
    //  Val (Message) - Message data
    Messages        map[int]*Message
    MessagesBind    map[int]binding.String
    // Chat metadata
    Metadata        ChatMetadata
    // Islands[RenderedIsland] = The island which contains messages currently being rendered
    RenderedIsland  int
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

// An island is one continuous line of loaded messages, defined by their earlier and later bounds. Two groups of loaded messages separated by unloaded messages between the former's later bound and the latter's earlier bound are two separate islands
type Island struct {
    // ID of the earliest message in the island
    EarlierBound    int
    // ID of the latest message in the island
    LaterBound      int
    // List of IDs of the messages in an island (earlier bound = MessageIDs[0], later bound = MessageIDs[len(MessageIDs)-1])
    MessageIDs      []int
}

// Merges earlierIsland and laterIsland into one island by connecting them with bridgeMsgID. bridgeMsgID should be one more than earlierIsland's later bound and one less than laterIsland's earlier bound
func MergeIslands(earlierIsland *Island, bridgeMsgID int, laterIsland *Island) *Island {
    return &Island{
        EarlierBound:   earlierIsland.EarlierBound,
        LaterBound:     laterIsland.LaterBound,
        MessageIDs:     slices.Concat(earlierIsland.MessageIDs, []int{ bridgeMsgID }, laterIsland.MessageIDs),
    }
}

// Organizes the IDs of msgsRaw into islands
func (chat *Chat) SortIntoIslands(msgsRaw []prot.MessageRaw) {
    for _, msgRaw := range msgsRaw {
        // Move to next message if the message is already in cache
        if _, ok := chat.Messages[msgRaw.ID]; ok { continue }

        // For each island do...
        for i, island := range chat.Islands {
            // If message ID is one less than this island's earlier bound...
            if msgRaw.ID == island.EarlierBound - 1 {
                // And is also one more than previous island's later bound, merge these two islands with the message ID as the bridge
                if i != 0 { if chat.Islands[i-1].LaterBound == msgRaw.ID + 1 {
                    // Replace one of the halves of the new islands with the merged one
                    chat.Islands[i] = MergeIslands(chat.Islands[i-1], msgRaw.ID, island)
                    // Remove the other half
                    chat.Islands = slices.Delete(chat.Islands, i-1, i)
                    if chat.RenderedIsland == i-1 {
                        chat.RenderedIsland = i
                    }
                // Otherwise, add this message ID to the early end of the island
                }} else {
                    island.EarlierBound--
                    island.MessageIDs = slices.Concat([]int{ msgRaw.ID }, island.MessageIDs)
                }
                break
            // At the last island, if the message ID is one more than the last island's later bound, add it to the late end of the island
            } else if msgRaw.ID == island.LaterBound + 1 && i == len(chat.Islands)-1 {
                island.LaterBound++
                island.MessageIDs = append(island.MessageIDs, msgRaw.ID)
            }
        }
    }
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
        Islands:            make([]*Island, 0),
        Messages:           make(map[int]*Message),
        MessagesBind:       make(map[int]binding.String),
        Metadata:           metadata,
        RenderedIsland:     prot.NO_ID,
        ReplyingTo:         make([]int, 0),
    }
}

func NewChatCache() *ChatCache {
    return &ChatCache{
        Chats: make(map[int]*Chat),
    }
}
