package main

import (
	"net"
	"time"
	"vncproxy/client"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
	listeners "vncproxy/tee-listeners"
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

	//rec := listeners.NewRecorder("c:/Users/betzalel/recording.rbs")
	rec := listeners.NewRecorder("/Users/amitbet/vncRec/recording.rbs")

	clientConn, err := client.NewClientConn(nc,
		&client.ClientConfig{
			Auth: authArr,
			Exclusive: true,
		})

	clientConn.Listeners.AddListener(rec)
	clientConn.Connect()

	if err != nil {
		logger.Errorf("error creating client: %s", err)
	}
	// err = clientConn.FramebufferUpdateRequest(false, 0, 0, 1024, 768)
	// if err != nil {
	// 	logger.Errorf("error requesting fb update: %s", err)
	// }
	encs := []common.Encoding{
		&encodings.TightEncoding{},
		//&encodings.TightPngEncoding{},
		//rre := encodings.RREEncoding{},
		//zlib := encodings.ZLibEncoding{},
		//zrle := encodings.ZRLEEncoding{},
		//&encodings.CopyRectEncoding{},
		//coRRE := encodings.CoRREEncoding{},
		//hextile := encodings.HextileEncoding{},
		&encodings.PseudoEncoding{int32(common.EncJPEGQualityLevelPseudo9)},
	}

	// file, _ := os.OpenFile("stam.bin", os.O_CREATE|os.O_RDWR, 0755)
	// defer file.Close()

	//tight.SetOutput(file)
	clientConn.SetEncodings(encs)

	clientConn.FramebufferUpdateRequest(false, 0, 0, 1280, 800)
	// clientConn.SetPixelFormat(&common.PixelFormat{
	// 	BPP:        32,
	// 	Depth:      24,
	// 	BigEndian:  0,
	// 	TrueColor:  1,
	// 	RedMax:     255,
	// 	GreenMax:   255,
	// 	BlueMax:    255,
	// 	RedShift:   16,
	// 	GreenShift: 8,
	// 	BlueShift:  0,
	// })
	start := getNowMillisec()
	go func() {
		for {
			if getNowMillisec()-start >= 10000 {
				break
			}

			err = clientConn.FramebufferUpdateRequest(true, 0, 0, 1280, 800)
			if err != nil {
				logger.Errorf("error requesting fb update: %s", err)
			}
			time.Sleep(250 * time.Millisecond)
		}
		clientConn.Close()
	}()

	//go func() {
	// for msg := range vncSrvMessagesChan {
	// 	logger.Debugf("message type: %d, content: %v\n", msg.Type(), msg)
	// }
	for {
		time.Sleep(time.Minute)
	}
	//}()

	//clientConn.Close()
}
func getNowMillisec() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}
