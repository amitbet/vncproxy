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

type WsHandler func(io.ReadWriter, *ServerConfig, string)

// func checkOrigin(config *websocket.Config, req *http.Request) (err error) {
// 	config.Origin, err = websocket.Origin(config, req)
// 	if err == nil && config.Origin == nil {
// 		return fmt.Errorf("null origin")
// 	}
// 	return err
// }

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

	// http.HandleFunc(url.Path,
	// 	func(w http.ResponseWriter, req *http.Request) {
	// 		sessionId := req.URL.Query().Get("sessionId")
	// 		s := websocket.Server{Handshake: checkOrigin, Handler: websocket.Handler(
	// 			func(ws *websocket.ServerConn) {
	// 				ws.PayloadType = websocket.BinaryFrame
	// 				handlerFunc(ws, wsServer.cfg, sessionId)
	// 			})}
	// 		s.ServeHTTP(w, req)
	// 	})

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
