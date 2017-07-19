package server

import (
	"io"
	"net/http"
	"net/url"
	"vncproxy/logger"

	"golang.org/x/net/websocket"
)

type WsServer struct {
	cfg *ServerConfig
}

type WsHandler func(io.ReadWriter, *ServerConfig, string)


func (wsServer *WsServer) Listen(urlStr string, handlerFunc WsHandler) {

	if urlStr == "" {
		urlStr = "/"
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		logger.Errorf("error while parsing url: ", err)
	}

	http.Handle(url.Path, websocket.Handler(
		func(ws *websocket.Conn) {
			path := ws.Request().URL.Path
			var sessionId string
			if path != "" {
				sessionId = path[1:]
			}

			ws.PayloadType = websocket.BinaryFrame
			handlerFunc(ws, wsServer.cfg, sessionId)
		}))

	err = http.ListenAndServe(url.Host, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
