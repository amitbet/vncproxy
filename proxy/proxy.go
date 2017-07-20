package proxy

import (
	"log"
	"net"
	"path"
	"strconv"
	"time"
	"vncproxy/client"
	"vncproxy/common"
	"vncproxy/encodings"
	"vncproxy/logger"
	"vncproxy/server"
	listeners "vncproxy/tee-listeners"
)

type VncProxy struct {
	tcpListeningUrl       string      // empty = not listening on tcp
	wsListeningUrl        string      // empty = not listening on ws
	recordingDir          string      // empty = no recording
	proxyPassword         string      // empty = no auth
	targetServersPassword string      //empty = no auth
	SingleSession         *VncSession // to be used when not using sessions
	UsingSessions         bool        //false = single session - defined in the var above
	sessionManager        *SessionManager
}

func (vp *VncProxy) createClientConnection(targetServerUrl string) (*client.ClientConn, error) {
	nc, err := net.Dial("tcp", targetServerUrl)

	if err != nil {
		logger.Errorf("error connecting to vnc server: %s", err)
		return nil, err
	}

	var noauth client.ClientAuthNone
	authArr := []client.ClientAuth{&client.PasswordAuth{Password: vp.targetServersPassword}, &noauth}

	//vncSrvMessagesChan := make(chan common.ServerMessage)

	clientConn, err := client.NewClientConn(nc,
		&client.ClientConfig{
			Auth:      authArr,
			Exclusive: true,
		})
	//clientConn.Listener = split

	if err != nil {
		logger.Errorf("error creating client: %s", err)
		return nil, err
	}

	return clientConn, nil
}

// if sessions not enabled, will always return the configured target server (only one)
func (vp *VncProxy) getTargetServerFromSession(sessionId string) (*VncSession, error) {

	if !vp.UsingSessions {
		if vp.SingleSession == nil {
			logger.Errorf("SingleSession is empty, use sessions or populate the SingleSession member of the VncProxy struct.")
		}
		return vp.SingleSession, nil
	}
	return vp.sessionManager.GetSession(sessionId)
}

func (vp *VncProxy) newServerConnHandler(cfg *server.ServerConfig, sconn *server.ServerConn) error {

	recFile := "recording" + strconv.FormatInt(time.Now().Unix(), 10) + ".rbs"
	recPath := path.Join(vp.recordingDir, recFile)
	rec, err := listeners.NewRecorder(recPath)
	if err != nil {
		logger.Errorf("Proxy.newServerConnHandler can't open recorder save path: %s", recPath)
		return err
	}

	session, err := vp.getTargetServerFromSession(sconn.SessionId)
	if err != nil {
		logger.Errorf("Proxy.newServerConnHandler can't get session: %d", sconn.SessionId)
		return err
	}

	// for _, l := range rfbListeners {
	// 	sconn.Listeners.AddListener(l)
	// }
	sconn.Listeners.AddListener(rec)

	//clientSplitter := &common.MultiListener{}

	cconn, err := vp.createClientConnection(session.TargetHostname + ":" + session.TargetPort)
	if err != nil {
		logger.Errorf("Proxy.newServerConnHandler error creating connection: %s", err)
		return err
	}
	cconn.Listeners.AddListener(rec)
	//cconn.Listener = clientSplitter

	//creating cross-listeners between server and client parts to pass messages through the proxy:

	// gets the bytes from the actual vnc server on the env (client part of the proxy)
	// and writes them through the server socket to the vnc-client
	serverUpdater := &ServerUpdater{sconn}
	cconn.Listeners.AddListener(serverUpdater)

	// // serverMsgRepeater := &listeners.WriteTo{sconn, "vnc-client-bound"}
	// // cconn.Listeners.AddListener(serverMsgRepeater)

	// gets the messages from the server part (from vnc-client),
	// and write through the client to the actual vnc-server
	//clientMsgRepeater := &listeners.WriteTo{cconn, "vnc-server-bound"}
	clientUpdater := &ClientUpdater{cconn}
	sconn.Listeners.AddListener(clientUpdater)

	err = cconn.Connect()
	if err != nil {
		logger.Errorf("Proxy.newServerConnHandler error connecting to client: %s", err)
		return err
	}

	encs := []common.IEncoding{
		&encodings.RawEncoding{},
		&encodings.TightEncoding{},
		&encodings.EncCursorPseudo{},
		//encodings.TightPngEncoding{},
		&encodings.RREEncoding{},
		&encodings.ZLibEncoding{},
		&encodings.ZRLEEncoding{},
		&encodings.CopyRectEncoding{},
		&encodings.CoRREEncoding{},
		&encodings.HextileEncoding{},
	}
	cconn.Encs = encs
	//err = cconn.MsgSetEncodings(encs)
	if err != nil {
		logger.Errorf("Proxy.newServerConnHandler error connecting to client: %s", err)
		return err
	}
	return nil
}

func (vp *VncProxy) StartListening() {

	//chServer := make(chan common.ClientMessage)
	chClient := make(chan common.ServerMessage)

	secHandlers := []server.SecurityHandler{&server.ServerAuthNone{}}

	if vp.proxyPassword != "" {
		secHandlers = []server.SecurityHandler{&server.ServerAuthVNC{vp.proxyPassword}}
	}
	cfg := &server.ServerConfig{
		SecurityHandlers: secHandlers,
		Encodings:        []common.IEncoding{&encodings.RawEncoding{}, &encodings.TightEncoding{}, &encodings.CopyRectEncoding{}},
		PixelFormat:      common.NewPixelFormat(32),
		ClientMessages:   server.DefaultClientMessages,
		DesktopName:      []byte("workDesk"),
		Height:           uint16(768),
		Width:            uint16(1024),
		NewConnHandler:   vp.newServerConnHandler,
		UseDummySession:  !vp.UsingSessions,
		// func(cfg *server.ServerConfig, conn *server.ServerConn) error {
		// 	vp.newServerConnHandler(cfg, conn)
		// 	return nil
		// },
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
