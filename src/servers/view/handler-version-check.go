package view

import (
	"encoding/hex"

	"github.com/bitvalues/ruby-lobby/src/sessions"
	"github.com/davecgh/go-spew/spew"
)

func (s *ViewServer) handleVersionCheck(buffer []byte, session *sessions.Session) {
	// make sure the buffer is of correct size
	if len(buffer) < 28 {
		return
	}

	// debugging
	spew.Dump("PARSED HASH", hex.EncodeToString(buffer[12:28]))
}
