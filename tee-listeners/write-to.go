package listeners

import (
	"errors"
	"io"
	"vncproxy/common"
)

type WriteToListener struct {
	io.Writer
}

func (p *WriteToListener) Consume(seg *common.RfbSegment) error {
	switch seg.SegmentType {
	case common.SegmentMessageSeparator:
	case common.SegmentRectSeparator:
	case common.SegmentBytes:
		_, err := p.Writer.Write(seg.Bytes)
		return err

	default:
		return errors.New("undefined RfbSegment type")
	}
	return nil
}
