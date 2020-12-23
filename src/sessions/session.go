package sessions

import (
	"net"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	id           uuid.UUID
	ip           string
	lastActivity time.Time

	authSocket net.Conn
	dataSocket net.Conn
	viewSocket net.Conn
}

func CreateNewSession(ip string) Session {
	session := Session{
		id:           uuid.NewV4(),
		ip:           ip,
		lastActivity: time.Now(),
	}

	return session
}

func (s *Session) GetID() string {
	return s.id.String()
}

func (s *Session) GetIP() string {
	return s.ip
}

func (s *Session) GetLastActivity() time.Time {
	return s.lastActivity
}

func (s *Session) SetAuthSocket(socket net.Conn) {
	s.authSocket = socket
}

func (s *Session) GetAuthSocket() net.Conn {
	return s.authSocket
}

func (s *Session) SetDataSocket(socket net.Conn) {
	s.dataSocket = socket
}

func (s *Session) GetDataSocket() net.Conn {
	return s.dataSocket
}

func (s *Session) SetViewSocket(socket net.Conn) {
	s.viewSocket = socket
}

func (s *Session) GetViewSocket() net.Conn {
	return s.viewSocket
}
