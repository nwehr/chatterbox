package chatterbox

import "testing"

func TestIdentityHost(t *testing.T) {
	ident := Identity("@nate.errorcode.io")

	if ident.Host() != "errorcode.io" {
		t.Errorf("expected 'errorcodelio'; got '%s'", ident.Host())
	}

	if ident.Name() != "@nate" {
		t.Errorf("expected 'nate'; got '%s'", ident.Name())
	}
}
