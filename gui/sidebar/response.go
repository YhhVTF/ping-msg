package sside

import (
    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/protocol"
)

func (g *ScreenSidebar) ChatMetadataRespGet(
    resp prot.ChatMetadataResponse, c *chat.ChatCache,
) {
    for _, chatMD := range resp.Metadata {
        c.Chats[chatMD.ID] = chat.NewChat(&chatMD)
    }
}
