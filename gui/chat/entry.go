package schat

import (
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/global"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func entryMessageOnSubmitted(g *ScreenChat, text string, c *chat.ChatCache, u *user.UserCache) {
	log.Info.Printf("Widget EntryMessage submitted (%s)\n", text)

    // If there's no text in the entry or Ping isn't connected to the server, don't continue
	if text == "" || !ping.Connected { return }

    // Send new ADD chat request net.serverSend
	req := prot.ChatRequest{
		ChatID:         1,
		MessageContent: text,
		MessageID:      prot.NONE_INT,
        RepliedIDs:     c.ThisChat.ReplyingTo,
		Type:           prot.REQ_ADD,
		UserID:         u.ThisUserID,
	}
	g.OutgoingRequests <- req

    // Clear the entry
	g.Widgets.EntryMessage.SetText("")
    // Reset ScreenChat.ReplyingTo
    c.ThisChat.ReplyingTo = make([]int, 0)
}

func createEntryMessage(g *ScreenChat, u *user.UserCache, opt *options.Options) *widget.Entry {
	// Initialize message entry
    entry := widget.NewEntry()
	entry.PlaceHolder = opt.GUIText.EntryMessage.Label
	entry.OnSubmitted = func(text string) {
        entryMessageOnSubmitted(g, text, ping.ChatCache, u)
	}
    return entry
}
