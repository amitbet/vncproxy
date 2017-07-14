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

// FramebufferUpdateMessage consists of a sequence of rectangles of
// pixel data that the client should put into its framebuffer.
type FramebufferUpdateMessage struct {
	Rectangles []common.Rectangle
}

func (m *FramebufferUpdateMessage) String() string {
	str := fmt.Sprintf("FramebufferUpdateMessage (type=%d) Rects: \n", m.Type())
	for _, rect := range m.Rectangles {
		str += rect.String() + "\n"
	}
	return str
}

func (*FramebufferUpdateMessage) Type() uint8 {
	return 0
}

func (fbm *FramebufferUpdateMessage) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	reader := common.NewRfbReadHelper(r)
	writeTo := &listeners.WriteTo{w, "FramebufferUpdateMessage.CopyTo"}
	reader.Listeners.AddListener(writeTo)
	_, err := fbm.Read(c, reader)
	return err
}

func (fbm *FramebufferUpdateMessage) Read(c common.IClientConn, r *common.RfbReadHelper) (common.ServerMessage, error) {

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
	encMap := make(map[int32]common.Encoding)
	for _, enc := range c.Encodings() {
		encMap[enc.Type()] = enc
	}

	// We must always support the raw encoding
	rawEnc := new(encodings.RawEncoding)
	encMap[rawEnc.Type()] = rawEnc
	logger.Debugf("numrects= %d", numRects)

	rects := make([]common.Rectangle, numRects)
	for i := uint16(0); i < numRects; i++ {
		logger.Debugf("###############rect################: %d\n", i)

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

		logger.Debugf("rect hdr data: enctype=%s, data: %s\n", encType, string(jBytes))
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
			} else {
				logger.Errorf("unsupported encoding type: %d, %s", encodingTypeInt, encType)
				return nil, fmt.Errorf("unsupported encoding type: %d, %s", encodingTypeInt, encType)
			}
		}
	}

	return &FramebufferUpdateMessage{rects}, nil
}

// SetColorMapEntriesMessage is sent by the server to set values into
// the color map. This message will automatically update the color map
// for the associated connection, but contains the color change data
// if the consumer wants to read it.
//
// See RFC 6143 Section 7.6.2
type SetColorMapEntriesMessage struct {
	FirstColor uint16
	Colors     []common.Color
}

func (fbm *SetColorMapEntriesMessage) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	reader := &common.RfbReadHelper{Reader: r}
	writeTo := &listeners.WriteTo{w, "SetColorMapEntriesMessage.CopyTo"}
	reader.Listeners.AddListener(writeTo)
	_, err := fbm.Read(c, reader)
	return err
}
func (m *SetColorMapEntriesMessage) String() string {
	return fmt.Sprintf("SetColorMapEntriesMessage (type=%d) first:%d colors: %v: ", m.Type(), m.FirstColor, m.Colors)
}

func (*SetColorMapEntriesMessage) Type() uint8 {
	return 1
}

func (*SetColorMapEntriesMessage) Read(c common.IClientConn, r *common.RfbReadHelper) (common.ServerMessage, error) {
	// Read off the padding
	var padding [1]byte
	if _, err := io.ReadFull(r, padding[:]); err != nil {
		return nil, err
	}

	var result SetColorMapEntriesMessage
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
type BellMessage byte

func (fbm *BellMessage) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	return nil
}
func (m *BellMessage) String() string {
	return fmt.Sprintf("BellMessage (type=%d)", m.Type())
}

func (*BellMessage) Type() uint8 {
	return 2
}

func (*BellMessage) Read(common.IClientConn, *common.RfbReadHelper) (common.ServerMessage, error) {
	return new(BellMessage), nil
}

// ServerCutTextMessage indicates the server has new text in the cut buffer.
//
// See RFC 6143 Section 7.6.4
type ServerCutTextMessage struct {
	Text string
}

func (fbm *ServerCutTextMessage) CopyTo(r io.Reader, w io.Writer, c common.IClientConn) error {
	reader := &common.RfbReadHelper{Reader: r}
	writeTo := &listeners.WriteTo{w, "ServerCutTextMessage.CopyTo"}
	reader.Listeners.AddListener(writeTo)
	_, err := fbm.Read(c, reader)
	return err
}
func (m *ServerCutTextMessage) String() string {
	return fmt.Sprintf("ServerCutTextMessage (type=%d)", m.Type())
}

func (*ServerCutTextMessage) Type() uint8 {
	return 3
}

func (*ServerCutTextMessage) Read(conn common.IClientConn, r *common.RfbReadHelper) (common.ServerMessage, error) {
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
	//var textLength uint32
	// if err := binary.Read(r, binary.BigEndian, &textLength); err != nil {
	// 	return nil, err
	// }

	// textBytes := make([]uint8, textLength)
	// if err := binary.Read(r, binary.BigEndian, &textBytes); err != nil {
	// 	return nil, err
	// }

	return &ServerCutTextMessage{string(textBytes)}, nil
}
