package encodings

import (
	"encoding/binary"
	"io"
	"vncproxy/common"
)

type RREEncoding struct {
	//Colors []Color
	numSubRects     uint32
	backgroundColor []byte
	subRectData     []byte
}

func (z *RREEncoding) WriteTo(w io.Writer) (n int, err error) {
	binary.Write(w, binary.BigEndian, z.numSubRects)
	if err != nil {
		return 0, err
	}

	w.Write(z.backgroundColor)
	if err != nil {
		return 0, err
	}

	w.Write(z.subRectData)

	if err != nil {
		return 0, err
	}
	b := len(z.backgroundColor) + len(z.subRectData) + 4
	return b, nil
}

func (z *RREEncoding) Type() int32 {
	return 2
}
func (z *RREEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	bytesPerPixel := int(pixelFmt.BPP / 8)
	numOfSubrectangles, err := r.ReadUint32()
	if err != nil {
		return nil, err
	}
	z.numSubRects = numOfSubrectangles

	//read whole-rect background color
	z.backgroundColor, err = r.ReadBytes(bytesPerPixel)
	if err != nil {
		return nil, err
	}

	//read all individual rects (color=bytesPerPixel + x=16b + y=16b + w=16b + h=16b)
	z.subRectData, err = r.ReadBytes(int(numOfSubrectangles) * (bytesPerPixel + 8)) // x+y+w+h=8 bytes
	if err != nil {
		return nil, err
	}
	return z, nil
}
