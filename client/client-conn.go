package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"unicode"

	"github.com/amitbet/vncproxy/common"
	"github.com/amitbet/vncproxy/logger"
)

// A ServerMessage implements a message sent from the server to the client.

// A ClientAuth implements a method of authenticating with a remote server.
type ClientAuth interface {
	// SecurityType returns the byte identifier sent by the server to
	// identify this authentication scheme.
	SecurityType() uint8

	// Handshake is called when the authentication handshake should be
	// performed, as part of the general RFB handshake. (see 7.2.1)
	Handshake(io.ReadWriteCloser) error
}

type ClientConn struct {
	conn io.ReadWriteCloser

	//c      net.IServerConn
	config *ClientConfig

	// If the pixel format uses a color map, then this is the color
	// map that is used. This should not be modified directly, since
	// the data comes from the server.
	ColorMap common.ColorMap

	// Encodings supported by the client. This should not be modified
	// directly. Instead, SetEncodings should be used.
	Encs []common.IEncoding

	// Width of the frame buffer in pixels, sent from the server.
	FrameBufferWidth uint16

	// Height of the frame buffer in pixels, sent from the server.
	FrameBufferHeight uint16

	// Name associated with the desktop, sent from the server.
	DesktopName string

	// The pixel format associated with the connection. This shouldn't
	// be modified. If you wish to set a new pixel format, use the
	// SetPixelFormat method.
	PixelFormat common.PixelFormat

	Listeners *common.MultiListener
}

// A ClientConfig structure is used to configure a ClientConn. After
// one has been passed to initialize a connection, it must not be modified.
type ClientConfig struct {
	// A slice of ClientAuth methods. Only the first instance that is
	// suitable by the server will be used to authenticate.
	Auth []ClientAuth

	// Exclusive determines whether the connection is shared with other
	// clients. If true, then all other clients connected will be
	// disconnected when a connection is established to the VNC server.
	Exclusive bool

	// A slice of supported messages that can be read from the server.
	// This only needs to contain NEW server messages, and doesn't
	// need to explicitly contain the RFC-required messages.
	ServerMessages []common.ServerMessage
}

func NewClientConn(c net.Conn, cfg *ClientConfig) (*ClientConn, error) {
	conn := &ClientConn{
		conn:      c,
		config:    cfg,
		Listeners: &common.MultiListener{},
	}
	return conn, nil
}

func (conn *ClientConn) Connect() error {

	if err := conn.handshake(); err != nil {
		logger.Errorf("ClientConn.Connect error: %v", err)
		conn.Close()
		return err
	}

	go conn.mainLoop()

	return nil
}

func (c *ClientConn) Close() error {
	return c.conn.Close()
}

func (c *ClientConn) Encodings() []common.IEncoding {
	return c.Encs
}

func (c *ClientConn) CurrentPixelFormat() *common.PixelFormat {
	return &c.PixelFormat
}

func (c *ClientConn) Write(bytes []byte) (n int, err error) {
	return c.conn.Write(bytes)
}

func (c *ClientConn) Read(bytes []byte) (n int, err error) {
	return c.conn.Read(bytes)
}

// func (c *ClientConn) CurrentColorMap() *common.ColorMap {
// 	return &c.ColorMap
// }

// CutText tells the server that the client has new text in its cut buffer.
// The text string MUST only contain Latin-1 characters. This encoding
// is compatible with Go's native string format, but can only use up to
// unicode.MaxLatin values.
//
// See RFC 6143 Section 7.5.6
func (c *ClientConn) CutText(text string) error {
	var buf bytes.Buffer

	// This is the fixed size data we'll send
	fixedData := []interface{}{
		uint8(6),
		uint8(0),
		uint8(0),
		uint8(0),
		uint32(len(text)),
	}

	for _, val := range fixedData {
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return err
		}
	}

	for _, char := range text {
		if char > unicode.MaxLatin1 {
			return fmt.Errorf("Character '%v' is not valid Latin-1", char)
		}

		if err := binary.Write(&buf, binary.BigEndian, uint8(char)); err != nil {
			return err
		}
	}

	dataLength := 8 + len(text)
	if _, err := c.conn.Write(buf.Bytes()[0:dataLength]); err != nil {
		return err
	}

	return nil
}

// Requests a framebuffer update from the server. There may be an indefinite
// time between the request and the actual framebuffer update being
// received.
//
// See RFC 6143 Section 7.5.3
func (c *ClientConn) FramebufferUpdateRequest(incremental bool, x, y, width, height uint16) error {
	var buf bytes.Buffer
	var incrementalByte uint8 = 0

	if incremental {
		incrementalByte = 1
	}

	data := []interface{}{
		uint8(3),
		incrementalByte,
		x, y, width, height,
	}

	for _, val := range data {
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return err
		}
	}

	if _, err := c.conn.Write(buf.Bytes()[0:10]); err != nil {
		return err
	}

	return nil
}

