package vnc_rec

type SessionManager struct {
	sessions map[string]*VncSession
}

func (s *SessionManager) GetSession(sessionId string) (*VncSession, error) {
	return s.sessions[sessionId], nil
}

func (s *SessionManager) SetSession(sessionId string, session *VncSession) error {
	s.sessions[sessionId] = session
	return nil
}

func (s *SessionManager) DeleteSession(sessionId string) error {
	delete(s.sessions, sessionId)
	return nil
}
