package sside

import (
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func (g *ScreenSidebar) ChatMetadataRequestCreate(u *user.UserCache) {
    g.OutgoingReqChatMetadata <- prot.ChatMetadataRequest{
        ChatID:     nil,
        Type:       prot.ChatMetadataRequestCreate,
        Username:   u.ThisUser.Username,
    }
}

func (g *ScreenSidebar) ChatMetadataRequestGet(chatIDs []int, u *user.UserCache) {
    g.OutgoingReqChatMetadata <- prot.ChatMetadataRequest{
        ChatID:     chatIDs,
        Type:       prot.ChatMetadataRequestGet,
        Username:   u.ThisUser.Username,
    }
}

func (g *ScreenSidebar) ChatMetadataRequestJoin(chatID int, u *user.UserCache) {
    g.OutgoingReqChatMetadata <- prot.ChatMetadataRequest{
        ChatID:     []int{ chatID },
        Type:       prot.ChatMetadataRequestJoin,
        Username:   u.ThisUser.Username,
    }
}

func (g *ScreenSidebar) MemberRequestGet(u *user.UserCache) {
    g.OutgoingReqMember <- prot.MemberRequest{
        ChatID:     prot.NONE_INT,
        Type:       prot.REQ_GET,
        Username:   u.ThisUser.Username,
    }
}

func (g *ScreenSidebar) MemberRequestToggle(chatID int, u *user.UserCache) {
    g.OutgoingReqMember <- prot.MemberRequest{
        ChatID:     chatID,
        Type:       prot.REQ_GET,
        Username:   u.ThisUser.Username,
    }
}
