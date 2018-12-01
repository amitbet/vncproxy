package server

import (
	"log"
	"testing"
	"github.com/amitbet/vncproxy/common"
	"github.com/amitbet/vncproxy/encodings"
)

func newServerConnHandler(cfg *ServerConfig, conn *ServerConn) error {

	return nil
}

func TestServer(t *testing.T) {

	//chServer := make(chan common.ClientMessage)
	chClient := make(chan common.ServerMessage)

	cfg := &ServerConfig{
		//SecurityHandlers: []SecurityHandler{&ServerAuthNone{}, &ServerAuthVNC{}},
		SecurityHandlers: []SecurityHandler{&ServerAuthVNC{"Ch_#!T@8"}},
		Encodings:        []common.IEncoding{&encodings.RawEncoding{}, &encodings.TightEncoding{}, &encodings.CopyRectEncoding{}},
		PixelFormat:      common.NewPixelFormat(32),
		ClientMessages:   DefaultClientMessages,
		DesktopName:      []byte("workDesk"),
		Height:           uint16(768),
		Width:            uint16(1024),
		NewConnHandler:   newServerConnHandler,
	}
	url := "http://localhost:8091/"
	go WsServe(url, cfg)
	go TcpServe(":5903", cfg)
	// Process messages coming in on the ClientMessage channel.
	for {
		msg := <-chClient
		switch msg.Type() {
		default:
			log.Printf("Received message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}
