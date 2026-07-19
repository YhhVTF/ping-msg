package main

import (
    "net"
    "time"
)

func StartNet(gui *GUI) {
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
