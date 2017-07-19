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
	Enc    IEncoding
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("(%d,%d) (width: %d, height: %d), Enc= %d", r.X, r.Y, r.Width, r.Height, r.Enc.Type())
}
