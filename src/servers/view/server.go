package view

import (
	"fmt"
	"net"
	"strings"

	"github.com/bitvalues/ruby-lobby/src/config"
	"github.com/bitvalues/ruby-lobby/src/sessions"
	"github.com/sirupsen/logrus"
)

type ViewServer struct {
	cfg    config.Config
	sm     *sessions.SessionManager
	log    *logrus.Entry
	socket net.Listener
}

func NewServer(cfg config.Config, sm *sessions.SessionManager, log *logrus.Logger) ViewServer {
	socket, err := net.Listen("tcp4", fmt.Sprintf(":%d", cfg.ViewServerPort))
	if err != nil {
		log.WithError(err).Fatal("Could not create view server listener")
	}

	return ViewServer{
		cfg: cfg,
		sm:  sm,
		log: log.WithFields(logrus.Fields{
			"server": "view",
			"port":   cfg.ViewServerPort,
		}),
		socket: socket,
	}
}

func (s *ViewServer) Startup() {
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

func (s *ViewServer) Shutdown() {
	s.log.Info("Shutting down")
	s.socket.Close()
}

func (s *ViewServer) getSessionFromConnection(conn net.Conn) sessions.Session {
	return sessions.Session{}
}
