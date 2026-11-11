package chat

import (
    "fyne.io/fyne/v2/data/binding"

    "errors"
    "fmt"
    "io"
    "slices"

    "github.com/YhhVTF/ping-msg/log"
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

// Organizes the IDs of newMsgRaws into islands of chat
func (chat *Chat) SortIntoIslands(newMsgRaws []prot.MessageRaw) error {
    newMsgIDs := make([]int, len(newMsgRaws))

    // Check if newMsgRaws a continuous line of messages, with the IDs incrementing by 1 as the index of the slice does the same. Return an error if it doesn't match this pattern
    for i, msgRaw := range newMsgRaws {
        newMsgIDs[i] = msgRaw.ID

        if i == 0 { continue }

        if msgRaw.ID - newMsgRaws[i-1].ID != 1 {
            return errors.New("Raw messages are not in continuous ascending order")
        }
    }

    msgRawsEarlierBound := newMsgRaws[0].ID
    msgRawsLaterBound := newMsgRaws[len(newMsgRaws)-1].ID

    // Check if any of newMsgRaws are already in cache, and return an error if they are
    _, earlierPresent := chat.Messages[msgRawsEarlierBound]
    _, laterPresent := chat.Messages[msgRawsLaterBound]
    if earlierPresent || laterPresent {
        return errors.New("Part or all of the raw messages are already loaded")
    }

    // If there are currently no islands, simply create one for the new message IDs
    if len(chat.Islands) == 0 {
        newIsland := &Island{
            EarlierBound:   msgRawsEarlierBound,
            LaterBound:     msgRawsLaterBound,
            MessageIDs:     newMsgIDs,
        }
        chat.Islands = []*Island{ newIsland }
        chat.RenderedIsland = 0

        log.Info.Printf("Chat %d: new island %d: %d\n", chat.Metadata.ID, 0, newIsland.MessageIDs)
        return nil
    }

    // For each island do...
    for i, island := range chat.Islands {
        var nextIsland *Island
        if i != len(chat.Islands) - 1 { nextIsland = chat.Islands[i+1] }

        // If the earlier bound of newMsgRaws is one less than this island's later bound...
        if msgRawsEarlierBound == island.LaterBound + 1 {
            // And the later bound is one less than the next island's earlier bound, merge the two islands with the new message IDs in the middle
            if i != len(chat.Islands)-1 { if msgRawsLaterBound == nextIsland.EarlierBound - 1 {
                mergedIsland := &Island{
                    EarlierBound:   island.EarlierBound,
                    LaterBound:     nextIsland.LaterBound,
                    MessageIDs:     slices.Concat(island.MessageIDs, newMsgIDs, nextIsland.MessageIDs),
                }
                // Replace this island with the merged island and remove the old next island
                chat.Islands[i] = mergedIsland
                chat.Islands = slices.Delete(chat.Islands, i+1, i+2)

                log.Info.Printf("Chat %d: new merged island %d: %d\n", chat.Metadata.ID, i, mergedIsland.MessageIDs)
                break
            }}
            // Otherwise, add the new message IDs to the later end of this island
            island.LaterBound = msgRawsLaterBound
            island.MessageIDs = slices.Concat(island.MessageIDs, newMsgIDs)
            
            log.Info.Printf("Chat %d: expanded island %d: %d\n", chat.Metadata.ID, i, island.MessageIDs)
            break
        // Or if not, and the earlier bound of newMsgRaws is still more than the later bound of this island, and the later bound is still less than the next island's earlier bound, create a new island for the new IDs between them
        } else if msgRawsEarlierBound > island.LaterBound {
            newIsland := &Island{
                EarlierBound:   msgRawsEarlierBound,
                LaterBound:     msgRawsLaterBound,
                MessageIDs:     newMsgIDs,
            }
            if i != len(chat.Islands)-1 { if msgRawsLaterBound < nextIsland.EarlierBound {
                chat.Islands = slices.Insert(chat.Islands, i+1, newIsland)

                log.Info.Printf("Chat %d: new island %d: %d\n", chat.Metadata.ID, i+1, newIsland.MessageIDs)
                break
            }}
            chat.Islands = append(chat.Islands, newIsland)

            log.Info.Printf("Chat %d: new island %d: %d\n", chat.Metadata.ID, i+1, newIsland.MessageIDs)
            break
        }

        // If the later bound of newMsgRaws is one less than the earlier bound of this island, add the new message IDs to the earlier end of this island
        if msgRawsLaterBound == island.EarlierBound - 1 {
            island.EarlierBound = msgRawsEarlierBound
            island.MessageIDs = slices.Concat(newMsgIDs, island.MessageIDs)

            log.Info.Printf("Chat %d: expanded island %d: %d\n", chat.Metadata.ID, i, island.MessageIDs)
            break
        // If not, but the later bound of newMsgRaws is still less than the earlier bound of this island, create a new island earlier than this one for the new message IDs
        } else if i == 0 && msgRawsLaterBound < island.EarlierBound {
            newIsland := &Island{
                EarlierBound:   msgRawsEarlierBound,
                LaterBound:     msgRawsLaterBound,
                MessageIDs:     newMsgIDs,
            }
            chat.Islands = slices.Concat([]*Island{ newIsland }, chat.Islands)

            log.Info.Printf("Chat %d: new island %d: %d\n", chat.Metadata.ID, i, newIsland.MessageIDs)
            break
        }
    }
    return nil
}

type Message struct {
    Content     string
    ID          int
    RepliedIDs  []int
    Time        int64
    Username    *string
}

func (chat *Chat) CacheMessages(messagesRaw []prot.MessageRaw, u *user.UserCache) error {
    // Sort new messages into islands
    err := chat.SortIntoIslands(messagesRaw)
    if err != nil {
        return err
    }

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
