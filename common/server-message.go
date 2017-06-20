package common

type IClientConn interface {
	CurrentPixelFormat() *PixelFormat
	CurrentColorMap() *ColorMap
	Encodings() []Encoding
}

type ServerMessage interface {
	// The type of the message that is sent down on the wire.
	Type() uint8
	String() string
	// Read reads the contents of the message from the reader. At the point
	// this is called, the message type has already been read from the reader.
	// This should return a new ServerMessage that is the appropriate type.
	Read(IClientConn, *RfbReadHelper) (ServerMessage, error)
}
type ServerMessageType int8

const (
	FramebufferUpdate ServerMessageType = iota
	SetColourMapEntries
	Bell
	ServerCutText
)
