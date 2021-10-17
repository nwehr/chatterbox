package chatterbox

type Password string

func (p Password) String() string {
	return string(p)
}
