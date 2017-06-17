package encodings

import (
	"vncproxy/common"
)

type CoRREEncoding struct {
	//Colors []Color
}

func (z *CoRREEncoding) Type() int32 {
	return 4
}

func (z *CoRREEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.Encoding, error) {
	bytesPerPixel := int(pixelFmt.BPP / 8)
	numOfSubrectangles, _ := r.ReadUint32()

	//read whole rect background color
	r.ReadBytes(bytesPerPixel)

	//read all individual rects (color=BPP + x=16b + y=16b + w=16b + h=16b)
	_, err := r.ReadBytes(int(numOfSubrectangles) * (bytesPerPixel + 4))

	if err != nil {
		return nil, err
	}
	return z, nil

	//int nSubrects = rfb.readU32();

	//byte[] bg_buf = new byte[bytesPerPixel];
	//rfb.readFully(bytesPerPixel);
	//Color pixel;
	// if (bytesPixel == 1) {
	//   pixel = colors[bg_buf[0] & 0xFF];
	// } else {
	//   pixel = new Color(bg_buf[2] & 0xFF, bg_buf[1] & 0xFF, bg_buf[0] & 0xFF);
	// }
	// memGraphics.setColor(pixel);
	// memGraphics.fillRect(x, y, w, h);

	// byte[] buf = new byte[nSubrects * (bytesPixel + 4)];
	// rfb.readFully(buf);

}
