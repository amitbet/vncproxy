package player

import (
	"log"
	"testing"
	"time"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
	"vncproxy/server"
	"encoding/binary"
)

type ServerMessageHandler struct {
	Conn         *server.ServerConn
	Fbs          *FbsReader
	firstSegDone bool
	startTime    int
}

func (handler *ServerMessageHandler) Consume(seg *common.RfbSegment) error {

	switch seg.SegmentType {
	case common.SegmentFullyParsedClientMessage:
		clientMsg := seg.Message.(common.ClientMessage)
		logger.Debugf("ClientUpdater.Consume:(vnc-server-bound) got ClientMessage type=%s", clientMsg.Type())
		switch clientMsg.Type() {

		case common.FramebufferUpdateRequestMsgType:
			if !handler.firstSegDone {
				handler.firstSegDone = true
				handler.startTime = int(time.Now().UnixNano() / int64(time.Millisecond))
			}
			handler.sendFbsMessage()
		}
		// server.FramebufferUpdateRequest:
	}
	return nil
}

func (h *ServerMessageHandler) sendFbsMessage() {
	var messageType uint8
	fbs := h.Fbs
	//conn := h.Conn
	binary.Read(fbs,binary.BigEndian,&messageType)
	bytes := messages[messageType].Read(fbs)
	h.Conn.Write(bytes)

	//seg, err := fbs.ReadSegment()
	//
	//now := int(time.Now().UnixNano() / int64(time.Millisecond))
	//if err != nil {
	//	logger.Error("TestServer.NewConnHandler: Error in reading FBS segment: ", err)
	//	return
	//}
	//timeSinceStart := now - h.startTime
	//
	//timeToWait := timeSinceStart - int(seg.timestamp)
	//
	//if timeToWait > 0 {
	//	time.Sleep(time.Duration(timeToWait) * time.Millisecond)
	//}
	//fmt.Printf("bytes: %v", seg.bytes)
	//conn.Write(seg.bytes)
}	

func loadFbsFile(filename string, conn *server.ServerConn) (*FbsReader, error) {
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
	chClient := make(chan common.ServerMessage)

	cfg := &server.ServerConfig{
		//SecurityHandlers: []SecurityHandler{&ServerAuthNone{}, &ServerAuthVNC{}},
		SecurityHandlers: []server.SecurityHandler{&server.ServerAuthNone{}},
		Encodings:        []common.Encoding{&encodings.RawEncoding{}, &encodings.TightEncoding{}, &encodings.CopyRectEncoding{}},
		PixelFormat:      common.NewPixelFormat(32),
		//ClientMessageCh:  chServer,
		//ServerMessageCh: chClient,
		ClientMessages: server.DefaultClientMessages,
		DesktopName:    []byte("workDesk"),
		Height:         uint16(768),
		Width:          uint16(1024),
		//NewConnHandler:  serverNewConnHandler,
	}

	cfg.NewConnHandler = func(cfg *server.ServerConfig, conn *server.ServerConn) error {
		fbs, err := loadFbsFile("/Users/amitbet/Dropbox/recording.rbs", conn)
		if err != nil {
			logger.Error("TestServer.NewConnHandler: Error in loading FBS: ", err)
			return err
		}
		conn.Listeners.AddListener(&ServerMessageHandler{conn, fbs, false, 0})
		return nil
	}

	url := "http://localhost:7777/"
	go server.WsServe(url, cfg)
	go server.TcpServe(":5904", cfg)

	// fbs, err := loadFbsFile("/Users/amitbet/vncRec/recording.rbs", cfg)
	// if err != nil {
	// 	logger.Error("TestServer.NewConnHandler: Error in loading FBS: ", err)
	// 	return
	// }

	// Process messages coming in on the ClientMessage channel.

	for {
		msg := <-chClient
		switch msg.Type() {
		default:
			log.Printf("Received message type:%v msg:%v\n", msg.Type(), msg)

		}
	}

}
