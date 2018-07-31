package main

import "vncproxy/proxy"
import "flag"
import "vncproxy/logger"
import "os"

func main() {
	//create default session if required
	var tcpPort = flag.String("tcpPort", "", "tcp port")
	var wsPort = flag.String("wsPort", "", "websocket port")
	var vncPass = flag.String("vncPass", "", "password on incoming vnc connections to the proxy, defaults to no password")
	var recordDir = flag.String("recDir", "", "path to save FBS recordings WILL NOT RECORD if not defined.")
	var targetVncPort = flag.String("targPort", "", "target vnc server port")
	var targetVncHost = flag.String("targHost", "", "target vnc server host")
	var targetVncPass = flag.String("targPass", "", "target vnc password")

	flag.Parse()

	if *tcpPort == "" && *wsPort == "" {
		logger.Error("no listening port defined")
		flag.Usage()
		os.Exit(1)
	}

	if *targetVncPort == "" {
		logger.Error("no target vnc server port defined")
		flag.Usage()
		os.Exit(1)
	}

	if *vncPass == "" {
		logger.Warn("proxy will have no password")
	}
	if *recordDir == "" {
		logger.Warn("FBS recording is turned off")
	}

	tcpUrl := ""
	if *tcpPort != "" {
		tcpUrl = ":" + string(*tcpPort)
	}

	proxy := &proxy.VncProxy{
		WsListeningUrl:   "http://0.0.0.0:" + string(*wsPort) + "/", // empty = not listening on ws
		RecordingDir:     *recordDir,                                //"/Users/amitbet/vncRec",                     // empty = no recording
		TcpListeningUrl:  tcpUrl,
		ProxyVncPassword: *vncPass, //empty = no auth
		SingleSession: &proxy.VncSession{
			TargetHostname: *targetVncHost,
			TargetPort:     *targetVncPort,
			TargetPassword: *targetVncPass, //"vncPass",
			ID:             "dummySession",
			Status:         proxy.SessionStatusInit,
			Type:           proxy.SessionTypeRecordingProxy,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	proxy.StartListening()
}
