package player

import (
	"encoding/binary"

	"time"
	"vncproxy/client"
	"vncproxy/common"

	"vncproxy/logger"
	"vncproxy/server"
)

type FBSPlayListener struct {
	Conn             *server.ServerConn
	Fbs              *FbsReader
	serverMessageMap map[uint8]common.ServerMessage
	firstSegDone     bool
	startTime        int
}

func NewFBSPlayListener(conn *server.ServerConn, r *FbsReader) *FBSPlayListener {
	h := &FBSPlayListener{Conn: conn, Fbs: r}
	cm := client.BellMessage(0)
	h.serverMessageMap = make(map[uint8]common.ServerMessage)
	h.serverMessageMap[0] = &client.FramebufferUpdateMessage{}
	h.serverMessageMap[1] = &client.SetColorMapEntriesMessage{}
	h.serverMessageMap[2] = &cm
	h.serverMessageMap[3] = &client.ServerCutTextMessage{}

	return h
}
func (handler *FBSPlayListener) Consume(seg *common.RfbSegment) error {

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

func (h *FBSPlayListener) sendFbsMessage() {
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
