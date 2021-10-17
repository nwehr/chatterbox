package chatterbox

import "strings"

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
