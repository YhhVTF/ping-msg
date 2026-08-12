package gui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/gui/chat"
    "github.com/YhhVTF/ping-msg/gui/dialogs"
    "github.com/YhhVTF/ping-msg/gui/options"
    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/gui/sidebar"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/user"
)

type GUI struct {
    // Base container that splits the window into the sidebar and main screen
    Base *container.Split
    // Chat screen
    Chat *schat.ScreenChat
    // All dialogs used by Ping
    Dialogs dialogs.DialogTable
    // Manages inner windows
    InnerWindows *container.MultipleWindows
    // Options screen
    Options *sopt.ScreenOptions
    // Manages screens
    ScreenManager *screen.ScreenManager
    // Sidebar
    Sidebar *sside.ScreenSidebar
    // The main window
    Window fyne.Window
}

// InitGUI: Initializes the main window and all objects in it, closes the loading window and then shows the main window
// Parameters:
//
//	a (fyne.App) - The fyne application the window will be initialized in
//  c (*chat.ChatCache) - Data/cache for chats
//  u (*user.UserCache) - Information pertaining to users
//	loadingWindow (fyne.Window) - The loading window
//  opt (*options.Options) - Options/settings
//
// Returns:
//
//	*GUI - The main window and all its objects
func InitGUI(
    a fyne.App, c *chat.ChatCache,
    u *user.UserCache, loadingWindow fyne.Window, opt *options.Options,
) *GUI {
	log.Info.Printf("Creating GUI\n")

	g := &GUI{}
	g.Window = a.NewWindow(opt.GUIText.Window.Title)
    g.Window.Resize(fyne.NewSize(opt.GUI.Window.Size[0], opt.GUI.Window.Size[1]))
    g.InnerWindows = container.NewMultipleWindows()

    // Start screen management
    g.ScreenManager = screen.NewScreenManager()
    go g.manageScreens(g.ScreenManager, c, u, opt)

    c.Chats[1] = chat.NewChatCache()
    c.ThisChat = c.Chats[1]

    // Initialize sidebar
    log.Info.Printf("Initializing sidebar\n")
    g.Sidebar = sside.InitScreenSidebar(g.Window, g.ScreenManager, c)
    // Create split and add sidebar as the leading child
    g.Base = container.NewHSplit(g.Sidebar.Base, container.NewStack())

    // Initialize chat screen
    g.ScreenManager.ScreenChatFull()

    // Set window content as the split
    g.Window.SetContent(g.Base)

	// Close loading window, set the main window as master window and show it
	loadingWindow.Close()
	g.Window.SetMaster()
	g.Window.Show()

    return g
}

func (g *GUI) manageScreens(s *screen.ScreenManager, c *chat.ChatCache, u *user.UserCache, opt *options.Options) {
    for {
        switch <-s.Chan {
        case screen.S_CHAT_FOCUS_DEFAULT:
            fyne.Do(func() {
                g.Window.Canvas().Focus(g.Chat.Widgets.EntryMessage)
            })
        case screen.S_CHAT_FULL:
            log.Info.Printf("Opening chat screen on the main window\n")

            fyne.Do(func() {
                // Initialize the chat screen
                g.Chat = schat.InitScreenChat(s, g.Window, c, u, opt)
	            // Set trailing child of split as the base container of the chat screen
                g.Base.Trailing.(*fyne.Container).Add(g.Chat.Containers.Base)
            })

        case screen.S_OPTIONS_FLOAT:
            log.Info.Printf("Opening options screen as floating window\n")

            fyne.Do(func() {
                g.Options = sopt.InitScreenOptions(g.Window, opt)
                g.Options.Float(g.Window, g.InnerWindows, opt)
            })
        }
    }
}
