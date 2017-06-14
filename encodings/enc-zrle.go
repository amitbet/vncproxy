package encodings

import "io"
import "vncproxy/common"

type ZRLEEncoding struct {
	//Colors []Color
}

func (z *ZRLEEncoding) Type() int32 {
	return 16
}
func (z *ZRLEEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r io.Reader) (common.Encoding, error) {
	conn := common.RfbReadHelper{r}
	//conn := &DataSource{conn: conn.c, PixelFormat: conn.PixelFormat}
	//bytesPerPixel := c.PixelFormat.BPP / 8
	len, _ := conn.ReadUint32()
	_, err := conn.ReadBytes(int(len))

	if err != nil {
		return nil, err
	}
	return z, nil
}
