package encodings

import (
	"io"
	"math"
	"vncproxy/common"
)

type EncCursorPseudo struct {
}

func (pe *EncCursorPseudo) Type() int32 {
	return int32(common.EncCursorPseudo)
}
func (z *EncCursorPseudo) WriteTo(w io.Writer) (n int, err error) {
	return 0, nil
}
func (pe *EncCursorPseudo) Read(pf *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	if rect.Width*rect.Height == 0 {
		return pe, nil
	}

	bytesPixel := int(pf.BPP / 8) //calcTightBytePerPixel(pf)
	r.ReadBytes(int(rect.Width*rect.Height) * bytesPixel)
	mask := ((rect.Width + 7) / 8) * rect.Height
	r.ReadBytes(int(math.Floor(float64(mask))))
	return pe, nil
}
