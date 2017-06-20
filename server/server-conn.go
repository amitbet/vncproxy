package server

import (
	"bufio"
	"net"
	"sync"
	"vncproxy/common"
)

type ServerConn struct {
	c        net.Conn
	cfg      *ServerConfig
	br       *bufio.Reader
	bw       *bufio.Writer
	protocol string
	m        sync.Mutex
	// If the pixel format uses a color map, then this is the color
	// map that is used. This should not be modified directly, since
	// the data comes from the server.
	// Definition in §5 - Representation of Pixel Data.
	colorMap *common.ColorMap

	// Name associated with the desktop, sent from the server.
	desktopName string

	// Encodings supported by the client. This should not be modified
	// directly. Instead, SetEncodings() should be used.
	encodings []common.Encoding

	// Height of the frame buffer in pixels, sent to the client.
	fbHeight uint16

	// Width of the frame buffer in pixels, sent to the client.
	fbWidth uint16

	// The pixel format associated with the connection. This shouldn't
	// be modified. If you wish to set a new pixel format, use the
	// SetPixelFormat method.
	pixelFormat *common.PixelFormat

	quit chan struct{}
}

func (c *ServerConn) UnreadByte() error {
	return c.br.UnreadByte()
}

func (c *ServerConn) Conn() net.Conn {
	return c.c
}

func (c *ServerConn) SetEncodings(encs []common.EncodingType) error {
	encodings := make(map[int32]common.Encoding)
	for _, enc := range c.cfg.Encodings {
		encodings[enc.Type()] = enc
	}
	for _, encType := range encs {
		if enc, ok := encodings[int32(encType)]; ok {
			c.encodings = append(c.encodings, enc)
		}
	}
	return nil
}

func (c *ServerConn) SetProtoVersion(pv string) {
	c.protocol = pv
}

func (c *ServerConn) Flush() error {
	//	c.m.Lock()
	//	defer c.m.Unlock()
	return c.bw.Flush()
}

func (c *ServerConn) Close() error {
	return c.c.Close()
}

/*
func (c *ServerConn) Input() chan *ServerMessage {
	return c.cfg.ServerMessageCh
}

func (c *ServerConn) Output() chan *ClientMessage {
	return c.cfg.ClientMessageCh
}
*/
func (c *ServerConn) Read(buf []byte) (int, error) {
	return c.br.Read(buf)
}

func (c *ServerConn) Write(buf []byte) (int, error) {
	//	c.m.Lock()
	//	defer c.m.Unlock()
	return c.bw.Write(buf)
}

func (c *ServerConn) ColorMap() *common.ColorMap {
	return c.colorMap
}

func (c *ServerConn) SetColorMap(cm *common.ColorMap) {
	c.colorMap = cm
}
func (c *ServerConn) DesktopName() string {
	return c.desktopName
}
func (c *ServerConn) PixelFormat() *common.PixelFormat {
	return c.pixelFormat
}
func (c *ServerConn) SetDesktopName(name string) {
	c.desktopName = name
}
func (c *ServerConn) SetPixelFormat(pf *common.PixelFormat) error {
	c.pixelFormat = pf
	return nil
}
func (c *ServerConn) Encodings() []common.Encoding {
	return c.encodings
}
func (c *ServerConn) Width() uint16 {
	return c.fbWidth
}
func (c *ServerConn) Height() uint16 {
	return c.fbHeight
}
func (c *ServerConn) Protocol() string {
	return c.protocol
}

// TODO send desktopsize pseudo encoding
func (c *ServerConn) SetWidth(w uint16) {
	c.fbWidth = w
}
func (c *ServerConn) SetHeight(h uint16) {
	c.fbHeight = h
}
