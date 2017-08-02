package main

import (
	"net"
	"time"
	"vncproxy/client"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
	"vncproxy/recorder"
)

func main() {

	//nc, err := net.Dial("tcp", "192.168.1.101:5903")
	nc, err := net.Dial("tcp", "localhost:5903")

	if err != nil {
		logger.Errorf("error connecting to vnc server: %s", err)
	}
	var noauth client.ClientAuthNone
	authArr := []client.ClientAuth{&client.PasswordAuth{Password: "Ch_#!T@8"}, &noauth}

	//vncSrvMessagesChan := make(chan common.ServerMessage)

	//rec, err := recorder.NewRecorder("c:/Users/betzalel/recording.rbs")
	rec, err := recorder.NewRecorder("/Users/amitbet/vncRec/recording.rbs")
	if err != nil {
		logger.Errorf("error creating recorder: %s", err)
		return
	}

	clientConn, err := client.NewClientConn(nc,
		&client.ClientConfig{
			Auth:      authArr,
			Exclusive: true,
		})

	clientConn.Listeners.AddListener(rec)
	clientConn.Listeners.AddListener(&recorder.RfbRequester{Conn: clientConn, Name: "Rfb Requester"})
	clientConn.Connect()

	if err != nil {
		logger.Errorf("error creating client: %s", err)
		return
	}
	// err = clientConn.FramebufferUpdateRequest(false, 0, 0, 1024, 768)
	// if err != nil {
	// 	logger.Errorf("error requesting fb update: %s", err)
	// }
	encs := []common.IEncoding{
		&encodings.TightEncoding{},
		//&encodings.TightPngEncoding{},
		//rre := encodings.RREEncoding{},
		//zlib := encodings.ZLibEncoding{},
		//zrle := encodings.ZRLEEncoding{},
		//&encodings.CopyRectEncoding{},
		//coRRE := encodings.CoRREEncoding{},
		//hextile := encodings.HextileEncoding{},
		&encodings.PseudoEncoding{int32(common.EncJPEGQualityLevelPseudo8)},
	}

	clientConn.SetEncodings(encs)
	//width := uint16(1280)
	//height := uint16(800)

	//clientConn.FramebufferUpdateRequest(false, 0, 0, width, height)

	// // clientConn.SetPixelFormat(&common.PixelFormat{
	// // 	BPP:        32,
	// // 	Depth:      24,
	// // 	BigEndian:  0,
	// // 	TrueColor:  1,
	// // 	RedMax:     255,
	// // 	GreenMax:   255,
	// // 	BlueMax:    255,
	// // 	RedShift:   16,
	// // 	GreenShift: 8,
	// // 	BlueShift:  0,
	// // })

	// start := getNowMillisec()
	// go func() {
	// 	for {
	// 		if getNowMillisec()-start >= 10000 {
	// 			break
	// 		}

	// 		err = clientConn.FramebufferUpdateRequest(true, 0, 0, 1280, 800)
	// 		if err != nil {
	// 			logger.Errorf("error requesting fb update: %s", err)
	// 		}
	// 		time.Sleep(250 * time.Millisecond)
	// 	}
	// 	clientConn.Close()
	// }()

	for {
		time.Sleep(time.Minute)
	}
}
func getNowMillisec() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}
