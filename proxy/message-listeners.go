package proxy

import (
	"vncproxy/client"
	"vncproxy/common"
	"vncproxy/logger"
	"vncproxy/server"
)

type ClientUpdater struct {
	conn *client.ClientConn
}

// Consume recieves vnc-server-bound messages (Client messages) and updates the server part of the proxy
func (cc *ClientUpdater) Consume(seg *common.RfbSegment) error {
	//logger.Debugf("ClientUpdater.Consume (vnc-server-bound): got segment type=%s bytes: %v", seg.SegmentType, seg.Bytes)
	switch seg.SegmentType {

	case common.SegmentFullyParsedClientMessage:
		clientMsg := seg.Message.(common.ClientMessage)
		logger.Debugf("ClientUpdater.Consume:(vnc-server-bound) got ClientMessage type=%s", clientMsg.Type())
		switch clientMsg.Type() {

		case common.SetPixelFormatMsgType:
			// update pixel format
			logger.Debugf("ClientUpdater.Consume: updating pixel format")
			pixFmtMsg := clientMsg.(*server.SetPixelFormat)
			cc.conn.PixelFormat = pixFmtMsg.PF
		}

		err := clientMsg.Write(cc.conn)
		if err != nil {
			logger.Errorf("ClientUpdater.Consume (vnc-server-bound, SegmentFullyParsedClientMessage): problem writing to port: %s", err)
		}
		return err
	}
	return nil
}
