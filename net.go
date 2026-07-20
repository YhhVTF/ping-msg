package main

import (
    "net"
    "time"
)

// StartNet: Connect to the server and show an error dialog if it fails
// Parameters:
//  gui (*GUI) - GUI elements
func StartNet(gui *GUI) {
    // Until Ping has been quit...
    for !PingQuit {
        // Connect to the server
        Info.Printf("Connecting to server\n")
        _, err := net.Dial("tcp", "127.0.0.1:5555")
        if err != nil {
            Error.Printf("Failed to connect to server: %s\n", err)
            // Show an error dialog if connecting fails
            if gui.Dialogs.ConnectionIssues == nil {
                gui.DialogConnectionIssues(err)
            }
            // Try connecting again after some time
            time.Sleep(30*time.Second)
            continue
        }
    }
}