// KeyEvent indiciates a key press or release and sends it to the server.
// The key is indicated using the X Window System "keysym" value. Use
// Google to find a reference of these values. To simulate a key press,
// you must send a key with both a down event, and a non-down event.
//
// See 7.5.4.
func (c *ClientConn) KeyEvent(keysym uint32, down bool) error {
	var downFlag uint8 = 0
	if down {
		downFlag = 1
	}

	data := []interface{}{
		uint8(4),
		downFlag,
		uint8(0),
		uint8(0),
		keysym,
	}

	for _, val := range data {
		if err := binary.Write(c.conn, binary.BigEndian, val); err != nil {
			return err
		}
	}

	return nil
}

// PointerEvent indicates that pointer movement or a pointer button
// press or release.
//
// The mask is a bitwise mask of various ButtonMask values. When a button
// is set, it is pressed, when it is unset, it is released.
//
// See RFC 6143 Section 7.5.5
func (c *ClientConn) PointerEvent(mask ButtonMask, x, y uint16) error {
	var buf bytes.Buffer

	data := []interface{}{
		uint8(5),
		uint8(mask),
		x,
		y,
	}

	for _, val := range data {
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return err
		}
	}

	if _, err := c.conn.Write(buf.Bytes()[0:6]); err != nil {
		return err
	}

	return nil
}

// SetEncodings sets the encoding types in which the pixel data can
// be sent from the server. After calling this method, the encs slice
// given should not be modified.
//
// See RFC 6143 Section 7.5.2
func (c *ClientConn) SetEncodings(encs []common.IEncoding) error {
	data := make([]interface{}, 3+len(encs))
	data[0] = uint8(2)
	data[1] = uint8(0)
	data[2] = uint16(len(encs))

	for i, enc := range encs {
		data[3+i] = int32(enc.Type())
	}

	var buf bytes.Buffer
	for _, val := range data {
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return err
		}
	}

	dataLength := 4 + (4 * len(encs))
	if _, err := c.conn.Write(buf.Bytes()[0:dataLength]); err != nil {
		return err
	}

	c.Encs = encs

	return nil
}

// SetPixelFormat sets the format in which pixel values should be sent
// in FramebufferUpdate messages from the server.
//
// See RFC 6143 Section 7.5.1
func (c *ClientConn) SetPixelFormat(format *common.PixelFormat) error {
	var keyEvent [20]byte
	keyEvent[0] = 0

	pfBytes, err := writePixelFormat(format)
	if err != nil {
		return err
	}

	// Copy the pixel format bytes into the proper slice location
	copy(keyEvent[4:], pfBytes)

	// Send the data down the connection
	if _, err := c.conn.Write(keyEvent[:]); err != nil {
		return err
	}

	// Reset the color map as according to RFC.
	var newColorMap common.ColorMap
	c.ColorMap = newColorMap

	return nil
}

const pvLen = 12 // ProtocolVersion message length.

func parseProtocolVersion(pv []byte) (uint, uint, error) {
	var major, minor uint

	if len(pv) < pvLen {
		return 0, 0, fmt.Errorf("ProtocolVersion message too short (%v < %v)", len(pv), pvLen)
	}

	l, err := fmt.Sscanf(string(pv), "RFB %d.%d\n", &major, &minor)
	if l != 2 {
		return 0, 0, fmt.Errorf("error parsing ProtocolVersion.")
	}
	if err != nil {
		return 0, 0, err
	}

	return major, minor, nil
}

