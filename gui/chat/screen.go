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

// A collection of all GUI elements to be used chat screen
type ScreenChat struct {
	// All containers
	Containers          ContainerTableChat
	// All widgets
	Widgets             WidgetTableChat
    // The main window
    Window              fyne.Window
	OutgoingMessages    chan prot.ChatRequest // Connects to net.go
    // IDs of messages that the next message sent will reply to
    ReplyingTo          []int
}

// All widgets to be used by the chat screen
type WidgetTableChat struct {
	// Button in the bottom bar that sends the contents of EntryMessage when pressed
	ButtonSend *widget.Button
	// Entry in the bottom bar used to type and send messages
	EntryMessage *widget.Entry
    // Button in the bottom left for attaching files to messages
    ButtonAttach *widget.Button
    // Button in the bottom right of the screen for opening a settings window
    ButtonOptions *widget.Button
    // Messages in the chat containers
    Messages map[int]*Message
    // Widgets that show replying to a message
    //  Key (int) - Message ID
    //  Val (*RepliedMessage) - Replied message
    RepliedMessages map[int]*RepliedMessage
}

func InitScreenChat(
    s *screen.ScreenManager, w fyne.Window, u *user.UserCache, opt *options.Options,
) *ScreenChat {
	log.Info.Printf("Initializing chat screen\n")

	g := &ScreenChat{}
    g.Window = w

	g.OutgoingMessages = make(chan prot.ChatRequest)
    g.ReplyingTo = make([]int, 0)

    g.Widgets.Messages = make(map[int]*Message)
    g.Widgets.RepliedMessages = make(map[int]*RepliedMessage)

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

	msgText := fmt.Sprintf("<User %d> %s", rawMsg.UserID, rawMsg.Content)

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
	}

	g.OutgoingMessages <- req
}
