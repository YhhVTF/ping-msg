package prot

const PROTOCOL_VER_LATEST = "5"

// Shows that an integer id field in a request or response is empty
const NO_ID = -1
// Shows that an error field in a request or response is empty
const NO_ERR = "##"
// Reserved username for server messages
const SERVER_USERNAME = "\n"

// What action is being requested. What these do differ depending on the RequestWhere
type RequestWhat string

// Information about a chat
type ChatMetadataRaw struct {
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

// ChatMetadataRequestCreate creates a new chat with the provided chat metadata, with the user who made the request initially being its sole member
const ChatMetadataRequestCreate = "CHATMD_CREATE"
// ChatMetadataRequestGet gives the client the metadata for the chats its user is a member of
const ChatMetadataRequestGet    = "CHATMD_GET"
// ChatMetadataRequestJoin makes the user of the client who made the request a member of the specified chat along with giving the client the metadata for that chat
const ChatMetadataRequestJoin   = "CHATMD_JOIN"
// ChatMetadataRequestLeave revokes the user of the client who made the resquest's membership of the specified chat
const ChatMetadataRequestLeave  = "CHATMD_LEAVE"
// ChatMetadataRequestStart tells the server to start the specified chat if it hasn't already
const ChatMetadataRequestStart  = "CHATMD_START"

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
    // IDs of the chats involved
    ChatID      []int               `json:"id"`
    Error       string              `json:"err"`
    // Chat metadata the client may have requested
    Metadata    []ChatMetadataRaw   `json:"metadata"`
    // Action that this response fulfilled
    Type        RequestWhat         `json:"req_what"`
}

// Collection of raw messages, used for saving and loading messages, not for communication between the client and server
type ChatRaw struct {
    // Protocol version
    Version string          `json:"ver"`
    // Array of all messages in chat
    Messages []MessageRaw   `json:"msgs"`
}

// ChatRequestAdd adds a new message to the specified chat
const ChatRequestAdd    = "CHAT_ADD"
// ChatRequestDelete deletes the desired message from the specified chat
const ChatRequestDelete = "CHAT_DEL"
// ChatRequestEdit edits to content or greetings of the desired message
const ChatRequestEdit   = "CHAT_EDIT"
// ChatRequestGet sends a block of messages from the specified chat to the client who made the request
const ChatRequestGet    = "CHAT_GET"

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
    Type RequestWhat        `json:"chatreq_what"`
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
    // Usernames of users who sent the messages in the Messages field
    Users []string          `json:"users"`
    // Action that this response fulfilled
    Type RequestWhat        `json:"chatresp_what"`
}

// Request to change user's membership status in a chat
type MemberRequest struct {
    // ID of the chat to change whether the user is a member of, is NONE_INT for REQ_GET
    ChatID      int         `json:"chat_id"`
    Type        RequestWhat `json:"memreq_what"`
    Username    string      `json:"username"`
}

// Response to MemberRequest
type MemberResponse struct {
    // IDs of the chats the user is a member of, only used with REQ_GET
    ChatIDs []int       `json:"chat_ids"`
    Error   string      `json:"err"`
    Type    RequestWhat `json:"memresp_what"`
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

// UserRaw - Includes profile data along with other data associated with the user, wwhich users can access what data is specified by the Visibility fields
type UserRaw struct {
    // String of text written by the uesr, usually to describe themselves - private by default
    Bio             string          `json:"bio"`
    // IDs of chats the user is a member of - private by default
    MemberOf        []int           `json:"memberof"`
    // Number of notifications the user has in each chat - private by default
    //  Key (int) - ChatID
    //  Val (int) - Number of notifications in the chat
    Notifications   map[int]int     `json:"notif"`
    // Whether or not the user currently has ping open - private by default
    Online          bool            `json:"online"`
    // Small string of text written by the user, generally to describe their current state - private by default
    Status          string          `json:"status"`
    // User's username - public
    Username        string          `json:"username"`
    // Defines which users can access each field of this user
    Visibility      UserVisibility  `json:"visibility"`
}

// UserRequestRegister registers a username on a new connection.
const UserRequestRegister = "USER_REG"

// UserRequest is sent before chat traffic to establish a connection identity.
type UserRequest struct {
    Type        RequestWhat `json:"userreq_what"`
    // Username of user, only used with REQ_ADD and REQ_EDIT
    Username    string      `json:"username"`
}

// UserResponse is returned for a UserRequest. Error is empty on success.
type UserResponse struct {
    Error   string      `json:"err"`
    Type    RequestWhat `json:"userresp_what"`
    User    UserRaw     `json:"user"`
}

const UserVisibilityPrivate = 0
const UserVisibilityFriendsOnly = 1
const UserVisibilityPublic = 2

// Defines which users should be able to access which fields, each can be either UserVisibilityPrivate, UserVisibilityFriendsOnly, or UserVisibilityPublic. Username is not present here as it must be public
type UserVisibility struct {
    Bio             int `json:"bio"`
    MemberOf        int `json:"memberof"`
    Notifications   int `json:"notif"`
    Online          int `json:"online"`
    Status          int `json:"status"`
}
