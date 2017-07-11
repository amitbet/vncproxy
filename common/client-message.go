package common

import (
	"io"

)

type ClientMessageType uint8

//go:generate stringer -type=ClientMessageType

// Client-to-Server message types.
const (
	SetPixelFormatMsgType           ClientMessageType = iota
	_
	SetEncodingsMsgType
	FramebufferUpdateRequestMsgType
	KeyEventMsgType
	PointerEventMsgType
	ClientCutTextMsgType
)

// Color represents a single color in a color map.
type Color struct {
	pf      *PixelFormat
	cm      *ColorMap
	cmIndex uint32 // Only valid if pf.TrueColor is false.
	R, G, B uint16
}

type ColorMap [256]Color

type ServerConn interface {
	io.ReadWriter
	//ServerConn() io.ReadWriter
	Protocol() string
	CurrentPixelFormat() *PixelFormat
	SetPixelFormat(*PixelFormat) error
	//ColorMap() *ColorMap
	SetColorMap(*ColorMap)
	Encodings() []Encoding
	SetEncodings([]EncodingType) error
	Width() uint16
	Height() uint16
	SetWidth(uint16)
	SetHeight(uint16)
	DesktopName() string
	SetDesktopName(string)
	//Flush() error
	SetProtoVersion(string)
	// Write([]byte) (int, error)
}

// ClientMessage is the interface
type ClientMessage interface {
	Type() ClientMessageType
	Read(io.Reader) (ClientMessage, error)
	Write(io.Writer) error
}

func (cmt ClientMessageType) String() string {
	switch  cmt {
	case SetPixelFormatMsgType:
		return "SetPixelFormat"
	case SetEncodingsMsgType:
		return "SetEncodings"
	case FramebufferUpdateRequestMsgType:
		return "FramebufferUpdateRequest"
	case KeyEventMsgType:
		return "KeyEvent"
	case PointerEventMsgType:
		return "PointerEvent"
	case ClientCutTextMsgType:
		return "ClientCutText"
	}
	return ""
}
