package main

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/amitbet/vncproxy/logger"
	vncproxy "github.com/amitbet/vncproxy/proxy"
)

func main() {
	//create default session if required
	var tcpPort = flag.String("tcpPort", "", "tcp port")
	var wsPort = flag.String("wsPort", "", "websocket port")
	var vncPass = flag.String("vncPass", "", "password on incoming vnc connections to the proxy, defaults to no password")
	var recordDir = flag.String("recDir", "", "path to save FBS recordings WILL NOT RECORD if not defined.")
	var targetVnc = flag.String("target", "", "target vnc server (host:port or /path/to/unix.socket)")
	var targetVncPort = flag.String("targPort", "", "target vnc server port (deprecated, use -target)")
	var targetVncHost = flag.String("targHost", "", "target vnc server host (deprecated, use -target)")
	var targetVncPass = flag.String("targPass", "", "target vnc password")
	var logLevel = flag.String("logLevel", "info", "change logging level")
	var reMarkable = flag.String("reMarkable", "", "reMarkable device ID (enable reMarkable 2.10+ support)")
	var overrideEncodings = flag.String("overrideEncodings", "", "force a specific set of encodings to be sent to the target vnc server (encoding types separated by comma)")
	var tls = flag.Bool("tls", false, "use TLS connection (turned on automatically if reMarkable)")
	flag.Parse()
	logger.SetLogLevel(*logLevel)

	if *tcpPort == "" && *wsPort == "" {
		logger.Error("no listening port defined")
		flag.Usage()
		os.Exit(1)
	}

	if *targetVnc == "" && *targetVncPort == "" {
		logger.Error("no target vnc server host/port or socket defined")
		flag.Usage()
		os.Exit(1)
	}

	if *vncPass == "" {
		logger.Warn("proxy will have no password")
	}

	tcpURL := ""
	if *tcpPort != "" {
		tcpURL = ":" + string(*tcpPort)
	}
	wsURL := ""
	if *wsPort != "" {
		wsURL = "http://0.0.0.0:" + string(*wsPort) + "/"
	}
	var overrideEncodingsList []uint32
	if *overrideEncodings != "" {
		for _, enc := range strings.Split(*overrideEncodings, ",") {
			encI, err := strconv.Atoi(enc)
			if err != nil {
				panic(err)
			}
			overrideEncodingsList = append(overrideEncodingsList, uint32(encI))
		}
	}
	proxy := &vncproxy.VncProxy{
		WsListeningURL:   wsURL, // empty = not listening on ws
		TCPListeningURL:  tcpURL,
		ProxyVncPassword: *vncPass, //empty = no auth
		SingleSession: &vncproxy.VncSession{
			Target:             *targetVnc,
			TargetHostname:     *targetVncHost,
			TargetPort:         *targetVncPort,
			TargetPassword:     *targetVncPass, //"vncPass",
			ID:                 "dummySession",
			Status:             vncproxy.SessionStatusInit,
			Type:               vncproxy.SessionTypeProxyPass,
			RemarkableDeviceId: *reMarkable,
			TLS:                *reMarkable != "" || *tls,
			OverrideEncodings:  overrideEncodingsList,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	if *recordDir != "" {
		fullPath, err := filepath.Abs(*recordDir)
		if err != nil {
			logger.Error("bad recording path: ", err)
		}
		logger.Info("FBS recording is turned on, writing to dir: ", fullPath)
		proxy.RecordingDir = fullPath
		proxy.SingleSession.Type = vncproxy.SessionTypeRecordingProxy
	} else {
		logger.Info("FBS recording is turned off")
	}
	if *reMarkable != "" {
		logger.Info("reMarkable 2.10+ support turned on")
		proxy.Remarkable = true
	}

	proxy.StartListening()
}
