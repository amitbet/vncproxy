package server

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"vncproxy/common"
)

var DefaultClientMessages = []common.ClientMessage{
	&SetPixelFormat{},
	&SetEncodings{},
	&FramebufferUpdateRequest{},
	&KeyEvent{},
	&PointerEvent{},
	&ClientCutText{},
}

type ServerInit struct {
	FBWidth, FBHeight uint16
	PixelFormat       common.PixelFormat
	NameLength        uint32
	NameText          []byte
}

//var _ Conn = (*ServerConn)(nil)

// ServerMessage represents a Client-to-Server RFB message type.
// type ServerMessageType uint8

// //go:generate stringer -type=ServerMessageType

// // Client-to-Server message types.
// const (
// 	FramebufferUpdateMsgType ServerMessageType = iota
// 	SetColorMapEntriesMsgType
// 	BellMsgType
// 	ServerCutTextMsgType
// )

// FramebufferUpdate holds a FramebufferUpdate wire format message.
type FramebufferUpdate struct {
	_       [1]byte             // pad
	NumRect uint16              // number-of-rectangles
	Rects   []*common.Rectangle // rectangles
}

// func (*FramebufferUpdate) Type() ServerMessageType {
// 	return FramebufferUpdateMsgType
// }

type ServerHandler func(*ServerConfig, *ServerConn) error

type ServerConfig struct {
	//VersionHandler    ServerHandler
	//SecurityHandler   ServerHandler
	SecurityHandlers []SecurityHandler
	//ClientInitHandler ServerHandler
	//ServerInitHandler ServerHandler
	Encodings       []common.Encoding
	PixelFormat     *common.PixelFormat
	ColorMap        *common.ColorMap
	ClientMessageCh chan common.ClientMessage
	ServerMessageCh chan common.ServerMessage
	ClientMessages  []common.ClientMessage
	DesktopName     []byte
	Height          uint16
	Width           uint16
}

func NewServerConn(c net.Conn, cfg *ServerConfig) (*ServerConn, error) {
	if cfg.ClientMessageCh == nil {
		return nil, fmt.Errorf("ClientMessageCh nil")
	}

	if len(cfg.ClientMessages) == 0 {
		return nil, fmt.Errorf("ClientMessage 0")
	}

	return &ServerConn{
		c:           c,
		br:          bufio.NewReader(c),
		bw:          bufio.NewWriter(c),
		cfg:         cfg,
		quit:        make(chan struct{}),
		encodings:   cfg.Encodings,
		pixelFormat: cfg.PixelFormat,
		fbWidth:     cfg.Width,
		fbHeight:    cfg.Height,
	}, nil
}

func Serve(ctx context.Context, ln net.Listener, cfg *ServerConfig) error {
	for {

		c, err := ln.Accept()
		if err != nil {
			continue
		}

		conn, err := NewServerConn(c, cfg)
		if err != nil {
			continue
		}

		if err := ServerVersionHandler(cfg, conn); err != nil {
			conn.Close()
			continue
		}

		if err := ServerSecurityHandler(cfg, conn); err != nil {
			conn.Close()
			continue
		}

		if err := ServerClientInitHandler(cfg, conn); err != nil {
			conn.Close()
			continue
		}

		if err := ServerServerInitHandler(cfg, conn); err != nil {
			conn.Close()
			continue
		}

		go conn.Handle()
	}
}

func (c *ServerConn) Handle() error {
	//var err error
	var wg sync.WaitGroup

	defer c.Close()

	//create a map of all message types
	clientMessages := make(map[common.ClientMessageType]common.ClientMessage)
	for _, m := range c.cfg.ClientMessages {
		clientMessages[m.Type()] = m
	}
	wg.Add(2)

	// server
	go func() error {
		defer wg.Done()
		for {
			select {
			case msg := <-c.cfg.ServerMessageCh:
				fmt.Printf("%v", msg)
				// if err = msg.Write(c); err != nil {
				// 	return err
				// }
			case <-c.quit:
				return nil
			}
		}
	}()

	// client
	go func() error {
		defer wg.Done()
		for {
			select {
			case <-c.quit:
				return nil
			default:
				var messageType common.ClientMessageType
				if err := binary.Read(c, binary.BigEndian, &messageType); err != nil {
					return err
				}
				msg, ok := clientMessages[messageType]
				if !ok {
					return fmt.Errorf("unsupported message-type: %v", messageType)

				}
				parsedMsg, err := msg.Read(c)
				if err != nil {
					fmt.Printf("srv err %s\n", err.Error())
					return err
				}
				fmt.Printf("message:%s, %v\n",parsedMsg.Type(), parsedMsg)
				//c.cfg.ClientMessageCh <- parsedMsg
			}
		}
	}()

	wg.Wait()
	return nil
}

// type ServerCutText struct {
// 	_      [1]byte
// 	Length uint32
// 	Text   []byte
// }

// func (*ServerCutText) Type() ServerMessageType {
// 	return ServerCutTextMsgType
// }

// func (*ServerCutText) Read(c common.Conn) (common.ServerMessage, error) {
// 	msg := ServerCutText{}

