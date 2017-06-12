package vnc

import "io"

type ZRLEEncoding struct {
	Colors []Color
}

func (z *ZRLEEncoding) Type() int32 {
	return 16
}
func (z *ZRLEEncoding) Read(conn *ClientConn, rect *Rectangle, r io.Reader) (Encoding, error) {
	//conn := &DataSource{conn: conn.c, PixelFormat: conn.PixelFormat}
	//bytesPerPixel := c.PixelFormat.BPP / 8
	len, _ := conn.readUint32()
	_, err := conn.readBytes(int(len))

	if err != nil {
		return nil, err
	}
	return z, nil
}
