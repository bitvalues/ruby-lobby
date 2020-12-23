package sessions

import (
	"time"

	"github.com/sirupsen/logrus"
)

type SessionManager struct {
	sessions []Session
	log      *logrus.Entry
}

func CreateNewSessionManager(log *logrus.Logger) *SessionManager {
	return &SessionManager{
		sessions: []Session{},
		log: log.WithFields(logrus.Fields{
			"component": "sessionManager",
		}),
	}
}

func (sm *SessionManager) AddSession(session Session) {
	sm.log.WithFields(logrus.Fields{
		"id": session.GetID(),
		"ip": session.GetIP(),
	}).Debug("Adding session")

	sm.sessions = append(sm.sessions, session)
}

func (sm *SessionManager) RemoveSession(session Session) {
	sm.log.WithFields(logrus.Fields{
		"id": session.GetID(),
		"ip": session.GetIP(),
	}).Debug("Removing session")

	// create a new array of sessions
	sessions := []Session{}

	// iterate through the existing sessions, only adding ones that don't
	// match the session we're avoiding
	for _, tmp := range sm.sessions {
		if tmp == session {
			continue
		}

		sessions = append(sessions, tmp)
	}

	// re-assign the sessions
	sm.sessions = sessions
}

func (sm *SessionManager) GetSessionByIP(ip string) *Session {
	if len(sm.sessions) == 0 {
		return nil
	}

	for _, session := range sm.sessions {
		if session.GetIP() == ip {
			return &session
		}
	}

	return nil
}

func (sm *SessionManager) StartCleanup() {
	for {
		// wait for 5 seconds
		time.Sleep(time.Second * 1)

		// make sure there are sessions to be cleaned up
		if len(sm.sessions) == 0 {
			continue
		}

		// purge old sessions
		sessions := []Session{}

		// iterate through the existing sessions, capturing ones that are
		// still relevant
		for _, session := range sm.sessions {
			if time.Now().Sub(session.GetLastActivity()).Seconds() < MaxSessionAge {
				sessions = append(sessions, session)
			} else {
				sm.log.WithFields(logrus.Fields{
					"id": session.GetID(),
					"ip": session.GetIP(),
				}).Debug("Removing old session")

				// close out the auth socket if it exists
				if socket := session.GetAuthSocket(); socket != nil {
					socket.Close()
				}

				// close out the data socket if it exists
				if socket := session.GetDataSocket(); socket != nil {
					socket.Close()
				}

				// close out the view socket if it exists
				if socket := session.GetViewSocket(); socket != nil {
					socket.Close()
				}
			}
		}

		// re-assign the sessions that are still active
		sm.sessions = sessions
	}
}
