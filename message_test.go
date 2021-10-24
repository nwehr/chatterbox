package chatterbox

import (
	"bytes"
	"encoding/base64"
	"io"
	"testing"

	"filippo.io/age"
)

func TestEncryptDecrypt(t *testing.T) {
	encoded := &bytes.Buffer{}

	{
		publicKey := "age1cy0su9fwf3gf9mw868g5yut09p6nytfmmnktexz2ya5uqg9vl9sss4euqm"
		recipient, err := age.ParseX25519Recipient(publicKey)
		if err != nil {
			t.Errorf("Failed to parse public key %q: %v", publicKey, err)
		}

		encoder := base64.NewEncoder(base64.RawStdEncoding, encoded)
		encrypter, err := age.Encrypt(encoder, recipient)
		if err != nil {
			t.Errorf("Failed to create encrypted file: %v", err)
		}

		if _, err := io.WriteString(encrypter, "Live free or die!"); err != nil {
			t.Errorf("Failed to write to encrypted file: %v", err)
		}

		encrypter.Close()
		encoder.Close()
	}

	{
		identity, err := age.ParseX25519Identity("AGE-SECRET-KEY-184JMZMVQH3E6U0PSL869004Y3U2NYV7R30EU99CSEDNPH02YUVFSZW44VU")
		if err != nil {
			t.Errorf("Failed to parse private key: %v", err)
		}

		decoder := base64.NewDecoder(base64.RawStdEncoding, encoded)
		decrypter, err := age.Decrypt(decoder, identity)
		if err != nil {
			t.Errorf("Failed to open encrypted file: %v", err)
		}

		decrypted := &bytes.Buffer{}
		if _, err := io.Copy(decrypted, decrypter); err != nil {
			t.Errorf("Failed to read encrypted file: %v", err)
		}

		if decryptedStr := decrypted.String(); decryptedStr != "Live free or die!" {
			t.Errorf("unexpected result: %s", decryptedStr)
		}
	}
}

func TestMsg(t *testing.T) {
	buf := new(bytes.Buffer)

	{
		msg := Msg(Identity("@nate.errorcodelio"), []Identity{"@kevpatt.errorcode.io", "@nate.errorcode.io"}, "Hello, World!")
		if _, err := msg.WriteTo(buf); err != nil {
			t.Error(err)
		}
	}

	{
		msg := Message{}
		if _, err := msg.ReadFrom(buf); err != nil {
			t.Error(err)
		}

		if msg.Type != "MSG" {
			t.Errorf("expected MSG; got %s", msg.Type)
		}

		if msg.Args["Recipients"][0] != "@kevpatt.errorcode.io" {
			t.Errorf("'%s' != '%s'", msg.Args["Recipients"][0], "@kevpatt.errorcode.io")
		}

		identity, err := age.ParseX25519Identity("AGE-SECRET-KEY-1SG5WQSCUEPGY9ZZZ0M74AEK6DQDRP774ZHLVXS4662YAHUQLWS8SCJCR0V")
		if err != nil {
			t.Error(err)
		}

		decoder := base64.NewDecoder(base64.RawStdEncoding, bytes.NewReader(msg.Data))
		decrypter, err := age.Decrypt(decoder, identity)
		if err != nil {
			t.Error(err)
		}

		decrypted := &bytes.Buffer{}
		io.Copy(decrypted, decrypter)

		if decryptedStr := decrypted.String(); decryptedStr != "Hello, World!" {
			t.Errorf("'%s' != '%s'", decryptedStr, "Hello, World!")
		}
	}
}

func TestLConv(t *testing.T) {
	buf := new(bytes.Buffer)

	{
		msg := LConv()
		if _, err := msg.WriteTo(buf); err != nil {
			t.Error(err)
		}
	}

	{
		msg := Message{}
		if _, err := msg.ReadFrom(buf); err != nil {
			t.Error(err)
		}

		if msg.Type != "LCONV" {
			t.Errorf("expected LCONV; got %s", msg.Type)
		}
	}
}

func TestConv(t *testing.T) {
	buf := new(bytes.Buffer)

	{
		recipients := Identities{"@kevpatt.errorcode.io", "@nate.errorcode.io"}

		msg := Conv(recipients.ConversationID(), recipients, 3)
		if _, err := msg.WriteTo(buf); err != nil {
			t.Error(err)
		}
	}

	{
		msg := Message{}
		if _, err := msg.ReadFrom(buf); err != nil {
			t.Error(err)
		}

		if msg.Type != "CONV" {
			t.Errorf("expected CONV; got %s", msg.Type)
		}

		if msg.Args["Id"][0] != "a83087325cec029b2c39464aa0d2ea94f9282531cb53ad41143ea8bcea78205a" {
			t.Errorf("'%s' != '%s'", msg.Args["Id"][0], "a83087325cec029b2c39464aa0d2ea94f9282531cb53ad41143ea8bcea78205a")
		}

		if msg.Args["Unread"][0] != "3" {
			t.Errorf("'%s' != '%s'", msg.Args["Unread"][0], "3")
		}
	}
}
