package vnc

import "io"

type ZLibEncoding struct {
	Colors []Color
}

func (z *ZLibEncoding) Type() int32 {
	return 6
}
func (z *ZLibEncoding) Read(conn *ClientConn, rect *Rectangle, r io.Reader) (Encoding, error) {
	//conn := &DataSource{conn: conn.c, PixelFormat: conn.PixelFormat}
	//bytesPerPixel := c.PixelFormat.BPP / 8
	len, _ := conn.readUint32()
	_, err := conn.readBytes(int(len))

	if err != nil {
		return nil, err
	}
	return z, nil
}
