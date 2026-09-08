package prot

import (
    "errors"
)

func CheckChatRequest(req *ChatRequest) error {
    if req.Username == SERVER_USERNAME || req.Username == "" {
        return errors.New("Invalid username")
    }
    if req.ChatID < 0 {
        return errors.New("Invalid chat ID")
    }

    switch req.Type {
    case ChatRequestAdd:
        {}
    case ChatRequestDelete:
        if req.MessageID < 0 { return errors.New("Invalid message ID") }
    case ChatRequestEdit:
        if req.MessageID < 0 { return errors.New("Invalid message ID") }
    case ChatRequestGet:
        if req.MessageID < 0 { return errors.New("Invalid message ID") }
    default:
        return errors.New("Invalid request type")
    }
    return nil
}

func CheckChatMetadataRequest(req *ChatMetadataRequest) error {
    if req.Username == SERVER_USERNAME || req.Username == "" {
        return errors.New("Invalid username")
    }

    switch req.Type {
    case ChatMetadataRequestCreate:
        {}
    case ChatMetadataRequestGet:
        if req.ChatID == nil { return errors.New("Chat ID(s) must be provided") }
    case ChatMetadataRequestJoin:
        if req.ChatID == nil { return errors.New("Chat ID(s) must be provided") }
    case ChatMetadataRequestLeave:
        if req.ChatID == nil { return errors.New("Chat ID(s) must be provided") }
    default:
        return errors.New("Invalid request type")
    }
    return nil
}

func CheckUserRequest(req *UserRequest) error {
    if req.Username == SERVER_USERNAME || req.Username == "" {
        return errors.New("Invalid username")
    }

    switch req.Type {
    case UserRequestRegister:
        {}
    default:
        return errors.New("Invalid request type")
    }
    return nil
}
