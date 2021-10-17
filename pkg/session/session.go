package session

import (
	"net"

	"github.com/nwehr/chatterbox"
)

type Session struct {
	ID       uint
	Identity chatterbox.Identity
	Conn     net.Conn
}

func NewSession(ident chatterbox.Identity, conn net.Conn) *Session {
	s := &Session{
		ID:       nextID(),
		Identity: ident,
		Conn:     conn,
	}

	sessions = append(sessions, s)

	return s
}

func QuitSession(id uint) {
	nextSessions := []*Session{}

	for _, s := range sessions {
		if s.ID != id {
			nextSessions = append(nextSessions, s)
		}
	}

	sessions = nextSessions
}

func SessionsForIdent(ident chatterbox.Identity) []*Session {
	identSessions := []*Session{}

	for _, s := range sessions {
		if s.Identity == ident {
			identSessions = append(identSessions, s)
		}
	}

	return identSessions
}

var id uint

func nextID() uint {
	id += 1
	return id
}

var sessions []*Session
