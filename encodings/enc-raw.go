package encodings

import (
	"bytes"
	"io"
	"github.com/amitbet/vncproxy/common"
)

// RawEncoding is raw pixel data sent by the server.
type RawEncoding struct {
	bytes []byte
}

func (*RawEncoding) Type() int32 {
	return 0
}
func (z *RawEncoding) WriteTo(w io.Writer) (n int, err error) {
	return w.Write(z.bytes)
}
func (*RawEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.IEncoding, error) {

	bytesPerPixel := int(pixelFmt.BPP / 8)

	bytes := &bytes.Buffer{}
	for y := uint16(0); y < rect.Height; y++ {
		for x := uint16(0); x < rect.Width; x++ {
			if bts, err := r.ReadBytes(bytesPerPixel); err != nil {
				StoreBytes(bytes, bts)
				return nil, err
			}
		}
	}

	return &RawEncoding{bytes.Bytes()}, nil
}
