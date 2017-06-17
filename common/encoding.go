package common

// An Encoding implements a method for encoding pixel data that is
// sent by the server to the client.
type Encoding interface {
	// The number that uniquely identifies this encoding type.
	Type() int32

	// Read reads the contents of the encoded pixel data from the reader.
	// This should return a new Encoding implementation that contains
	// the proper data.
	Read(*PixelFormat, *Rectangle, *RfbReadHelper) (Encoding, error)
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
