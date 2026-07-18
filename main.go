package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

var PingQuit = false

func main() {
    a := app.New()

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
