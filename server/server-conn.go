package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
	"vncproxy/common"
	"vncproxy/logger"
)

type ServerConn struct {
	c   io.ReadWriter
	cfg *ServerConfig
	//br       *bufio.Reader
	//bw       *bufio.Writer
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

	// a consumer for the parsed messages, to allow for recording and proxy
	Listeners *common.MultiListener

	SessionId string

	quit chan struct{}
}

// func (c *ServerConn) UnreadByte() error {
// 	return c.br.UnreadByte()
// }

func NewServerConn(c io.ReadWriter, cfg *ServerConfig) (*ServerConn, error) {
	// if cfg.ClientMessageCh == nil {
	// 	return nil, fmt.Errorf("ClientMessageCh nil")
	// }

	if len(cfg.ClientMessages) == 0 {
		return nil, fmt.Errorf("ClientMessage 0")
	}

	return &ServerConn{
		c: c,
		//br:          bufio.NewReader(c),
		//bw:          bufio.NewWriter(c),
		cfg:         cfg,
		quit:        make(chan struct{}),
		encodings:   cfg.Encodings,
		pixelFormat: cfg.PixelFormat,
		fbWidth:     cfg.Width,
		fbHeight:    cfg.Height,
		Listeners:   &common.MultiListener{},
	}, nil
}

func (c *ServerConn) Conn() io.ReadWriter {
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

// func (c *ServerConn) Flush() error {
// 	//	c.m.Lock()
// 	//	defer c.m.Unlock()
// 	return c.bw.Flush()
// }

func (c *ServerConn) Close() error {
	return c.c.(io.ReadWriteCloser).Close()
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
	return c.c.Read(buf)
}

func (c *ServerConn) Write(buf []byte) (int, error) {
	//	c.m.Lock()
	//	defer c.m.Unlock()
	return c.c.Write(buf)
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
func (c *ServerConn) CurrentPixelFormat() *common.PixelFormat {
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

func (c *ServerConn) handle() error {
	//var err error
	//var wg sync.WaitGroup

	//defer c.Close()
	defer func() {
		c.Listeners.Consume(&common.RfbSegment{
			SegmentType: common.SegmentConnectionClosed,
		})
	}()
	//create a map of all message types
	clientMessages := make(map[common.ClientMessageType]common.ClientMessage)
	for _, m := range c.cfg.ClientMessages {
		clientMessages[m.Type()] = m
	}
	//wg.Add(2)

	// server
	// go func() error {
	// 	//defer wg.Done()
	// 	for {
	// 		select {
	// 		case msg := <-c.cfg.ServerMessageCh:
	// 			logger.Debugf("%v", msg)
	// 			// if err = msg.Write(c); err != nil {
	// 			// 	return err
	// 			// }
	// 		case <-c.quit:
	// 			c.Close()
	// 			return nil
	// 		}
	// 	}
	// }()

	// client
	//go func() error {
	//defer wg.Done()
	for {
		select {
		case <-c.quit:
			return nil
		default:
			var messageType common.ClientMessageType
			if err := binary.Read(c, binary.BigEndian, &messageType); err != nil {
				logger.Errorf("ServerConn.handle error: %v", err)
				return err
			}
			msg, ok := clientMessages[messageType]
			if !ok {
				return fmt.Errorf("ServerConn.Handle: unsupported message-type: %v", messageType)
			}

			parsedMsg, err := msg.Read(c)

			//update connection for pixel format / color map changes
			switch parsedMsg.Type() {
			case common.SetPixelFormatMsgType:
				// update pixel format
				logger.Debugf("ClientUpdater.Consume: updating pixel format")
				pixFmtMsg := parsedMsg.(*SetPixelFormat)
				c.SetPixelFormat(&pixFmtMsg.PF)
				if pixFmtMsg.PF.TrueColor != 0 {
					c.SetColorMap(&common.ColorMap{})
				}
			}
			////////

			if err != nil {
				logger.Errorf("srv err %s", err.Error())
				return err
			}

			//logger.Debugf("ServerConn.Handle got client message, type=%s", parsedMsg.Type())
			logger.Debugf("ServerConn.Handle got ClientMessage: %s, %v", parsedMsg.Type(), parsedMsg)
			//parsedMsg.Type()
			seg := &common.RfbSegment{
				SegmentType: common.SegmentFullyParsedClientMessage,
				Message:     parsedMsg,
			}
			err = c.Listeners.Consume(seg)
			if err != nil {
				logger.Errorf("ServerConn.Handle: listener consume err %s", err.Error())
				return err
			}

			//c.cfg.ClientMessageCh <- parsedMsg
		}
	}
	//}()

	//wg.Wait()
	//return nil
}
