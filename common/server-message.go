package common

import (
	"io"
)

type IClientConn interface {
	CurrentPixelFormat() *PixelFormat
	CurrentColorMap() *ColorMap
	Encodings() []Encoding
}

type ServerMessage interface {
	// The type of the message that is sent down on the wire.
	Type() uint8
	String() string
	CopyTo(r io.Reader, w io.Writer, c IClientConn) error
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

func (typ ServerMessageType) String() string {
	switch typ {
	case FramebufferUpdate:
		return "FramebufferUpdate"
	case SetColourMapEntries:
		return "SetColourMapEntries"
	case Bell:
		return "Bell"
	case ServerCutText:
		return "ServerCutText"
	}
	return ""
}

type ServerInit struct {
	FBWidth, FBHeight uint16
	PixelFormat       PixelFormat
	NameLength        uint32
	NameText          []byte
}
