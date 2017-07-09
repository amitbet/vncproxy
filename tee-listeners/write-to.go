package listeners

import (
	"errors"
	"io"
	"vncproxy/common"
)

type WriteTo struct {
	Writer io.Writer
	Name   string
}

func (p *WriteTo) Consume(seg *common.RfbSegment) error {
	switch seg.SegmentType {
	case common.SegmentMessageSeparator:
	case common.SegmentRectSeparator:
	case common.SegmentBytes:
		_, err := p.Writer.Write(seg.Bytes)
		return err
	case common.SegmentFullyParsedClientMessage:
		clientMsg := seg.Message.(common.ClientMessage)
		clientMsg.Write(p.Writer)
	default:
		return errors.New("undefined RfbSegment type")
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
