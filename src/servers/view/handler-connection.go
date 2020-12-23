package view

import (
	"io"
	"net"
	"strings"

	"github.com/bitvalues/ruby-lobby/src/sessions"
	"github.com/bitvalues/ruby-lobby/src/tools"
	"github.com/sirupsen/logrus"
)

func (s *ViewServer) handleConnection(conn net.Conn) {
	// extract the ip from the remote address
	ip, _ := tools.ParseRemoteAddress(conn.RemoteAddr().String())

	// check if the ip address is banned
	if tools.IPAddressIsBanned(ip) {
		s.log.WithField("ip", ip).Info("Dropping banned connection")
		conn.Close()
		return
	}

	// create an empty session holder
	var session sessions.Session

	// create a buffer for receiving data
	buffer := make([]byte, 1024)

	// debugging
	s.log.WithField("ip", ip).Debug("New connection established")

	// attempt to read data from the connection
	for {
		// attempt to get a packet fromt the connection
		len, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else if strings.Contains(err.Error(), "connection was forcibly closed") {
				break
			}

			s.log.WithError(err).Debug("Could not read buffer from connection")
			break
		}

		// make sure the packet is over 9 bytes
		if len < 9 {
			continue
		}

		// parse the code from the packet
		code := buffer[8]

		// debugging
		s.log.WithFields(logrus.Fields{
			"size": len,
			"code": code,
		}).Debug("Read packet from connection")

		// react to the request
		switch code {
		case 0x26:
			s.handleVersionCheck(buffer, &session)
		}
	}
}
