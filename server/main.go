package server

import (
	"log"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
)

func newServerConnHandler(cfg *ServerConfig, conn *ServerConn) error {
	return nil
}

func main() {

	//chServer := make(chan common.ClientMessage)
	chClient := make(chan common.ServerMessage)

	cfg := &ServerConfig{
		//SecurityHandlers: []SecurityHandler{&ServerAuthNone{}, &ServerAuthVNC{}},
		SecurityHandlers: []SecurityHandler{&ServerAuthVNC{"Ch_#!T@8"}},
		Encodings:        []common.Encoding{&encodings.RawEncoding{}, &encodings.TightEncoding{}, &encodings.CopyRectEncoding{}},
		PixelFormat:      common.NewPixelFormat(32),
		//ClientMessageCh:  chServer,
		ServerMessageCh: chClient,
		ClientMessages:  DefaultClientMessages,
		DesktopName:     []byte("workDesk"),
		Height:          uint16(768),
		Width:           uint16(1024),
		NewConnHandler:  newServerConnHandler,
	}

	loadFbsFile("c:\\Users\\betzalel\\Dropbox\\recording.rbs", cfg)

	url := "http://localhost:8091/"
	go WsServe(url, cfg)
	go TcpServe(":5904", cfg)
	// Process messages coming in on the ClientMessage channel.
	for {
		msg := <-chClient
		switch msg.Type() {
		default:
			log.Printf("Received message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}

func loadFbsFile(filename string, cfg *ServerConfig) {
	fbs, err := NewFbsReader(filename)
	if err != nil {
		logger.Error("failed to open fbs reader:", err)
	}
	//NewFbsReader("/Users/amitbet/vncRec/recording.rbs")
	initMsg, err := fbs.ReadStartSession()
	if err != nil {
		logger.Error("failed to open read fbs start session:", err)
	}

	cfg.PixelFormat = &initMsg.PixelFormat
	cfg.Height = initMsg.FBHeight
	cfg.Width = initMsg.FBWidth
	cfg.DesktopName = initMsg.NameText
}
