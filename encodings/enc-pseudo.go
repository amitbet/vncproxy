package encodings

import (
	"io"
	"vncproxy/common"
)

//PseudoEncoding ...
type PseudoEncoding struct {
	Typ int32
}

//Type ...
func (pe *PseudoEncoding) Type() int32 {
	return pe.Typ
}

//WriteTo ...
func (pe *PseudoEncoding) WriteTo(w io.Writer) (n int, err error) {
	return 0, nil
}

//Read ...
func (pe *PseudoEncoding) Read(*common.PixelFormat, *common.Rectangle, *common.RfbReadHelper) (common.IEncoding, error) {
	return pe, nil
}
