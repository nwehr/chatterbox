package chatterbox

import (
	"bytes"
	"testing"
	// "github.com/nwehr/chatterbox"
)

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

	if string(msg.Data) != "Hello World" {
		t.Errorf("'%s' != '%s'", string(msg.Data), "Hello World")
	}
}

func TestSendWrite(t *testing.T) {
	buf := make([]byte, 512)
	w := bytes.NewBuffer(buf)

	msg := Send(Identity("@nate.errorcode.io"), []Identity{"@nate.errorcode.io", "@kevpatt.errorcode.io"}, "Hello, World!")
	msg.Write(w)

	// expected := "SEND To=@nate.errorcode.io;@kevpatt.errorcode.io From=@nate.errorcode.io Length=13\nHello, World!\n"

	// if len(w.String()) != len(expected) {
	// 	t.Errorf("expected %d; got %d", len(strings.TrimSpace(expected)), len(strings.TrimSpace(w.String())))
	// }

	// if w.String() != expected {
	// 	t.Errorf("got '%s'; expected '%s'", w.String(), expected)
	// }
}
