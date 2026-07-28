package gui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"

    "github.com/YhhVTF/ping-msg/gui/chat"
    "github.com/YhhVTF/ping-msg/gui/options"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/user"
)

type GUI struct {
    // Chat screen
    Chat *schat.ScreenChat
    // Manages inner windows
    InnerWindows *container.MultipleWindows
    // Options screen
    Options *sopt.ScreenOptions
    // The main window
    Window fyne.Window
}

// InitGUI: Initializes the main window and all objects in it, closes the loading window and then shows the main window
// Parameters:
//
//	a (fyne.App) - The fyne application the window will be initialized in
//  u (*user.UserData) - log.Information pertaining to users
//	loadingWindow (fyne.Window) - The loading window
//  opt (*options.Options) - Options/settings
//
// Returns:
//
//	*GUI - The main window and all its objects
func InitGUI(a fyne.App, u *user.UserData, loadingWindow fyne.Window, opt *options.Options) *GUI {
	log.Info.Printf("Creating GUI\n")

	g := &GUI{}
	g.Window = a.NewWindow(opt.GUIText.Window.Title)
    g.Window.Resize(fyne.NewSize(opt.GUI.Window.Size[0], opt.GUI.Window.Size[1]))
    g.InnerWindows = container.NewMultipleWindows()

    // Initialize the chat screen
    g.Chat = schat.InitScreenChat(g.Window, u, opt)
	// Set window content as the base container of the chat screen
	g.Window.SetContent(g.Chat.Containers.Base)

	// Close loading window, set the main window as master window and show it
	loadingWindow.Close()
	g.Window.SetMaster()
	g.Window.Show()

    return g
}
