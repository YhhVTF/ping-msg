package main

import (
    "encoding/json"
	"net"
	"time"
)

// StartNet: Connect to the server and show an error dialog if it fails
// Parameters:
//
//	gui (*GUI) - GUI elements
func StartNet(gui *GUI) {
	// Until Ping has been quit...
	for !PingQuit {
		// Connect to the server
		Info.Printf("Connecting to server\n")

		conn, err := net.Dial("tcp", "127.0.0.1:5555")
		if err != nil {
			Error.Printf("Failed to connect to server: %s\n", err)
			if gui.Dialogs.ConnectionIssues == nil {
				gui.DialogConnectionIssues(err)
			}
			time.Sleep(30 * time.Second)
			continue
		}

		Info.Printf("Successfully connected to server\n")

		connDone := make(chan bool)
		go HandleServerCommunication(conn, gui, connDone)

		<-connDone
		Error.Printf("Connection lost. Reconnecting in 5 seconds...\n")
		time.Sleep(5 * time.Second)
	}
}

func HandleServerCommunication(conn net.Conn, gui *GUI, connDone chan bool) {
	defer conn.Close()

	connectionFailed := make(chan bool)

	go serverRecieve(conn, gui, connectionFailed)

	go serverSend(conn, gui, connectionFailed)

	<-connectionFailed // Wait for communication failure
	connDone <- true   // tell StartNet connection died
}

func serverRecieve(conn net.Conn, gui *GUI, done chan bool) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			done <- true
			return
		}

		Info.Printf("Recieved %d bytes from server\n", n)

        response := ChatResponse{}
        err = json.Unmarshal(buffer, &response)
        if err != nil {
            Error.Printf("Failed to parse incoming data from server: %s\n", err)
            continue
        }

        gui.AddMessage(response.Messages.Messages[0])
	}
}

func serverSend(conn net.Conn, gui *GUI, done chan bool) {
	for {
		select {
		case msg := <-gui.OutgoingMessages:
			_, err := conn.Write(msg)
			if err != nil {
				done <- true
				return
			}
		case <-done:
			return
		}
	}
}
