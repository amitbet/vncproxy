package encodings

import (
	"io"
	"vncproxy/common"
)

type CopyRectEncoding struct {
	//Colors       []Color
	copyRectSrcX uint16
	copyRectSrcY uint16
}

func (z *CopyRectEncoding) Type() int32 {
	return 1
}
func (z *CopyRectEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r io.Reader) (common.Encoding, error) {
	conn := common.RfbReadHelper{r}
	//conn := &DataSource{conn: conn.c, PixelFormat: conn.PixelFormat}
	//bytesPerPixel := c.PixelFormat.BPP / 8
	z.copyRectSrcX, _ = conn.ReadUint16()
	z.copyRectSrcY, _ = conn.ReadUint16()
	return z, nil
}

//////////
