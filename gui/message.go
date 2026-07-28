package gui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "time"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func createMessage(g *GUI, msgRaw prot.MessageRaw, u *user.UserData) Message {
    log.Info.Printf("Creating new message widget\n")

    msg := Message{}

    msg.Username = widget.NewLabel(msgRaw.Username)
	msg.Username.Wrapping = fyne.TextWrapWord
    msg.Username.TextStyle.Bold = true

    msg.Time = widget.NewLabel(time.Unix(msgRaw.Time, 0).Format("3:04 PM"))

    buttonCopy := widget.NewButton("C", func() {
        log.Info.Printf("Copied message %d\n", msgRaw.ID)

        // Give focus back to the message entry when done
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)
        // Copy the message contents
        fyne.CurrentApp().Clipboard().SetContent(msg.Content.Text)
    })

    var c *fyne.Container

    // If the message was sent by the user of this client...
    if msgRaw.Username == u.ThisUser {
        // Add a delete button
        buttonDelete := widget.NewButton("D", func() {
            messageOnDelete(g, msgRaw.ID, u)
        })

        // Add an edit button, upon pressing...
        buttonEdit := widget.NewButton("E", func() {
            messageOnEdit(g, &msg, msgRaw.ID, u)
        })
        c = container.NewHBox(buttonCopy, buttonEdit, buttonDelete, msg.Time)
    } else {
        c = container.NewHBox(buttonCopy, msg.Time)
    }

    msg.Border = container.NewBorder(nil, nil, msg.Username, c, nil)

    msg.Content = widget.NewLabel(msgRaw.Content)
    msg.Content.Wrapping = fyne.TextWrapWord

    msg.VBox = container.NewVBox(msg.Border, msg.Content)

	msg.Base = widget.NewCard("", "", msg.VBox)

	return msg
}

func messageOnDelete(g *GUI, msgID int, u *user.UserData) {
    log.Info.Printf("Delete button on message %d pressed\n", msgID)

    // Give focus back to the message entry when done
    defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

    // Send new DEL chat request to net.serverSend
    req := prot.ChatRequest{
        ChatID: 1,
        MessageContent: prot.NONE_STRING,
        MessageID: msgID,
        Type: prot.REQ_DEL,
        Username: u.ThisUser,
    }
    g.OutgoingMessages <- req
}

func messageOnEdit(g *GUI, msg *Message, msgID int, u *user.UserData) {
    log.Info.Printf("Edit button on message %d pressed\n", msgID)

    // Replace the message content label with an entry to allow editing
    msg.Content.Hide()
    editEntry := widget.NewEntry()
    editEntry.Text = msg.Content.Text
    msg.VBox.Add(editEntry)

    // Focus on the edit entry
    g.Window.Canvas().Focus(editEntry)

    // On submission of entry...
    editEntry.OnSubmitted = func(text string) {
        log.Info.Printf("Edit entry on message %d submitted\n", msgID)

        // Give focus back to the message entry when done
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

        // Send edit request if there was an actual edit
        if text != msg.Content.Text {
            req := prot.ChatRequest{
                ChatID: 1,
                MessageContent: text,
                MessageID: msgID,
                Type: prot.REQ_EDIT,
                Username: u.ThisUser,
            }
            g.OutgoingMessages <- req
        }
        // Replace the edit entry with the message content label again
        editEntry.Hide()
        msg.Content.Show()
    }
}
