package proxy

type SessionStatus int
type SessionType int

const (
	SessionStatusInit SessionStatus = iota
	SessionStatusActive
	SessionStatusError
)

const (
	SessionTypeRecordingProxy SessionType = iota
	SessionTypeReplayServer
	SessionTypeProxyPass
)

type VncSession struct {
	TargetHostname string
	TargetPort     string
	TargetPassword string
	ID             string
	Status         SessionStatus
	Type           SessionType
	ReplayFilePath string
}
