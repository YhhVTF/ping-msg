package main

import (
	"encoding/json"
	"net"
    "os"
	"time"

	"fyne.io/fyne/v2"

    "ping/protocol"
)

var Connected = false

// StartNet: Connect to the server and show an error dialog if it fails
// Parameters:
//
//	gui (*GUI) - GUI elements
//  u (*UserData) - Information pertaining to users
func StartNet(gui *GUI, u *UserData) {
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

		Info.Printf("Received response from server\n")

        switch resp.Type {
        case prot.REQ_ADD:
            fyne.Do(func() {
                for _, msg := range resp.Messages {
                    msgWidget := gui.NewMessage(
                        msg.Content, msg.Username,
                        time.Unix(msg.Time, 0).Format("3:04 PM"),
                        msg.Username == u.ThisUser,
                        func() {
                            req := prot.ChatRequest{
                                ChatID: 1,
                                MessageContent: prot.NONE_STRING,
                                MessageID: msg.ID,
                                Type: prot.REQ_DEL,
                                Username: u.ThisUser,
                            }
                            gui.OutgoingMessages <- req
                        },
                        func(newText string) {
                            req := prot.ChatRequest{
                                ChatID: 1,
                                MessageContent: newText,
                                MessageID: msg.ID,
                                Type: prot.REQ_EDIT,
                                Username: u.ThisUser,
                            }
                            gui.OutgoingMessages <- req
                        },
                    )
                    gui.Widgets.Messages[msg.ID] = msgWidget
                    gui.Containers.Chat.VBox.Add(msgWidget.Base)
                }
                gui.Containers.Chat.VBox.Refresh()
                gui.Containers.Chat.VScroll.ScrollToBottom()
            })
        case prot.REQ_DEL:
            fyne.Do(func() {
                if _, exists := gui.Widgets.Messages[resp.MessageID]; exists {
                    gui.Widgets.Messages[resp.MessageID].Base.Hide()
                    delete(gui.Widgets.Messages, resp.MessageID)
                }
            })
        case prot.REQ_EDIT:
            fyne.Do(func() {
                if _, exists := gui.Widgets.Messages[resp.MessageID]; exists {
                    gui.Widgets.Messages[resp.MessageID].Content.Text =
                        resp.Messages[0].Content
                    gui.Widgets.Messages[resp.MessageID].Content.Refresh()
                }
            })
        }
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
