package proxy

import "testing"

func TestProxy(t *testing.T) {
	//create default session if required

	proxy := &VncProxy{
		WsListeningUrl:  "http://localhost:8183/", // empty = not listening on ws
		RecordingDir:    "",                       //"/Users/amitbet/vncRec",  // empty = no recording
		TcpListeningUrl: ":5905",
		//recordingDir:          "C:\\vncRec", // empty = no recording
		ProxyVncPassword: "", //empty = no auth
		SingleSession: &VncSession{
			TargetHostname: "192.168.1.101",
			TargetPort:     "5900",
			TargetPassword: "ancient1", //"Ch_#!T@8", //
			ID:             "dummySession",
			Status:         SessionStatusInit,
			Type:           SessionTypeRecordingProxy,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	proxy.StartListening()
}
