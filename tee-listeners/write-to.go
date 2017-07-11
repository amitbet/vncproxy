package listeners

import (
	"io"
	"vncproxy/common"
	"vncproxy/logger"
)

type WriteTo struct {
	Writer io.Writer
	Name   string
}

func (p *WriteTo) Consume(seg *common.RfbSegment) error {

	logger.Debugf("WriteTo.Consume ("+p.Name+"): got segment type=%s bytes: %v", seg.SegmentType, seg.Bytes)
	switch seg.SegmentType {
	case common.SegmentMessageSeparator:
	case common.SegmentRectSeparator:
	case common.SegmentBytes:
		_, err := p.Writer.Write(seg.Bytes)
		if (err != nil) {
			logger.Errorf("WriteTo.Consume ("+p.Name+" SegmentBytes): problem writing to port: %s", err)
		}
		return err
	case common.SegmentFullyParsedClientMessage:

		clientMsg := seg.Message.(common.ClientMessage)
		logger.Debugf("WriteTo.Consume ("+p.Name+"): got ClientMessage type=%s", clientMsg.Type())
		err := clientMsg.Write(p.Writer)
		if (err != nil) {
			logger.Errorf("WriteTo.Consume ("+p.Name+" SegmentFullyParsedClientMessage): problem writing to port: %s", err)
		}
		return err
	default:
		//return errors.New("WriteTo.Consume: undefined RfbSegment type")
	}
	return nil
}



// type SendToClientMessageChan struct {
// 	Channel chan *common.ClientMessage
// }

// func (p *SendToClientMessageChan) Consume(seg *common.RfbSegment) error {
// 	switch seg.SegmentType {
// 	case common.SegmentMessageSeparator:
// 	case common.SegmentRectSeparator:
// 	case common.SegmentBytes:
// 	case common.SegmentFullyParsedClientMessage:
// 		p.Channel <- seg.Message.(*common.ClientMessage)
// 		//_, err := p.Writer.Write(seg.Bytes)
// 		//return err

// 	default:
// 		//return errors.New("undefined RfbSegment type")
// 	}
// 	return nil
// }
