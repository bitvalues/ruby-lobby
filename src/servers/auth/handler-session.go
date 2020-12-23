package auth

import (
	"bytes"
	"io"
	"strings"

	"github.com/bitvalues/ruby-lobby/src/sessions"
)

func (s *AuthServer) handleSession(session sessions.Session) {
	buffer := make([]byte, 33)

	for {
		// try reading a chunk of data
		_, err := session.GetAuthSocket().Read(buffer)
		if err != nil {
			// if the socket closed, just break
			if err == io.EOF {
				break
			}

			// if the socket was forcibly closed, just break
			if strings.Contains(err.Error(), "connection was forcibly closed") {
				break
			}

			// otherwise, log the error
			s.log.WithError(err).Debug("could not read buffer from connection")
			break
		}

		// extract information from the buffer
		username := strings.TrimSpace(string(bytes.Trim(buffer[:16], "\x00")))
		password := strings.TrimSpace(string(bytes.Trim(buffer[16:32], "\x00")))
		action := int8(buffer[32])

		switch action {
		case ActionLogin:
			s.handleActionLogin(session, username, password)
		case ActionCreateAccount:
			s.handleActionCreateAccount(session, username, password)
		case ActionPasswordChange:
			s.handleActionPasswordChange(session, username, password)
		}
	}
}
