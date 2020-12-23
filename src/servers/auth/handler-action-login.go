package auth

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/bitvalues/ruby-lobby/src/sessions"
	"github.com/bitvalues/ruby-lobby/src/tools"
	"github.com/davecgh/go-spew/spew"
	uuid "github.com/satori/go.uuid"
)

func (s *AuthServer) handleActionLogin(session sessions.Session, username, password string) {
	logger := s.log.WithField("sessionID", session.GetID())
	logger.Debug("Attempting to login")

	// declare an empty response
	var response []byte

	// attempt to login
	if accountID, found := tools.GetUserAccountID(username, password); found {
		// debugging
		logger.WithField("accountID", accountID).Debug("Login succeeded")

		// set the session account ID
		session.SetAccountID(accountID)

		// create the packet to be sent back
		response = make([]byte, 33)
		response[0] = ActionLoginSucceeded
		copy(response[1:], tools.Uint32ToByteArray(accountID))

		// generate a session hash and copy it into the packet
		hash := md5.Sum(uuid.NewV4().Bytes())
		copy(response[5:], hash[:])
		spew.Dump("GENERATED HASH", hex.EncodeToString(hash[:]))

	} else {
		// debugging
		logger.Debug("Login failed")

		// create the packet to be sent back
		response = []byte{ActionLoginFailed}
	}

	// finally, write the response to the socket
	session.GetAuthSocket().Write(response)
}
