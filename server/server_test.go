package server

import (
	"context"
	"log"
	"net"
	"testing"
	"vncproxy/common"
	"vncproxy/encodings"
)

func TestServer(t *testing.T) {
	ln, err := net.Listen("tcp", ":5903")
	if err != nil {
		log.Fatalf("Error listen. %v", err)
	}

	chServer := make(chan common.ClientMessage)
	chClient := make(chan common.ServerMessage)

	cfg := &ServerConfig{
		//SecurityHandlers: []SecurityHandler{&ServerAuthNone{}, &ServerAuthVNC{}},
		SecurityHandlers: []SecurityHandler{&ServerAuthVNC{}},
		Encodings:        []common.Encoding{&encodings.RawEncoding{}, &encodings.TightEncoding{}, &encodings.CopyRectEncoding{}},
		PixelFormat:      common.NewPixelFormat(32),
		ClientMessageCh:  chServer,
		ServerMessageCh:  chClient,
		ClientMessages:   DefaultClientMessages,
		DesktopName:      []byte("workDesk"),
		Height:           uint16(768),
		Width:            uint16(1024),

	}
	go Serve(context.Background(), ln, cfg)

	// Process messages coming in on the ClientMessage channel.
	for {
		msg := <-chClient
		switch msg.Type() {
		default:
			log.Printf("Received message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}
