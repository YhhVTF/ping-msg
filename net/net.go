package net

import (
	"encoding/json"
	"net"
    "os"
	"time"

    "github.com/YhhVTF/ping-msg/global"
    "github.com/YhhVTF/ping-msg/gui"
    "github.com/YhhVTF/ping-msg/log"
    "github.com/YhhVTF/ping-msg/opt"
    "github.com/YhhVTF/ping-msg/protocol"
    "github.com/YhhVTF/ping-msg/user"
)

// StartNet: Connect to the server and show an error dialog if it fails
// Parameters:
//
//	gui (*gui.GUI) - GUI elements
//  u (*user.UserData) - Information pertaining to users
//  opt (*options.Options) - Options/settings
func StartNet(gui *gui.GUI, u *user.UserData, opt *options.Options) {
    // Prompt for a username
    gui.Chat.DialogLogin(u, opt)
    for gui.Chat.Dialogs.Login != nil {}

        addr := "127.0.0.1:5555"
        if len(os.Args) > 1 {
            addr = os.Args[1]
        }

	// Until Ping has been quit...
	for !ping.Quit {
		// Connect to the server
		log.Info.Printf("Connecting to server\n")

        conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Error.Printf("Failed to connect to server: %s\n", err)
			if gui.Chat.Dialogs.ConnectionIssues == nil {
				gui.Chat.DialogConnectionIssues(err, opt)
			}
			time.Sleep(30 * time.Second)
			continue
		}

        ping.Connected = true
		log.Info.Printf("Successfully connected to server\n")
        // Dismiss connection issues dialog after connecting successfully
        if gui.Chat.Dialogs.ConnectionIssues != nil {
            gui.Chat.Dialogs.ConnectionIssues.Dismiss()
            gui.Chat.Dialogs.ConnectionIssues = nil
        }

		connDone := make(chan bool)
		go HandleServerCommunication(conn, gui, u, connDone)

		<-connDone
        ping.Connected = false
		log.Error.Printf("Connection lost. Reconnecting in 5 seconds...\n")
		time.Sleep(5 * time.Second)
	}
}

func HandleServerCommunication(conn net.Conn, gui *gui.GUI, u *user.UserData, connDone chan bool) {
	defer conn.Close()

	connectionFailed := make(chan bool)

	go serverRecieve(conn, gui, u, connectionFailed)

	go serverSend(conn, gui, u, connectionFailed)

	<-connectionFailed // Wait for communication failure
	connDone <- true   // tell StartNet connection died
}

func serverRecieve(conn net.Conn, gui *gui.GUI, u *user.UserData, done chan bool) {
	decoder := json.NewDecoder(conn)
	for {
		var resp prot.ChatResponse
		err := decoder.Decode(&resp)
		if err != nil {
			done <- true
			return
		}

		if resp.Error != prot.NONE_STRING {
			log.Error.Printf("Server returned error: %s\n", resp.Error)
			continue
		}

		log.Info.Printf("Received %s response from server\n", resp.Type)

        switch resp.Type {
        case prot.REQ_ADD:
            gui.Chat.RespAdd(&resp, u)

        case prot.REQ_DEL:
            gui.Chat.RespDel(&resp)

        case prot.REQ_EDIT:
            gui.Chat.RespEdit(&resp)
        }
	}
}

func serverSend(conn net.Conn, gui *gui.GUI, u *user.UserData, done chan bool) {
	for {
		select {
		case req := <-gui.Chat.OutgoingMessages:
			reqBytes, err := json.Marshal(req)
			if err != nil {
				log.Error.Printf("Failed to marshal outgoing request\n")
				continue
			}
			reqBytes = append(reqBytes, '\n')

            log.Info.Printf("Sending %s request to server\n", req.Type)

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
		log.Error.Printf("Failed to marshal chat request: %s\n", err)
		return nil
	}
	return append(bytes, '\n')
}
