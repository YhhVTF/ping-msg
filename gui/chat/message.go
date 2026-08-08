package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "fmt"
    "slices"
    "time"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/global"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

type Message struct {
    // Base VBox container for the message that contains the reply section if there is one, and the message card
    Base *fyne.Container
    // A VBox with all replied messages, added to Base first if there is any
    RepliedSection *fyne.Container
    // Surrounds message in a card, added to Base first, second if there's replies
    Card *widget.Card
    // Base VBox container for card content
    VBox *fyne.Container
    // Border container for message metadata, added to VBox first, second if there's replied messages
    Border *fyne.Container
    // Label for message content, added to VBox second, third of there are replied messages
    Content *widget.Label
    // Label for username, added to left side of Border
    Username *widget.Label
    // Label for time, added to right side of Border
    Time *widget.Label
}

// messageUsername returns the best available display name for a message.
// MessageRaw identifies senders by UserID; usernames are held separately in
// the local user cache.
func messageUsername(msgRaw *prot.MessageRaw, u *user.UserCache) string {
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

func createMessage(g *ScreenChat, msgRaw *prot.MessageRaw, u *user.UserCache) *Message {
    log.Info.Printf("Creating new message widget\n")

    msg := &Message{}

    if len(msgRaw.RepliedIDs) > 0 {
        // Replied messages are added to ScreenChat.RepliedMessages in createRepliedSection
        msg.RepliedSection = createRepliedSection(g, msgRaw)
    }

    msg.Username = widget.NewLabel(messageUsername(msgRaw, u))
	msg.Username.Wrapping = fyne.TextWrapWord
    msg.Username.TextStyle.Bold = true

    msg.Time = widget.NewLabel(time.Unix(msgRaw.Time, 0).Format("3:04 PM"))

    // Add a reply button
    buttonReply := widget.NewButton("R", func() {
        messageOnReply(g, msgRaw.ID, ping.ChatCache)
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
            messageOnEdit(g, msg, msgRaw.ID, u)
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

	msg.Card = widget.NewCard("", "", msg.VBox)

    if msg.RepliedSection == nil {
        msg.Base = container.NewVBox(msg.Card)
    } else {
        msg.Base = container.NewVBox(msg.RepliedSection, msg.Card)
    }
	return msg
}

func createRepliedSection(g *ScreenChat, rawMsg *prot.MessageRaw) *fyne.Container {
    vbox := container.NewVBox()

    // Get the username and content of the messages being replied to via the messageids in rawMsg.RepliedMessages and add each one to the replied section
    for _, repliedID := range rawMsg.RepliedIDs {
        if _, exists := g.Widgets.RepliedMessages[repliedID]; !exists {
            g.Widgets.RepliedMessages[repliedID] = createRepliedMessage(g, repliedID)
        }
        vbox.Add(g.Widgets.RepliedMessages[repliedID].Base)
    }
    c := container.NewPadded(vbox)
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

func messageOnReply(g *ScreenChat, msgID int, c *chat.ChatCache) {
    log.Info.Printf("Reply button on message %d pressed\n", msgID)
    // Give focus back to the message entry when done
    defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

    // Add the message ID to c.ThisChat.ReplyingTo if it isn't already there or remove it from there if it is
    if slices.Contains(c.ThisChat.ReplyingTo, msgID) {
        i := slices.Index(c.ThisChat.ReplyingTo, msgID)
        c.ThisChat.ReplyingTo = slices.Delete(c.ThisChat.ReplyingTo, i, i+1)
    } else {
        c.ThisChat.ReplyingTo = append(c.ThisChat.ReplyingTo, msgID)
    }

    log.Info.Printf("The next message sent will reply to messages %d\n", c.ThisChat.ReplyingTo)
}
