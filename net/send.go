package net

import (
    "encoding/json"
    "net"

    "github.com/YhhVTF/ping-msg/gui"
    "github.com/YhhVTF/ping-msg/log"
)

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
