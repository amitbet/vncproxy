package vnc

import "io"

type RREEncoding struct {
	Colors []Color
}

func (z *RREEncoding) Type() int32 {
	return 2
}
func (z *RREEncoding) Read(conn *ClientConn, rect *Rectangle, r io.Reader) (Encoding, error) {

	bytesPerPixel := int(conn.PixelFormat.BPP / 8)
	numOfSubrectangles, _ := conn.readUint32()

	//read whole rect background color
	conn.readBytes(bytesPerPixel)

	//read all individual rects (color=BPP + x=16b + y=16b + w=16b + h=16b)
	_, err := conn.readBytes(int(numOfSubrectangles) * (bytesPerPixel + 8))

	if err != nil {
		return nil, err
	}
	return z, nil
}

type CoRREEncoding struct {
	Colors []Color
}

func (z *CoRREEncoding) Type() int32 {
	return 4
}

func (z *CoRREEncoding) Read(conn *ClientConn, rect *Rectangle, r io.Reader) (Encoding, error) {

	bytesPerPixel := int(conn.PixelFormat.BPP / 8)
	numOfSubrectangles, _ := conn.readUint32()

	//read whole rect background color
	conn.readBytes(bytesPerPixel)

	//read all individual rects (color=BPP + x=16b + y=16b + w=16b + h=16b)
	_, err := conn.readBytes(int(numOfSubrectangles) * (bytesPerPixel + 4))

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
