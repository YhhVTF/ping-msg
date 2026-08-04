package net

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"

	"github.com/coder/websocket"

	ping "github.com/YhhVTF/ping-msg/global"
	"github.com/YhhVTF/ping-msg/gui"
	"github.com/YhhVTF/ping-msg/gui/dialogs"
	"github.com/YhhVTF/ping-msg/log"
	options "github.com/YhhVTF/ping-msg/opt"
	prot "github.com/YhhVTF/ping-msg/protocol"
	"github.com/YhhVTF/ping-msg/user"
)

// StartNet connects to the server, registers the selected username, and then
// starts the chat request/response loops.
func StartNet(gui *gui.GUI, u *user.UserCache, opt *options.Options) {
	fyne.DoAndWait(func() {
        gui.Dialogs.Login = dialogs.InitDialogLogin(gui.Window, u, opt)
    })
	for gui.Dialogs.Login.Dialog != nil {
		time.Sleep(10 * time.Millisecond)
	}
    gui.Dialogs.Login = nil

	endpoint := "wss://ping.da5h1n.uk:5555/ws"
	if len(os.Args) > 1 {
		endpoint = websocketEndpoint(os.Args[1])
	}

	for !ping.Quit {
		log.Info.Printf("Connecting to server\n")
		wsConn, _, err := websocket.Dial(context.Background(), endpoint, nil)
		if err == nil {
			conn := websocket.NetConn(context.Background(), wsConn, websocket.MessageText)
			decoder, registerErr := registerUser(conn, u)
			if registerErr == nil {
				ping.Connected = true
				log.Info.Printf("Successfully connected as user %d\n", u.ThisUserID)
				if gui.Dialogs.ConnectionIssues != nil {
					fyne.DoAndWait(func() {
                        gui.Dialogs.ConnectionIssues.Dialog.Dismiss()
                    })
					gui.Dialogs.ConnectionIssues = nil
				}
				connDone := make(chan bool)
				go HandleServerCommunication(conn, decoder, gui, u, connDone)
				<-connDone
				ping.Connected = false
				log.Error.Printf("Connection lost. Reconnecting in 5 seconds...\n")
				time.Sleep(5 * time.Second)
				continue
			}
			err = registerErr
			conn.Close()
		}
		log.Error.Printf("Failed to connect or register with server: %s\n", err)
		if gui.Dialogs.ConnectionIssues == nil {
			fyne.DoAndWait(func() {
                dialogs.InitDialogConnIssues(gui.Window, err, opt)
            })
		}
		time.Sleep(5 * time.Second)
	}
}

func websocketEndpoint(value string) string {
	if !strings.HasPrefix(value, "ws://") && !strings.HasPrefix(value, "wss://") {
		value = "ws://" + value
	}
	addressAndPath := strings.TrimPrefix(strings.TrimPrefix(value, "wss://"), "ws://")
	if !strings.Contains(addressAndPath, "/") {
		value += "/ws"
	}
	return value
}

func registerUser(conn net.Conn, u *user.UserCache) (*json.Decoder, error) {
	request := prot.UserRequest{RequestType: prot.UserRequestRegister, Username: u.ThisUsername}
	if err := json.NewEncoder(conn).Encode(request); err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(conn)
	var response prot.UserResponse
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, fmt.Errorf("%s", response.Error)
	}
	if response.User.ID == prot.SERVER_USER_ID || response.User.Username == "" {
		return nil, fmt.Errorf("server returned an invalid user registration")
	}

	u.ThisUserID = response.User.ID
	u.ThisUsername = response.User.Username
	if u.Users == nil {
		u.Users = make(map[int]user.User)
	}
	u.Users[response.User.ID] = user.User{ID: response.User.ID, Username: response.User.Username}
	return decoder, nil
}

func HandleServerCommunication(conn net.Conn, decoder *json.Decoder, gui *gui.GUI, u *user.UserCache, connDone chan bool) {
	defer conn.Close()

	done := make(chan struct{})
	var once sync.Once
	signalDone := func() { once.Do(func() { close(done) }) }

	go serverRecieve(decoder, gui, u, signalDone)
	go serverSend(conn, gui, done, signalDone)

	<-done
	connDone <- true
}

func cacheUsers(u *user.UserCache, users []prot.User) {
	if u.Users == nil {
		u.Users = make(map[int]user.User)
	}
	for _, profile := range users {
		u.Users[profile.ID] = user.User{ID: profile.ID, Username: profile.Username}
	}
}

func serverRecieve(decoder *json.Decoder, gui *gui.GUI, u *user.UserCache, signalDone func()) {
	for {
		var resp prot.ChatResponse
		if err := decoder.Decode(&resp); err != nil {
			signalDone()
			return
		}

		if resp.Error != prot.NONE_STRING {
			log.Error.Printf("Server returned error: %s\n", resp.Error)
			continue
		}

		cacheUsers(u, resp.Users)

		switch resp.Type {
		case prot.REQ_ADD:
			fyne.Do(func() { gui.Chat.RespAdd(&resp, u) })
		case prot.REQ_DEL:
			fyne.Do(func() { gui.Chat.RespDel(&resp) })
		case prot.REQ_EDIT:
			fyne.Do(func() { gui.Chat.RespEdit(&resp) })
		}
	}
}

func serverSend(conn net.Conn, gui *gui.GUI, done <-chan struct{}, signalDone func()) {
	encoder := json.NewEncoder(conn)
	for {
		select {
		case req := <-gui.Chat.OutgoingMessages:
			log.Info.Printf("Sending %s request to server\n", req.Type)
			if err := encoder.Encode(req); err != nil {
				log.Error.Printf("Failed to send outgoing request: %s\n", err)
				signalDone()
				return
			}
		case <-done:
			return
		}
	}
}

func CreateChatRequest(chatID int, reqType prot.RequestWhat, userID int, messageContent string, messageID int) []byte {
	req := prot.ChatRequest{ChatID: chatID, Type: reqType, UserID: userID, MessageContent: messageContent, MessageID: messageID}
	bytes, err := json.Marshal(req)
	if err != nil {
		log.Error.Printf("Failed to marshal chat request: %s\n", err)
		return nil
	}
	return append(bytes, '\n')
}
