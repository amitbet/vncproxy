package server

import (
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

//var _ ServerConn = (*ServerConn)(nil)

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
	Encodings   []common.Encoding
	PixelFormat *common.PixelFormat
	ColorMap    *common.ColorMap
	//ClientMessageCh chan common.ClientMessage
	//ServerMessageCh chan common.ServerMessage
	ClientMessages  []common.ClientMessage
	DesktopName     []byte
	Height          uint16
	Width           uint16
	UseDummySession bool
	//handler to allow for registering for messages, this can't be a channel
	//because of the websockets handler function which will kill the connection on exit if conn.handle() is run on another thread
	NewConnHandler ServerHandler
}

func wsHandlerFunc(ws io.ReadWriter, cfg *ServerConfig, sessionId string) {
	// header := ws.Request().Header
	// url := ws.Request().URL
	// //stam := header.Get("Origin")
	// logger.Debugf("header: %v\nurl: %v", header, url)
	// io.Copy(ws, ws)

	err := attachNewServerConn(ws, cfg, sessionId)
	if err != nil {
		log.Fatalf("Error attaching new connection. %v", err)
	}
}

func WsServe(url string, cfg *ServerConfig) error {
	//server := WsServer1{cfg}
	server := WsServer{cfg}
	server.Listen(url, WsHandler(wsHandlerFunc))
	return nil
}

func TcpServe(url string, cfg *ServerConfig) error {
	ln, err := net.Listen("tcp", url)
	if err != nil {
		log.Fatalf("Error listen. %v", err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			return err
		}
		go attachNewServerConn(c, cfg, "dummySession")
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}

func attachNewServerConn(c io.ReadWriter, cfg *ServerConfig, sessionId string) error {

	conn, err := NewServerConn(c, cfg)
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

	conn.SessionId = sessionId
	if cfg.UseDummySession {
		conn.SessionId = "dummySession"
	}
	cfg.NewConnHandler(cfg, conn)

	//go here will kill ws connections
	conn.handle()

	return nil
}
