package chatterbox

import (
	"bytes"
	"testing"
	// "github.com/nwehr/chatterbox"
)

// func TestRequestMainLine(t *testing.T) {
// 	r := Request{
// 		Type: "LOGIN",
// 		Args: map[string]string{},
// 	}

// 	r.Args.Set("Identity", "@nate.errorcode.io")
// 	r.Args.Set("Password", "abc123")

// 	expected := "LOGIN Identity=@nate.errorcode.io Password=abc123"

// 	if r.mainLine() != expected {
// 		t.Errorf("unexpected result: %s; got %s", r.mainLine(), expected)
// 	}
// }

// func TestLoginMainLine(t *testing.T) {
// 	r := LoginRequest(Identity("@nate.errorcode.io"), "abc123")

// 	expected := "LOGIN Identity=@nate.errorcode.io Password=abc123"

// 	if r.mainLine() != expected {
// 		t.Errorf("unexpected result: %s; got %s", r.mainLine(), expected)
// 	}
// }

// func TestSendMainLine(t *testing.T) {
// 	r := SendRequest(Identity("@nate.errorcode.io"), "Hello, World")

// 	expected := "SEND To=@nate.errorcode.io"

// 	if r.mainLine() != expected {
// 		t.Errorf("unexpected result: %s; got %s", r.mainLine(), expected)
// 	}
// }

// func TestReqestParseMainLine(t *testing.T) {
// 	{
// 		mainLine := []byte("LOGIN Identity=@nate Password=abc123")

// 		req := Request{}
// 		req.parseMainLine(mainLine)

// 		if req.Type != "LOGIN" {
// 			t.Errorf("expected type LOGIN; got %s", req.Type)
// 		}
// 	}

// 	{
// 		mainLine := []byte("SEND To=@nate.errorcode.io")

// 		req := Request{}
// 		req.parseMainLine(mainLine)

// 		if req.Type != "SEND" {
// 			t.Errorf("expected type LOGIN; got %s", req.Type)
// 		}

// 		if req.Args["To"] != "@nate.errorcode.io" {
// 			t.Errorf("expected type @nate.errorcode.io; got %s", req.Type)
// 		}
// 	}
// }

func TestSendRead(t *testing.T) {
	buf := []byte("SEND To=@nate.errorcode.io From=@nate.errorcode.io Length=11\nHello World\n")
	rd := bytes.NewReader(buf)

	msg := Message{}
	msg.Read(rd)

	if msg.Type != "SEND" {
		t.Errorf("'%s' != '%s'", msg.Type, "SEND")
	}

	if msg.Args["To"][0] != "@nate.errorcode.io" {
		t.Errorf("'%s' != '%s'", msg.Args["To"][0], "@nate.errorcode.io")
	}

	// if req.Args["To"][1] != "@kevpatt.errorcode.io" {
	// 	t.Errorf("'%s' != '%s'", req.Args["To"][1], "@evpatt.errorcode.io")
	// }

	if string(msg.Data) != "Hello World" {
		t.Errorf("'%s' != '%s'", string(msg.Data), "Hello World")
	}
}

func TestSendWrite(t *testing.T) {
	buf := make([]byte, 512)
	w := bytes.NewBuffer(buf)

	msg := Send(Identity("@nate.errorcode.io"), []Identity{"@nate.errorcode.io", "@kevpatt.errorcode.io"}, "Hello, world!")
	msg.Write(w)

	if w.String() != "" {
		// t.Error(w.String())
	}
}
