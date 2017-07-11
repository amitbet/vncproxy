package encodings

import "vncproxy/common"

type CopyRectEncoding struct {
	//Colors       []Color
	copyRectSrcX uint16
	copyRectSrcY uint16
}

func (z *CopyRectEncoding) Type() int32 {
	return 1
}
func (z *CopyRectEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	z.copyRectSrcX, _ = r.ReadUint16()
	z.copyRectSrcY, _ = r.ReadUint16()
	return z, nil
}

//////////
