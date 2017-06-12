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
