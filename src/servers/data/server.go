package data

import (
	"fmt"
	"net"
	"strings"

	"github.com/bitvalues/ruby-lobby/src/config"
	"github.com/bitvalues/ruby-lobby/src/sessions"
	"github.com/sirupsen/logrus"
)

type DataServer struct {
	cfg    config.Config
	sm     *sessions.SessionManager
	log    *logrus.Entry
	socket net.Listener
}

func NewServer(cfg config.Config, sm *sessions.SessionManager, log *logrus.Logger) DataServer {
	socket, err := net.Listen("tcp4", fmt.Sprintf(":%d", cfg.DataServerPort))
	if err != nil {
		log.WithError(err).Fatal("Could not create data server listener")
	}

	return DataServer{
		cfg: cfg,
		sm:  sm,
		log: log.WithFields(logrus.Fields{
			"server": "data",
			"port":   cfg.DataServerPort,
		}),
		socket: socket,
	}
}

func (s *DataServer) Startup() {
	s.log.Info("Starting up")

	for {
		// attempt to accept a connection
		conn, err := s.socket.Accept()
		if err != nil {
			// prevent the closed network connection spam
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			// otherwise debug what the issue is
			s.log.WithError(err).Debug("Cannot accept new connection")
			continue
		}

		s.getSessionFromConnection(conn)
	}
}

func (s *DataServer) Shutdown() {
	s.log.Info("Shutting down")
	s.socket.Close()
}

func (s *DataServer) getSessionFromConnection(conn net.Conn) sessions.Session {
	return sessions.Session{}
}
