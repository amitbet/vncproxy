package encodings

import "vncproxy/common"

type PseudoEncoding struct {
	Typ int32
}

func (pe *PseudoEncoding) Type() int32 {
	return pe.Typ
}

func (pe *PseudoEncoding) Read(*common.PixelFormat, *common.Rectangle, *common.RfbReadHelper) (common.Encoding, error) {
	return pe, nil
}
