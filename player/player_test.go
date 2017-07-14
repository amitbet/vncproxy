package player

import (
	"encoding/binary"
	"log"
	"testing"
	"time"
	"vncproxy/client"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
	"vncproxy/server"
)

type ServerMessageHandler struct {
	Conn             *server.ServerConn
	Fbs              *FbsReader
	serverMessageMap map[uint8]common.ServerMessage
	firstSegDone     bool
	startTime        int
}

func NewServerMessageHandler(conn *server.ServerConn, r *FbsReader) *ServerMessageHandler {
	h := &ServerMessageHandler{Conn: conn, Fbs: r}
	cm := client.BellMessage(0)
	h.serverMessageMap = make(map[uint8]common.ServerMessage)
	h.serverMessageMap[0] = &client.FramebufferUpdateMessage{}
	h.serverMessageMap[1] = &client.SetColorMapEntriesMessage{}
	h.serverMessageMap[2] = &cm
	h.serverMessageMap[3] = &client.ServerCutTextMessage{}

	return h
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
	//messages := make(map[uint8]common.ServerMessage)
	fbs := h.Fbs
	//conn := h.Conn
	err := binary.Read(fbs, binary.BigEndian, &messageType)
	if err != nil {
		logger.Error("TestServer.NewConnHandler: Error in reading FBS segment: ", err)
		return
	}
	//common.IClientConn{}
	binary.Write(h.Conn, binary.BigEndian, messageType)
	msg := h.serverMessageMap[messageType]
	if msg == nil {
		logger.Error("TestServer.NewConnHandler: Error unknown message type: ", messageType)
		return
	}
	timeSinceStart := int(time.Now().UnixNano()/int64(time.Millisecond)) - h.startTime
	timeToSleep := fbs.currentTimestamp - timeSinceStart
	if timeToSleep > 0 {
		time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	}

	err = msg.CopyTo(fbs, h.Conn, fbs)
	if err != nil {
		logger.Error("TestServer.NewConnHandler: Error in reading FBS segment: ", err)
		return
	}
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
		conn.Listeners.AddListener(NewServerMessageHandler(conn, fbs))
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
