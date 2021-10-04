package proxy

import (
	"github.com/amitbet/vncproxy/client"
	"github.com/amitbet/vncproxy/common"
	"github.com/amitbet/vncproxy/logger"
	"github.com/amitbet/vncproxy/server"
)

type ClientUpdater struct {
	conn                   *client.ClientConn
	suppressedMessageTypes []common.ClientMessageType
	overrideEncodings      []common.EncodingType
}

// Consume recieves vnc-server-bound messages (Client messages) and updates the server part of the proxy
func (cc *ClientUpdater) Consume(seg *common.RfbSegment) error {
	logger.Tracef("ClientUpdater.Consume (vnc-server-bound): got segment type=%s bytes: %v", seg.SegmentType, seg.Bytes)
	switch seg.SegmentType {

	case common.SegmentFullyParsedClientMessage:
		clientMsg := seg.Message.(common.ClientMessage)
		logger.Debugf("ClientUpdater.Consume:(vnc-server-bound) got ClientMessage type=%s", clientMsg.Type())

		switch clientMsg.Type() {
		case common.SetEncodingsMsgType:
			if len(cc.overrideEncodings) > 0 {
				logger.Debugf("ClientUpdater.Consume:(vnc-server-bound) overriding supported encodings with %v", cc.overrideEncodings)
				encodingsMsg := clientMsg.(*server.MsgSetEncodings)
				encodingsMsg.EncNum = uint16(len(cc.overrideEncodings))
				encodingsMsg.Encodings = cc.overrideEncodings
			}
		case common.SetPixelFormatMsgType:
			// update pixel format
			logger.Debugf("ClientUpdater.Consume: updating pixel format")
			pixFmtMsg := clientMsg.(*server.MsgSetPixelFormat)
			cc.conn.PixelFormat = pixFmtMsg.PF
		}

		suppressMessage := false
		for _, suppressed := range cc.suppressedMessageTypes {
			if suppressed == clientMsg.Type() {
				suppressMessage = true
				break
			}
		}
		if suppressMessage {
			logger.Infof("ClientUpdater.Consume:(vnc-server-bound) Suppressing client message type=%s", clientMsg.Type())
			return nil
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

	logger.Debugf("WriteTo.Consume (ServerUpdater): got segment type=%s, object type:%d", seg.SegmentType, seg.UpcomingObjectType)
	switch seg.SegmentType {
	case common.SegmentMessageStart:
	case common.SegmentRectSeparator:
	case common.SegmentServerInitMessage:
		serverInitMessage := seg.Message.(*common.ServerInit)
		p.conn.SetHeight(serverInitMessage.FBHeight)
		p.conn.SetWidth(serverInitMessage.FBWidth)
		p.conn.SetDesktopName(string(serverInitMessage.NameText))
		p.conn.SetPixelFormat(&serverInitMessage.PixelFormat)

	case common.SegmentBytes:
		logger.Debugf("WriteTo.Consume (ServerUpdater SegmentBytes): got bytes len=%d", len(seg.Bytes))
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
