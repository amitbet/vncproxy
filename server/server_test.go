package server

import (
	"context"
	"log"
	"testing"
	"vncproxy/common"
	"vncproxy/encodings"
)

func TestServer(t *testing.T) {

	chServer := make(chan common.ClientMessage)
	chClient := make(chan common.ServerMessage)

	cfg := &ServerConfig{
		//SecurityHandlers: []SecurityHandler{&ServerAuthNone{}, &ServerAuthVNC{}},
		SecurityHandlers: []SecurityHandler{&ServerAuthVNC{"Ch_#!T@8"}},
		Encodings:        []common.Encoding{&encodings.RawEncoding{}, &encodings.TightEncoding{}, &encodings.CopyRectEncoding{}},
		PixelFormat:      common.NewPixelFormat(32),
		ClientMessageCh:  chServer,
		ServerMessageCh:  chClient,
		ClientMessages:   DefaultClientMessages,
		DesktopName:      []byte("workDesk"),
		Height:           uint16(768),
		Width:            uint16(1024),
	}
	url := "http://localhost:8091/"
	go WsServe(url, context.Background(), cfg)
	go TcpServe(":5903", context.Background(), cfg)
	// Process messages coming in on the ClientMessage channel.
	for {
		msg := <-chClient
		switch msg.Type() {
		default:
			log.Printf("Received message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}
