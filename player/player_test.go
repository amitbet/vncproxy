package player

import (
	"testing"
	"time"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
	"vncproxy/server"
)

func connectFbsFile(filename string, conn *server.ServerConn) (*FbsReader, error) {
	fbs, err := NewFbsReader(filename)
	if err != nil {
		logger.Error("failed to open fbs reader:", err)
		return nil, err
	}
	//NewFbsReader("/Users/amitbet/vncRec/recording.rbs")
	initMsg, err := fbs.ReadStartSession()
	if err != nil {
		logger.Error("failed to open read fbs start session:", err)
		return nil, err
	}
	conn.SetPixelFormat(&initMsg.PixelFormat)
	conn.SetHeight(initMsg.FBHeight)
	conn.SetWidth(initMsg.FBWidth)
	conn.SetDesktopName(string(initMsg.NameText))

	return fbs, nil
}

func TestServer(t *testing.T) {

	//chServer := make(chan common.ClientMessage)
	//chClient := make(chan common.ServerMessage)

	encs := []common.IEncoding{
		&encodings.RawEncoding{},
		&encodings.TightEncoding{},
		&encodings.EncCursorPseudo{},
		//encodings.TightPngEncoding{},
		&encodings.RREEncoding{},
		&encodings.ZLibEncoding{},
		&encodings.ZRLEEncoding{},
		&encodings.CopyRectEncoding{},
		&encodings.CoRREEncoding{},
		&encodings.HextileEncoding{},
	}

	cfg := &server.ServerConfig{
		//SecurityHandlers: []SecurityHandler{&ServerAuthNone{}, &ServerAuthVNC{}},
		SecurityHandlers: []server.SecurityHandler{&server.ServerAuthNone{}},
		Encodings:        encs,
		PixelFormat:      common.NewPixelFormat(32),
		ClientMessages:   server.DefaultClientMessages,
		DesktopName:      []byte("workDesk"),
		Height:           uint16(768),
		Width:            uint16(1024),
	}

	cfg.NewConnHandler = func(cfg *server.ServerConfig, conn *server.ServerConn) error {
		//fbs, err := loadFbsFile("/Users/amitbet/Dropbox/recording.rbs", conn)
		//fbs, err := loadFbsFile("/Users/amitbet/vncRec/recording.rbs", conn)
		fbs, err := connectFbsFile("/Users/amitbet/vncRec/recording1500554806.rbs", conn)

		if err != nil {
			logger.Error("TestServer.NewConnHandler: Error in loading FBS: ", err)
			return err
		}
		conn.Listeners.AddListener(NewFBSPlayListener(conn, fbs))
		return nil
	}

	url := "http://localhost:7777/"
	go server.WsServe(url, cfg)
	go server.TcpServe(":5904", cfg)

	for {
		time.Sleep(time.Minute)
	}

}
