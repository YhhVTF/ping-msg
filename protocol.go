package main

var PROTOCOL_VER_LATEST = "1"

// What action is being requested. What these do differ depending on the RequestWhere
type RequestWhat int
// Adds new data provided by the client
//  REQ_CHAT - Adds a message to the chat
//  REQ_CHATMETADATA - Creates a new chat
//  REQ_USER - Registers a new user
const REQ_ADD   RequestWhat = 0
// Deletes existing data
//  REQ_CHAT - Deletes the a message
//  REQ_CHATMETADATA - Deletes a chat
//  REQ_USER - Deletes a user
const REQ_DEL   RequestWhat = 1
// Edits existing data
//  REQ_CHAT - Edits a message
//  REQ_CHATMETADATA - Edits a chat's metadata
//  REQ_USER - Edits user information
const REQ_EDIT  RequestWhat = 2
// Sends existing data to the client that requested it
//  REQ_CHAT - Sends a chat block to the client
//  REQ_CHATMETADATA - Sends a chat's metadata to the client
//  REQ_USER - Sends user information to the client
const REQ_GET  RequestWhat = 3

// What data is the request asking for the action to be done upon
type RequestWhere int
// Request for chat
const REQ_CHAT          RequestWhere = 0
// Request for chat metadata
const REQ_CHATMETADATA  RequestWhere = 1
// Request for user
const REQ_USER          RequestWhere = 2

// Collection of raw messages, used for sending and receicing messages to and from the server
type ChatRaw struct {
    // Server version the chat was encoded in
    Version string          `json:"version"`
    // Array of all messages in chat
    Messages []MessageRaw   `json:"messages"`
}

// A request to change something about or get information from a chat (e.g., add a message, delete a message, load and receive a message)
type ChatRequest struct {
    // ID of the chat involved
    ChatID int              `json:"chat_id"`
    // Content to be assigned to the given message for editing and adding messages, is empty for `REQ_DEL` and `REQ_GET`
    MessageContent string   `json:"content"`
    // ID of the message involved, is -1 for `REQ_ADD`
    MessageID int           `json:"message_id"`
    // What the request is (e.g., message deletion, editing)
    Type RequestWhat        `json:"chat_request_type"`
    // Username of who sent the request
    Username string         `json:"username"`
}

// A response to a chat request after the request is fulfuilled
type ChatResponse struct {
    // ID of the chat involved
    ChatID int          `json:"chat_id"`
    // An error that prevented the request from being fulfilled. Is empty if no error occurred
    Error string        `json:"error"`
    // Messages that the client may have requested
    Messages ChatRaw    `json:"messages"`
    // Action that this response fulfilled
    Type RequestWhat    `json:"chat_response_type"`
}

// Message data sent to and received from the server
type MessageRaw struct {
    // What the message says
    Content string  `json:"content"`
    // Unique message identifier (assigned incrementally, is -1 when a message is received from a client until the server assigns it an ID)
    ID int          `json:"id"`
    // Time at which the message was sent (unix format)
    Time int64      `json:"time"`
    // Username of who sent the message
    Username string `json:"username"`
}
