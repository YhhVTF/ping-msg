package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/data/binding"
    "fyne.io/fyne/v2/widget"

    "slices"
    "time"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/global"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
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

func createMessage(g *ScreenChat, msg *chat.Message, cacheBind binding.String, chat *chat.Chat, u *user.UserCache, opt *options.Options) *Message {
    log.Info.Printf("Creating new message widget\n")

    msgWidget := &Message{}

    if len(msg.RepliedIDs) > 0 {
        // Replied messages are added to ScreenChat.RepliedMessages in createRepliedSection
        msgWidget.RepliedSection = createRepliedSection(g, msg.RepliedIDs, chat, u, opt)
    }

    msgWidget.Username = widget.NewLabelWithData(u.UsersBind[*msg.Username].Username)

	msgWidget.Username.Wrapping = fyne.TextWrapWord
    msgWidget.Username.TextStyle.Bold = true

    msgWidget.Time = widget.NewLabel(time.Unix(msg.Time, 0).Format("3:04 PM"))

    // Add a reply button
    buttonReply := widget.NewButton("R", func() {
        messageOnReply(g, msg.ID, chat)
    })

    // Add a copy button
    buttonCopy := widget.NewButton("C", func() {
        log.Info.Printf("Copied message %d\n", msg.ID)

        // Give focus back to the message entry when done
        defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)
        // Copy the message contents
        fyne.CurrentApp().Clipboard().SetContent(msgWidget.Content.Text)
    })

    var c *fyne.Container

    // If the message was sent by the user of this client...
    if *msg.Username == u.ThisUser.Username {
        // Add a delete button
        buttonDelete := widget.NewButton("D", func() {
            messageOnDelete(g, msg.ID, ping.ChatCache, u)
        })

        // Add an edit button, upon pressing...
        buttonEdit := widget.NewButton("E", func() {
            messageOnEdit(g, msgWidget, msg.ID, ping.ChatCache, u)
        })
        c = container.NewHBox(
            buttonReply, buttonCopy, buttonEdit, buttonDelete, msgWidget.Time,
        )
    } else {
        c = container.NewHBox(buttonReply, buttonCopy, msgWidget.Time)
    }

    msgWidget.Border = container.NewBorder(nil, nil, msgWidget.Username, c, nil)

    msgWidget.Content = widget.NewLabelWithData(chat.MessagesBind[msg.ID])
    msgWidget.Content.Wrapping = fyne.TextWrapWord

    msgWidget.VBox = container.NewVBox(msgWidget.Border, msgWidget.Content)

	msgWidget.Card = widget.NewCard("", "", msgWidget.VBox)

    if msgWidget.RepliedSection == nil {
        msgWidget.Base = container.NewVBox(msgWidget.Card)
    } else {
        msgWidget.Base = container.NewVBox(msgWidget.RepliedSection, msgWidget.Card)
    }
	return msgWidget
}

func createRepliedSection(
    g *ScreenChat, repliedIDs []int, chat *chat.Chat,
    u *user.UserCache, opt *options.Options,
) *fyne.Container {
    vbox := container.NewVBox()

    // Get the username and content of the messages being replied to, create a replied message widget for each one and add each one to the replied section
    for _, repliedID := range repliedIDs {
        if _, exists := g.Widgets.RepliedMessages[repliedID]; !exists {
            g.Widgets.RepliedMessages[repliedID] =
                createRepliedMessage(chat.Messages[repliedID], repliedID, u, opt)
        }
        vbox.Add(g.Widgets.RepliedMessages[repliedID].Base)
    }
    c := container.NewPadded(vbox)
    return c
}

func messageOnDelete(g *ScreenChat, msgID int, c *chat.ChatCache, u *user.UserCache) {
    log.Info.Printf("Delete button on message %d pressed\n", msgID)

    // Give focus back to the message entry when done
    defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

    // Send new DEL chat request to net.serverSend
    g.ChatRequestDelete(msgID, c, u)
}

func messageOnEdit(g *ScreenChat, msg *Message, msgID int, c *chat.ChatCache, u *user.UserCache) {
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
            g.ChatRequestEdit(msgID, text, c, u)
        }
        // Replace the edit entry with the message content label again
        editEntry.Hide()
        msg.Content.Show()
    }
}

func messageOnReply(g *ScreenChat, msgID int, chatCache *chat.Chat) {
    log.Info.Printf("Reply button on message %d pressed\n", msgID)
    // Give focus back to the message entry when done
    defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

    // Add the message ID to chatCache.ReplyingTo if it isn't already there or remove it from there if it is
    if slices.Contains(chatCache.ReplyingTo, msgID) {
        i := slices.Index(chatCache.ReplyingTo, msgID)
        chatCache.ReplyingTo = slices.Delete(chatCache.ReplyingTo, i, i+1)
    } else {
        chatCache.ReplyingTo = append(chatCache.ReplyingTo, msgID)
    }

    log.Info.Printf("The next message sent will reply to messages %d\n", chatCache.ReplyingTo)
}
