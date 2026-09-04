package net

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
    "slices"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"

	"github.com/coder/websocket"

    "github.com/YhhVTF/ping-msg/chat"
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
func StartNet(
    gui *gui.GUI, u *user.UserCache, c *chat.ChatCache, opt *options.Options,
) {
	fyne.DoAndWait(func() {
        gui.Dialogs.Login =
            dialogs.InitDialogLogin(gui.Window, gui.ScreenManager, u, opt)
    })
	for gui.Dialogs.Login.Dialog != nil {
		time.Sleep(10 * time.Millisecond)
	}
    gui.Dialogs.Login = nil

	endpoint := "wss://ping.da5h1n.uk/ws"
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
				log.Info.Printf("Successfully connected as user %s\n", u.ThisUser.Username)

                fyne.DoAndWait(func() { gui.Sidebar.InitChatsList(c, u) })

				if gui.Dialogs.ConnectionIssues != nil {
					fyne.DoAndWait(func() {
                        gui.Dialogs.ConnectionIssues.Dialog.Dismiss()
                    })
					gui.Dialogs.ConnectionIssues = nil
				}
				connDone := make(chan bool)
				go HandleServerCommunication(conn, decoder, gui, c, u, opt, connDone)
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
                gui.Dialogs.ConnectionIssues =
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
	request := prot.UserRequest{Type: prot.UserRequestRegister, Username: u.ThisUser.Username}
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
	if response.User.Username == prot.NONE_STRING || response.User.Username == "" {
		return nil, fmt.Errorf("server returned an invalid user registration")
	}

    u.CacheUserBack(response.User)
    u.SetThisUser(response.User.Username)
	return decoder, nil
}

func HandleServerCommunication(conn net.Conn, decoder *json.Decoder, gui *gui.GUI, c *chat.ChatCache, u *user.UserCache, opt *options.Options, connDone chan bool) {
	defer conn.Close()

	done := make(chan struct{})
	var once sync.Once
	signalDone := func() { once.Do(func() { close(done) }) }

	go serverRecieve(decoder, gui, c, u, opt, signalDone)
	go serverSend(conn, gui, done, signalDone)

	<-done
	connDone <- true
}

func serverRecieve(
    decoder *json.Decoder, gui *gui.GUI, c *chat.ChatCache, u *user.UserCache,
    opt *options.Options, signalDone func(),
) {
	for {
        var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			signalDone()
			return
		}

        var chatResp prot.ChatResponse
        if err := json.Unmarshal(raw, &chatResp); err == nil && chatResp.Type != "" {
            log.Info.Printf("Received %s from server\n", chatResp.Type)

            if chatResp.Error != prot.NONE_STRING && chatResp.Error != "" {
                log.Error.Printf("Server returned error: %s\n", chatResp.Error)
                continue
            }

            switch chatResp.Type {
            case prot.REQ_ADD:
                u.CacheUserFront(chatResp.Users) // Cache the usernames of users involved
                fyne.Do(func() { gui.Chat.RespAdd(&chatResp, ping.ChatCache, u, opt) })
            case prot.REQ_DEL:
                fyne.Do(func() { gui.Chat.RespDel(&chatResp, ping.ChatCache, opt) })
            case prot.REQ_EDIT:
                fyne.Do(func() { gui.Chat.RespEdit(&chatResp, ping.ChatCache, u) })
            }
        }
        var userResp prot.UserResponse
        if err := json.Unmarshal(raw, &userResp); err == nil && userResp.Type != "" {
            log.Info.Printf("Received %s from server\n", userResp.Type)

            if userResp.Error != prot.NONE_STRING && userResp.Error != "" {
                log.Error.Printf("Server returned error: %s\n", chatResp.Error)
                continue
            }
        }
        var chatMDResp prot.ChatMetadataResponse
        if err := json.Unmarshal(raw, &chatMDResp); err == nil && chatMDResp.Type != "" {
            log.Info.Printf("Received %s from server\n", chatMDResp.Type)

            if chatMDResp.Error != prot.NONE_STRING && chatMDResp.Error != "" {
                log.Error.Printf("Server returned error: %s\n", chatResp.Error)
                continue
            }

            switch chatMDResp.Type {
            case prot.ChatMetadataRequestCreate:
                c.Chats[chatMDResp.ChatID[0]] = chat.NewChat(&chatMDResp.Metadata[0])
                u.ThisUser.MemberOf = append(u.ThisUser.MemberOf, chatMDResp.ChatID[0])
                fyne.Do(func() { gui.Sidebar.ChatMetadataRespCreate(u) })
            case prot.ChatMetadataRequestGet:
                for _, chatMD := range chatMDResp.Metadata {
                    c.Chats[chatMD.ID] = chat.NewChat(&chatMD)
                }
                fyne.Do(func() { gui.Sidebar.Widgets.ChatsList.Refresh() })
            case prot.ChatMetadataRequestJoin:
                c.Chats[chatMDResp.ChatID[0]] = chat.NewChat(&chatMDResp.Metadata[0])
                u.ThisUser.MemberOf = append(u.ThisUser.MemberOf, chatMDResp.ChatID[0])
                fyne.Do(func() { gui.Sidebar.ChatMetadataRespJoin(u) })
            case prot.ChatMetadataRequestLeave:
                delete(c.Chats, chatMDResp.ChatID[0])
                i := slices.Index(u.ThisUser.MemberOf, chatMDResp.ChatID[0])
                u.ThisUser.MemberOf = slices.Delete(u.ThisUser.MemberOf, i, i+1)

                fyne.Do(func() { gui.Sidebar.Widgets.ChatsList.Refresh() })
            }
        }
	}
}

func serverSend(conn net.Conn, gui *gui.GUI, done <-chan struct{}, signalDone func()) {
	encoder := json.NewEncoder(conn)
	for {
		select {
		case req := <-gui.OutgoingRequests:
			if err := encoder.Encode(req); err != nil {
				log.Error.Printf("Failed to send outgoing request %s: %s\n", req.Type, err)
				signalDone()
				return
			}
            log.Info.Printf("Sent request %s to server\n", req.Type)
        case req := <-gui.Sidebar.OutgoingReqChatMetadata:
            if err := encoder.Encode(req); err != nil {
                log.Error.Printf("Failed to send outgoing request %s: %s\n", req.Type, err)
                signalDone()
                return
            }
            log.Info.Printf("Sent request %s to server\n", req.Type)
		case <-done:
			return
		}
	}
}

func CreateChatRequest(chatID int, reqType prot.RequestWhat, username string, messageContent string, messageID int) []byte {
	req := prot.ChatRequest{ChatID: chatID, Type: reqType, Username: username, MessageContent: messageContent, MessageID: messageID}
	bytes, err := json.Marshal(req)
	if err != nil {
		log.Error.Printf("Failed to marshal chat request: %s\n", err)
		return nil
	}
	return append(bytes, '\n')
}
