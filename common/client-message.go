package common

import "io"

type ClientMessageType uint8

//go:generate stringer -type=ClientMessageType

// Client-to-Server message types.
const (
	SetPixelFormatMsgType ClientMessageType = iota
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

type Conn interface {
	io.ReadWriter
	Conn() io.ReadWriter
	Protocol() string
	PixelFormat() *PixelFormat
	SetPixelFormat(*PixelFormat) error
	ColorMap() *ColorMap
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
}

// ClientMessage is the interface
type ClientMessage interface {
	Type() ClientMessageType
	Read(Conn) (ClientMessage, error)
	Write(Conn) error
}
