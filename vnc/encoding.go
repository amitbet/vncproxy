package vnc

import "io"

// An Encoding implements a method for encoding pixel data that is
// sent by the server to the client.
type Encoding interface {
	// The number that uniquely identifies this encoding type.
	Type() int32

	// Read reads the contents of the encoded pixel data from the reader.
	// This should return a new Encoding implementation that contains
	// the proper data.
	Read(*ClientConn, *Rectangle, io.Reader) (Encoding, error)
}

const (
	EncodingRaw      = 0
	EncodingCopyRect = 1
	EncodingRRE      = 2
	EncodingCoRRE    = 4
	EncodingHextile  = 5
	EncodingZlib     = 6
	EncodingTight    = 7
	EncodingZRLE     = 16
)

type CopyRectEncoding struct {
	Colors       []Color
	copyRectSrcX uint16
	copyRectSrcY uint16
}

func (z *CopyRectEncoding) Type() int32 {
	return 1
}
func (z *CopyRectEncoding) Read(conn *ClientConn, rect *Rectangle, r io.Reader) (Encoding, error) {
	//conn := &DataSource{conn: conn.c, PixelFormat: conn.PixelFormat}
	//bytesPerPixel := c.PixelFormat.BPP / 8
	z.copyRectSrcX, _ = conn.readUint16()
	z.copyRectSrcY, _ = conn.readUint16()
	return z, nil
}

//////////
