package auth

import "github.com/bitvalues/ruby-lobby/src/sessions"

func (s *AuthServer) handleActionPasswordChange(session sessions.Session, username, password string) {
	s.log.WithField("sessionID", session.GetID()).Debug("Attempting to change password")
}
