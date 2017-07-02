package encodings

import (
	"fmt"
	"vncproxy/common"
)

type TightPngEncoding struct {
}

func (*TightPngEncoding) Type() int32 { return int32(common.EncTightPng) }

func (t *TightPngEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	bytesPixel := calcTightBytePerPixel(pixelFmt)

	//var subencoding uint8
	compctl, err := r.ReadUint8()
	if err != nil {
		fmt.Printf("error in handling tight encoding: %v\n", err)
		return nil, err
	}
	fmt.Printf("bytesPixel= %d, subencoding= %d\n", bytesPixel, compctl)

	//move it to position (remove zlib flush commands)
	compType := compctl >> 4 & 0x0F

	fmt.Printf("afterSHL:%d\n", compType)
	switch compType {
	case TightPNG:
		len, err := r.ReadCompactLen()
		_, err = r.ReadBytes(len)

		if err != nil {
			return t, err
		}

	case TightFill:
		r.ReadBytes(int(bytesPixel))
	default:
		return nil, fmt.Errorf("unknown tight compression %d", compType)
	}
	return t, nil
}
