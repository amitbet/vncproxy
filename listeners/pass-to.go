package listeners

import "vncproxy/common"
import "io"

type PassListener struct {
	io.Writer
}

func (*PassListener) Consume(seg *common.RfbSegment) error {
	return nil
}
