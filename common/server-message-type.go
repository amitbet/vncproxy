package common

import (
	"io"
)

//ServerMessage ...
type ServerMessage interface {
	Type() uint8
	String() string
	CopyTo(r io.Reader, w io.Writer, c IClientConn) error
	Read(IClientConn, *RfbReadHelper) (ServerMessage, error)
}

//ServerMessageType ...
type ServerMessageType int8

//ServerMessageType ...
const (
	FramebufferUpdate ServerMessageType = iota
	SetColourMapEntries
	Bell
	ServerCutText
	ServerFence = 248
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

//ServerInit ...
type ServerInit struct {
	FBWidth, FBHeight uint16
	PixelFormat       PixelFormat
	NameLength        uint32
	NameText          []byte
}