func (c *ClientConn) handshake() error {
	var protocolVersion [pvLen]byte

	// 7.1.1, read the ProtocolVersion message sent by the server.
	if _, err := io.ReadFull(c.conn, protocolVersion[:]); err != nil {
		return err
	}

	maxMajor, maxMinor, err := parseProtocolVersion(protocolVersion[:])
	if err != nil {
		return err
	}
	if maxMajor < 3 {
		return fmt.Errorf("unsupported major version, less than 3: %d", maxMajor)
	}
	if maxMinor < 8 {
		return fmt.Errorf("unsupported minor version, less than 8: %d", maxMinor)
	}

	// Respond with the version we will support
	if _, err = c.conn.Write([]byte("RFB 003.008\n")); err != nil {
		return err
	}

	// 7.1.2 Security Handshake from server
	var numSecurityTypes uint8
	if err = binary.Read(c.conn, binary.BigEndian, &numSecurityTypes); err != nil {
		return fmt.Errorf("Error reading security types: %v", err)
		return err
	}

	if numSecurityTypes == 0 {
		return fmt.Errorf("Error: no security types: %s", c.readErrorReason())
	}

	securityTypes := make([]uint8, numSecurityTypes)
	if err = binary.Read(c.conn, binary.BigEndian, &securityTypes); err != nil {
		return err
	}

	clientSecurityTypes := c.config.Auth
	if clientSecurityTypes == nil {
		clientSecurityTypes = []ClientAuth{new(ClientAuthNone)}
	}

	var auth ClientAuth
FindAuth:
	for _, curAuth := range clientSecurityTypes {
		for _, securityType := range securityTypes {
			if curAuth.SecurityType() == securityType {
				// We use the first matching supported authentication
				auth = curAuth
				break FindAuth
			}
		}
	}

	if auth == nil {
		return fmt.Errorf("no suitable auth schemes found. server supported: %#v", securityTypes)
	}

	// Respond back with the security type we'll use
	if err = binary.Write(c.conn, binary.BigEndian, auth.SecurityType()); err != nil {
		return err
	}

	if err = auth.Handshake(c.conn); err != nil {
		return err
	}

	// 7.1.3 SecurityResult Handshake
	var securityResult uint32
	if err = binary.Read(c.conn, binary.BigEndian, &securityResult); err != nil {
		return err
	}

	if securityResult == 1 {
		return fmt.Errorf("security handshake failed: %s", c.readErrorReason())
	}

	// 7.3.1 ClientInit
	var sharedFlag uint8 = 1
	if c.config.Exclusive {
		sharedFlag = 0
	}

	if err = binary.Write(c.conn, binary.BigEndian, sharedFlag); err != nil {
		return err
	}

	// 7.3.2 ServerInit
	if err = binary.Read(c.conn, binary.BigEndian, &c.FrameBufferWidth); err != nil {
		return err
	}

	if err = binary.Read(c.conn, binary.BigEndian, &c.FrameBufferHeight); err != nil {
		return err
	}

	// Read the pixel format
	if err = readPixelFormat(c.conn, &c.PixelFormat); err != nil {
		return err
	}

	var nameLength uint32
	if err = binary.Read(c.conn, binary.BigEndian, &nameLength); err != nil {
		return err
	}

	nameBytes := make([]uint8, nameLength)
	if err = binary.Read(c.conn, binary.BigEndian, &nameBytes); err != nil {
		return err
	}

	c.DesktopName = string(nameBytes)
	srvInit := common.ServerInit{
		NameLength:  nameLength,
		NameText:    nameBytes,
		FBHeight:    c.FrameBufferHeight,
		FBWidth:     c.FrameBufferWidth,
		PixelFormat: c.PixelFormat,
	}
	rfbSeg := &common.RfbSegment{SegmentType: common.SegmentServerInitMessage, Message: &srvInit}

	return c.Listeners.Consume(rfbSeg)
}

// mainLoop reads messages sent from the server and routes them to the
// proper channels for users of the client to read.
func (c *ClientConn) mainLoop() {
	defer c.Close()

	reader := &common.RfbReadHelper{Reader: c.conn, Listeners: c.Listeners}
	// Build the map of available server messages
	typeMap := make(map[uint8]common.ServerMessage)

	defaultMessages := []common.ServerMessage{
		new(MsgFramebufferUpdate),
		new(MsgSetColorMapEntries),
		new(MsgBell),
		new(MsgServerCutText),
		new(MsgServerFence),
	}

	for _, msg := range defaultMessages {
		typeMap[msg.Type()] = msg
	}

	if c.config.ServerMessages != nil {
		for _, msg := range c.config.ServerMessages {
			typeMap[msg.Type()] = msg
		}
	}

	defer func() {
		logger.Warn("ClientConn.MainLoop: exiting!")
		c.Listeners.Consume(&common.RfbSegment{
			SegmentType: common.SegmentConnectionClosed,
		})
	}()

	for {
		var messageType uint8
		if err := binary.Read(c.conn, binary.BigEndian, &messageType); err != nil {
			logger.Errorf("ClientConn.MainLoop: error reading messagetype, %s", err)
			break
		}

		msg, ok := typeMap[messageType]
		if !ok {
			logger.Errorf("ClientConn.MainLoop: bad message type, %d", messageType)
			// Unsupported message type! Bad!
			break
		}
		logger.Debugf("ClientConn.MainLoop: got ServerMessage:%s", common.ServerMessageType(messageType))
		reader.SendMessageStart(common.ServerMessageType(messageType))
		reader.PublishBytes([]byte{byte(messageType)})

		parsedMsg, err := msg.Read(c, reader)
		if err != nil {
			logger.Errorf("ClientConn.MainLoop: error parsing message, %s", err)
			break
		}
		logger.Debugf("ClientConn.MainLoop: read & parsed ServerMessage:%d, %s", parsedMsg.Type(), parsedMsg)
	}
}

func (c *ClientConn) readErrorReason() string {
	var reasonLen uint32
	if err := binary.Read(c.conn, binary.BigEndian, &reasonLen); err != nil {
		return "<error>"
	}

	reason := make([]uint8, reasonLen)
	if err := binary.Read(c.conn, binary.BigEndian, &reason); err != nil {
		return "<error>"
	}

	return string(reason)
}
