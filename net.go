package main

import (
	"encoding/json"
	"net"
    "os"
	"time"

    "ping/protocol"
)

var Connected = false

// StartNet: Connect to the server and show an error dialog if it fails
// Parameters:
//
//	gui (*GUI) - GUI elements
//  u (*UserData) - Information pertaining to users
//  opt (*Options) - Options/settings
func StartNet(gui *GUI, u *UserData, opt *Options) {
    // Prompt for a username
    gui.DialogLogin(u)
    for gui.Dialogs.Login != nil {}

        addr := "127.0.0.1:5555"
        if len(os.Args) > 1 {
            addr = os.Args[1]
        }

	// Until Ping has been quit...
	for !PingQuit {
		// Connect to the server
		Info.Printf("Connecting to server\n")

        conn, err := net.Dial("tcp", addr)
		if err != nil {
			Error.Printf("Failed to connect to server: %s\n", err)
			if gui.Dialogs.ConnectionIssues == nil {
				gui.DialogConnectionIssues(err)
			}
			time.Sleep(30 * time.Second)
			continue
		}

        Connected = true
		Info.Printf("Successfully connected to server\n")
        // Dismiss connection issues dialog after connecting successfully
        if gui.Dialogs.ConnectionIssues != nil {
            gui.Dialogs.ConnectionIssues.Dismiss()
            gui.Dialogs.ConnectionIssues = nil
        }

		connDone := make(chan bool)
		go HandleServerCommunication(conn, gui, u, connDone)

		<-connDone
        Connected = false
		Error.Printf("Connection lost. Reconnecting in 5 seconds...\n")
		time.Sleep(5 * time.Second)
	}
}

func HandleServerCommunication(conn net.Conn, gui *GUI, u *UserData, connDone chan bool) {
	defer conn.Close()

	connectionFailed := make(chan bool)

	go serverRecieve(conn, gui, u, connectionFailed)

	go serverSend(conn, gui, u, connectionFailed)

	<-connectionFailed // Wait for communication failure
	connDone <- true   // tell StartNet connection died
}

func serverRecieve(conn net.Conn, gui *GUI, u *UserData, done chan bool) {
	decoder := json.NewDecoder(conn)
	for {
		var resp prot.ChatResponse
		err := decoder.Decode(&resp)
		if err != nil {
			done <- true
			return
		}

		if resp.Error != prot.NONE_STRING {
			Error.Printf("Server returned error: %s\n", resp.Error)
			continue
		}

		Info.Printf("Received %s response from server\n", resp.Type)

        switch resp.Type {
        case prot.REQ_ADD:
            gui.RespAdd(&resp, u)

        case prot.REQ_DEL:
            gui.RespDel(&resp)

        case prot.REQ_EDIT:
            gui.RespEdit(&resp)
        }
	}
}

func serverSend(conn net.Conn, gui *GUI, u *UserData, done chan bool) {
	for {
		select {
		case req := <-gui.OutgoingMessages:
			reqBytes, err := json.Marshal(req)
			if err != nil {
				Error.Printf("Failed to marshal outgoing request\n")
				continue
			}
			reqBytes = append(reqBytes, '\n')

            Info.Printf("Sending %s request to server\n", req.Type)

			_, err = conn.Write(reqBytes)
			if err != nil {
				done <- true
				return
			}
		case <-done:
			return
		}
	}
}

func CreateChatRequest(chatID int, reqType prot.RequestWhat, username string, messageContent string, messageID int) []byte {
	req := prot.ChatRequest{
		ChatID:         chatID,
		Type:           reqType,
		Username:       username,
		MessageContent: messageContent,
		MessageID:      messageID,
	}
	bytes, err := json.Marshal(req)
	if err != nil {
		Error.Printf("Failed to marshal chat request: %s\n", err)
		return nil
	}
	return append(bytes, '\n')
}
