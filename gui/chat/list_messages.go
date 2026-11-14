package schat

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"

    "time"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

func createListMessages(g *ScreenChat, chat *chat.Chat, c *chat.ChatCache, u *user.UserCache) *widget.List {
    list := widget.NewList(
        func() int {
            if len(chat.Islands) == 0 {
                return 0
            }
            return len(chat.Islands[chat.RenderedIsland].MessageIDs)
        },
        func() fyne.CanvasObject {
            return createMessageCardTemplate()
        },
        func(itemID widget.ListItemID, o fyne.CanvasObject) {
            msgID := chat.Islands[chat.RenderedIsland].MessageIDs[itemID]
            msg, exists := chat.Messages[msgID]
            if !exists {
                o.(*fyne.Container).Hide()
                return
            }

            card := o.(*fyne.Container).Objects[1].(*widget.Card)
            wTopOfMsg := card.Content.(*fyne.Container).Objects[0].(*fyne.Container)

            wUsername := wTopOfMsg.Objects[0].(*widget.Label)
            wUsername.SetText(*msg.Username)

            wContent := card.Content.(*fyne.Container).Objects[1].(*widget.RichText)
            wContent.ParseMarkdown(msg.Content)

            // If this message should be edited, and it hasn't already been done...
            if chat.Editing == msg.ID && !wContent.Hidden {
                // Replace the message content label with an entry to allow editing
                wContent.Hide()
                editEntry := widget.NewEntry()
                editEntry.Text = msg.Content
                card.Content.(*fyne.Container).Add(editEntry)

                // Move cursor to the end of the text
                editEntry.CursorRow = len(msg.Content)

                // Focus on the edit entry
                g.Window.Canvas().Focus(editEntry)

                // On submission of entry...
                editEntry.OnSubmitted = func(text string) {
                    log.Info.Printf("Edit entry on message %d submitted\n", msgID)

                    // Give focus back to the message entry when done
                    defer g.Window.Canvas().Focus(g.Widgets.EntryMessage)

                    // Send edit request if there was an actual edit
                    if text != msg.Content {
                        g.ChatRequestEdit(msgID, text, c, u)
                    }
                    // Replace the edit entry with the message content label again
                    editEntry.Hide()
                    wContent.Show()

                    // Mark no message as being edited
                    chat.Editing = prot.NO_ID
                }
            }

            hbox := card.Content.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container)
            var wTime *widget.Label


            if len(hbox.Objects) > 1 {
                wTime = hbox.Objects[4].(*widget.Label)
            // If it hasn't already been done...
            } else {
                wTime = hbox.Objects[0].(*widget.Label)

                // And the message has been highlighted...
                if chat.Selections.GeneralSelection == msg.ID {
                    hbox.RemoveAll()

                    // Create and add the action buttons (edit, copy, etc.) beside the time label
                    hbox.Add(createButtonReply(g, msgID, chat))
                    hbox.Add(createButtonCopy(g, itemID, chat))
                    hbox.Add(createButtonEdit(g, itemID, chat))
                    hbox.Add(createButtonDelete(g, msgID, c, u))
                    hbox.Add(wTime)
                }
            }
            wTime.SetText(time.Unix(msg.Time, 0).Format("3:04 PM"))
        },
    )

    list.OnHighlighted = func(itemID widget.ListItemID) {
        msgID := chat.Islands[chat.RenderedIsland].MessageIDs[itemID]
        chat.Selections.GeneralSelection = msgID
        list.Refresh()
    }
    return list
}

func createMessageCardTemplate() *fyne.Container {
    labelUsername := widget.NewLabel("")
    labelUsername.Wrapping = fyne.TextWrapWord
    labelUsername.TextStyle.Bold = true

    // Base VBox
    return container.NewVBox(
        // VBox for replied messages
        container.NewVBox(),

        // Main message card
        widget.NewCard("", "",
            // VBox for contents of message card
            container.NewVBox(
                // Border container for message metadata
                container.NewBorder(
                    nil, nil,
                    // Username label (left of border)
                    labelUsername,
                    // HBox for time label along with buttons to perform actions on the message (the latter appears when message is highlighted, right of border)
                    container.NewHBox(
                        // Time label
                        widget.NewLabel(""),
                    ),
                    nil,
                ),

                // Message content
                widget.NewRichTextWithText(""),
            ),
        ),
    )
}
