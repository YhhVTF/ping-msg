package schat

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fmt"
	"image/color"

    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

// All containers to be used by chat screen
type ContainerTableChat struct {
	// Highest level container that contains all other objects in the window
	Base *fyne.Container
	// Container for objects at the bottom of the screen (like the message entry, send button, etc.)
	BottomBar *fyne.Container
    // Container for buttons at the bottom of the screen to the left of the message entry
    BottomLeftCluster *fyne.Container
	// Containers containing the messages in the chat
	Chat Chat
}

// All dialogs to be used by chat screen
type DialogTable struct {
	// Informs the user that there are issues with connecting to the server
	ConnectionIssues *dialog.CustomDialog
    // Asks the user for information to log in
    Login *dialog.CustomDialog
}

// A collection of all GUI elements to be used chat screen
type ScreenChat struct {
	// All containers
	Containers ContainerTableChat
	// All dialogs
	Dialogs DialogTable
	// All widgets
	Widgets          WidgetTable
    // The main window
    Window fyne.Window
	OutgoingMessages chan prot.ChatRequest // Connects to net.go
}

// All widget to be used by Ping
type WidgetTable struct {
	// Button in the bottom bar that sends the contents of EntryMessage when pressed
	ButtonSend *widget.Button
	// Entry in the bottom bar used to type and send messages
	EntryMessage *widget.Entry
    // Button in the bottom left for attaching files to messages
    ButtonAttach *widget.Button
    // Button in the bottom right of the screen for opening a settings window
    ButtonOptions *widget.Button
    // Messages in the chat containers
    Messages map[int]Message
}

// DialogConnectionIssues: Creates and shows a dialog set to the default size that informs the user that there are connection issues, with a user friendly message and a technical message
// Parameters:
//
//	err (error) - The error that occurred. This will be used as the technical error message
//  opt (*options.Options) - Options/settings
func (g *ScreenChat) DialogConnectionIssues(err error, opt *options.Options) {
	log.Info.Printf("Creating Dialog ConnectionIssues\n")

	// Create the user friendly error message as a label
	uxErrMsg := widget.NewLabel(opt.GUIText.DialogConnIssues.Prompt)
	// Create the technical error message as a label, make it selectable and a low importance widget
	technicalErrMsg := widget.NewLabel(err.Error())
	technicalErrMsg.Selectable = true
	technicalErrMsg.Importance = widget.LowImportance
	// add then to a new vbox
	c := container.NewVBox(uxErrMsg, technicalErrMsg)

	// Create a dialog with the vbox as the content
	dialog := dialog.NewCustom(opt.GUIText.DialogConnIssues.Title, "", c, g.Window)
	// add an ok button
	dialog.SetButtons([]fyne.CanvasObject{
		widget.NewButton(opt.GUIText.DialogConnIssues.Buttons[0].Label, func() {
			log.Info.Printf("Dialog ConnectionsIssues dismissed\n")
			// Dismiss the dialog and set it to nil in the dialog table
			dialog.Dismiss()
			g.Dialogs.ConnectionIssues = nil

            // Set focus on message entry now that dialog is dismissed
            g.Window.Canvas().Focus(g.Widgets.EntryMessage)
		}),
	})
	// Resize to default dialog size and show the dialog
	dialog.Resize(fyne.NewSize(opt.GUI.DialogConnIssues.Size[0], opt.GUI.DialogConnIssues.Size[1]))
	dialog.Show()

	// add the dialog to the dialog table
	g.Dialogs.ConnectionIssues = dialog
}

