package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fmt"
	"image/color"
)

// All containers to be used by Ping
type ContainerTable struct {
	// Highest level container that contains all other objects in the window
	Base *fyne.Container
	// Container for objects at the bottom of the screen (like the message entry, send button, etc.)
	BottomBar *fyne.Container
	// Containers containing the messages in the chat
	Chat Chat
}

// All dialogs to be used by Ping
type DialogTable struct {
	// Informs the user that there are issues with connecting to the server
	ConnectionIssues *dialog.CustomDialog
}

// A collection of all GUI elements to be used
type GUI struct {
	// The main window
	Window fyne.Window
	// All containers
	Containers ContainerTable
	// All dialogs
	Dialogs DialogTable
	// All widgets
	Widgets          WidgetTable
	OutgoingMessages chan ChatRequest // Connects to net.go
}

// All widget to be used by Ping
type WidgetTable struct {
	// Button in the bottom bar that sends the contents of BottomBarEntry when pressed
	BottomBarButtonSend *widget.Button
	// Entry in the bottom bar used to type and send messages
	BottomBarEntry *widget.Entry
}

// DialogConnectionIssues: Creates and shows a dialog set to the default size that informs the user that there are connection issues, with a user friendly message and a technical message
// Parameters:
//
//	err (error) - The error that occurred. This will be used as the technical error message
func (g *GUI) DialogConnectionIssues(err error) {
	Info.Printf("Creating Dialog ConnectionIssues\n")

	// Create the user friendly error message as a label
	uxErrMsg := widget.NewLabel("Failed to connect to server. The server may be down or your internet connection could be unstable")
	// Create the technical error message as a label, make it selectable and a low importance widget
	technicalErrMsg := widget.NewLabel(err.Error())
	technicalErrMsg.Selectable = true
	technicalErrMsg.Importance = widget.LowImportance
	// add then to a new vbox
	c := container.NewVBox(uxErrMsg, technicalErrMsg)

	// Create a dialog with the vbox as the content
	dialog := dialog.NewCustom("Connection Issues", "", c, g.Window)
	// add an ok button
	dialog.SetButtons([]fyne.CanvasObject{
		widget.NewButton("Ok", func() {
			Info.Printf("Dialog ConnectionsIssues dismissed\n")
			// Dismiss the dialog and set it to nil in the dialog table
			dialog.Dismiss()
			g.Dialogs.ConnectionIssues = nil
		}),
	})
	// Resize to default dialog size and show the dialog
	dialog.Resize(fyne.NewSize(350, 200))
	dialog.Show()

	// add the dialog to the dialog table
	g.Dialogs.ConnectionIssues = dialog
}

// InitGUI: Initializes the main window and all objects in it, closes the loading window and then shows the main window
// Parameters:
//
//	a (fyne.App) - The fyne application the window will be initialized in
//	loadingWindow (fyne.Window) - The loading window
//
// Returns:
//
//	*GUI - The main window and all its objects
func InitGUI(a fyne.App, loadingWindow fyne.Window) *GUI {
	Info.Printf("Creating GUI\n")

	g := &GUI{}
	g.OutgoingMessages = make(chan ChatRequest)
	g.Window = a.NewWindow("Ping")

	// Initialize message entry
	g.Widgets.BottomBarEntry = widget.NewEntry()
	g.Widgets.BottomBarEntry.PlaceHolder = "Type a message..."
	g.Widgets.BottomBarEntry.OnSubmitted = func(text string) {
		Info.Printf("Widget BottomBarEntry submitted (%s)\n", text)
		if text == "" {
			return
		}

		req := ChatRequest{
			ChatID:         1,
			MessageContent: text,
			MessageID:      -1,
			Type:           REQ_ADD,
			Username:       "TestUser",
		}
		g.OutgoingMessages <- req

		g.Widgets.BottomBarEntry.SetText("")
	}

	// Initialize send button
	g.Widgets.BottomBarButtonSend = widget.NewButton("Send", func() {
		Info.Printf("Widget BottomBarButtonSend pressed\n")
		text := g.Widgets.BottomBarEntry.Text
		if text == "" {
			return
		}

		req := ChatRequest{
			ChatID:         1,
			MessageContent: text,
			MessageID:      -1,
			Type:           REQ_ADD,
			Username:       "TestUser",
		}
		g.OutgoingMessages <- req

		g.Widgets.BottomBarEntry.SetText("")
	})

	// Initialize chat containers
	g.Containers.Chat = NewChat()

	// Initialize bottom bar container and add the message entry and send button
	g.Containers.BottomBar = container.NewBorder(
		// top, bottom, left, right, center
		nil, nil, nil, g.Widgets.BottomBarButtonSend, g.Widgets.BottomBarEntry,
	)

	// Initialize the base container and add the chat scroll container and bottom bar
	g.Containers.Base = container.NewBorder(
		// top, bottom, left, right, center
		nil, g.Containers.BottomBar, nil, nil, g.Containers.Chat.Base,
	)
	// Set window content as the base container
	g.Window.SetContent(g.Containers.Base)

	// Close loading window, set the main window as master window and show it
	loadingWindow.Close()
	g.Window.SetMaster()
	g.Window.Show()

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
func (g *GUI) NewDialog(title, content string) *dialog.CustomDialog {
	dialog := dialog.NewCustom(title, "", widget.NewLabel(content), g.Window)
	dialog.Resize(fyne.NewSize(350, 200))
	dialog.Show()
	return dialog
}

func (g *GUI) ReceiveMessage(rawMsg MessageRaw) {
	Info.Printf("Received message\n")

	msgText := fmt.Sprintf("<%s> %s", rawMsg.Username, rawMsg.Content)

	msg := canvas.NewText(msgText, color.NRGBA{255, 255, 255, 255})
	g.Containers.Chat.VBox.Add(msg)
	g.Containers.Chat.VScroll.ScrollToBottom()
}

// UNUSED BUT EXPANDABLE LATER??

// oopie - dash

// huh? Twt

// oopie :P - entie
func (g *GUI) SendMessage() {
	Info.Printf("Sending message\n")

	req := ChatRequest{
		ChatID:         0,
		MessageContent: g.Widgets.BottomBarEntry.Text,
		MessageID:      -1,
		Type:           REQ_ADD,
		Username:       "buh",
	}

	g.OutgoingMessages <- req
}
