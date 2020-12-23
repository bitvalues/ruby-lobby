package auth

import (
	"github.com/bitvalues/ruby-lobby/src/sessions"
	"github.com/bitvalues/ruby-lobby/src/tools"
)

func (s *AuthServer) handleActionLogin(session sessions.Session, username, password string) {
	logger := s.log.WithField("sessionID", session.GetID())
	logger.Debug("Attempting to login")

	// declare an empty response
	var response []byte

	// attempt to login
	if accountID, found := s.attemptToLogin(username, password); found {
		// succeeded
		logger.WithField("accountID", accountID).Debug("Login succeeded")
		response = make([]byte, 33)
		response[0] = ActionLoginSucceeded
		copy(response[1:], tools.Uint32ToByteArray(accountID))
	} else {
		// failed
		logger.Debug("Login failed")
		response = []byte{ActionLoginFailed}
	}

	// finally, write the response to the socket
	session.GetAuthSocket().Write(response)
}

func (s *AuthServer) attemptToLogin(username, password string) (uint32, bool) {
	// TODO: Actually implement logic for this
	return 999, true
}
