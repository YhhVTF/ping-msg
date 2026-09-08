package net

import (
    "fyne.io/fyne/v2"

    "encoding/json"
    "slices"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/gui"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func handleChatMDResponse(resp *prot.ChatMetadataResponse, gui *gui.GUI, c *chat.ChatCache, u *user.UserCache) {
    switch resp.Type {
    case prot.ChatMetadataRequestCreate:
        c.Chats[resp.ChatID[0]] = chat.NewChat(&resp.Metadata[0])
        u.ThisUser.MemberOf = append(u.ThisUser.MemberOf, resp.ChatID[0])

        fyne.Do(func() { gui.Sidebar.ChatMetadataRespCreate(u) })

    case prot.ChatMetadataRequestGet:
        for _, chatMD := range resp.Metadata {
            c.Chats[chatMD.ID] = chat.NewChat(&chatMD)
        }
        fyne.Do(func() { gui.Sidebar.Widgets.ChatsList.Refresh() })

    case prot.ChatMetadataRequestJoin:
        c.Chats[resp.ChatID[0]] = chat.NewChat(&resp.Metadata[0])
        u.ThisUser.MemberOf = append(u.ThisUser.MemberOf, resp.ChatID[0])

        fyne.Do(func() { gui.Sidebar.ChatMetadataRespJoin(u) })

    case prot.ChatMetadataRequestLeave:
        delete(c.Chats, resp.ChatID[0])
        i := slices.Index(u.ThisUser.MemberOf, resp.ChatID[0])
        u.ThisUser.MemberOf = slices.Delete(u.ThisUser.MemberOf, i, i+1)

        fyne.Do(func() { gui.Sidebar.Widgets.ChatsList.Refresh() })
    }
}

func handleChatResponse(
    resp *prot.ChatResponse, gui *gui.GUI, c *chat.ChatCache, u *user.UserCache, opt *options.Options,
) {
    switch resp.Type {
    case prot.ChatRequestAdd:
        u.CacheUserFront(resp.Users) // Cache the usernames of users involved
        fyne.Do(func() { gui.Chat.RespAdd(resp, c, u, opt) })
    case prot.ChatRequestDelete:
        fyne.Do(func() { gui.Chat.RespDel(resp, c, opt) })
    case prot.ChatRequestEdit:
        fyne.Do(func() { gui.Chat.RespEdit(resp, c, u) })
    }
}

func serverRecieve(
    decoder *json.Decoder, gui *gui.GUI, c *chat.ChatCache, u *user.UserCache,
    opt *options.Options, signalDone func(),
) {
	for {
        var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			signalDone()
			return
		}

        var chatResp prot.ChatResponse
        if err := json.Unmarshal(raw, &chatResp); err == nil && chatResp.Type != "" {
            log.Info.Printf("Received %s from server\n", chatResp.Type)

            if chatResp.Error != prot.NO_ERR {
                log.Error.Printf("Server returned error: %s\n", chatResp.Error)
                continue
            }
            handleChatResponse(&chatResp, gui, c, u, opt)
        }
        var userResp prot.UserResponse
        if err := json.Unmarshal(raw, &userResp); err == nil && userResp.Type != "" {
            log.Info.Printf("Received %s from server\n", userResp.Type)

            if userResp.Error != prot.NO_ERR { 
                log.Error.Printf("Server returned error: %s\n", chatResp.Error)
                continue
            }
        }
        var chatMDResp prot.ChatMetadataResponse
        if err := json.Unmarshal(raw, &chatMDResp); err == nil && chatMDResp.Type != "" {
            log.Info.Printf("Received %s from server\n", chatMDResp.Type)

            if chatMDResp.Error != prot.NO_ERR {
                log.Error.Printf("Server returned error: %s\n", chatResp.Error)
                continue
            }
            handleChatMDResponse(&chatMDResp, gui, c, u)
        }
	}
}
