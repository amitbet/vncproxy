package encodings

import (
	"encoding/binary"
	"io"
	"vncproxy/common"
)

type CoRREEncoding struct {
	numSubRects     uint32
	backgroundColor []byte
	subRectData     []byte
}

func (z *CoRREEncoding) Type() int32 {
	return 4
}

func (z *CoRREEncoding) WriteTo(w io.Writer) (n int, err error) {
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

func (z *CoRREEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.IEncoding, error) {
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

	//read all individual rects (color=BPP + x=16b + y=16b + w=16b + h=16b)
	z.subRectData, err = r.ReadBytes(int(numOfSubrectangles) * (bytesPerPixel + 4))
	if err != nil {
		return nil, err
	}

	return z, nil
}
