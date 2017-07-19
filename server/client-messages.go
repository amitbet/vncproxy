package server

import (
	"encoding/binary"
	"io"
	"vncproxy/common"
)

// Key represents a VNC key press.
type Key uint32

//go:generate stringer -type=Key

// Keys is a slice of Key values.
type Keys []Key

// MsgSetPixelFormat holds the wire format message.
type MsgSetPixelFormat struct {
	_  [3]byte            // padding
	PF common.PixelFormat // pixel-format
	_  [3]byte            // padding after pixel format
}

func (*MsgSetPixelFormat) Type() common.ClientMessageType {
	return common.SetPixelFormatMsgType
}

func (msg *MsgSetPixelFormat) Write(c io.Writer) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}

	if err := binary.Write(c, binary.BigEndian, msg); err != nil {
		return err
	}

	//pf := c.CurrentPixelFormat()
	// Invalidate the color map.
	// if pf.TrueColor {
	// 	c.SetColorMap(&common.ColorMap{})
	// }

	return nil
}

func (*MsgSetPixelFormat) Read(c io.Reader) (common.ClientMessage, error) {
	msg := MsgSetPixelFormat{}
	if err := binary.Read(c, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// MsgSetEncodings holds the wire format message, sans encoding-type field.
type MsgSetEncodings struct {
	_         [1]byte // padding
	EncNum    uint16  // number-of-encodings
	Encodings []common.EncodingType
}

func (*MsgSetEncodings) Type() common.ClientMessageType {
	return common.SetEncodingsMsgType
}

func (*MsgSetEncodings) Read(c io.Reader) (common.ClientMessage, error) {
	msg := MsgSetEncodings{}
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
	c.(common.IServerConn).SetEncodings(msg.Encodings)
	return &msg, nil
}

func (msg *MsgSetEncodings) Write(c io.Writer) error {
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
	return nil
}

// MsgFramebufferUpdateRequest holds the wire format message.
type MsgFramebufferUpdateRequest struct {
	Inc           uint8  // incremental
	X, Y          uint16 // x-, y-position
	Width, Height uint16 // width, height
}

func (*MsgFramebufferUpdateRequest) Type() common.ClientMessageType {
	return common.FramebufferUpdateRequestMsgType
}

func (*MsgFramebufferUpdateRequest) Read(c io.Reader) (common.ClientMessage, error) {
	msg := MsgFramebufferUpdateRequest{}
	if err := binary.Read(c, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (msg *MsgFramebufferUpdateRequest) Write(c io.Writer) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}
	if err := binary.Write(c, binary.BigEndian, msg); err != nil {
		return err
	}
	return nil
}

// MsgKeyEvent holds the wire format message.
type MsgKeyEvent struct {
	Down uint8   // down-flag
	_    [2]byte // padding
	Key  Key     // key
}

func (*MsgKeyEvent) Type() common.ClientMessageType {
	return common.KeyEventMsgType
}

func (*MsgKeyEvent) Read(c io.Reader) (common.ClientMessage, error) {
	msg := MsgKeyEvent{}
	if err := binary.Read(c, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (msg *MsgKeyEvent) Write(c io.Writer) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}
	if err := binary.Write(c, binary.BigEndian, msg); err != nil {
		return err
	}
	return nil
}

// PointerEventMessage holds the wire format message.
type MsgPointerEvent struct {
	Mask uint8  // button-mask
	X, Y uint16 // x-, y-position
}

func (*MsgPointerEvent) Type() common.ClientMessageType {
	return common.PointerEventMsgType
}

func (*MsgPointerEvent) Read(c io.Reader) (common.ClientMessage, error) {
	msg := MsgPointerEvent{}
	if err := binary.Read(c, binary.BigEndian, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (msg *MsgPointerEvent) Write(c io.Writer) error {
	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
		return err
	}
	if err := binary.Write(c, binary.BigEndian, msg); err != nil {
		return err
	}
	return nil
}

type MsgClientFence struct {
}

func (*MsgClientFence) Type() common.ClientMessageType {
	return common.ClientFenceMsgType
}

func (cf *MsgClientFence) Read(c io.Reader) (common.ClientMessage, error) {
	bytes := make([]byte, 3)
	c.Read(bytes)
	if _, err := c.Read(bytes); err != nil {
		return nil, err
	}
	var flags uint32
	if err := binary.Read(c, binary.BigEndian, &flags); err != nil {
		return nil, err
	}

	var length uint8
	if err := binary.Read(c, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	bytes = make([]byte, length)
	if _, err := c.Read(bytes); err != nil {
		return nil, err
	}
	return cf, nil
}

func (msg *MsgClientFence) Write(c io.Writer) error {
	panic("not implemented!")
}

// MsgClientCutText holds the wire format message, sans the text field.
type MsgClientCutText struct {
	_      [3]byte // padding
	Length uint32  // length
	Text   []byte
}

func (*MsgClientCutText) Type() common.ClientMessageType {
	return common.ClientCutTextMsgType
}

func (*MsgClientCutText) Read(c io.Reader) (common.ClientMessage, error) {
	msg := MsgClientCutText{}
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

func (msg *MsgClientCutText) Write(c io.Writer) error {
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

	return nil
}
