package identity

import "strings"

type Identity string

func (ident Identity) Host() string {
	return strings.Join(strings.Split(string(ident), ".")[1:], ".")
}
