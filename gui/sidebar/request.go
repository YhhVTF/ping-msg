package sside

import (
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func (g *ScreenSidebar) ChatMetadataRequestGet(chatIDs []int, u *user.UserCache) {
    g.outgoingReqChatMetadata <- prot.ChatMetadataRequest{
        ChatID:     chatIDs,
        Type:       prot.ChatMetadataRequestGet,
        Username:   u.ThisUser.Username,
    }
}

func (g *ScreenSidebar) MemberRequestGet(u *user.UserCache) {
    g.outgoingReqMember <- prot.MemberRequest{
        ChatID:     prot.NONE_INT,
        Type:       prot.REQ_GET,
        Username:   u.ThisUser.Username,
    }
}

func (g *ScreenSidebar) MemberRequestToggle(chatID int, u *user.UserCache) {
    g.outgoingReqMember <- prot.MemberRequest{
        ChatID:     chatID,
        Type:       prot.REQ_GET,
        Username:   u.ThisUser.Username,
    }
}