// 	var pad [1]byte
// 	if err := binary.Read(c, binary.BigEndian, &pad); err != nil {
// 		return nil, err
// 	}

// 	if err := binary.Read(c, binary.BigEndian, &msg.Length); err != nil {
// 		return nil, err
// 	}

// 	msg.Text = make([]byte, msg.Length)
// 	if err := binary.Read(c, binary.BigEndian, &msg.Text); err != nil {
// 		return nil, err
// 	}
// 	return &msg, nil
// }

// func (msg *ServerCutText) Write(c common.Conn) error {
// 	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
// 		return err
// 	}
// 	var pad [1]byte
// 	if err := binary.Write(c, binary.BigEndian, pad); err != nil {
// 		return err
// 	}

// 	if msg.Length < uint32(len(msg.Text)) {
// 		msg.Length = uint32(len(msg.Text))
// 	}
// 	if err := binary.Write(c, binary.BigEndian, msg.Length); err != nil {
// 		return err
// 	}

// 	if err := binary.Write(c, binary.BigEndian, msg.Text); err != nil {
// 		return err
// 	}
// 	return nil
// }

// type Bell struct{}

// func (*Bell) Type() ServerMessageType {
// 	return BellMsgType
// }

// func (*Bell) Read(c common.Conn) (common.ServerMessage, error) {
// 	return &Bell{}, nil
// }

// func (msg *Bell) Write(c common.Conn) error {
// 	return binary.Write(c, binary.BigEndian, msg.Type())
// }

// type SetColorMapEntries struct {
// 	_          [1]byte
// 	FirstColor uint16
// 	ColorsNum  uint16
// 	Colors     []common.Color
// }

// func (*SetColorMapEntries) Type() ServerMessageType {
// 	return SetColorMapEntriesMsgType
// }

// func (*SetColorMapEntries) Read(c common.Conn) (common.ServerMessage, error) {
// 	msg := SetColorMapEntries{}
// 	var pad [1]byte
// 	if err := binary.Read(c, binary.BigEndian, &pad); err != nil {
// 		return nil, err
// 	}

// 	if err := binary.Read(c, binary.BigEndian, &msg.FirstColor); err != nil {
// 		return nil, err
// 	}

// 	if err := binary.Read(c, binary.BigEndian, &msg.ColorsNum); err != nil {
// 		return nil, err
// 	}

// 	msg.Colors = make([]common.Color, msg.ColorsNum)
// 	colorMap := c.ColorMap()

// 	for i := uint16(0); i < msg.ColorsNum; i++ {
// 		color := &msg.Colors[i]
// 		if err := binary.Read(c, binary.BigEndian, &color); err != nil {
// 			return nil, err
// 		}
// 		colorMap[msg.FirstColor+i] = *color
// 	}
// 	c.SetColorMap(colorMap)
// 	return &msg, nil
// }

// func (msg *SetColorMapEntries) Write(c common.Conn) error {
// 	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
// 		return err
// 	}
// 	var pad [1]byte
// 	if err := binary.Write(c, binary.BigEndian, &pad); err != nil {
// 		return err
// 	}

// 	if err := binary.Write(c, binary.BigEndian, msg.FirstColor); err != nil {
// 		return err
// 	}

// 	if msg.ColorsNum < uint16(len(msg.Colors)) {
// 		msg.ColorsNum = uint16(len(msg.Colors))
// 	}
// 	if err := binary.Write(c, binary.BigEndian, msg.ColorsNum); err != nil {
// 		return err
// 	}

// 	for i := 0; i < len(msg.Colors); i++ {
// 		color := msg.Colors[i]
// 		if err := binary.Write(c, binary.BigEndian, color); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (*FramebufferUpdate) Read(cliInfo common.IClientConn, c *common.RfbReadHelper) (common.ServerMessage, error) {
// 	msg := FramebufferUpdate{}
// 	var pad [1]byte
// 	if err := binary.Read(c, binary.BigEndian, &pad); err != nil {
// 		return nil, err
// 	}

// 	if err := binary.Read(c, binary.BigEndian, &msg.NumRect); err != nil {
// 		return nil, err
// 	}
// 	for i := uint16(0); i < msg.NumRect; i++ {
// 		rect := &common.Rectangle{}
// 		if err := rect.Read(c); err != nil {
// 			return nil, err
// 		}
// 		msg.Rects = append(msg.Rects, rect)
// 	}
// 	return &msg, nil
// }

// func (msg *FramebufferUpdate) Write(c common.Conn) error {
// 	if err := binary.Write(c, binary.BigEndian, msg.Type()); err != nil {
// 		return err
// 	}
// 	var pad [1]byte
// 	if err := binary.Write(c, binary.BigEndian, pad); err != nil {
// 		return err
// 	}
// 	if err := binary.Write(c, binary.BigEndian, msg.NumRect); err != nil {
// 		return err
// 	}
// 	for _, rect := range msg.Rects {
// 		if err := rect.Write(c); err != nil {
// 			return err
// 		}
// 	}
// 	return c.Flush()
// }
