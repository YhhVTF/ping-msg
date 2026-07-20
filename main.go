package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"os"
)

// Whether or not the user has quit Ping
var PingQuit = false

// StartPing: Wrapper function that initializes the main window and gui and then connects to the server by calling InitGUI and StartNet respectively
// Parameters:
//
//	a (fyne.App) - argument for InitGUI
//	loadingWindow (fyne.Window) - argument for InitGUI
func StartPing(a fyne.App, loadingWindow fyne.Window) {
	StartNet(InitGUI(a, loadingWindow))
}

func main() {
	InitLog(os.Stderr, os.Stdout, os.Stdout)
	Info.Printf("Starting Ping\n")

	a := app.New()

	iconData, err := os.ReadFile("assets/icon.png")
	if err == nil {
		a.SetIcon(fyne.NewStaticResource("icon.png", iconData))
	}

	// Create a loading window
	// This window will be shown while Ping starts and will be closed when it is ready
	loadingWindow := a.NewWindow("Launching Ping")
	c := container.NewCenter(widget.NewLabel("Loading..."))
	fyne.Do(func() {
		loadingWindow.SetContent(c)
	})
	fyne.Do(loadingWindow.Show)
	//fyne.Do(loadingWindow.SetMaster)

	go StartPing(a, loadingWindow)

	a.Run()
}
