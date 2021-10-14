package session

import (
	"net"

	"github.com/nwehr/chatterbox/pkg/identity"
)

type Session struct {
	Identity identity.Identity
	Conn     net.Conn
}

// func Start(conn net.Conn) Session {
// 	conn.
// }
