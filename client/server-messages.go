package client

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
	listeners "vncproxy/tee-listeners"
)

// MsgFramebufferUpdate consists of a sequence of rectangles of
// pixel data that the client should put into its framebuffer.
type MsgFramebufferUpdate struct {
	Rectangles []common.Rectangle
}

func (m *MsgFramebufferUpdate) String() string {
	str := fmt.Sprintf("MsgFramebufferUpdate (type=%d) Rects: ", m.Type())
	for _, rect := range m.Rectangles {
		str += rect.String() + "\n"
		//if this is the last rect, break the loop
		if rect.Enc.Type() == int32(common.EncLastRectPseudo) {
			break
		}
	}
	return str
}

func (*MsgFramebufferUpdate) Type() uint8 {
	return 0
}

func (fbm *MsgFramebufferUpdate) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	reader := common.NewRfbReadHelper(r)
	writeTo := &listeners.WriteTo{w, "MsgFramebufferUpdate.CopyTo"}
	reader.Listeners.AddListener(writeTo)
	_, err := fbm.Read(c, reader)
	return err
}

func (fbm *MsgFramebufferUpdate) Read(c common.IClientConn, r *common.RfbReadHelper) (common.ServerMessage, error) {

	// Read off the padding
	var padding [1]byte
	if _, err := io.ReadFull(r, padding[:]); err != nil {
		return nil, err
	}

	var numRects uint16
	if err := binary.Read(r, binary.BigEndian, &numRects); err != nil {
		return nil, err
	}

	// Build the map of encodings supported
	encMap := make(map[int32]common.IEncoding)
	for _, enc := range c.Encodings() {
		encMap[enc.Type()] = enc
	}

	// We must always support the raw encoding
	rawEnc := new(encodings.RawEncoding)
	encMap[rawEnc.Type()] = rawEnc
	logger.Infof("MsgFramebufferUpdate.Read: numrects= %d", numRects)

	rects := make([]common.Rectangle, numRects)
	for i := uint16(0); i < numRects; i++ {
		logger.Debugf("MsgFramebufferUpdate.Read: ###############rect################: %d", i)

		var encodingTypeInt int32
		r.SendRectSeparator(-1)
		rect := &rects[i]
		data := []interface{}{
			&rect.X,
			&rect.Y,
			&rect.Width,
			&rect.Height,
			&encodingTypeInt,
		}

		for _, val := range data {
			if err := binary.Read(r, binary.BigEndian, val); err != nil {
				logger.Errorf("err: %v", err)
				return nil, err
			}
		}
		jBytes, _ := json.Marshal(data)

		encType := common.EncodingType(encodingTypeInt)

		logger.Infof("MsgFramebufferUpdate.Read: rect# %d, rect hdr data: enctype=%s, data: %s", i, encType, string(jBytes))
		enc, supported := encMap[encodingTypeInt]
		if supported {
			var err error
			rect.Enc, err = enc.Read(c.CurrentPixelFormat(), rect, r)
			if err != nil {
				return nil, err
			}
		} else {
			if strings.Contains(encType.String(), "Pseudo") {
				rect.Enc = &encodings.PseudoEncoding{encodingTypeInt}

				//if this is the last rect, break the for loop
				if rect.Enc.Type() == int32(common.EncLastRectPseudo) {
					break
				}
			} else {
				logger.Errorf("MsgFramebufferUpdate.Read: unsupported encoding type: %d, %s", encodingTypeInt, encType)
				return nil, fmt.Errorf("MsgFramebufferUpdate.Read: unsupported encoding type: %d, %s", encodingTypeInt, encType)
			}
		}
	}

	return &MsgFramebufferUpdate{rects}, nil
}

// MsgSetColorMapEntries is sent by the server to set values into
// the color map. This message will automatically update the color map
// for the associated connection, but contains the color change data
// if the consumer wants to read it.
//
// See RFC 6143 Section 7.6.2
type MsgSetColorMapEntries struct {
	FirstColor uint16
	Colors     []common.Color
}

