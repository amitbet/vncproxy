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
			pixFmtMsg := clientMsg.(*server.MsgSetPixelFormat)
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

type ServerUpdater struct {
	conn *server.ServerConn
}

func (p *ServerUpdater) Consume(seg *common.RfbSegment) error {

	logger.Debugf("WriteTo.Consume (ServerUpdater): got segment type=%s", seg.SegmentType)
	switch seg.SegmentType {
	case common.SegmentMessageSeparator:
	case common.SegmentRectSeparator:
	case common.SegmentServerInitMessage:
		serverInitMessage := seg.Message.(*common.ServerInit)
		p.conn.SetHeight(serverInitMessage.FBHeight)
		p.conn.SetWidth(serverInitMessage.FBWidth)
		p.conn.SetDesktopName(string(serverInitMessage.NameText))
		p.conn.SetPixelFormat(&serverInitMessage.PixelFormat)

	case common.SegmentBytes:
		_, err := p.conn.Write(seg.Bytes)
		if err != nil {
			logger.Errorf("WriteTo.Consume (ServerUpdater SegmentBytes): problem writing to port: %s", err)
		}
		return err
	case common.SegmentFullyParsedClientMessage:

		clientMsg := seg.Message.(common.ClientMessage)
		logger.Debugf("WriteTo.Consume (ServerUpdater): got ClientMessage type=%s", clientMsg.Type())
		err := clientMsg.Write(p.conn)
		if err != nil {
			logger.Errorf("WriteTo.Consume (ServerUpdater SegmentFullyParsedClientMessage): problem writing to port: %s", err)
		}
		return err
	default:
		//return errors.New("WriteTo.Consume: undefined RfbSegment type")
	}
	return nil
}
