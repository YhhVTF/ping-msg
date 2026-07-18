package main

import (
    "fyne.io/fyne/v2"

    "net"
    "time"
)

func StartPing(a fyne.App, loadingWindow fyne.Window) {
    gui := InitGUI(a)

    loadingWindow.Close()
    gui.Window.SetMaster()
    gui.Window.Show()

    for !PingQuit {
        _, err := net.Dial("tcp", "127.0.0.1:5555")
        if err != nil {
            if gui.Dialogs.ConnectionIssues == nil {
                gui.DialogConnectionIssues(err)
            }
            time.Sleep(30*time.Second)
            continue
        }
    }
}
