package gui

import (
    "fyne.io/fyne/v2/widget"

    "github.com/YhhVTF/ping-msg/global"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func buttonSendOnPressed(g *GUI, u *user.UserData) {
	log.Info.Printf("Widget ButtonSend pressed\n")

    // If there's no text in the message entry or Ping is not connected to the server, don't continue
	text := g.Widgets.EntryMessage.Text
	if text == "" || !ping.Connected { return }

    // Send a new ADD chat request to net.serverSend
	req := prot.ChatRequest{
		ChatID:         1,
		MessageContent: text,
		MessageID:      prot.NONE_INT,
		Type:           prot.REQ_ADD,
		Username:       u.ThisUser,
	}
	g.OutgoingMessages <- req

    // Clear the text in the message entry and give focus back to it
	g.Widgets.EntryMessage.SetText("")
    g.Window.Canvas().Focus(g.Widgets.EntryMessage)
}

func createButtonSend(g *GUI, u *user.UserData, opt *options.Options) *widget.Button {
	// Initialize send button
    button := widget.NewButton(opt.GUIText.ButtonSend.Label, func() {
        buttonSendOnPressed(g, u)
	})
    // Set importance to high
    g.Widgets.ButtonSend.Importance = widget.HighImportance

    return button
}
