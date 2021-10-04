package proxy

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"net"
	"path"
	"strconv"
	"time"

	"github.com/amitbet/vncproxy/client"
	"github.com/amitbet/vncproxy/common"
	"github.com/amitbet/vncproxy/encodings"
	"github.com/amitbet/vncproxy/logger"
	"github.com/amitbet/vncproxy/player"
	listeners "github.com/amitbet/vncproxy/recorder"
	"github.com/amitbet/vncproxy/server"
)

type VncProxy struct {
	TCPListeningURL  string      // empty = not listening on tcp
	WsListeningURL   string      // empty = not listening on ws
	RecordingDir     string      // empty = no recording
	ProxyVncPassword string      //empty = no auth
	SingleSession    *VncSession // to be used when not using sessions
	UsingSessions    bool        //false = single session - defined in the var above
	Remarkable       bool
	sessionManager   *SessionManager
}

func (vp *VncProxy) createClientConnection(target string, vncPass string, tlsEnabled bool, reMarkableDeviceId string) (*client.ClientConn, error) {
	var (
		nc  net.Conn
		err error
	)

	if target[0] == '/' {
		nc, err = net.Dial("unix", target)
	} else {
		nc, err = net.Dial("tcp", target)
		if tlsEnabled {
			logger.Info("Upgrading to TLS connection...")
			config := tls.Config{
				InsecureSkipVerify: true,
			}
			tc := tls.Client(nc, &config)
			err = tc.Handshake()
			if err != nil {
				return nil, err
			}

			nc = tc
		}
	}

	if err != nil {
		logger.Errorf("error connecting to vnc server: %s", err)
		return nil, err
	}

	var noauth client.ClientAuthNone
	authArr := []client.ClientAuth{&client.PasswordAuth{Password: vncPass}, &noauth}
	var rmTimestamp uint64
	if reMarkableDeviceId != "" {
		logger.Info("Enabling reMarkable 2.10+ support")

		authC, err := net.ListenPacket("udp", ":5901")
		if err != nil {
			return nil, err
		}
		defer authC.Close()

		buffer := make([]byte, 128) // datagram size 51 bytes
		n, addr, err := authC.ReadFrom(buffer)
		logger.Debugf("Received datagram from reMarkable (%v, %v bytes)", addr, n)

		pb := bytes.NewBuffer(buffer)
		err = binary.Read(pb, binary.BigEndian, &rmTimestamp)
		if err != nil {
			return nil, err
		}
		logger.Debugf("reMarkable timestamp is %v", rmTimestamp)
		authArr = append(authArr, &client.RemarkableAuth{
			RemarkableTimestamp: rmTimestamp,
			RemarkableDeviceId:  reMarkableDeviceId,
		})
	}

	clientConn, err := client.NewClientConn(nc,
		&client.ClientConfig{
			Auth:      authArr,
			Exclusive: true,
		})

	if err != nil {
		logger.Errorf("error creating client: %s", err)
		return nil, err
	}

	return clientConn, nil
}

// if sessions not enabled, will always return the configured target server (only one)
func (vp *VncProxy) getProxySession(sessionId string) (*VncSession, error) {

	if !vp.UsingSessions {
		if vp.SingleSession == nil {
			logger.Errorf("SingleSession is empty, use sessions or populate the SingleSession member of the VncProxy struct.")
		}
		return vp.SingleSession, nil
	}
	return vp.sessionManager.GetSession(sessionId)
}

