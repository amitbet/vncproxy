package encodings

import (
	"encoding/binary"
	"io"
	"github.com/amitbet/vncproxy/common"
)

type CopyRectEncoding struct {
	//Colors       []Color
	copyRectSrcX uint16
	copyRectSrcY uint16
}

func (z *CopyRectEncoding) Type() int32 {
	return 1
}
func (z *CopyRectEncoding) WriteTo(w io.Writer) (n int, err error) {
	binary.Write(w, binary.BigEndian, z.copyRectSrcX)
	if err != nil {
		return 0, err
	}
	binary.Write(w, binary.BigEndian, z.copyRectSrcY)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (z *CopyRectEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.IEncoding, error) {
	z.copyRectSrcX, _ = r.ReadUint16()
	z.copyRectSrcY, _ = r.ReadUint16()
	return z, nil
}

//////////
