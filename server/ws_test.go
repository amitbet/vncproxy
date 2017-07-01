package server

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"testing"

// 	"golang.org/x/net/websocket"
// )

// func TestWsServer(t *testing.T) {
// 	server := WsServer{}
// 	server.Listen(":8090")
// }

// // Echo the data received on the WebSocket.
// func EchoHandler(ws *websocket.Conn) {
// 	header := ws.Request().Header
// 	url := ws.Request().URL
// 	//stam := header.Get("Origin")
// 	fmt.Printf("header: %v\nurl: %v\n", header, url)
// 	io.Copy(ws, ws)
// }

// // This example demonstrates a trivial echo server.
// func TestGoWsServer(t *testing.T) {
// 	http.Handle("/", websocket.Handler(EchoHandler))
// 	err := http.ListenAndServe(":11111", nil)
// 	if err != nil {
// 		panic("ListenAndServe: " + err.Error())
// 	}
// }
