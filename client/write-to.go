package client

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

	logger.Debugf("WriteTo.Consume ("+p.Name+"): got segment type=%s", seg.SegmentType)
	switch seg.SegmentType {
	case common.SegmentMessageStart:
	case common.SegmentRectSeparator:
	case common.SegmentBytes:
		_, err := p.Writer.Write(seg.Bytes)
		if err != nil {
			logger.Errorf("WriteTo.Consume ("+p.Name+" SegmentBytes): problem writing to port: %s", err)
		}
		return err
	case common.SegmentFullyParsedClientMessage:

		clientMsg := seg.Message.(common.ClientMessage)
		logger.Debugf("WriteTo.Consume ("+p.Name+"): got ClientMessage type=%s", clientMsg.Type())
		err := clientMsg.Write(p.Writer)
		if err != nil {
			logger.Errorf("WriteTo.Consume ("+p.Name+" SegmentFullyParsedClientMessage): problem writing to port: %s", err)
		}
		return err
	default:
		//return errors.New("WriteTo.Consume: undefined RfbSegment type")
	}
	return nil
}
