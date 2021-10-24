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

func TestConversationID(t *testing.T) {
	id := Identities{"@nate.errorcode.io", "@kevpatt.errorcode.io"}.ConversationID()
	if id != "a83087325cec029b2c39464aa0d2ea94f9282531cb53ad41143ea8bcea78205a" {
		t.Errorf("unexpected hash value %s", id)
	}
}
