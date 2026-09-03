package sside

import (
    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func (g *ScreenSidebar) ChatMetadataRespCreate(u *user.UserCache) {
    newChatCardItemID := len(u.ThisUser.MemberOf)-1
    g.Widgets.ChatsList.Refresh()
    g.Widgets.ChatsList.ScrollTo(newChatCardItemID)
    g.Widgets.ChatsList.Select(newChatCardItemID)
}

func (g *ScreenSidebar) ChatMetadataRespGet(
    resp prot.ChatMetadataResponse, c *chat.ChatCache,
) {
    for _, chatMD := range resp.Metadata {
        c.Chats[chatMD.ID] = chat.NewChat(&chatMD)
    }
}

func (g *ScreenSidebar) ChatMetadataRespJoin(u *user.UserCache) {
    newChatCardItemID := len(u.ThisUser.MemberOf)-1
    g.Widgets.ChatsList.Refresh()
    g.Widgets.ChatsList.ScrollTo(newChatCardItemID)
    g.Widgets.ChatsList.Select(newChatCardItemID)
}
