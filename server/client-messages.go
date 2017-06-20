package server

import (
	"encoding/binary"
	"vncproxy/common"
)

// SetPixelFormat holds the wire format message.
type SetPixelFormat struct {
	_  [3]byte            // padding
	PF common.PixelFormat // pixel-format
	_  [3]byte            // padding after pixel format
}

// Key represents a VNC key press.
type Key uint32

//go:generate stringer -type=Key

// Keys is a slice of Key values.
type Keys []Key

func (*SetPixelFormat) Type() common.ClientMessageType {
	return common.SetPixelFormatMsgType
}

func (msg *SetPixelFormat) Write(c common.Conn) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}

	if err := binary.Write(c, binary.BigEndian, msg); err != nil {
		return err
	}

	pf := c.PixelFormat()
	// Invalidate the color map.
	if pf.TrueColor {
		c.SetColorMap(&common.ColorMap{})
	}

	return c.Flush()
}

func (*SetPixelFormat) Read(c common.Conn) (common.ClientMessage, error) {
	msg := SetPixelFormat{}
	if err := binary.Read(c, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// SetEncodings holds the wire format message, sans encoding-type field.
type SetEncodings struct {
	_         [1]byte // padding
	EncNum    uint16  // number-of-encodings
	Encodings []common.EncodingType
}

func (*SetEncodings) Type() common.ClientMessageType {
	return common.SetEncodingsMsgType
}

func (*SetEncodings) Read(c common.Conn) (common.ClientMessage, error) {
	msg := SetEncodings{}
	var pad [1]byte
	if err := binary.Read(c, binary.BigEndian, &pad); err != nil {
		return nil, err
	}

	if err := binary.Read(c, binary.BigEndian, &msg.EncNum); err != nil {
		return nil, err
	}
	var enc common.EncodingType
	for i := uint16(0); i < msg.EncNum; i++ {
		if err := binary.Read(c, binary.BigEndian, &enc); err != nil {
			return nil, err
		}
		msg.Encodings = append(msg.Encodings, enc)
	}
	c.SetEncodings(msg.Encodings)
	return &msg, nil
}

func (msg *SetEncodings) Write(c common.Conn) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}

	var pad [1]byte
	if err := binary.Write(c, binary.BigEndian, pad); err != nil {
		return err
	}

	if uint16(len(msg.Encodings)) > msg.EncNum {
		msg.EncNum = uint16(len(msg.Encodings))
	}
	if err := binary.Write(c, binary.BigEndian, msg.EncNum); err != nil {
		return err
	}
	for _, enc := range msg.Encodings {
		if err := binary.Write(c, binary.BigEndian, enc); err != nil {
			return err
		}
	}
	return c.Flush()
}

// FramebufferUpdateRequest holds the wire format message.
type FramebufferUpdateRequest struct {
	Inc           uint8  // incremental
	X, Y          uint16 // x-, y-position
	Width, Height uint16 // width, height
}

func (*FramebufferUpdateRequest) Type() common.ClientMessageType {
	return common.FramebufferUpdateRequestMsgType
}

func (*FramebufferUpdateRequest) Read(c common.Conn) (common.ClientMessage, error) {
	msg := FramebufferUpdateRequest{}
	if err := binary.Read(c, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (msg *FramebufferUpdateRequest) Write(c common.Conn) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}
	if err := binary.Write(c, binary.BigEndian, msg); err != nil {
		return err
	}
	return c.Flush()
}

// KeyEvent holds the wire format message.
type KeyEvent struct {
	Down uint8   // down-flag
	_    [2]byte // padding
	Key  Key     // key
}

func (*KeyEvent) Type() common.ClientMessageType {
	return common.KeyEventMsgType
}

func (*KeyEvent) Read(c common.Conn) (common.ClientMessage, error) {
	msg := KeyEvent{}
	if err := binary.Read(c, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (msg *KeyEvent) Write(c common.Conn) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}
	if err := binary.Write(c, binary.BigEndian, msg); err != nil {
		return err
	}
	return c.Flush()
}

// PointerEventMessage holds the wire format message.
type PointerEvent struct {
	Mask uint8  // button-mask
	X, Y uint16 // x-, y-position
}

func (*PointerEvent) Type() common.ClientMessageType {
	return common.PointerEventMsgType
}

func (*PointerEvent) Read(c common.Conn) (common.ClientMessage, error) {
	msg := PointerEvent{}
	if err := binary.Read(c, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (msg *PointerEvent) Write(c common.Conn) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}
	if err := binary.Write(c, binary.BigEndian, msg); err != nil {
		return err
	}
	return c.Flush()
}

// ClientCutText holds the wire format message, sans the text field.
type ClientCutText struct {
	_      [3]byte // padding
	Length uint32  // length
	Text   []byte
}

func (*ClientCutText) Type() common.ClientMessageType {
	return common.ClientCutTextMsgType
}

func (*ClientCutText) Read(c common.Conn) (common.ClientMessage, error) {
	msg := ClientCutText{}
	var pad [3]byte
	if err := binary.Read(c, binary.BigEndian, &pad); err != nil {
		return nil, err
	}

	if err := binary.Read(c, binary.BigEndian, &msg.Length); err != nil {
		return nil, err
	}

	msg.Text = make([]byte, msg.Length)
	if err := binary.Read(c, binary.BigEndian, &msg.Text); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (msg *ClientCutText) Write(c common.Conn) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}

	var pad [3]byte
	if err := binary.Write(c, binary.BigEndian, &pad); err != nil {
		return err
	}

	if uint32(len(msg.Text)) > msg.Length {
		msg.Length = uint32(len(msg.Text))
	}

	if err := binary.Write(c, binary.BigEndian, msg.Length); err != nil {
		return err
	}

	if err := binary.Write(c, binary.BigEndian, msg.Text); err != nil {
		return err
	}

	return c.Flush()
}
