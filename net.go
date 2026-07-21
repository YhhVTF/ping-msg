package main

import (
	"encoding/json"
	"net"
	"time"

	"fyne.io/fyne/v2"
)

// StartNet: Connect to the server and show an error dialog if it fails
// Parameters:
//
//	gui (*GUI) - GUI elements
//  u (*UserData) - Information pertaining to users
func StartNet(gui *GUI, u *UserData) {
    // Prompt for a username
    gui.DialogLogin(u)
    for gui.Dialogs.Login != nil {}

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
		go HandleServerCommunication(conn, gui, u, connDone)

		<-connDone
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
		var resp ChatResponse
		err := decoder.Decode(&resp)
		if err != nil {
			done <- true
			return
		}

		if resp.Error != "" {
			Error.Printf("Server returned error: %s\n", resp.Error)
			continue
		}

		Info.Printf("Successfully recieved chat update for Chat ID %d. Total messages: %d\n", resp.ChatID, len(resp.Messages.Messages))

		fyne.Do(func() {
			for _, msg := range resp.Messages.Messages {
				card := NewMessage(msg.Content, msg.Username)
				gui.Containers.Chat.VBox.Add(card)
			}
			gui.Containers.Chat.VBox.Refresh()
			gui.Containers.Chat.VScroll.ScrollToBottom()
		})
	}
}

func serverSend(conn net.Conn, gui *GUI, u *UserData, done chan bool) {
	for {
		select {
		case msg := <-gui.OutgoingMessages:
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				Error.Printf("Failed to marshal outgoing request\n")
				continue
			}
			msgBytes = append(msgBytes, '\n')

			_, err = conn.Write(msgBytes)
			if err != nil {
				done <- true
				return
			}
		case <-done:
			return
		}
	}
}

func CreateChatRequest(chatID int, reqType RequestWhat, username string, messageContent string, messageID int) []byte {
	req := ChatRequest{
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
