package chatterbox

import (
	"testing"
)

func TestRequestMainLine(t *testing.T) {
	r := Request{
		Type: "LOGIN",
		Args: map[string]string{},
	}

	r.Args.Set("Identity", "@nate.errorcode.io")
	r.Args.Set("Password", "abc123")

	expected := "LOGIN Identity=@nate.errorcode.io Password=abc123"

	if r.mainLine() != expected {
		t.Errorf("unexpected result: %s; got %s", r.mainLine(), expected)
	}
}

func TestLoginMainLine(t *testing.T) {
	r := LoginRequest(Identity("@nate.errorcode.io"), "abc123")

	expected := "LOGIN Identity=@nate.errorcode.io Password=abc123"

	if r.mainLine() != expected {
		t.Errorf("unexpected result: %s; got %s", r.mainLine(), expected)
	}
}

func TestReqestParseMainLine(t *testing.T) {
	mainLine := []byte("LOGIN Identity=@nate Password=abc123")

	req := Request{}
	req.parseMainLine(mainLine)

	if req.Type != "LOGIN" {
		t.Errorf("expected type LOGIN; got %s", req.Type)
	}
}
