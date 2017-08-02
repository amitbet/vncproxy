package proxy

import "testing"

func TestProxy(t *testing.T) {
	//create default session if required

	proxy := &VncProxy{
		WsListeningUrl:  "http://localhost:7777/", // empty = not listening on ws
		RecordingDir:    "/Users/amitbet/vncRec",  // empty = no recording
		TcpListeningUrl: ":5904",
		//recordingDir:          "C:\\vncRec", // empty = no recording
		ProxyVncPassword: "1234", //empty = no auth
		SingleSession: &VncSession{
			TargetHostname: "localhost",
			TargetPort:     "5903",
			TargetPassword: "Ch_#!T@8",
			ID:             "dummySession",
			Status:         SessionStatusInit,
			Type:           SessionTypeRecordingProxy,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	proxy.StartListening()
}
