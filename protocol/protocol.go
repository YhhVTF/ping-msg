package prot

const PROTOCOL_VER_LATEST = "4"

// Shows that an integer field in a request or response is empty
const NONE_INT = -1
// Shows that a string field in a request or response is empty
const NONE_STRING = "##"

// What action is being requested. What these do differ depending on the RequestWhere
type RequestWhat string
// Adds new data provided by the client
//  REQ_CHAT - Adds a message to the chat
//  REQ_CHATMETADATA - Creates a new chat
//  REQ_MEMBER - Makes the user who requested it a member of the specified chat
//  REQ_USER - Registers a new user
const REQ_ADD   RequestWhat = "ADD"
// Deletes existing data
//  REQ_CHAT - Deletes the a message
//  REQ_CHATMETADATA - Deletes a chat
//  REQ_MEMBER - Removes the user's membership of the specified chat
//  REQ_USER - Deletes a user
const REQ_DEL   RequestWhat = "DEL"
// Edits existing data
//  REQ_CHAT - Edits a message
//  REQ_CHATMETADATA - Edits a chat's metadata
//  REQ_USER - Edits user information
const REQ_EDIT  RequestWhat = "EDIT"
// Sends existing data to the client that requested it
//  REQ_CHAT - Sends a chat block to the client
//  REQ_CHATMETADATA - Sends a chat's metadata to the client
//  REQ_MEMBER - Gets IDs of all the chats the user is a member of
//  REQ_USER - Sends user information to the client
const REQ_GET   RequestWhat = "GET"

// What data is the request asking for the action to be done upon
type RequestWhere string
// Request for chat
const REQ_CHAT          RequestWhere = "CHAT"
// Request for chat metadata
const REQ_CHATMETADATA  RequestWhere = "CHATMETADATA"
// Request for user's membership of chats
const REQ_MEMBER        RequestWhere = "MEMBER"
// Request for user
const REQ_USER          RequestWhere = "USER"

// Information about a chat
type ChatMetadata struct {
	// Number of messages there are in a block
	BlockSize int `json:"blk_size"`
	// Username of who created the chat
	Creator string `json:"creator"`
	// A description of what the chat is about
	Description string `json:"desc"`
	// ID of the chat
	ID int `json:"id"`
	// Is the chat public?
	IsPublic bool `json:"public"`
	// ID of the last block in the chat. Should be set to -1 and not 0 when creating a chat to prevent and off by 1 error
	LastBlockID int `json:"last_blk"`
	// ID of the last message in the chat. Should be set to -1 and not 0 when creating a chat to prevent and off by 1 error

	LastMessageID int `json:"last_msg"`
	// Usernames of everyone who has access to the chat
	Members []string `json:"members"`
	// The name of the chat
	Name string `json:"name"`
	// Total number of blocks in the chat. This should not be updated while the chat is loaded and instead should be calculated from the field `LastBlockID` when saving metadata
	NumberOfBlocks int `json:"blks"`
	// Total number of messages in the chat. This should not be updated while the chat is loaded and instead should be calculated from the field `LastMessageID` when saving metadata
	NumberOfMessages int `json:"msgs"`
}

type ChatMetadataRequest struct {
    // IDs of the chats involved
    ChatID      []int       `json:"id"`
    // New chat description
    Description string      `json:"desc"`
    // New chat viewability
    IsPublic    bool        `json:"public"`
    // New chat name
    Name        string      `json:"name"`
    // What the request is (e.g., creating new chat, getting metadata, etc.)
    Type        RequestWhat `json:"req_what"`
    // ID of user the request is from
    Username    string      `json:"username"`
}

type ChatMetadataResponse struct {
    // Chat metadata the client may have requested
    Metadata    []ChatMetadata  `json:"metadata"`
    // Action that this response fulfilled
    Type        RequestWhat     `json:"req_what"`
}

