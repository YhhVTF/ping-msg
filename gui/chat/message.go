package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"

    "fmt"
    "time"

    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

// messageUsername returns the best available display name for a message.
// MessageRaw identifies senders by UserID; usernames are held separately in
// the local user cache.
func messageUsername(msgRaw prot.MessageRaw, u *user.UserCache) string {
    if u != nil {
        if msgRaw.UserID == u.ThisUserID && u.ThisUsername != "" {
            return u.ThisUsername
        }
        if cachedUser, exists := u.Users[msgRaw.UserID]; exists && cachedUser.Username != "" {
            return cachedUser.Username
        }
    }

    return fmt.Sprintf("User %d", msgRaw.UserID)
}

func createMessage(g *ScreenChat, msgRaw prot.MessageRaw, u *user.UserCache) Message {
    log.Info.Printf("Creating new message widget\n")

    msg := Message{}

    msg.Username = widget.NewLabel(messageUsername(msgRaw, u))
	msg.Username.Wrapping = fyne.TextWrapWord
    msg.Username.TextStyle.Bold = true

    msg.Time = widget.NewLabel(time.Unix(msgRaw.Time, 0).Format("3:04 PM"))

    // Add a reply button
    buttonReply := widget.NewButton("R", func() {
        messageOnReply(g, msgRaw.ID)
    })

    // Add a copy button
    buttonCopy := widget.NewButton("C", func() {
        log.Info.Printf("Copied message %d\n", msgRaw.ID)

        // Give focus back to the message entry when done
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)
        // Copy the message contents
        fyne.CurrentApp().Clipboard().SetContent(msg.Content.Text)
    })

    var c *fyne.Container

    // If the message was sent by the user of this client...
    if msgRaw.UserID == u.ThisUserID {
        // Add a delete button
        buttonDelete := widget.NewButton("D", func() {
            messageOnDelete(g, msgRaw.ID, u)
        })

        // Add an edit button, upon pressing...
        buttonEdit := widget.NewButton("E", func() {
            messageOnEdit(g, &msg, msgRaw.ID, u)
        })
        c = container.NewHBox(
            buttonReply, buttonCopy, buttonEdit, buttonDelete, msg.Time,
        )
    } else {
        c = container.NewHBox(buttonReply, buttonCopy, msg.Time)
    }

    msg.Border = container.NewBorder(nil, nil, msg.Username, c, nil)

    msg.Content = widget.NewLabel(msgRaw.Content)
    msg.Content.Wrapping = fyne.TextWrapWord

    msg.VBox = container.NewVBox(msg.Border, msg.Content)

	msg.Base = widget.NewCard("", "", msg.VBox)

	return msg
}

func createRepliedSection(g *ScreenChat, rawMsg *prot.MessageRaw) *fyne.Container {
    c := container.NewVBox()

    // Get the username and content of the messages being replied to via the messageids in rawMsg.RepliedMessages and add each one to the replied section
    for _, repliedID := range rawMsg.RepliedIDs {
        repliedMsg, exists := g.Widgets.Messages[repliedID]
        var repliedText string
        if !exists {
            repliedText = "Message could not be loaded"
        } else {
            repliedText = 
                fmt.Sprintf("%s: %s", repliedMsg.Username.Text, repliedMsg.Content.Text)
        }

        replyIcon := widget.NewIcon(theme.Current().Icon(theme.IconNameMailReply))
        repliedLabel := widget.NewLabel(repliedText)

        cRepliedMsg := container.NewHBox(replyIcon, repliedLabel)
        c.Add(cRepliedMsg)
    }
    return c
}

func messageOnDelete(g *ScreenChat, msgID int, u *user.UserCache) {
    log.Info.Printf("Delete button on message %d pressed\n", msgID)

    // Give focus back to the message entry when done
    defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

    // Send new DEL chat request to net.serverSend
    req := prot.ChatRequest{
        ChatID: 1,
        MessageContent: prot.NONE_STRING,
        MessageID: msgID,
        Type: prot.REQ_DEL,
        UserID: u.ThisUserID,
    }
    g.OutgoingMessages <- req
}

func messageOnEdit(g *ScreenChat, msg *Message, msgID int, u *user.UserCache) {
    log.Info.Printf("Edit button on message %d pressed\n", msgID)

    // Replace the message content label with an entry to allow editing
    msg.Content.Hide()
    editEntry := widget.NewEntry()
    editEntry.Text = msg.Content.Text
    msg.VBox.Add(editEntry)
    
    // Move cursor to the end of the text
    editEntry.CursorRow = len(msg.Content.Text)

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
                UserID: u.ThisUserID,
            }
            g.OutgoingMessages <- req
        }
        // Replace the edit entry with the message content label again
        editEntry.Hide()
        msg.Content.Show()
    }
}

func messageOnReply(g *ScreenChat, msgID int) {
    log.Info.Printf("Reply button on message %d pressed\n", msgID)
    // Give focus back to the message entry when done
    defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

    if g.RepliedMessages == nil {
        g.RepliedMessages = make(map[int]bool)
    }

    // Add the message ID to g.RepliedMessages if it isn't already there or delete it from there if it is
    if g.RepliedMessages[msgID] {
        delete(g.RepliedMessages, msgID)
        log.Info.Printf("The next message sent will NOT reply to message %d\n", msgID)
    } else {
        g.RepliedMessages[msgID] = true
        log.Info.Printf("The next message sent will reply to message %d\n", msgID)
    }
}
