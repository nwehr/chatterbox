package chatterbox

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

type Identity string

func (ident Identity) Host() string {
	return strings.Join(strings.Split(string(ident), ".")[1:], ".")
}

func (ident Identity) Name() string {
	return strings.Split(string(ident), ".")[0]
}

func (ident Identity) String() string {
	return string(ident)
}

type Identities []Identity

func NewIdentities(strs []string) Identities {
	idents := Identities{}

	for _, str := range strs {
		idents = append(idents, Identity(str))
	}

	return idents
}

func (i Identities) Strings() []string {
	strs := []string{}

	for _, ident := range i {
		strs = append(strs, ident.String())
	}

	sort.Strings(strs)

	return strs
}

func (i Identities) ConversationID() string {
	h := sha256.New()
	h.Write([]byte(strings.TrimSpace(strings.Join(i.Strings(), ";"))))

	return fmt.Sprintf("%x", h.Sum(nil))
}
