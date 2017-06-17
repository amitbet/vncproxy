package encodings

import "vncproxy/common"

type RREEncoding struct {
	//Colors []Color
}

func (z *RREEncoding) Type() int32 {
	return 2
}
func (z *RREEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	//conn := common.RfbReadHelper{Reader:r}
	bytesPerPixel := int(pixelFmt.BPP / 8)
	numOfSubrectangles, _ := r.ReadUint32()

	//read whole rect background color
	r.ReadBytes(bytesPerPixel)

	//read all individual rects (color=BPP + x=16b + y=16b + w=16b + h=16b)
	_, err := r.ReadBytes(int(numOfSubrectangles) * (bytesPerPixel + 8))

	if err != nil {
		return nil, err
	}
	return z, nil
}