// Collection of raw messages, used for saving and loading messages, not for communication between the client and server
type ChatRaw struct {
    // Protocol version
    Version string          `json:"ver"`
    // Array of all messages in chat
    Messages []MessageRaw   `json:"msgs"`
}

// A request to change something about or get information from a chat (e.g., add a message, delete a message, load and receive a message)
type ChatRequest struct {
    // ID of the chat involved
    ChatID int              `json:"chat_id"`
    // Content to be assigned to the given message for editing and adding messages, is empty for `REQ_DEL` and `REQ_GET`
    MessageContent string   `json:"content"`
    // ID of the message involved, is -1 for `REQ_ADD`
    MessageID int           `json:"msg_id"`
    // IDs of the messages the message sent in this chat request is replying to, is non nil only for REQ_ADD and REQ_EDIT
    RepliedIDs []int        `json:"replied_ids"`
    // What the request is (e.g., message deletion, editing)
    Type RequestWhat        `json:"req_what"`
    // Username of who sent the request
    Username string         `json:"username"`
}

// A response to a chat request after the request is fulfuilled
type ChatResponse struct {
    // ID of the chat involved
    ChatID int              `json:"chat_id"`
    // An error that prevented the request from being fulfilled. Is empty if no error occurred
    Error string            `json:"err"`
    // ID of the message involved
    MessageID int           `json:"msg_id"`
    // Messages that the client may have requested
    Messages []MessageRaw   `json:"msgs"`
    // Public user profiles needed to render message senders
    Users []UserPublicRaw   `json:"users"`
    // Action that this response fulfilled
    Type RequestWhat        `json:"resp_what"`
}

// Request to change user's membership status in a chat
type MemberRequest struct {
    // ID of the chat to change whether the user is a member of, is NONE_INT for REQ_GET
    ChatID      int         `json:"chat_id"`
    Type        RequestWhat `json:"req_what"`
    Username    string      `json:"username"`
}

// Response to MemberRequest
type MemberResponse struct {
    // IDs of the chats the user is a member of, only used with REQ_GET
    ChatIDs []int       `json:"chat_ids"`
    Error   string      `json:"err"`
    Type    RequestWhat `json:"resp_what"`
}

// Message data sent to and received from the server
type MessageRaw struct {
    // What the message says
    Content string      `json:"content"`
    // Unique message identifier (assigned incrementally, is -1 when a message is received from a client until the server assigns it an ID)
    ID int              `json:"id"`
    // IDs of the message(s) this message is replying to
    RepliedIDs []int    `json:"replied_ids"`
    // Time at which the message was sent (unix format)
    Time int64          `json:"time"`
    // Username of who sent the message
    Username string     `json:"username"`
}
// UserRequestRegister registers a username on a new connection. The server
// assigns the numeric ID; clients must not generate their own IDs.
const UserRequestRegister = "REGISTER"

// UserRaw - Includes public profile data, along with other data that other users shouldn't access. A user can only get their own private data and no one else's
type UserRaw struct {
    // String of text written by the uesr, usually to describe themselves - private by default
    Bio             string      `json:"bio"`
    // IDs of chats the user is a member of - private by default
    MemberOf        []int       `json:"memberof"`
    // Number of notifications the user has in each chat - private by default
    //  Key (int) - ChatID
    //  Val (int) - Number of notifications in the chat
    Notifications   map[int]int `json:"notif"`
    // Whether or not the user currently has ping open - private by default
    Online          bool        `json:"online"`
    // Small string of text written by the user, generally to describe their current state - private by default
    Status          string      `json:"status"`
    // User's username - public
    Username        string      `json:"username"`
}

// UserRequest is sent before chat traffic to establish a connection identity.
type UserRequest struct {
    Type        RequestWhat `json:"req_what"`
    // Username of user, only used with REQ_ADD and REQ_EDIT
    Username    string      `json:"username"`
}

// UserResponse is returned for a UserRequest. Error is empty on success.
type UserResponse struct {
    Error   string  `json:"err"`
    User    UserRaw `json:"user"`
}
