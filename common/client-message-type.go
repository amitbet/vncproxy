package common

import (
	"io"
)

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
	ClientFenceMsgType          = 248
	QEMUExtendedKeyEventMsgType = 255
)

// Color represents a single color in a color map.
type Color struct {
	pf      *PixelFormat
	cm      *ColorMap
	cmIndex uint32 // Only valid if pf.TrueColor is false.
	R, G, B uint16
}

type ColorMap [256]Color

// ClientMessage is the interface
type ClientMessage interface {
	Type() ClientMessageType
	Read(io.Reader) (ClientMessage, error)
	Write(io.Writer) error
}

func (cmt ClientMessageType) String() string {
	switch cmt {
	case SetPixelFormatMsgType:
		return "SetPixelFormat"
	case SetEncodingsMsgType:
		return "SetEncodings"
	case FramebufferUpdateRequestMsgType:
		return "FramebufferUpdateRequest"
	case KeyEventMsgType:
		return "KeyEvent"
	case QEMUExtendedKeyEventMsgType:
		return "QEMUExtendedKeyEvent"
	case PointerEventMsgType:
		return "PointerEvent"
	case ClientCutTextMsgType:
		return "ClientCutText"
	}
	return ""
}
