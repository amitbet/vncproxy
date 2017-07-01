package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/websocket"
)

type WsServer struct {
	cfg *ServerConfig
}

type WsHandler func(io.ReadWriter, *ServerConfig)

// This example demonstrates a trivial echo server.
func (wsServer *WsServer) Listen(urlStr string, handlerFunc WsHandler) {
	//http.Handle("/", websocket.Handler(EchoHandler))
	if urlStr == "" {
		urlStr = "/"
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("error while parsing url: ", err)
	}

	http.Handle(url.Path, websocket.Handler(func(ws *websocket.Conn) {
		// header := ws.Request().Header
		// url := ws.Request().URL
		// //stam := header.Get("Origin")
		// fmt.Printf("header: %v\nurl: %v\n", header, url)
		// io.Copy(ws, ws)
		ws.PayloadType = websocket.BinaryFrame
		handlerFunc(ws, wsServer.cfg)
	}))

	err = http.ListenAndServe(url.Host, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
