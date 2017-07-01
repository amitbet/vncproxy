package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"bytes"

	"github.com/gorilla/websocket"
)

type WsServer1 struct {
	cfg *ServerConfig
}

type WsHandler1 func(io.ReadWriter, *ServerConfig)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

type WsConnection struct {
	Reader WsReader
	Writer WsWriter
}

func NewWsConnection(c *websocket.Conn) *WsConnection {
	return &WsConnection{
		WsReader{},
		WsWriter{c},
	}
}

type WsWriter struct {
	conn *websocket.Conn
}

type WsReader struct {
	Buff bytes.Buffer
}

func (wr WsReader) Read(p []byte) (n int, err error) {
	return wr.Buff.Read(p)
}

func (wr WsWriter) Write(p []byte) (int, error) {
	err := wr.conn.WriteMessage(websocket.BinaryMessage, p)
	return len(p), err
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Print("got connection:", r.URL)
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	myConn := NewWsConnection(c)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.BinaryMessage {
			myConn.Reader.Buff.Write(message)
		}
		log.Printf("recv: %s", message)
		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

// This example demonstrates a trivial echo server.
func (wsServer *WsServer1) Listen(urlStr string, handlerFunc WsHandler) {
	//http.Handle("/", websocket.Handler(EchoHandler))
	if urlStr == "" {
		urlStr = "/"
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("error while parsing url: ", err)
	}

	http.HandleFunc(url.Path, handleConnection)

	err = http.ListenAndServe(url.Host, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