func (vp *VncProxy) newServerConnHandler(cfg *server.ServerConfig, sconn *server.ServerConn) error {
	var err error
	session, err := vp.getProxySession(sconn.SessionId)
	if err != nil {
		logger.Errorf("Proxy.newServerConnHandler can't get session: %d", sconn.SessionId)
		return err
	}

	var rec *listeners.Recorder

	if session.Type == SessionTypeRecordingProxy {
		recFile := "recording" + strconv.FormatInt(time.Now().Unix(), 10) + ".rbs"
		recPath := path.Join(vp.RecordingDir, recFile)
		rec, err = listeners.NewRecorder(recPath)
		if err != nil {
			logger.Errorf("Proxy.newServerConnHandler can't open recorder save path: %s", recPath)
			return err
		}

		sconn.Listeners.AddListener(rec)
	}

	session.Status = SessionStatusInit
	if session.Type == SessionTypeProxyPass || session.Type == SessionTypeRecordingProxy {
		target := session.Target
		if session.TargetHostname != "" && session.TargetPort != "" {
			target = session.TargetHostname + ":" + session.TargetPort
		}

		cconn, err := vp.createClientConnection(target, session.TargetPassword, session.TLS, session.RemarkableDeviceId)
		if err != nil {
			session.Status = SessionStatusError
			logger.Errorf("Proxy.newServerConnHandler error creating connection: %s", err)
			return err
		}
		if session.Type == SessionTypeRecordingProxy {
			cconn.Listeners.AddListener(rec)
		}

		//creating cross-listeners between server and client parts to pass messages through the proxy:

		// gets the bytes from the actual vnc server on the env (client part of the proxy)
		// and writes them through the server socket to the vnc-client
		serverUpdater := &ServerUpdater{sconn}
		cconn.Listeners.AddListener(serverUpdater)

		// gets the messages from the server part (from vnc-client),
		// and write through the client to the actual vnc-server
		var clientUpdater *ClientUpdater
		clientUpdater = &ClientUpdater{conn: cconn}
		if session.RemarkableDeviceId != "" {
			clientUpdater.suppressedMessageTypes = []common.ClientMessageType{common.SetPixelFormatMsgType}
		}
		if len(session.OverrideEncodings) > 0 {
			for _, enc := range session.OverrideEncodings {
				clientUpdater.overrideEncodings = append(clientUpdater.overrideEncodings, common.EncodingType(enc))
			}
		}
		sconn.Listeners.AddListener(clientUpdater)

		err = cconn.Connect()
		if err != nil {
			session.Status = SessionStatusError
			logger.Errorf("Proxy.newServerConnHandler error connecting to client: %s", err)
			return err
		}

		encs := []common.IEncoding{
			&encodings.RawEncoding{},
			&encodings.TightEncoding{},
			&encodings.EncCursorPseudo{},
			&encodings.EncLedStatePseudo{},
			&encodings.TightPngEncoding{},
			&encodings.RREEncoding{},
			&encodings.ZLibEncoding{},
			&encodings.ZRLEEncoding{},
			&encodings.CopyRectEncoding{},
			&encodings.CoRREEncoding{},
			&encodings.HextileEncoding{},
		}
		cconn.Encs = encs

		if err != nil {
			session.Status = SessionStatusError
			logger.Errorf("Proxy.newServerConnHandler error connecting to client: %s", err)
			return err
		}
	}

	if session.Type == SessionTypeReplayServer {
		fbs, err := player.ConnectFbsFile(session.ReplayFilePath, sconn)

		if err != nil {
			logger.Error("TestServer.NewConnHandler: Error in loading FBS: ", err)
			return err
		}
		sconn.Listeners.AddListener(player.NewFBSPlayListener(sconn, fbs))
		return nil

	}

	session.Status = SessionStatusActive
	return nil
}

func (vp *VncProxy) StartListening() {

	secHandlers := []server.SecurityHandler{&server.ServerAuthNone{}}

	if vp.ProxyVncPassword != "" {
		secHandlers = []server.SecurityHandler{&server.ServerAuthVNC{vp.ProxyVncPassword}}
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
	}
	if vp.Remarkable {
		cfg.DesktopName = []byte("reMarkable rfb")
	}

	if vp.TCPListeningURL != "" && vp.WsListeningURL != "" {
		logger.Infof("running two listeners: tcp port: %s, ws url: %s", vp.TCPListeningURL, vp.WsListeningURL)

		go server.WsServe(vp.WsListeningURL, cfg)
		server.TcpServe(vp.TCPListeningURL, cfg)
	}

	if vp.WsListeningURL != "" {
		logger.Infof("running ws listener url: %s", vp.WsListeningURL)
		server.WsServe(vp.WsListeningURL, cfg)
	}
	if vp.TCPListeningURL != "" {
		logger.Infof("running tcp listener on port: %s", vp.TCPListeningURL)
		server.TcpServe(vp.TCPListeningURL, cfg)
	}
}