func (fbm *MsgSetColorMapEntries) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	reader := &common.RfbReadHelper{Reader: r}
	writeTo := &listeners.WriteTo{w, "MsgSetColorMapEntries.CopyTo"}
	reader.Listeners.AddListener(writeTo)
	_, err := fbm.Read(c, reader)
	return err
}
func (m *MsgSetColorMapEntries) String() string {
	return fmt.Sprintf("MsgSetColorMapEntries (type=%d) first:%d colors: %v: ", m.Type(), m.FirstColor, m.Colors)
}

func (*MsgSetColorMapEntries) Type() uint8 {
	return 1
}

func (*MsgSetColorMapEntries) Read(c common.IClientConn, r *common.RfbReadHelper) (common.ServerMessage, error) {
	// Read off the padding
	var padding [1]byte
	if _, err := io.ReadFull(r, padding[:]); err != nil {
		return nil, err
	}

	var result MsgSetColorMapEntries
	if err := binary.Read(r, binary.BigEndian, &result.FirstColor); err != nil {
		return nil, err
	}

	var numColors uint16
	if err := binary.Read(r, binary.BigEndian, &numColors); err != nil {
		return nil, err
	}

	result.Colors = make([]common.Color, numColors)
	for i := uint16(0); i < numColors; i++ {

		color := &result.Colors[i]
		data := []interface{}{
			&color.R,
			&color.G,
			&color.B,
		}

		for _, val := range data {
			if err := binary.Read(r, binary.BigEndian, val); err != nil {
				return nil, err
			}
		}
		cmap := c.CurrentColorMap()
		// Update the connection's color map
		cmap[result.FirstColor+i] = *color
	}

	return &result, nil
}

// Bell signals that an audible bell should be made on the client.
//
// See RFC 6143 Section 7.6.3
type MsgBell byte

func (fbm *MsgBell) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	return nil
}
func (m *MsgBell) String() string {
	return fmt.Sprintf("MsgBell (type=%d)", m.Type())
}

func (*MsgBell) Type() uint8 {
	return 2
}

func (*MsgBell) Read(common.IClientConn, *common.RfbReadHelper) (common.ServerMessage, error) {
	return new(MsgBell), nil
}

type MsgServerFence byte

func (fbm *MsgServerFence) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	return nil
}
func (m *MsgServerFence) String() string {
	return fmt.Sprintf("MsgServerFence (type=%d)", m.Type())
}

func (*MsgServerFence) Type() uint8 {
	return uint8(common.ServerFence)
}

func (sf *MsgServerFence) Read(info common.IClientConn, c *common.RfbReadHelper) (common.ServerMessage, error) {
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
	return sf, nil
}

// MsgServerCutText indicates the server has new text in the cut buffer.
//
// See RFC 6143 Section 7.6.4
type MsgServerCutText struct {
	Text string
}

func (fbm *MsgServerCutText) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	reader := &common.RfbReadHelper{Reader: r}
	writeTo := &listeners.WriteTo{w, "MsgServerCutText.CopyTo"}
	reader.Listeners.AddListener(writeTo)
	_, err := fbm.Read(c, reader)
	return err
}
func (m *MsgServerCutText) String() string {
	return fmt.Sprintf("MsgServerCutText (type=%d)", m.Type())
}

func (*MsgServerCutText) Type() uint8 {
	return 3
}

func (*MsgServerCutText) Read(conn common.IClientConn, r *common.RfbReadHelper) (common.ServerMessage, error) {
	//reader := common.RfbReadHelper{Reader: r}

	// Read off the padding
	var padding [3]byte
	if _, err := io.ReadFull(r, padding[:]); err != nil {
		return nil, err
	}
	textLength, err := r.ReadUint32()
	if err != nil {
		return nil, err
	}
	textBytes, err := r.ReadBytes(int(textLength))
	if err != nil {
		return nil, err
	}

	return &MsgServerCutText{string(textBytes)}, nil
}
