package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"os"

    "github.com/YhhVTF/ping-msg/chat"
    "github.com/YhhVTF/ping-msg/global"
    "github.com/YhhVTF/ping-msg/gui"
    "github.com/YhhVTF/ping-msg/gui/screen"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/net"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/user"
)

// StartPing: Wrapper function that initializes the main window and gui and then connects to the server by calling InitGUI and StartNet respectively
// Parameters:
//
//	a (fyne.App) - argument for InitGUI
//	loadingWindow (fyne.Window) - argument for InitGUI
//  opt (*options.Options) - options/settings
func StartPing(a fyne.App, loadingWindow fyne.Window, opt *options.Options) {
    ping.ChatCache = &chat.ChatCache{Chats: make(map[int]*chat.Chat)}
    ping.Connected = false
    ping.Quit = false
    ping.ScreenManager = screen.NewScreenManager()
    ping.UserCache = &user.UserCache{Users: make(map[int]user.User)}

    var g *gui.GUI
    fyne.DoAndWait(func() {
        g = gui.InitGUI(a, ping.ChatCache, ping.UserCache, loadingWindow, opt)
    })

	net.StartNet(g, ping.UserCache, opt)
}

func main() {
	log.InitLog(os.Stderr, os.Stdout, os.Stdout)
	log.Info.Printf("Starting Ping\n")

	a := app.New()

    // Create a loading window
    // This window will be shown while Ping starts and will be closed when it is ready
    loadingWindow := a.NewWindow("Launching Ping")
    c := container.NewCenter(widget.NewLabel("Loading..."))
    fyne.Do(func() {
        loadingWindow.SetContent(c)
        loadingWindow.Resize(fyne.NewSize(500, 500))
    })
    fyne.Do(loadingWindow.Show)

    var err error
    ping.Options, err = options.LoadOptions(".ping/options")
    if err != nil {
        return
    }

    log.Info.Printf("Loading assets\n")
	iconData, err := os.ReadFile(".ping/assets/icons/ping.png")
	if err == nil {
		a.SetIcon(fyne.NewStaticResource("ping.png", iconData))
	} else {
        log.Error.Printf("Failed to load asset assets/ping.png: %s\n", err)
    }

	go StartPing(a, loadingWindow, ping.Options)

	a.Run()
}
