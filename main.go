package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/vnc"
)

func main() {
	//fmt.Println("")
	//nc, err := net.Dial("tcp", "192.168.1.101:5903")
	nc, err := net.Dial("tcp", "localhost:5903")

	if err != nil {
		fmt.Printf(";error connecting to vnc server: %s", err)
	}
	var noauth vnc.ClientAuthNone
	authArr := []vnc.ClientAuth{&vnc.PasswordAuth{Password: "Ch_#!T@8"}, &noauth}

	vncSrvMessagesChan := make(chan vnc.ServerMessage)
	clientConn, err := vnc.Client(nc, &vnc.ClientConfig{Auth: authArr, ServerMessageCh: vncSrvMessagesChan, Exclusive: true})
	if err != nil {
		fmt.Printf("error creating client: %s", err)
	}
	// err = clientConn.FramebufferUpdateRequest(false, 0, 0, 1024, 768)
	// if err != nil {
	// 	fmt.Printf("error requesting fb update: %s\n", err)
	// }

	tight := encodings.TightEncoding{}
	//rre := encodings.RREEncoding{}
	//zlib := encodings.ZLibEncoding{}
	//zrle := encodings.ZRLEEncoding{}
	cpyRect := encodings.CopyRectEncoding{}
	//coRRE := encodings.CoRREEncoding{}
	//hextile := encodings.HextileEncoding{}
	file, _ := os.OpenFile("stam.bin", os.O_CREATE|os.O_RDWR, 0755)
	defer file.Close()

	tight.SetOutput(file)
	clientConn.SetEncodings([]common.Encoding{&cpyRect, &tight})

	go func() {
		for {
			err = clientConn.FramebufferUpdateRequest(true, 0, 0, 1280, 800)
			if err != nil {
				fmt.Printf("error requesting fb update: %s\n", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	//go func() {
	for msg := range vncSrvMessagesChan {
		fmt.Printf("message type: %d, content: %v\n", msg.Type(), msg)
	}
	//}()

	//clientConn.Close()
}
