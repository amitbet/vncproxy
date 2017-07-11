package proxy

import "testing"

func TestProxy(t *testing.T) {
	//create default session if required

	proxy := &VncProxy{
		wsListeningUrl:        "http://localhost:7777/", // empty = not listening on ws
		recordingDir:          "/Users/amitbet/vncRec",  // empty = no recording
		targetServersPassword: "Ch_#!T@8",               //empty = no auth
		SingleSession: &VncSession{
			TargetHostname: "localhost",
			TargetPort:     "5903",
			TargetPassword: "vncPass",
			ID:             "dummySession",
			Status:         SessionStatusActive,
			Type:           SessionTypeRecordingProxy,
		}, // to be used when not using sessions
		UsingSessions: false, //false = single session - defined in the var above
	}

	proxy.StartListening()
}
