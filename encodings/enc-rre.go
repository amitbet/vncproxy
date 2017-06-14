package encodings

import "io"
import "vncproxy/common"

type RREEncoding struct {
	//Colors []Color
}

func (z *RREEncoding) Type() int32 {
	return 2
}
func (z *RREEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r io.Reader) (common.Encoding, error) {
	conn := common.RfbReadHelper{r}
	bytesPerPixel := int(pixelFmt.BPP / 8)
	numOfSubrectangles, _ := conn.ReadUint32()

	//read whole rect background color
	conn.ReadBytes(bytesPerPixel)

	//read all individual rects (color=BPP + x=16b + y=16b + w=16b + h=16b)
	_, err := conn.ReadBytes(int(numOfSubrectangles) * (bytesPerPixel + 8))

	if err != nil {
		return nil, err
	}
	return z, nil
}
