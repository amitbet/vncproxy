package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
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

func newServerConn(c io.ReadWriter, cfg *ServerConfig) (*ServerConn, error) {
	if cfg.ClientMessageCh == nil {
		return nil, fmt.Errorf("ClientMessageCh nil")
	}

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
	}, nil
}
func wsHandlerFunc(ws io.ReadWriter, cfg *ServerConfig) {
	// header := ws.Request().Header
	// url := ws.Request().URL
	// //stam := header.Get("Origin")
	// fmt.Printf("header: %v\nurl: %v\n", header, url)
	// io.Copy(ws, ws)

	err := attachNewServerConn(ws, cfg)
	if err != nil {
		log.Fatalf("Error attaching new connection. %v", err)
	}
}

func WsServe(url string, ctx context.Context, cfg *ServerConfig) error {
	//server := WsServer1{cfg}
	server := WsServer{cfg}
	server.Listen(url, WsHandler(wsHandlerFunc))
	return nil
}

func TcpServe(url string, ctx context.Context, cfg *ServerConfig) error {
	ln, err := net.Listen("tcp", ":5903")
	if err != nil {
		log.Fatalf("Error listen. %v", err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			return err
		}
		go attachNewServerConn(c, cfg)
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}

func attachNewServerConn(c io.ReadWriter, cfg *ServerConfig) error {

	conn, err := newServerConn(c, cfg)
	if err != nil {
		return err
	}

	if err := ServerVersionHandler(cfg, conn); err != nil {
		fmt.Errorf("err: %v\n", err)
		conn.Close()
		return err
	}

	if err := ServerSecurityHandler(cfg, conn); err != nil {
		conn.Close()
		return err
	}

	if err := ServerClientInitHandler(cfg, conn); err != nil {
		conn.Close()
		return err
	}

	if err := ServerServerInitHandler(cfg, conn); err != nil {
		conn.Close()
		return err
	}

	//go
	conn.handle()

	return nil
}

func (c *ServerConn) handle() error {
	//var err error
	//var wg sync.WaitGroup

	//defer c.Close()

	//create a map of all message types
	clientMessages := make(map[common.ClientMessageType]common.ClientMessage)
	for _, m := range c.cfg.ClientMessages {
		clientMessages[m.Type()] = m
	}
	//wg.Add(2)

	// server
	// go func() error {
	// 	defer wg.Done()
	// 	for {
	// 		select {
	// 		case msg := <-c.cfg.ServerMessageCh:
	// 			fmt.Printf("%v", msg)
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
				fmt.Printf("Error: %v\n", err)
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
			fmt.Printf("message:%s, %v\n", parsedMsg.Type(), parsedMsg)
			//c.cfg.ClientMessageCh <- parsedMsg
		}
	}
	//}()

	//wg.Wait()
	//return nil
}
