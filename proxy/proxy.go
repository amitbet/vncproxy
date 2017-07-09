package proxy

import (
	"fmt"
	"log"
	"net"
	"path"
	"strconv"
	"time"
	"vncproxy/client"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/server"
	listeners "vncproxy/tee-listeners"
)

type VncProxy struct {
	tcpListeningUrl       string      // empty = not listening on tcp
	wsListeningUrl        string      // empty = not listening on ws
	recordingDir          string      // empty = no recording
	proxyPassword         string      // empty = no auth
	targetServersPassword string      //empty = no auth
	singleSession         *VncSession // to be used when not using sessions
	usingSessions         bool        //false = single session - defined in the var above
	sessionManager        *SessionManager
}

func (vp *VncProxy) connectToVncServer(targetServerUrl string) (*client.ClientConn, error) {
	nc, err := net.Dial("tcp", targetServerUrl)

	if err != nil {
		fmt.Printf("error connecting to vnc server: %s", err)
		return nil, err
	}

	var noauth client.ClientAuthNone
	authArr := []client.ClientAuth{&client.PasswordAuth{Password: vp.targetServersPassword}, &noauth}

	vncSrvMessagesChan := make(chan common.ServerMessage)

	//rec := listeners.NewRecorder("recording.rbs")

	// split := &listeners.MultiListener{}
	// for _, listener := range rfbListeners {
	// 	split.AddListener(listener)
	// }

	clientConn, err := client.Client(nc,
		&client.ClientConfig{
			Auth:            authArr,
			ServerMessageCh: vncSrvMessagesChan,
			Exclusive:       true,
		})
	//clientConn.Listener = split

	if err != nil {
		fmt.Printf("error creating client: %s", err)
		return nil, err
	}

	tight := encodings.TightEncoding{}
	tightPng := encodings.TightPngEncoding{}
	rre := encodings.RREEncoding{}
	zlib := encodings.ZLibEncoding{}
	zrle := encodings.ZRLEEncoding{}
	cpyRect := encodings.CopyRectEncoding{}
	coRRE := encodings.CoRREEncoding{}
	hextile := encodings.HextileEncoding{}

	clientConn.SetEncodings([]common.Encoding{&cpyRect, &tightPng, &tight, &hextile, &coRRE, &rre, &zlib, &zrle})
	return clientConn, nil
}

// if sessions not enabled, will always return the configured target server (only one)
func (vp *VncProxy) getTargetServerFromSession(sessionId string) (*VncSession, error) {

	if !vp.usingSessions {
		return vp.singleSession, nil
	}
	return vp.sessionManager.GetSession(sessionId)
}

func (vp *VncProxy) newServerConnHandler(cfg *server.ServerConfig, sconn *server.ServerConn, rfbListeners []common.SegmentConsumer) error {

	recFile := "recording" + strconv.FormatInt(time.Now().Unix(), 10) + ".rbs"
	recPath := path.Join(vp.recordingDir, recFile)
	rec := listeners.NewRecorder(recPath)
	session, err := vp.getTargetServerFromSession(sconn.SessionId)
	if err != nil {
		fmt.Printf("Proxy.newServerConnHandler can't get session: %d\n", sconn.SessionId)
		return err
	}

	serverSplitter := &listeners.MultiListener{}
	for _, l := range rfbListeners {
		serverSplitter.AddListener(l)
	}
	serverSplitter.AddListener(rec)
	sconn.Listener = serverSplitter

	clientSplitter := &listeners.MultiListener{}
	clientSplitter.AddListener(rec)

	cconn, err := vp.connectToVncServer(session.TargetHostname + ":" + session.TargetPort)
	cconn.Listener = clientSplitter

	//creating cross-listeners between server and client parts to pass messages through the proxy:

	// gets the bytes from the actual vnc server on the env (client part of the proxy)
	// and writes them through the server socket to the vnc-client
	serverMsgRepeater := &listeners.WriteTo{sconn, "vnc-client bound"}
	clientSplitter.AddListener(serverMsgRepeater)

	// gets the messages from the server part (from vnc-client),
	// and write through the client to the actual vnc-server
	clientMsgRepeater := &listeners.WriteTo{cconn, "vnc-server bound"}
	serverSplitter.AddListener(clientMsgRepeater)
	return nil
}

func (vp *VncProxy) StartListening(rfbListeners []common.SegmentConsumer) {

	//chServer := make(chan common.ClientMessage)
	chClient := make(chan common.ServerMessage)

	secHandlers := []server.SecurityHandler{&server.ServerAuthNone{}}

	if vp.proxyPassword != "" {
		secHandlers = []server.SecurityHandler{&server.ServerAuthVNC{vp.proxyPassword}}
	}
	cfg := &server.ServerConfig{
		SecurityHandlers: secHandlers,
		Encodings:        []common.Encoding{&encodings.RawEncoding{}, &encodings.TightEncoding{}, &encodings.CopyRectEncoding{}},
		PixelFormat:      common.NewPixelFormat(32),
		ServerMessageCh:  chClient,
		ClientMessages:   server.DefaultClientMessages,
		DesktopName:      []byte("workDesk"),
		Height:           uint16(768),
		Width:            uint16(1024),
		NewConnHandler: func(cfg *server.ServerConfig, conn *server.ServerConn) error {
			vp.newServerConnHandler(cfg, conn, rfbListeners)
			return nil
		},
	}

	if vp.wsListeningUrl != "" {
		go server.WsServe(vp.wsListeningUrl, cfg)
	}
	if vp.tcpListeningUrl != "" {
		go server.TcpServe(vp.tcpListeningUrl, cfg)
	}

	// Process messages coming in on the ClientMessage channel.
	for {
		msg := <-chClient
		switch msg.Type() {
		default:
			log.Printf("Received message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}
