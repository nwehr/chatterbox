package message

import (
	"github.com/nwehr/chatterbox/pkg/identity"
)

type Message struct {
	From identity.Identity
	To   []identity.Identity
	Body string
}