func (g *ScreenChat) DialogLogin(u *user.UserData, opt *options.Options) {
	log.Info.Printf("Creating Dialog Login\n")

    // Text telling the user what to do
	prompt := widget.NewLabel(opt.GUIText.DialogLogin.Prompt)
    // Entry for username
    entry := widget.NewEntry()
    entry.SetPlaceHolder(opt.GUIText.EntryUsername.Label)

	// add them to a new vbox
	c := container.NewVBox(prompt, entry)

	// Create a dialog with the vbox as the content
	dialog := dialog.NewCustom(opt.GUIText.DialogLogin.Title, "", c, g.Window)
	// add a login button
	dialog.SetButtons([]fyne.CanvasObject{
		widget.NewButton(opt.GUIText.DialogLogin.Buttons[0].Label, func() {
            // Set the text in the entry as the username if it isn't a reserved username or empty
            if entry.Text == "" { return }
            if entry.Text == prot.SERVER_USERNAME || entry.Text == prot.NONE_STRING {
                prompt.SetText(opt.GUIText.DialogLoginAltPrompt)
                return
            }
            u.ThisUser = entry.Text

            // Dismiss the dialog and set it as nil in the dialog table
            dialog.Dismiss()
            g.Dialogs.Login = nil

            // Set focus on message entry now that dialog is dismissed
            g.Window.Canvas().Focus(g.Widgets.EntryMessage)

            log.Info.Printf("Username set as %s\n", u.ThisUser)
		    log.Info.Printf("Dialog Login dismissed\n")
		}),
	})
	// Resize to default dialog size and show the dialog
	dialog.Resize(fyne.NewSize(opt.GUI.DialogConnIssues.Size[0], opt.GUI.DialogConnIssues.Size[1]))
	dialog.Show()

    entry.OnSubmitted = func(text string) {
        // Set the text in the entry as the username if it isn't a reserved username or empty
        if entry.Text == "" { return }
        if entry.Text == prot.SERVER_USERNAME || entry.Text == prot.NONE_STRING {
            prompt.SetText(opt.GUIText.DialogLoginAltPrompt)
            return
        }
        u.ThisUser = entry.Text

        // Dismiss the dialog and set it as nil in the dialog table
        dialog.Dismiss()
        g.Dialogs.Login = nil

        // Set focus on message entry now that dialog is dismissed
        g.Window.Canvas().Focus(g.Widgets.EntryMessage)

        log.Info.Printf("Username set as %s\n", u.ThisUser)
		log.Info.Printf("Dialog Login dismissed\n")
    }

	// add the dialog to the dialog table
	g.Dialogs.Login = dialog
}

func InitScreenChat(
    s *screen.ScreenManager, w fyne.Window, u *user.UserData, opt *options.Options,
) *ScreenChat {
	log.Info.Printf("Initializing chat screen\n")

	g := &ScreenChat{}
    g.Window = w

	g.OutgoingMessages = make(chan prot.ChatRequest)

    g.Widgets.Messages = make(map[int]Message)

    g.Widgets.ButtonAttach = createButtonAttach(g, opt)

    g.Widgets.ButtonSend = createButtonSend(g, u, opt)

    g.Containers.BottomLeftCluster =
        container.NewHBox(g.Widgets.ButtonAttach, g.Widgets.ButtonSend)

    g.Widgets.EntryMessage = createEntryMessage(g, u, opt)
    // Set focus on message entry
    g.Window.Canvas().Focus(g.Widgets.EntryMessage)

    g.Widgets.ButtonOptions = createButtonOptions(g, s, opt)

	// Initialize chat containers
	g.Containers.Chat = NewChat()

	// Initialize bottom bar container and add the message entry and send button
	g.Containers.BottomBar = container.NewBorder(
		// top, bottom, left, right, center
		nil, nil, g.Containers.BottomLeftCluster, g.Widgets.ButtonOptions, g.Widgets.EntryMessage,
	)

	// Initialize the base container and add the chat scroll container and bottom bar
	g.Containers.Base = container.NewBorder(
		// top, bottom, left, right, center
		nil, g.Containers.BottomBar, nil, nil, g.Containers.Chat.Base,
	)
	return g
}

// NewDialog: Creates a new custom dialog with the specified title and content set to the default dialog size
// Parameters:
//
//	title (string) - Title of the dialog
//	content (string) - Text to be set as the dialog content
//
// Returns:
//
//	*dialog.CustomDialog - The new dialog
func (g *ScreenChat) NewDialog(title, content string) *dialog.CustomDialog {
	dialog := dialog.NewCustom(title, "", widget.NewLabel(content), g.Window)
	dialog.Resize(fyne.NewSize(350, 200))
	dialog.Show()
	return dialog
}

func (g *ScreenChat) ReceiveMessage(rawMsg prot.MessageRaw) {
	log.Info.Printf("Received message\n")

	msgText := fmt.Sprintf("<%s> %s", rawMsg.Username, rawMsg.Content)

	msg := canvas.NewText(msgText, color.NRGBA{255, 255, 255, 255})
	g.Containers.Chat.VBox.Add(msg)
	g.Containers.Chat.VScroll.ScrollToBottom()
}

// UNUSED BUT EXPANDABLE LATER??

// oopie - dash

// huh? Twt

// oopie :P - entie
func (g *ScreenChat) SendMessage() {
	log.Info.Printf("Sending message\n")

	req := prot.ChatRequest{
		ChatID:         0,
		MessageContent: g.Widgets.EntryMessage.Text,
		MessageID:      prot.NONE_INT,
		Type:           prot.REQ_ADD,
		Username:       "buh",
	}

	g.OutgoingMessages <- req
}
