package common

import (
	"fmt"
)

// Rectangle represents a rectangle of pixel data.
type Rectangle struct {
	X      uint16
	Y      uint16
	Width  uint16
	Height uint16
	Enc    Encoding
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("(%d,%d) (width: %d, height: %d), Enc= %d", r.X, r.Y, r.Width, r.Height, r.Enc.Type())
}

// PixelFormat describes the way a pixel is formatted for a VNC connection.
//
// See RFC 6143 Section 7.4 for information on each of the fields.
type PixelFormat struct {
	BPP        uint8
	Depth      uint8
	BigEndian  bool
	TrueColor  bool
	RedMax     uint16
	GreenMax   uint16
	BlueMax    uint16
	RedShift   uint8
	GreenShift uint8
	BlueShift  uint8
}
